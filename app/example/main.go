package main

import (
	"context"
	"fmt"
	"log"
	"task/internal/config"
	"task/internal/services/floodControl"
	"task/internal/storage/redis"
)

// Example of usage as cli
func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	redisClient, err := redis.NewClient(cfg.Redis)
	if err != nil {
		log.Fatalf("Failed to initialize Redis client: %v", err)
	}
	defer redisClient.Close()

	floodCtrlService := floodControl.NewService(redisClient, cfg.FloodControl)

	ctx := context.Background()
	userID := int64(3) // ID example
	allowed, err := floodCtrlService.Check(ctx, userID)
	if err != nil {
		log.Fatalf("Failed to check flood control: %v", err)
	}

	fmt.Printf("Request allowed: %v\n", allowed)
}
