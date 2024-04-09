package main

import (
	"github.com/ankush-003/alerts-simulation-and-remediation/middleware/sim/alerts"
	"github.com/ankush-003/alerts-simulation-and-remediation/middleware/sim/kafka"

	"asmr/store"
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func main() {
	// loading .env file
	err_load := godotenv.Load()
	if err_load != nil {
		log.Fatalf("Error loading .env file: %s\n", err_load)
	}

	NodeID := uuid.New()
	logger := log.New(os.Stdout, fmt.Sprintf("Node %s:", NodeID.String()), log.LstdFlags)

	// config
	time_limit := 10 // 2 min

	broker := os.Getenv("KAFKA_BROKER")
	redis_addr := os.Getenv("REDIS_ADDR")

	if broker == "" {
		broker = "localhost:9092"
		logger.Println("KAFKA_BROKER not set, using default %s\n", broker)
	}

	brokers := []string{broker}
	// Create a new Sarama configuration
	username := os.Getenv("KAFKA_USERNAME")
	password := os.Getenv("KAFKA_PASSWORD")

	if username == "" || password == "" {
		// logger.Fatalf("KAFKA_USERNAME or KAFKA_PASSWORD not set\n")
		logger.Println("KAFKA_USERNAME or KAFKA_PASSWORD not set, using KAFKA locally")
	}

	config := kafka.NewConfig(username, password)

	producer, err := kafka.NewProducer(brokers, config, logger)
	if err != nil {
		logger.Fatalf("Error creating producer: %s\n", err)
	}
	defer producer.Close()

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

	alertsConfigChan := make(chan *alerts.AlertConfig)

	signalChan := make(chan os.Signal, 2)
	signal.Notify(signalChan, os.Interrupt)

	// Creating Alerts
	logger.Println("Creating alerts")

	go func() {
		// set lower limit of 1min and upper of 1min + time_limit
		interval := time.Duration(rand.Intn(time_limit)+1) * time.Minute
		ticker := time.NewTicker(interval)
		for {
			select {
			case <-ticker.C:
				alertConfig, err := redis.GetRandomAlertConfig(ctx)
				if err != nil {
					logger.Printf("Error getting random alert: %s\n", err)
					continue
				}
				alertsConfigChan <- &alertConfig

			case <-signalChan:
				logger.Printf("Interrupted\n")
				close(signalChan)
				return
			}
		}
	}()

	var wg sync.WaitGroup

	for {
		select {
		case alertConfig := <-alertsConfigChan:
			wg.Add(1)
			go func(alertConfig *alerts.AlertConfig) {
				defer wg.Done()
				newALert := alerts.NewAlert(alertConfig, NodeID, "demoSimulator")
				producer.SendAlert("alerts", newALert)
				err := redis.PublishAlerts(ctx, newALert)
				if err != nil {
					logger.Printf("Error publishing alert: %s\n", err)
				}
				logger.Printf("Alert sent: %s\n", newALert)
			}(alertConfig)

		case <-signalChan:
			logger.Printf("Stopping Simulator %s\n", NodeID.String())
			wg.Wait()
			return
		}
	}
}
