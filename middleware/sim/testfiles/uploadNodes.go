package main

import (
	"os"
	"time"
	"github.com/ankush-003/alerts-simulation-and-remediation/middleware/rule_engine_v2/mongo"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
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

	client, ctx, cancel, err := mongo.Connect(os.Getenv("MONGO_URI"))
	if err != nil {
		panic(err)
	}
	defer mongo.Close(client, ctx, cancel)
	db := client.Database("AlertSimAndRemediation")
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
