package domain

import "context"

type ProductUsecase interface {
	Create(ctx context.Context, product *Product) error
	GetAll(ctx context.Context) ([]Product, error)
	GetBySKU(ctx context.Context, sku string) (*Product, error)
	Update(ctx context.Context, sku string, product *Product) (*Product, error)
	Delete(ctx context.Context, sku string) error
	Purchase(ctx context.Context, sku string, qty int) error
}
