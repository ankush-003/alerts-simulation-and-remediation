package main

import (
	"context"
	"github.com/ankush-003/alerts-simulation-and-remediation/middleware/sim/alerts"
	"github.com/ankush-003/alerts-simulation-and-remediation/middleware/sim/kafka"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx := context.Background()

	err_load := godotenv.Load()
	if err_load != nil {
		log.Println("Error loading .env file")
	}

	logger := log.New(os.Stdout, "[kafka-group-consumer] ", log.LstdFlags)

	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		broker = "localhost:9092"
		logger.Println("KAFKA_BROKER not set, using default %s\n", broker)
	}

	logger.Println("Kafka broker: %s\n", broker)

	brokers := []string{broker}
	username := os.Getenv("KAFKA_USERNAME")
	password := os.Getenv("KAFKA_PASSWORD")
	config := kafka.NewConfig(username, password)
	groupID := "test-group"
	topics := []string{"alerts"}

	consumerGroup, err := kafka.NewConsumerGroup(brokers, groupID, topics, config, logger)
	if err != nil {
		logger.Fatalf("Error creating consumer group: %s\n", err)
	}
	defer consumerGroup.Close(logger)

	alertsChan := make(chan alerts.AlertInput)
	doneChan := make(chan struct{})

	go consumerGroup.Consume(ctx, topics, alertsChan, doneChan)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-signals:
			logger.Println("Interrupted")
			return

		case alert := <-alertsChan:
			logger.Printf("Received alert: %s\n", alert)
		}
	}

}
