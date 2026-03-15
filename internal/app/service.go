package app

import (
	"context"

	"github.com/HelenaBlack/anti-bruteforce/internal/config"
	"github.com/HelenaBlack/anti-bruteforce/internal/domain"
)

type AntiBruteforceService struct {
	limiter domain.RateLimiter
	ipRepo  domain.IPService
	cfg     *config.Config
}

func NewAntiBruteforceService(
	limiter domain.RateLimiter,
	ipRepo domain.IPService,
	cfg *config.Config,
) *AntiBruteforceService {
	return &AntiBruteforceService{
		limiter: limiter,
		ipRepo:  ipRepo,
		cfg:     cfg,
	}
}

func (s *AntiBruteforceService) Check(ctx context.Context, login, password, ip string) (bool, error) {
	// 1. Check Whitelist
	whitelisted, err := s.ipRepo.IsWhitelisted(ctx, ip)
	if err != nil {
		return false, err
	}
	if whitelisted {
		return true, nil
	}

	// 2. Check Blacklist
	blacklisted, err := s.ipRepo.IsBlacklisted(ctx, ip)
	if err != nil {
		return false, err
	}
	if blacklisted {
		return false, nil
	}

	// 3. Check Rate Limits
	// Login limit
	ok, err := s.limiter.Allow(ctx, domain.LimitLogin, login, s.cfg.LimitN)
	if err != nil || !ok {
		return false, err
	}

	// Password limit
	ok, err = s.limiter.Allow(ctx, domain.LimitPassword, password, s.cfg.LimitM)
	if err != nil || !ok {
		return false, err
	}

	// IP limit
	ok, err = s.limiter.Allow(ctx, domain.LimitIP, ip, s.cfg.LimitK)
	if err != nil || !ok {
		return false, err
	}

	return true, nil
}

func (s *AntiBruteforceService) Reset(ctx context.Context, login, ip string) error {
	return s.limiter.Reset(ctx, login, ip)
}

func (s *AntiBruteforceService) AddToBlacklist(ctx context.Context, subnet string) error {
	return s.ipRepo.AddToBlacklist(ctx, subnet)
}

func (s *AntiBruteforceService) RemoveFromBlacklist(ctx context.Context, subnet string) error {
	return s.ipRepo.RemoveFromBlacklist(ctx, subnet)
}

func (s *AntiBruteforceService) AddToWhitelist(ctx context.Context, subnet string) error {
	return s.ipRepo.AddToWhitelist(ctx, subnet)
}

func (s *AntiBruteforceService) RemoveFromWhitelist(ctx context.Context, subnet string) error {
	return s.ipRepo.RemoveFromWhitelist(ctx, subnet)
}
