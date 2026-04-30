package main

import (
	"Inventory-Management-Erajaya/config"
	"Inventory-Management-Erajaya/delivery/http/handler"
	"Inventory-Management-Erajaya/delivery/http/router"
	"Inventory-Management-Erajaya/migration"
	"Inventory-Management-Erajaya/pkg/logger"
	repoPostgres "Inventory-Management-Erajaya/repository/postgres"
	repoRedis "Inventory-Management-Erajaya/repository/redis"
	"Inventory-Management-Erajaya/usecase"
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"

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

	// Initialize layers (Clean Architecture)
	productRepo := repoPostgres.NewProductRepository(db)
	stockCache := repoRedis.NewStockCache(rdb)
	productUsecase := usecase.NewProductUsecase(productRepo, stockCache)
	productHandler := handler.NewProductHandler(productUsecase)

	// Setup router and start server
	r := router.NewRouter(productHandler)
	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Infof("server running on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}

}
