package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ankush-003/alerts-simulation-and-remediation/middleware/sim/store"
	"github.com/joho/godotenv"
)

func main() {
	err_load := godotenv.Load()
	if err_load != nil {
		log.Fatalf("Error loading .env file: %s\n", err_load)
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

	node_id, close_id, err := mongo_client.GetNodeId(ctx)
	defer close_id()

	if err != nil {
		log.Fatalf("Error getting node id: %s\n", err)
	}

	fmt.Printf("Node id: %s\n", node_id)
	time.Sleep(30 * time.Second)
}
