package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"task/internal/config"
	"task/internal/services/floodControl"
	"task/internal/storage/redis"
)

type RequestPayload struct {
	UserID int64 `json:"userId"`
}

// Request handler
func floodHandler(floodCtrl *floodControl.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var reqBody RequestPayload
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		if err != nil {
			http.Error(w, "Failed to decode request body", http.StatusBadRequest)
			return
		}

		allowed, err := floodCtrl.Check(r.Context(), reqBody.UserID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to check flood control: %v", err), http.StatusInternalServerError)
			return
		}

		if allowed {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Request allowed"))
		} else {
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte("Request rate exceeded"))
		}
	}
}

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

	http.HandleFunc("/flood", floodHandler(floodCtrlService))

	port := 8080
	log.Printf("Server started at :%d", port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}
