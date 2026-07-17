package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	domainmetrics "laboratory-internet-ai-test/internal/domain/metrics"
	"log/slog"
	"time"
)

type Repository struct {
	client *redis.Client
	logger *slog.Logger
}

func New(client *redis.Client, logger *slog.Logger) *Repository {
	return &Repository{
		client: client,
		logger: logger,
	}
}

func (r *Repository) GetActionsCountIp(ctx context.Context, ip string, timeRequest time.Time, ttl int) (int, error) {

	key := fmt.Sprintf("ratelimit:%s", ip)

	ttlDuration := time.Duration(ttl) * time.Second
	timeUnix := timeRequest.Unix()
	minScore := timeUnix - int64(ttl)

	r.logger.InfoContext(ctx, "GetActionsCountIp", "timeUnix", timeUnix, "ttlUnix", ttl)

	pipe := r.client.Pipeline()
	r.logger.InfoContext(ctx, "GetActionsCountIp", "ZRemRangeByScore", fmt.Sprintf("%d", minScore))
	pipe.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", minScore))

	r.client.ZAdd(ctx, key, redis.Z{
		Score:  float64(timeUnix),
		Member: fmt.Sprintf("%d", timeUnix),
	})

	pipe.Expire(ctx, key, ttlDuration)

	countCmd := pipe.ZCard(ctx, key)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return 0, err
	}

	exists, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	if exists == 0 {
		err := r.client.Set(ctx, key, 1, ttlDuration).Err()
		if err != nil {
			return 0, err
		}
		return 1, nil
	}

	return int(countCmd.Val()), nil

}

func (r *Repository) IncrStatusMetric(ctx context.Context, dateKey string, status domainmetrics.Status) {
	r.client.HIncrBy(ctx, dateKey, string(status), 1)
	return
}

func (r *Repository) GetMetric(ctx context.Context, dateKey string) (domainmetrics.Metric, error) {
	data, err := r.client.HGetAll(ctx, dateKey).Result()
	if err != nil {
		return domainmetrics.Metric{}, err
	}

	r.logger.Info("Repo GetMetric", "data", data)

	metric := domainmetrics.Metric{
		DateKey:           dateKey,
		CriticalStatus:    data[string(domainmetrics.CriticalStatus)],
		OKStatus:          data[string(domainmetrics.OKStatus)],
		BadReqStatus:      data[string(domainmetrics.BadReqStatus)],
		ManyRequestStatus: data[string(domainmetrics.ManyRequestStatus)],
		AnotherStatus:     data[string(domainmetrics.AnotherStatus)],
	}

	return metric, nil
}
