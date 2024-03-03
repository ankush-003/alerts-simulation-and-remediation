package main

import (
	"asmr/alerts"
	"asmr/store"
	"context"
	// "encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	// loading .env file
	err_load := godotenv.Load()
	if err_load != nil {
		log.Fatalf("error loading .env file: %s\n", err_load)
	}

	// redis_addr := os.Getenv("REDIS_ADDR")
	mongo_uri := os.Getenv("MONGO_URI")
	if mongo_uri == "" {
		mongo_uri = "mongodb://localhost:27017"
		log.Println("MONGO_URI not set, using default %s\n", mongo_uri)
	}

	database := "asmr"
	coll := "alertConfigs"

	ctx := context.Background()
	// Create a new MongoStore
	mongoStore, err := store.NewMongoStore(ctx, mongo_uri, database, coll)
	if err != nil {
		panic(err)
	}

	// Fetch alert configs from the store
	alertConfigs, err := mongoStore.FetchAlertConfigs(ctx)
	if err != nil {
		panic(err)
	}

	if(len(alertConfigs) == 0) {
		log.Println("No alert configs found")
		log.Println("Creating alert configs")
		// Create predefined alert configs
		descriptions := []string{"High CPU usage", "Low disk space", "Network down", "Service unavailable"}
		severities := []string{"critical", "warning", "info"}

		for _, desc := range descriptions {
			for _, sev := range severities {
				alertConfig := alerts.NewAlertConfig(desc, sev)
				alertConfigs = append(alertConfigs, *alertConfig)
			}
		}

		// Insert the alert configs into the store
		err = mongoStore.InsertAlertConfigs(ctx, alertConfigs)
		if err != nil {
			panic(err)
		}
	}

	// Update Redis with the alert configs
	// redisStore := store.NewRedisStore(redis_addr)
	// for _, alertConfig := range alertConfigs {
	// 	err = redisStore.StoreAlertConfig(ctx, &alertConfig)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }

	for _, alertConfig := range alertConfigs {
		fmt.Printf("Alert Config: %v\n", alertConfig)
	}


	if err != nil {
		panic(err)
	}
}