package main

import (
	"asmr/alerts"
	"asmr/kafka"
	"asmr/rule_engine"
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
	time_limit := 60 // 1 minute

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

	redis := store.NewRedisStore(redis_addr)
	defer redis.Close()

	alertsConfigChan := make(chan *alerts.AlertConfig)

	signalChan := make(chan os.Signal, 2)
	signal.Notify(signalChan, os.Interrupt)

	ctx := context.Background()
	
	// Creating Alerts
	logger.Println("Creating alerts")

	go func() {
		interval := time.Duration(rand.Intn(time_limit)) * time.Second
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
				// fmt.Println(alertConfig.ID)
				defer wg.Done()

				

				
				
				


				testInput:= rule_engine.AlertInput{
					"alert-1",
					"Memory",
					"router",
					"datacenter",
					rule_engine.Memory{
						Usage:      1024,
						PageFaults: 10,
						SwapUsage:  512,
					},
					time.Now(),
					false,
				}
				

				producer.SendAlert("alerts", testInput)
			}(alertConfig)

		case <-signalChan:
			logger.Printf("Stopping Simulator %s\n", NodeID.String())
			wg.Wait()
			return
		}
	}
}