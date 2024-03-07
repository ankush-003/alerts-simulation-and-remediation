package main

import (
	"asmr/alerts"
	"asmr/kafka"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {

	err_load := godotenv.Load()
	if err_load != nil {
		log.Fatalf("Error loading .env file: %v\n", err_load)
	}

	logger := log.New(os.Stdout, "sse-server: ", log.LstdFlags)

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

	alertsChan := make(chan alerts.Alerts)
	doneChan := make(chan struct{})

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to the SSE server",
		})
	})

	router.GET("/stream", func(c *gin.Context) {
		logger.Println("Consuming alerts !")

		go consumer.ConsumeAlerts("alerts", alertsChan, doneChan)
	consumerLoop:

		for {
			select {
			case alert := <-alertsChan:
				logger.Printf("Received alert: alrtID: %s, NodeID: %s, Description: %s, Severity: %s, Source: %s, CreatedAt: %s\n", alert.ID.String(), alert.NodeID.String(), alert.Description, alert.Severity, alert.Source, alert.CreatedAt)
				c.SSEvent("alert", map[string]interface{} {
					"alertID": alert.ID.String(),
					"nodeID": alert.NodeID.String(),
					"description": alert.Description,
					"severity": alert.Severity,
					"source": alert.Source,
					"createdAt": alert.CreatedAt,
				})

				c.Writer.Flush()

			case <-doneChan:
				break consumerLoop
			}
		}

		c.Writer.Flush()
	})

	router.Run(":8080")
}
