package main

import (
	"Inventory-Management-Erajaya/config"
	"Inventory-Management-Erajaya/logger"
	"database/sql"
)

func main() {
	cfg := config.Load()
	logger.Init()
	log := logger.Log

	// Connect to PostgreSQL
	db, err := sql.Open("postgres", cfg.DSN())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

}
