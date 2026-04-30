package domain

import "context"

type StockCache interface {
	// InitStock overwrites stock in cache (used for create/update).
	InitStock(ctx context.Context, sku string, qty int) error
	// InitStockIfNotExists sets stock only if not yet cached (race-safe for purchase).
	InitStockIfNotExists(ctx context.Context, sku string, qty int) (bool, error)
	// DecrementStock atomically decrements stock. Returns error if insufficient.
	DecrementStock(ctx context.Context, sku string, qty int) (int, error)
	DeleteStock(ctx context.Context, sku string) error
	GetStock(ctx context.Context, sku string) (int, error)
}
