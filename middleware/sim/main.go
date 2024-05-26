package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
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
		log.Println("Error loading .env file")
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
	log.Println("Received Node ID: ", NodeID)
	defer close_id()
	if err != nil {
		log.Fatalf("Error getting node id: %s\n", err)
	}
	logger := log.New(os.Stdout, fmt.Sprintf("[Node %s] ", NodeID), log.LstdFlags)
	
	broker := os.Getenv("KAFKA_BROKER")
	redis_addr := os.Getenv("REDIS_ADDR")
	
	if broker == "" {
		broker = "localhost:9092"
		logger.Println("KAFKA_BROKER not set, using default", broker)
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
	logger.Println("Kafka Producer created")
	
	if redis_addr == "" {
		logger.Println("REDIS_ADDR not set, using default localhost:6379")
		redis_addr = "localhost:6379"
	}

	redis, redisErr := store.NewRedisStore(ctx, redis_addr)

	if redisErr != nil {
		logger.Fatalf("Error creating redis store: %s\n", redisErr)
	}
	defer redis.Close()
	logger.Println("Redis store created")

	logger.Println("Setup Complete")
	
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	heartbeatDoneChan := make(chan struct{})
	
	stream := os.Getenv("STREAM")
	if stream == "" {
		stream = "nodeAlerts"
	}
	
	// Creating Alerts
	logger.Println("Creating alerts")
	
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	var wg sync.WaitGroup
	wg.Add(1)

	go sendHeartBeatToRedis(ctx, redis, NodeID, heartbeatDoneChan, logger, &wg)
	defer wg.Wait()
	// defer redis.KillHeartBeat(ctx, NodeID, logger)

	for {
		select {
		case <-ticker.C:
			alert := alerts.GenRandomAlert(NodeID)
			times := rand.Int() % 6
			for i := 0; i < times; i++ {
				if err := producer.SendAlert("alerts", &alert); err != nil {
					logger.Printf("Error sending alert: %s\n", err)
				}
				logger.Printf("Sent alert: %v\n", alert)
				redis.PublishAlertInputs(ctx, &alert, stream)
			}
			newDuration := time.Duration(rand.Intn(30)+1) * time.Second // Random duration between 1 and 30 seconds
			ticker.Stop()
			ticker = time.NewTicker(newDuration)

		case <-signalChan:
			logger.Println("Received signal to stop")
			close(heartbeatDoneChan)
			return
		}
	}
}

func sendHeartBeatToRedis(ctx context.Context, redis *store.RedisStore, NodeID string, doneChan chan struct{}, logger *log.Logger, wg *sync.WaitGroup) {
	defer wg.Done()
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// err := redis.StoreHeartBeat(ctx, NodeID, alerts.NewRuntimeMetrics(), logger)
			err := redis.StreamHeartBeat(ctx, NodeID, alerts.NewRuntimeMetrics(), logger)

			if err != nil {
				logger.Printf("Error sending heartbeat: %s\n", err)
			}
			logger.Printf("Sent heartbeat: %s\n", NodeID)

		case <-doneChan:
			logger.Println("Received signal to stop Heartbeat")
			return
		}
	}
}
