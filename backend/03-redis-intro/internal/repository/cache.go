package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/e6a5/learning/backend/03-redis-intro/internal/models"
)

// CacheRepository handles Redis cache operations
type CacheRepository struct {
	client *redis.Client
	ctx    context.Context
}

// NewCacheRepository creates a new cache repository
func NewCacheRepository(client *redis.Client) *CacheRepository {
	return &CacheRepository{
		client: client,
		ctx:    context.Background(),
	}
}

// Get retrieves a value from Redis by key
func (r *CacheRepository) Get(key string) (*models.KeyValue, error) {
	val, err := r.client.Get(r.ctx, key).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("key not found: %s", key)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get key %s: %w", key, err)
	}

	return models.NewKeyValue(key, val, 0), nil
}

// Set stores a key-value pair in Redis with optional TTL
func (r *CacheRepository) Set(key, value string, ttl int) error {
	var expiration time.Duration
	if ttl > 0 {
		expiration = time.Duration(ttl) * time.Second
	}

	err := r.client.Set(r.ctx, key, value, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to set key %s: %w", key, err)
	}

	return nil
}

// Delete removes a key from Redis
func (r *CacheRepository) Delete(key string) error {
	deleted, err := r.client.Del(r.ctx, key).Result()
	if err != nil {
		return fmt.Errorf("failed to delete key %s: %w", key, err)
	}

	if deleted == 0 {
		return fmt.Errorf("key not found: %s", key)
	}

	return nil
}

// GetAllKeys retrieves all keys matching a pattern
func (r *CacheRepository) GetAllKeys(pattern string) ([]string, error) {
	if pattern == "" {
		pattern = "*"
	}

	keys, err := r.client.Keys(r.ctx, pattern).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get keys with pattern %s: %w", pattern, err)
	}

	return keys, nil
}

// GetTTL returns the time to live for a key
func (r *CacheRepository) GetTTL(key string) (time.Duration, error) {
	ttl, err := r.client.TTL(r.ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get TTL for key %s: %w", key, err)
	}

	return ttl, nil
}

// SetExpire sets the TTL for an existing key
func (r *CacheRepository) SetExpire(key string, ttl int) error {
	success, err := r.client.Expire(r.ctx, key, time.Duration(ttl)*time.Second).Result()
	if err != nil {
		return fmt.Errorf("failed to set expire for key %s: %w", key, err)
	}

	if !success {
		return fmt.Errorf("key not found: %s", key)
	}

	return nil
}

// Ping checks if Redis is accessible
func (r *CacheRepository) Ping() error {
	_, err := r.client.Ping(r.ctx).Result()
	if err != nil {
		return fmt.Errorf("redis ping failed: %w", err)
	}

	return nil
}
