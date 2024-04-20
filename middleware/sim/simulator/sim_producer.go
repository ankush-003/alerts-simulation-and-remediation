package main

import (
	"log"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"github.com/ankush-003/alerts-simulation-and-remediation/middleware/rule_engine_v2/mongo"
	"github.com/ankush-003/alerts-simulation-and-remediation/middleware/sim/alerts"
	"github.com/ankush-003/alerts-simulation-and-remediation/middleware/sim/kafka"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	client, ctx, cancel, err := mongo.Connect(os.Getenv("MONGO_URI"))
	if err != nil {
		panic(err)
	}
	defer mongo.Close(client, ctx, cancel)
	db := client.Database("AlertSimAndRemediation")
	collection := db.Collection("Nodes")
	cursor, err := collection.Find(ctx, bson.M{}, options.Find().SetProjection(bson.M{"node": 1}))
	if err != nil {
		panic(err)
	}
	var nodes []string
	for cursor.Next(ctx) {
		var document bson.M
		if err = cursor.Decode(&document); err != nil {
			panic(err)
		}
		nodes = append(nodes, document["node"].(string))
	}

	// create a timer to send the alert every 5 seconds
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	for {
		select {
		case <-ticker.C:
			alert := alerts.GenRandomAlert(nodes)
			if err := producer.SendAlert("alerts", &alert); err != nil {
				logger.Printf("Error sending alert: %s\n", err)
			}
			logger.Printf("Sent alert: %v\n", alert)
			newDuration := time.Duration(rand.Intn(5)+1) * time.Second // Random duration between 1 and 5 seconds
			ticker.Stop()
			ticker = time.NewTicker(newDuration)
		case <-signalChan:
			logger.Printf("Interrupted\n")
			return
		}
	}
}
