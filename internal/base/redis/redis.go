package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/Al-Khaimah/khaimah-golang-backend/config"
	"github.com/redis/go-redis/v9"
)

var (
	redisClient *redis.Client
)

// InitRedis initializes the Redis client
func InitRedis() error {
	redisAddr := config.GetEnv("REDIS_ADDR", "localhost:6379")
	redisPassword := config.GetEnv("REDIS_PASSWORD", "")
	redisDB := 0

	redisClient = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return nil
}

// GetRedisClient returns the Redis client
func GetRedisClient() *redis.Client {
	return redisClient
}

// SetWithTTL sets a key-value pair with a TTL
func SetWithTTL(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return redisClient.Set(ctx, key, value, ttl).Err()
}

// Get retrieves a value by key
func Get(ctx context.Context, key string) (string, error) {
	return redisClient.Get(ctx, key).Result()
}

// Delete removes a key
func Delete(ctx context.Context, key string) error {
	return redisClient.Del(ctx, key).Err()
}
