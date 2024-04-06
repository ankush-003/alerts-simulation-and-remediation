package main

import (
	"asmr/kafka"
	"log"
	"os"

	rule_engine "github.com/ankush-003/alerts-simulation-and-remediation/middleware/rule_engine_v2/engine"
	"github.com/joho/godotenv"
)

func main() {

	// loading .env file
	err_load := godotenv.Load()
	if err_load != nil {
		log.Fatalf("Error loading .env file: %v\n", err_load)
	}

	logger := log.New(os.Stdout, "kafka-consumer: ", log.LstdFlags)
	
	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		broker = "localhost:9092"
		logger.Println("KAFKA_BROKER not set, using default %s\n", broker)
	}
	brokers := []string{broker}
	username := os.Getenv("KAFKA_USERNAME")
	password := os.Getenv("KAFKA_PASSWORD")
	
	config := kafka.NewConfig(username, password) 

	consumer, err := kafka.NewConsumer(brokers, config, logger)
	if err != nil {
		logger.Fatalf("Error creating consumer: %s\n", err)
	}

	defer consumer.Close()

	alertsChan := make(chan rule_engine.AlertInput)
	doneChan := make(chan struct{})

	logger.Println("Consuming alerts !")

	go consumer.ConsumeAlerts("alerts", alertsChan, doneChan)
consumerLoop:

	for {
		select {
		case alert := <-alertsChan:
			// fmt.Println(alert)
			logger.Printf("Received alert: alrtID: %s, NodeID: %s, Description: %s, Severity: %s, Source: %s, CreatedAt: %s\n", alert.ID.String(), alert.NodeID.String(), alert.Description, alert.Severity, alert.Source, alert.CreatedAt)
			// mailserver.SendEmail(alert, nil)
		case <-doneChan:
			break consumerLoop
		}
	}
}