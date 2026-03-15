package limiter

import (
	"context"
	"fmt"
	"time"

	"github.com/HelenaBlack/anti-bruteforce/internal/domain"
	"github.com/redis/go-redis/v9"
)

type RedisLimiter struct {
	client *redis.Client
}

func NewRedisLimiter(client *redis.Client) *RedisLimiter {
	return &RedisLimiter{client: client}
}

func (r *RedisLimiter) Allow(ctx context.Context, limitType domain.RateLimitType, key string, limit int) (bool, error) {
	redisKey := fmt.Sprintf("limiter:%s:%s", limitType, key)
	now := time.Now().Unix()
	window := int64(60) // 1 minute window

	pipe := r.client.Pipeline()
	// Remove old entries
	pipe.ZRemRangeByScore(ctx, redisKey, "0", fmt.Sprintf("%d", now-window))
	// Add current attempt
	pipe.ZAdd(ctx, redisKey, redis.Z{Score: float64(now), Member: fmt.Sprintf("%d", time.Now().UnixNano())})
	// Count attempts in window
	pipe.ZCard(ctx, redisKey)
	// Set TTL
	pipe.Expire(ctx, redisKey, time.Duration(window)*time.Second)

	cmders, err := pipe.Exec(ctx)
	if err != nil {
		return false, err
	}

	count := cmders[2].(*redis.IntCmd).Val()
	return count <= int64(limit), nil
}

func (r *RedisLimiter) Reset(ctx context.Context, login, ip string) error {
	loginKey := fmt.Sprintf("limiter:%s:%s", domain.LimitLogin, login)
	ipKey := fmt.Sprintf("limiter:%s:%s", domain.LimitIP, ip)

	return r.client.Del(ctx, loginKey, ipKey).Err()
}
