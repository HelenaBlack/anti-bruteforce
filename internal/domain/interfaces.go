package domain

import (
	"context"
	"net"
)

type RateLimitType string

const (
	LimitLogin    RateLimitType = "login"
	LimitPassword RateLimitType = "password"
	LimitIP       RateLimitType = "ip"
)

type RateLimiter interface {
	Allow(ctx context.Context, limitType RateLimitType, key string, limit int) (bool, error)
	Reset(ctx context.Context, login, ip string) error
}

type IPService interface {
	IsBlacklisted(ctx context.Context, ip string) (bool, error)
	IsWhitelisted(ctx context.Context, ip string) (bool, error)
	AddToBlacklist(ctx context.Context, subnet string) error
	RemoveFromBlacklist(ctx context.Context, subnet string) error
	AddToWhitelist(ctx context.Context, subnet string) error
	RemoveFromWhitelist(ctx context.Context, subnet string) error
}

func IsInSubnet(ipStr, subnetStr string) (bool, error) {
	ip := net.ParseIP(ipStr)
	_, subnet, err := net.ParseCIDR(subnetStr)
	if err != nil {
		return false, err
	}
	return subnet.Contains(ip), nil
}
