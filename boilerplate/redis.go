package boilerplate

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisClient struct {
	RedisHost       string
	RedisPassword   string
	RedisDatabase   int // default is 0
	redisClient     *redis.Client
	RedisClientOnce sync.Once
}

// Validate redis configs
func (r *RedisClient) Validate() error {
	if r.RedisHost == "" {
		return errors.New("Missing redis host")
	}
	if r.RedisPassword == "" {
		return errors.New("Missing redis password")
	}
	return nil
}

// GetSetupRedisClient get redis client
func (r *RedisClient) GetSetupRedisClient() *redis.Client {
	r.RedisClientOnce.Do(func() {
		r.redisClient = redis.NewClient(&redis.Options{
			Addr:     r.RedisHost,
			Password: r.RedisPassword, // no password set
			DB:       r.RedisDatabase, // use default DB
		})
	})
	return r.redisClient
}

// SetKV set key value to redis
func (r *RedisClient) StoreKV(key string, value string, timeoutCtx context.Context, ttl time.Duration) error {
	err := r.GetSetupRedisClient().Set(timeoutCtx, key, value, ttl).Err()
	return err
}

// GetKV get key value from redis
func (r *RedisClient) GetKV(key string, timeoutCtx context.Context) (string, error) {
	val, err := r.GetSetupRedisClient().Get(timeoutCtx, key).Result()
	return val, err
}

// DeleteKV delete key value from redis
func (r *RedisClient) DeleteKV(key string, timeoutCtx context.Context) error {
	err := r.GetSetupRedisClient().Del(timeoutCtx, key).Err()
	if err == redis.Nil {
		return nil
	}
	return err
}
