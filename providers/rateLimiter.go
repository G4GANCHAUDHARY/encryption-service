package providers

import (
	"context"
	"time"
)

type RateLimiter struct {
	Redis *RedisLib
}

func GetRateLimiter(redis *RedisLib) *RateLimiter {
	return &RateLimiter{
		Redis: redis,
	}
}

func (rl *RateLimiter) Allow(capacity int, key string) bool {
	ctx := context.Background()

	// increase counter per ip
	count, err := rl.Redis.Increment(ctx, key)
	if err != nil {
		return false
	}

	// set key expiration window
	if count == 1 {
		_, _ = rl.Redis.Expire(ctx, key, 24*60*time.Minute)
	}

	// false if exceeding limit
	return int(count) <= capacity
}
