package metrics

import (
	"context"
	"laboratory-internet-ai-test/config"
	domainmetrics "laboratory-internet-ai-test/internal/domain/metrics"
	"log/slog"
	"time"
)

type Service struct {
	repo      RepositoryRateLimiter
	rateLimit rateLimit
	logger    *slog.Logger
}

type rateLimit struct {
	ttl   int
	limit int
}

func New(repo RepositoryRateLimiter, config config.RateLimiter, logger *slog.Logger) *Service {

	return &Service{repo: repo, logger: logger, rateLimit: rateLimit{
		ttl:   config.TTL,
		limit: config.Limit,
	},
	}
}

func (s *Service) RateLimiterAllow(ctx context.Context, ip string, timeRequest time.Time) (bool, error) {

	count, err := s.repo.GetActionsCountIp(ctx, ip, timeRequest, s.rateLimit.ttl)
	if err != nil {
		s.logger.Error("RateLimitAllow", "error", err)
		return false, err
	}

	if count > s.rateLimit.limit {
		return false, nil
	}

	return true, nil
}

func (s *Service) SaveMetrics(ctx context.Context, metric Metric) error {
	dateKey := metric.Now.Format(time.DateOnly)

	var status domainmetrics.Status
	switch metric.StatusCode {
	case 500:
		status = domainmetrics.CriticalStatus
	case 200, 201:
		status = domainmetrics.OKStatus
	case 429:
		status = domainmetrics.ManyRequestStatus
	case 400:
		status = domainmetrics.BadReqStatus
	default:
		status = domainmetrics.AnotherStatus
	}

	s.logger.InfoContext(ctx, "SaveMetrics", "status", status)

	s.repo.IncrStatusMetric(ctx, dateKey, status)
	return nil
}

func (s *Service) GetMetric(ctx context.Context) (domainmetrics.Metric, error) {

	metric, err := s.repo.GetMetric(ctx, time.Now().Format(time.DateOnly))
	if err != nil {
		return domainmetrics.Metric{}, nil
	}
	return metric, nil
}
