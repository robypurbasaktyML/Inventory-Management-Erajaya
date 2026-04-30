package domain

import "errors"

var (
	ErrProductNotFound   = errors.New("product not found")
	ErrDuplicateSKU      = errors.New("product with this SKU already exists")
	ErrInsufficientStock = errors.New("insufficient stock")
	ErrInvalidQty        = errors.New("qty cannot be negative")
	ErrInvalidPrice      = errors.New("price must be greater than 0")
)

type Product struct {
	ID    uint    `json:"id"`
	SKU   string  `json:"sku"`
	Name  string  `json:"name"`
	Qty   int     `json:"qty"`
	Price float64 `json:"price"`
}

func (p *Product) Validate() error {
	if p.Qty < 0 {
		return ErrInvalidQty
	}
	if p.Price <= 0 {
		return ErrInvalidPrice
	}
	return nil
}
