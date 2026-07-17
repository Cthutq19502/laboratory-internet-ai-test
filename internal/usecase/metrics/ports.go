package metrics

import (
	"context"
	domainmetrics "laboratory-internet-ai-test/internal/domain/metrics"
	"time"
)

type RepositoryRateLimiter interface {
	GetActionsCountIp(ctx context.Context, ip string, timeRequest time.Time, ttl int) (int, error)
	IncrStatusMetric(ctx context.Context, dateKey string, status domainmetrics.Status)
	GetMetric(ctx context.Context, dateKey string) (domainmetrics.Metric, error)
}

type RateLimiterUsecase interface {
	RateLimiterAllow(ctx context.Context, ip string, timeRequest time.Time) (bool, error)
}

type MetricUsecase interface {
	SaveMetrics(ctx context.Context, metric Metric) error
	GetMetric(ctx context.Context) (domainmetrics.Metric, error)
}

type Metric struct {
	Now        time.Time
	StatusCode int
}
