package main

import (
	"github.com/ankush-003/alerts-simulation-and-remediation/middleware/sim/alerts"
	"github.com/ankush-003/alerts-simulation-and-remediation/middleware/sim/kafka"

	"io"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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
	alertsChan := make(chan alerts.Alerts)
	doneChan := make(chan struct{})
	consumer, err := kafka.NewConsumer(brokers, config, logger)
	if err != nil {
		logger.Fatalf("Error creating consumer: %s\n", err)
	}
	defer consumer.Close()

	go consumer.ConsumeAlerts("alerts", alertsChan, doneChan)

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to the SSE server!",
		})
	})

	router.GET("/stream", func(c *gin.Context) {
		// send message initially to avoid timeout
		c.SSEvent("alert", map[string]interface{}{
			"message": "Connected to the server",
		})
		streamer(c, alertsChan)
	})

	router.Run(":8090")
}

func streamer(c *gin.Context, alertsChan chan alerts.Alerts) {
	c.Stream(func(w io.Writer) bool {
		select {
		case alert := <-alertsChan:
			c.SSEvent("alert", map[string]interface{}{
				"alertID":     alert.ID.String(),
				"nodeID":      alert.NodeID.String(),
				"description": alert.Description,
				"severity":    alert.Severity,
				"source":      alert.Source,
				"createdAt":   alert.CreatedAt,
				"runtimeMetrics": map[string]interface{}{
					"numGoroutine": alert.RuntimeMetrics.NumGoroutine,
					"cpuUsage":     alert.RuntimeMetrics.CpuUsage,
					"ramUsage":     alert.RuntimeMetrics.RamUsage,
				},
			})
			return true
		case <-c.Writer.CloseNotify():
			return false
		}
	})

	c.Writer.Flush()
}
