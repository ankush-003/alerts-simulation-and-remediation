package main

import (
	"asmr/alerts"
	"context"
	"asmr/store"
	"log"
	"os"
	"github.com/joho/godotenv"
)

func main() {

	// loading .env file
	err_load := godotenv.Load()
	if err_load != nil {
		log.Fatalf("Error loading .env file: %v\n", err_load)
	}

	logger := log.New(os.Stdout, "redis-consumer: ", log.LstdFlags)
	
	// redis_addr := os.Getenv("REDIS_ADDR")
	redis_addr := ""

	if redis_addr == "" {
		logger.Println("REDIS_ADDR not set, using default localhost:6379")
		redis_addr = "localhost:6379"
	}

	ctx := context.Background()

	redis, redisErr := store.NewRedisStore(ctx, redis_addr) 

	if redisErr != nil {
		logger.Fatalf("Error creating redis store: %s\n", redisErr)
	}
	
	defer redis.Close()

	alertsChan := make(chan alerts.Alerts)
	doneChan := make(chan struct{})

	logger.Println("Consuming alerts !")

	stream := "alerts"
    groupName := "alerts-group"

	redis.ConsumeAlerts(ctx, alertsChan, doneChan, stream, groupName)

	consumerLoop:
		for {
			select {
			case alert := <-alertsChan:
				logger.Printf("Received alert: alrtID: %s, NodeID: %s, Description: %s, Severity: %s, Source: %s, CreatedAt: %s\t", alert.ID.String(), alert.NodeID.String(), alert.Description, alert.Severity, alert.Source, alert.CreatedAt)
				logger.Printf("RuntimeMetrics: NumGoroutine: %d, CpuUsage: %f, RamUsage: %f\n\n", alert.RuntimeMetrics.NumGoroutine, alert.RuntimeMetrics.CpuUsage, alert.RuntimeMetrics.RamUsage)
			case <-doneChan:
				break consumerLoop
			}
		}
}