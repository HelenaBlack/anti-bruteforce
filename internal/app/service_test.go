package app

import (
	"context"
	"errors"
	"testing"

	"github.com/HelenaBlack/anti-bruteforce/internal/config" //nolint:depguard
	"github.com/HelenaBlack/anti-bruteforce/internal/domain" //nolint:depguard
	"github.com/stretchr/testify/assert"                     //nolint:depguard
	"github.com/stretchr/testify/mock"                       //nolint:depguard
)

// Mocks.
type mockLimiter struct {
	mock.Mock
}

func (m *mockLimiter) Allow(ctx context.Context, limitType domain.RateLimitType, key string, limit int) (bool, error) {
	args := m.Called(ctx, limitType, key, limit)
	return args.Bool(0), args.Error(1)
}

func (m *mockLimiter) Reset(ctx context.Context, login, ip string) error {
	args := m.Called(ctx, login, ip)
	return args.Error(0)
}

type mockIPRepo struct {
	mock.Mock
}

func (m *mockIPRepo) IsBlacklisted(ctx context.Context, ip string) (bool, error) {
	args := m.Called(ctx, ip)
	return args.Bool(0), args.Error(1)
}

func (m *mockIPRepo) IsWhitelisted(ctx context.Context, ip string) (bool, error) {
	args := m.Called(ctx, ip)
	return args.Bool(0), args.Error(1)
}

func (m *mockIPRepo) AddToBlacklist(ctx context.Context, subnet string) error {
	args := m.Called(ctx, subnet)
	return args.Error(0)
}

func (m *mockIPRepo) RemoveFromBlacklist(ctx context.Context, subnet string) error {
	args := m.Called(ctx, subnet)
	return args.Error(0)
}

func (m *mockIPRepo) AddToWhitelist(ctx context.Context, subnet string) error {
	args := m.Called(ctx, subnet)
	return args.Error(0)
}

func (m *mockIPRepo) RemoveFromWhitelist(ctx context.Context, subnet string) error {
	args := m.Called(ctx, subnet)
	return args.Error(0)
}

func TestAntiBruteforceService_Check(t *testing.T) {
	cfg := &config.Config{
		LimitN: 10,
		LimitM: 100,
		LimitK: 1000,
	}

	ctx := context.Background()

	t.Run("Whitelisted IP should pass regardless of limits", func(t *testing.T) {
		limiter := new(mockLimiter)
		repo := new(mockIPRepo)
		service := NewAntiBruteforceService(limiter, repo, cfg)

		repo.On("IsWhitelisted", ctx, "1.1.1.1").Return(true, nil)

		ok, err := service.Check(ctx, "user", "pass", "1.1.1.1")
		assert.NoError(t, err)
		assert.True(t, ok)
		repo.AssertExpectations(t)
		limiter.AssertExpectations(t)
	})

	t.Run("Blacklisted IP should fail", func(t *testing.T) {
		limiter := new(mockLimiter)
		repo := new(mockIPRepo)
		service := NewAntiBruteforceService(limiter, repo, cfg)

		repo.On("IsWhitelisted", ctx, "2.2.2.2").Return(false, nil)
		repo.On("IsBlacklisted", ctx, "2.2.2.2").Return(true, nil)

		ok, err := service.Check(ctx, "user", "pass", "2.2.2.2")
		assert.NoError(t, err)
		assert.False(t, ok)
		repo.AssertExpectations(t)
	})

	t.Run("Rate limit exceeded for login", func(t *testing.T) {
		limiter := new(mockLimiter)
		repo := new(mockIPRepo)
		service := NewAntiBruteforceService(limiter, repo, cfg)

		repo.On("IsWhitelisted", ctx, "3.3.3.3").Return(false, nil)
		repo.On("IsBlacklisted", ctx, "3.3.3.3").Return(false, nil)
		limiter.On("Allow", ctx, domain.LimitLogin, "user", cfg.LimitN).Return(false, nil)

		ok, err := service.Check(ctx, "user", "pass", "3.3.3.3")
		assert.NoError(t, err)
		assert.False(t, ok)
		limiter.AssertExpectations(t)
	})

	t.Run("All checks passed", func(t *testing.T) {
		limiter := new(mockLimiter)
		repo := new(mockIPRepo)
		service := NewAntiBruteforceService(limiter, repo, cfg)

		repo.On("IsWhitelisted", ctx, "4.4.4.4").Return(false, nil)
		repo.On("IsBlacklisted", ctx, "4.4.4.4").Return(false, nil)
		limiter.On("Allow", ctx, domain.LimitLogin, "user", cfg.LimitN).Return(true, nil)
		limiter.On("Allow", ctx, domain.LimitPassword, "pass", cfg.LimitM).Return(true, nil)
		limiter.On("Allow", ctx, domain.LimitIP, "4.4.4.4", cfg.LimitK).Return(true, nil)

		ok, err := service.Check(ctx, "user", "pass", "4.4.4.4")
		assert.NoError(t, err)
		assert.True(t, ok)
		limiter.AssertExpectations(t)
	})

	t.Run("Error in repository", func(t *testing.T) {
		limiter := new(mockLimiter)
		repo := new(mockIPRepo)
		service := NewAntiBruteforceService(limiter, repo, cfg)

		repo.On("IsWhitelisted", ctx, "5.5.5.5").Return(false, errors.New("db error"))

		ok, err := service.Check(ctx, "user", "pass", "5.5.5.5")
		assert.Error(t, err)
		assert.False(t, ok)
	})
}
