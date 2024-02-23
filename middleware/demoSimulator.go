package main

import (
	"asmr/alerts"
	"asmr/kafka"
	"asmr/mailserver"
	"asmr/store"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
)

func main() {

	NodeID := uuid.New()
	logger := log.New(os.Stdout, fmt.Sprintf("Node %s:", NodeID.String()), log.LstdFlags)

	broker := os.Getenv("KAFKA_BROKER")

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

	redis := store.NewRedisStore("localhost:6379")
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
			fmt.Println(redis.GetAlertsByNodeID(ctx,NodeID.String()))
			mailserver.SendEmail(redis.GetAlertsByNodeID(ctx, NodeID.String()))
			wg.Wait()
			return
		}
	}
}
