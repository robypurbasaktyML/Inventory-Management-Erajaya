package redis

import (
	"context"
	"errors"
	"fmt"

	"Inventory-Management-Erajaya/domain"

	"github.com/redis/go-redis/v9"
)

type stockCache struct {
	client *redis.Client
}

func NewStockCache(client *redis.Client) domain.StockCache {
	return &stockCache{client: client}
}

func stockKey(sku string) string {
	return fmt.Sprintf("stock:%s", sku)
}

func (c *stockCache) InitStock(ctx context.Context, sku string, qty int) error {
	return c.client.Set(ctx, stockKey(sku), qty, 0).Err()
}

func (c *stockCache) InitStockIfNotExists(ctx context.Context, sku string, qty int) (bool, error) {
	ok, err := c.client.SetNX(ctx, stockKey(sku), qty, 0).Result()
	if err != nil {
		return false, err
	}
	return ok, nil
}

func (c *stockCache) DecrementStock(ctx context.Context, sku string, qty int) (int, error) {
	// Lua script for atomic decrement with floor check.
	// Returns new stock if success, -1 if insufficient.
	script := redis.NewScript(`
		local key = KEYS[1]
		local decr = tonumber(ARGV[1])
		local current = tonumber(redis.call('GET', key))
		if current == nil then
			return -2
		end
		if decr > 0 and current < decr then
			return -1
		end
		local newStock = current - decr
		redis.call('SET', key, newStock)
		return newStock
	`)

	result, err := script.Run(ctx, c.client, []string{stockKey(sku)}, qty).Int()
	if err != nil {
		return 0, err
	}
	if result == -2 {
		return 0, errors.New("stock not initialized in cache")
	}
	if result == -1 {
		return 0, domain.ErrInsufficientStock
	}
	return result, nil
}

func (c *stockCache) DeleteStock(ctx context.Context, sku string) error {
	return c.client.Del(ctx, stockKey(sku)).Err()
}

func (c *stockCache) GetStock(ctx context.Context, sku string) (int, error) {
	val, err := c.client.Get(ctx, stockKey(sku)).Int()
	if errors.Is(err, redis.Nil) {
		return -1, nil
	}
	if err != nil {
		return 0, err
	}
	return val, nil
}
