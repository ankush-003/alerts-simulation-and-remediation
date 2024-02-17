package main

import (
	"asmr/alerts"
	"asmr/kafka"
	"log"
	"os"
	"os/signal"
	"time"
	"github.com/IBM/sarama"
)

func main() {
	logger := log.New(os.Stdout, "kafka-consumer: ", log.LstdFlags)

	brokers := []string{"localhost:9092"}
	config := sarama.NewConfig()

	producer, err := kafka.NewProducer(brokers, config, logger)

	if err != nil {
		logger.Fatalf("Error creating producer: %s\n", err)
	}

	defer producer.Close()

	alert := alerts.NewAlert("This is a test alert", "critical", "kafka-test-producer")

	// create a timer to send the alert every 5 seconds
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	for {
		select {
		case <-ticker.C:
			if err := producer.SendAlert("alerts", alert); err != nil {
				logger.Printf("Error sending alert: %s\n", err)
			}
			logger.Printf("Sent alert: %v\n", alert)
		case <-signalChan:
			logger.Printf("Interrupted\n")
			return
		}
	}

}