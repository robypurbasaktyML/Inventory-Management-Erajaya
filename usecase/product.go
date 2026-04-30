package usecase

import (
	"context"
	"math/rand"

	"Inventory-Management-Erajaya/domain"
	"Inventory-Management-Erajaya/pkg/logger"
)

type productUsecase struct {
	repo  domain.ProductRepository
	cache domain.StockCache
}

func NewProductUsecase(repo domain.ProductRepository, cache domain.StockCache) domain.ProductUsecase {
	return &productUsecase{repo: repo, cache: cache}
}

func generateID() uint {
	return uint(rand.Uint32())
}

func (u *productUsecase) Create(ctx context.Context, product *domain.Product) error {
	if err := product.Validate(); err != nil {
		return err
	}
	product.ID = generateID()
	if err := u.repo.Create(ctx, product); err != nil {
		return err
	}
	_ = u.cache.InitStock(ctx, product.SKU, product.Qty)
	logger.Log.WithField("sku", product.SKU).Info("product created")
	return nil
}

func (u *productUsecase) GetAll(ctx context.Context) ([]domain.Product, error) {
	return u.repo.GetAll(ctx)
}

func (u *productUsecase) GetBySKU(ctx context.Context, sku string) (*domain.Product, error) {
	return u.repo.GetBySKU(ctx, sku)
}

func (u *productUsecase) Update(ctx context.Context, sku string, product *domain.Product) (*domain.Product, error) {
	if err := product.Validate(); err != nil {
		return nil, err
	}
	if err := u.repo.Update(ctx, sku, product); err != nil {
		return nil, err
	}
	_ = u.cache.InitStock(ctx, sku, product.Qty)
	// Fetch updated product with id and sku
	updated, err := u.repo.GetBySKU(ctx, sku)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (u *productUsecase) Delete(ctx context.Context, sku string) error {
	if err := u.repo.Delete(ctx, sku); err != nil {
		return err
	}
	_ = u.cache.DeleteStock(ctx, sku)
	return nil
}

// Purchase handles flash sale concurrency.
// Flow:
//  1. Load stock ke Redis pakai SETNX (hanya request pertama yang set, sisanya skip)
//  2. Atomic decrement di Redis via Lua script (single-threaded, no race condition)
//  3. Hanya request yang lolos Redis yang persist ke DB
//  4. Rollback Redis jika DB gagal
func (u *productUsecase) Purchase(ctx context.Context, sku string, qty int) error {
	if qty <= 0 {
		return domain.ErrInvalidQty
	}

	// Ensure stock is loaded in cache using SETNX (race-safe)
	cached, err := u.cache.GetStock(ctx, sku)
	if err != nil {
		return err
	}
	if cached == -1 {
		product, err := u.repo.GetBySKU(ctx, sku)
		if err != nil {
			return err
		}
		// SETNX: only first request sets the value, others skip
		_, err = u.cache.InitStockIfNotExists(ctx, sku, product.Qty)
		if err != nil {
			return err
		}
	}

	// Step 1: Atomic decrement in Redis (Lua script, no race condition)
	_, err = u.cache.DecrementStock(ctx, sku, qty)
	if err != nil {
		logger.Log.WithFields(map[string]interface{}{"sku": sku, "qty": qty}).Warn("purchase rejected: " + err.Error())
		return err
	}

	// Step 2: Persist to database
	if err := u.repo.DecrementStock(ctx, sku, qty); err != nil {
		// Rollback: re-add stock to Redis
		_, _ = u.cache.DecrementStock(ctx, sku, -qty)
		logger.Log.WithFields(map[string]interface{}{"sku": sku, "qty": qty}).Error("db decrement failed, redis rolled back")
		return err
	}

	logger.Log.WithFields(map[string]interface{}{"sku": sku, "qty": qty}).Info("purchase successful")
	return nil
}
