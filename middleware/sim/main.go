package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/ankush-003/alerts-simulation-and-remediation/middleware/sim/alerts"
	"github.com/ankush-003/alerts-simulation-and-remediation/middleware/sim/kafka"
	"github.com/ankush-003/alerts-simulation-and-remediation/middleware/sim/store"
	"github.com/joho/godotenv"
)

func main() {
	// loading .env file
	err_load := godotenv.Load()
	if err_load != nil {
		log.Fatalf("Error loading .env file: %s\n", err_load)
	}

	ctx := context.Background()

	// create a new mongo store
	mongo_uri := os.Getenv("MONGO_URI")
	if mongo_uri == "" {
		log.Fatalf("MONGO_URI not set\n")
	}
	mongo_client, err := store.NewMongoStore(ctx, mongo_uri, "AlertSimAndRemediation", "Nodes")
	if err != nil {
		log.Fatalf("Error creating mongo store: %s\n", err)
	}
	defer mongo_client.Close(ctx)

	// new uid from available nodes from mongo
	NodeID, close_id, err := mongo_client.GetNodeId(ctx)
	if err != nil {
		log.Fatalf("Error getting node id: %s\n", err)
	}
	defer close_id()

	logger := log.New(os.Stdout, fmt.Sprintf("Node %s:", NodeID), log.LstdFlags)

	// config
	// time_limit := 10 // 2 min

	broker := os.Getenv("KAFKA_BROKER")
	redis_addr := os.Getenv("REDIS_ADDR")

	if broker == "" {
		broker = "localhost:9092"
		logger.Println("KAFKA_BROKER not set, using default %s\n", broker)
	}

	brokers := []string{broker}
	// Create a new Sarama configuration
	username := os.Getenv("KAFKA_USERNAME")
	password := os.Getenv("KAFKA_PASSWORD")

	if username == "" || password == "" {
		// logger.Fatalf("KAFKA_USERNAME or KAFKA_PASSWORD not set\n")
		logger.Println("KAFKA_USERNAME or KAFKA_PASSWORD not set, using KAFKA locally")
	}

	config := kafka.NewConfig(username, password)

	producer, err := kafka.NewProducer(brokers, config, logger)
	if err != nil {
		logger.Fatalf("Error creating producer: %s\n", err)
	}
	defer producer.Close()

	if redis_addr == "" {
		logger.Println("REDIS_ADDR not set, using default localhost:6379")
		redis_addr = "localhost:6379"
	}

	redis, redisErr := store.NewRedisStore(ctx, redis_addr)

	if redisErr != nil {
		logger.Fatalf("Error creating redis store: %s\n", redisErr)
	}

	defer redis.Close()

	signalChan := make(chan os.Signal, 2)
	signal.Notify(signalChan, os.Interrupt)

	stream := os.Getenv("STREAM")
	if stream == "" {
		stream = "alerts"
	}

	// Creating Alerts
	logger.Println("Creating alerts")

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	var wg sync.WaitGroup
	wg.Add(1)

	go sendHeartBeatToRedis(ctx, redis, NodeID, signalChan, logger, &wg)

	for {
		select {
		case <-ticker.C:
			alert := alerts.GenRandomAlert(NodeID)
			if err := producer.SendAlert("alerts", &alert); err != nil {
				logger.Printf("Error sending alert: %s\n", err)
			}
			logger.Printf("Sent alert: %v\n", alert)
			redis.PublishAlertInputs(ctx, &alert, stream)

			newDuration := time.Duration(rand.Intn(30)) * time.Second // Random duration between 1 and 30 seconds
			ticker.Stop()
			ticker = time.NewTicker(newDuration)

		case <-signalChan:
			logger.Println("Received signal to stop")
			return
		}
	}

	defer wg.Wait()
}

func sendHeartBeatToRedis(ctx context.Context, redis *store.RedisStore, NodeID string, signalChan chan os.Signal, logger *log.Logger, wg *sync.WaitGroup) {
	defer wg.Done()
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := redis.StoreHeartBeat(ctx, NodeID, alerts.NewRuntimeMetrics(), logger)
			if err != nil {
				logger.Printf("Error sending heartbeat: %s\n", err)
			}
			logger.Printf("Sent heartbeat: %s\n", NodeID)

		case <-signalChan:
			logger.Println("Received signal to stop")
			return
		}
	}
}
