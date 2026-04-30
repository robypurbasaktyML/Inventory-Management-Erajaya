package postgres

import (
	"context"
	"database/sql"

	"Inventory-Management/domain"
)

type productRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) domain.ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(ctx context.Context, product *domain.Product) error {
	query := `INSERT INTO products (id, sku, name, qty, price) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.ExecContext(ctx, query, product.ID, product.SKU, product.Name, product.Qty, product.Price)
	if err != nil {
		if isDuplicateError(err) {
			return domain.ErrDuplicateSKU
		}
		return err
	}
	return nil
}

func (r *productRepository) GetAll(ctx context.Context) ([]domain.Product, error) {
	query := `SELECT id, sku, name, qty, price FROM products ORDER BY id`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []domain.Product
	for rows.Next() {
		var p domain.Product
		if err := rows.Scan(&p.ID, &p.SKU, &p.Name, &p.Qty, &p.Price); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, rows.Err()
}

func (r *productRepository) GetBySKU(ctx context.Context, sku string) (*domain.Product, error) {
	query := `SELECT id, sku, name, qty, price FROM products WHERE sku = $1`
	var p domain.Product
	err := r.db.QueryRowContext(ctx, query, sku).Scan(&p.ID, &p.SKU, &p.Name, &p.Qty, &p.Price)
	if err == sql.ErrNoRows {
		return nil, domain.ErrProductNotFound
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *productRepository) Update(ctx context.Context, sku string, product *domain.Product) error {
	query := `UPDATE products SET name = $1, qty = $2, price = $3 WHERE sku = $4`
	result, err := r.db.ExecContext(ctx, query, product.Name, product.Qty, product.Price, sku)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return domain.ErrProductNotFound
	}
	return nil
}

func (r *productRepository) Delete(ctx context.Context, sku string) error {
	query := `DELETE FROM products WHERE sku = $1`
	result, err := r.db.ExecContext(ctx, query, sku)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return domain.ErrProductNotFound
	}
	return nil
}

func (r *productRepository) DecrementStock(ctx context.Context, sku string, qty int) error {
	query := `UPDATE products SET qty = qty - $1 WHERE sku = $2 AND qty >= $1`
	result, err := r.db.ExecContext(ctx, query, qty, sku)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return domain.ErrInsufficientStock
	}
	return nil
}
