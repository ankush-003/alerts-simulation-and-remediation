package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	// "github.com/ankush-003/alerts-simulation-and-remediation/middleware/sim/alerts"
	// rule_engine "github.com/ankush-003/alerts-simulation-and-remediation/middleware/rule_engine_v2/engine"
	"github.com/ankush-003/alerts-simulation-and-remediation/middleware/sim/alerts"
	"github.com/ankush-003/alerts-simulation-and-remediation/middleware/sim/kafka"
	"github.com/joho/godotenv"
	// "github.com/google/uuid"
)

func main() {

	err_load := godotenv.Load()
	if err_load != nil {
		log.Fatalf("Error loading .env file: %v\n", err_load)
	}

	logger := log.New(os.Stdout, "kafka-producer: ", log.LstdFlags)

	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		broker = "localhost:9092"
		logger.Println("KAFKA_BROKER not set, using default ", broker)
	}
	brokers := []string{broker}
	username := os.Getenv("KAFKA_USERNAME")
	password := os.Getenv("KAFKA_PASSWORD")

	config := kafka.NewConfig(username, password)

	producer, err := kafka.NewProducer(brokers, config, logger)

	if err != nil {
		logger.Fatalf("Error creating producer: %s\n", err)
	}

	defer producer.Close()

	// alertConf := alerts.NewAlertConfig("High CPU usage", "critical")
	// alert := alerts.NewAlert(alertConf, uuid.New(), "CPU")
	alert := alerts.AlertInput{
		ID:        "ID1",
		Category:  "Memory",
		Source:    "Hardware",
		Origin:    "NodeA",
		Params:    &alerts.Memory{Usage: 76, PageFaults: 30, SwapUsage: 2},
		CreatedAt: time.Now(),
		Handled:   false,
	}
	// create a timer to send the alert every 5 seconds
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	for {
		select {
		case <-ticker.C:
			if err := producer.SendAlert("alerts", &alert); err != nil {
				logger.Printf("Error sending alert: %s\n", err)
			}
			logger.Printf("Sent alert: %v\n", alert)
		case <-signalChan:
			logger.Printf("Interrupted\n")
			return
		}
	}
}
