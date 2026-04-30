package migration

import (
	"database/sql"

	"Inventory-Management-Erajaya/pkg/logger"
)

func Run(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS products (
		id         BIGINT       PRIMARY KEY,
		sku        VARCHAR(100) NOT NULL UNIQUE,
		name       VARCHAR(255) NOT NULL,
		qty        INTEGER      NOT NULL DEFAULT 0 CHECK (qty >= 0),
		price      NUMERIC(15,2) NOT NULL CHECK (price > 0)
	);

	CREATE INDEX IF NOT EXISTS idx_products_sku ON products (sku);
	`

	_, err := db.Exec(query)
	if err != nil {
		logger.Log.Fatalf("migration failed: %v", err)
	}
	logger.Log.Info("migration completed successfully")
}
