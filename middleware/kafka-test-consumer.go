package main

import (
	"asmr/alerts"
	"asmr/kafka"
	"log"
	"os"
)

func main() {
	logger := log.New(os.Stdout, "kafka-consumer: ", log.LstdFlags)
	
	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		broker = "active-boar-11578-us1-kafka.upstash.io:9092"
		logger.Println("KAFKA_BROKER not set, using default %s\n", broker)
	}
	brokers := []string{broker}
	// username := os.Getenv("KAFKA_USERNAME")
	// password := os.Getenv("KAFKA_PASSWORD")
	// if username == "" || password == "" {
	// 	logger.Fatalf("KAFKA_USERNAME or KAFKA_PASSWORD not set\n")
	// }
	
	username := "test"
	password := "test"
	config := kafka.NewConfig(username, password) 

	consumer, err := kafka.NewConsumer(brokers, config, logger)
	if err != nil {
		logger.Fatalf("Error creating consumer: %s\n", err)
	}

	defer consumer.Close()

	alertsChan := make(chan alerts.Alerts)
	doneChan := make(chan struct{})

	logger.Println("Consuming alerts !")

	go consumer.ConsumeAlerts("alerts", alertsChan, doneChan)
consumerLoop:

	for {
		select {
		case alert := <-alertsChan:
			logger.Printf("Received alert: alrtID: %s, NodeID: %s, Description: %s, Severity: %s, Source: %s, CreatedAt: %s\n", alert.ID.String(), alert.NodeID.String(), alert.Description, alert.Severity, alert.Source, alert.CreatedAt)

		case <-doneChan:
			break consumerLoop
		}
	}
}