package floodControl_test

import (
	"context"
	"github.com/17HIERARCH70/flood-control/internal/config"
	"github.com/17HIERARCH70/flood-control/internal/services/floodControl"
	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCheck(t *testing.T) {
	// Set up a mock Redis server
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mr.Close()

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	// Test configuration
	cfg := config.FloodControlConfig{
		RequestLimit:  5,
		PeriodSeconds: 60,
	}

	service := floodControl.NewService(client, cfg)

	for i := 0; i < cfg.RequestLimit; i++ {
		allowed, err := service.Check(context.Background(), 1)
		assert.Nil(t, err)
		assert.True(t, allowed)
	}

	allowed, err := service.Check(context.Background(), 1)
	assert.Nil(t, err)
	assert.False(t, allowed)
}
