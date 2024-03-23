package floodControl

import (
	"context"
	"fmt"
	"github.com/17HIERARCH70/flood-control/internal/config"
	"github.com/go-redis/redis/v8"
	"time"
)

// Service a structure that implements the logic of controlling flooding.
type Service struct {
	client *redis.Client
	config config.FloodControlConfig
}

// NewService function to create a new Service instance.
func NewService(client *redis.Client, cfg config.FloodControlConfig) *Service {
	return &Service{
		client: client,
		config: cfg,
	}
}

// Check method to check if the user's request limit has been exceeded.
func (s *Service) Check(ctx context.Context, userID int64) (bool, error) {
	key := fmt.Sprintf("user:%d:requests", userID)
	now := time.Now()

	windowStart := now.Add(-time.Duration(s.config.PeriodSeconds) * time.Second).Unix()

	pipe := s.client.TxPipeline()

	pipe.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", windowStart))

	pipe.ZAdd(ctx, key, &redis.Z{
		Score:  float64(now.Unix()),
		Member: now.UnixNano(),
	})

	countCmd := pipe.ZCount(ctx, key, fmt.Sprintf("%d", windowStart), fmt.Sprintf("%d", now.Unix()))

	pipe.Expire(ctx, key, time.Duration(s.config.PeriodSeconds+1)*time.Second)

	if _, err := pipe.Exec(ctx); err != nil {
		return false, fmt.Errorf("failed to execute Redis transaction: %v", err)
	}

	count, err := countCmd.Result()
	if err != nil {
		return false, fmt.Errorf("failed to count requests: %v", err)
	}

	return count <= int64(s.config.RequestLimit), nil
}
