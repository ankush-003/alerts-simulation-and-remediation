package main

import (
	"github.com/ankush-003/alerts-simulation-and-remediation/middleware/sim/alerts"

	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/ankush-003/alerts-simulation-and-remediation/middleware/sim/store"

	"github.com/joho/godotenv"
)

func main() {
	// loading .env file
	err_load := godotenv.Load()
	if err_load != nil {
		log.Fatalf("Error loading .env file: %s\n", err_load)
	}

	redis_addr := os.Getenv("REDIS_ADDR")
	if redis_addr == "" {
		redis_addr = "localhost:6379"
		log.Println("REDIS_ADDR not set, using default localhost:6379")
	}
	ctx := context.Background()
	redis, _ := store.NewRedisStore(ctx, redis_addr)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	//predefined alert configs
	descriptions := []string{"High CPU usage", "Low disk space", "Network down", "Service unavailable"}
	severities := []string{"critical", "warning", "info"}

	// push predefined alert configs to redis
	for _, desc := range descriptions {
		for _, sev := range severities {
			alertConfig := alerts.NewAlertConfig(desc, sev)
			err := redis.StoreAlertConfig(ctx, alertConfig)
			if err != nil {
				log.Printf("Error storing alert config: %s\n", err)
			}
		}
	}
	fmt.Println("Alert configs stored")
}
