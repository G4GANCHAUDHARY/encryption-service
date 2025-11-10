package providers

import (
	"context"
	redis "github.com/redis/go-redis/v9"
	"time"
)

type RedisLib struct {
	client *redis.Client
}

func GetRedisClient(config AppConfig) *RedisLib {
	rdb := redis.NewClient(&redis.Options{
		Addr: config.RedisConfig.Address,
		DB:   config.RedisConfig.DbNumber,
	})
	return &RedisLib{client: rdb}
}

func (r *RedisLib) Close() error {
	if r.client != nil {
		return r.client.Close()
	}
	return nil
}

func (r *RedisLib) Set(ctx context.Context, key string, value string) error {
	return r.client.Set(ctx, key, value, 0).Err()
}

func (r *RedisLib) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *RedisLib) Increment(ctx context.Context, key string) (int64, error) {
	return r.client.Incr(ctx, key).Result()
}

func (r *RedisLib) Delete(ctx context.Context, key string) (int64, error) {
	return r.client.Del(ctx, key).Result()
}

func (r *RedisLib) Expire(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	return r.client.Expire(ctx, key, expiration).Result()
}
