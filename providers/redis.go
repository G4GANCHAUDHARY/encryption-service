package providers

import (
	"context"
	redis "github.com/redis/go-redis/v9"
)

type RedisLib struct {
	client *redis.Client
}

func NewRedisWrapper() *RedisLib {
	return &RedisLib{client: InitRedis()}
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

func InitRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // default DB
	})

	return rdb
}
