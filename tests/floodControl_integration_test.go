package floodControl_test

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"task/internal/config"
	"task/internal/services/floodControl"
	"task/internal/storage/redis"
)

func TestFloodControlIntegration(t *testing.T) {
	cfg := config.RedisConfig{
		Address:  "localhost:6379",
		Password: "", // Update this if your Redis setup requires a password
		DB:       0,
	}

	client, err := redis.NewClient(cfg)
	if err != nil {
		log.Fatalf("Failed to create Redis client: %v", err)
	}
	defer client.Close()

	err = client.FlushDB(context.Background()).Err()
	if err != nil {
		log.Fatalf("Failed to flush Redis: %v", err)
	}

	floodControlService := floodControl.NewService(client, config.FloodControlConfig{
		RequestLimit:  2,
		PeriodSeconds: 5, // Extending the period to ensure clear separation of test phases
	})

	ctx := context.Background()
	userID := int64(123)

	// First request - should be allowed
	allowed1, err := floodControlService.Check(ctx, userID)
	assert.NoError(t, err)
	assert.True(t, allowed1, "First request should be allowed")

	// Second request - should also be allowed
	allowed2, err := floodControlService.Check(ctx, userID)
	assert.NoError(t, err)
	assert.True(t, allowed2, "Second request should be allowed")
	time.Sleep(2 * time.Second)

	// Third request - should be blocked
	allowed3, err := floodControlService.Check(ctx, userID)
	assert.NoError(t, err)
	assert.False(t, allowed3, "Third request should not be allowed due to rate limit")

	time.Sleep(5 * time.Second)

	// Fourth request - should be allowed again as we're in a new rate limit window
	allowed4, err := floodControlService.Check(ctx, userID)
	assert.NoError(t, err)
	assert.True(t, allowed4, "Fourth request should be allowed as we're in a new rate limit window")

	if !allowed3 || !allowed4 {
		currentCount, err := client.ZCount(ctx, fmt.Sprintf("user:%d:requests", userID), "-inf", "+inf").Result()
		if err != nil {
			log.Printf("Error fetching current count from Redis: %v", err)
		} else {
			log.Printf("Current request count in Redis for userID %d: %d", userID, currentCount)
		}
	}
}
