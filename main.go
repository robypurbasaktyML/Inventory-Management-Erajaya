package main

import (
	"Inventory-Management-Erajaya/config"
	"Inventory-Management-Erajaya/migration"
	"Inventory-Management-Erajaya/pkg/logger"
	"context"
	"database/sql"
	"github.com/redis/go-redis/v9"
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

	log.Info("connected to PostgreSQL")

	// Run migration
	migration.Run(db)

	// Connect to Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("failed to connect to Redis: %v", err)
	}
	log.Info("connected to Redis")

}
