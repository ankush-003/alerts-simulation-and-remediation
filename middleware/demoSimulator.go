package main

import (
	"asmr/alerts"
	"asmr/kafka"
	"asmr/store"
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"
)

func main() {

	NodeID := uuid.New()
	logger := log.New(os.Stdout, fmt.Sprintf("Node %s:", NodeID.String()), log.LstdFlags)

	broker := os.Getenv("KAFKA_BROKER")
	redis_addr := os.Getenv("REDIS_ADDR")

	if broker == "" {
		logger.Println("KAFKA_BROKER not set, using default localhost:9092")
		broker = "localhost:9092"
	}

	brokers := []string{broker}
	config := sarama.NewConfig()

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

	logger.Println("Creating alerts")

	go func() {
		ticker := time.NewTicker(2 * time.Second)
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
				producer.SendAlert("alerts", alerts.NewAlert(alertConfig, NodeID, "demoSimulator"))
			}(alertConfig)

		case <-signalChan:
			logger.Printf("Stopping Simulator %s\n", NodeID.String())
			wg.Wait()
			return
		}
	}
}
