package limiter

import (
	"context"
	"testing"

	"github.com/HelenaBlack/anti-bruteforce/internal/domain"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRedisLimiter(t *testing.T) {
	// Skip if no redis available locally
	client := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		t.Skip("Redis not available")
	}
	defer func() { _ = client.Close() }()

	r := NewRedisLimiter(client)
	key := "test_login"
	limit := 3

	// Cleanup
	client.Del(ctx, "limiter:login:test_login")

	// 1. Allow up to limit
	for i := 0; i < limit; i++ {
		ok, err := r.Allow(ctx, domain.LimitLogin, key, limit)
		require.NoError(t, err)
		assert.True(t, ok)
	}

	// 2. Deny after limit
	ok, err := r.Allow(ctx, domain.LimitLogin, key, limit)
	require.NoError(t, err)
	assert.False(t, ok)

	// 3. Reset
	err = r.Reset(ctx, key, "")
	require.NoError(t, err)

	// 4. Allow again after reset
	ok, err = r.Allow(ctx, domain.LimitLogin, key, limit)
	require.NoError(t, err)
	assert.True(t, ok)
}
