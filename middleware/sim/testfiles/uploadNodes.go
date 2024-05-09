package main

import (
	"os"
	"time"
	"github.com/ankush-003/alerts-simulation-and-remediation/middleware/sim/store"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"log"
	"context"
)

type Node struct {
	Node      string `json:"node"`
	CreatedAt string `json:"createdAt"`
	Available bool   `json:"avaliable"`
}

func main() {
	err_load := godotenv.Load()
	if err_load != nil {
		panic(err_load)
	}

	ctx := context.Background()

	mongo_uri := os.Getenv("MONGO_URI")
	if mongo_uri == "" {
		log.Fatalf("MONGO_URI not set\n")
	}
	mongo_client, err := store.NewMongoStore(ctx, mongo_uri, "AlertSimAndRemediation", "Nodes")
	if err != nil {
		log.Fatalf("Error creating mongo store: %s\n", err)
	}
	defer mongo_client.Close(ctx)
	db := mongo_client.Client.Database("AlertSimAndRemediation")
	collection := db.Collection("Nodes")
	for i := 0; i < 10; i++ {
		node := Node{
			Node:        uuid.New().String(),
			CreatedAt: time.Now().Format(time.DateTime),
			Available:   true,
		}
		_, err := collection.InsertOne(ctx, node)
		if err != nil {
			panic(err)
		}
	}
}
