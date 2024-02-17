package main

import (
	"asmr/alerts"
	"asmr/kafka"
	"asmr/store"
	"log"
	"sync"
	"os"
	"os/signal"
	"context"
	"time"
	"github.com/IBM/sarama"
)

func main() {
	logger := log.New(os.Stdout, "Simulator: ", log.LstdFlags)

	brokers := []string{"localhost:9092"}
	config := sarama.NewConfig()

	producer, err := kafka.NewProducer(brokers, config, logger)
	if err != nil {
		logger.Fatalf("Error creating producer: %s\n", err)
	}
	defer producer.Close()

	redis := store.NewRedisStore("localhost:6379")
	defer redis.Close()

	alertsChan := make(chan alerts.Alerts)
	
	signalChan := make(chan os.Signal, 2)
	signal.Notify(signalChan, os.Interrupt)

	ctx := context.Background()
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		for {
			select {
			case <-ticker.C:
				alert, err := redis.GetRandomAlert(ctx)
				if err != nil {
					logger.Printf("Error getting random alert: %s\n", err)
					continue
				}
				alertsChan <- alert

			case <-signalChan:
				logger.Printf("Interrupted\n")
				return
			}
		}
	} ()


	var wg sync.WaitGroup
	
	for {
		select {
		case alert := <-alertsChan:
			wg.Add(1)
			go func(alert alerts.Alerts) {
				defer wg.Done()
				producer.SendAlert("alerts", &alert)
			} (alert)

		case <-signalChan:
			logger.Printf("Interrupted\n")
			wg.Wait()
			return
		}
	}
}