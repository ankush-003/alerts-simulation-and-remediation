package mongo

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// connect is a userdefined function returns mongo.Client,
// context.Context, context.CancelFunc and error.
// context.Context is used to set a deadline for process
// context.CancelFunc is used to cancel context and its resources

func Connect(uri string) (*mongo.Client, context.Context,
	context.CancelFunc, error,
) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	// ctx, cancel := context.WithTimeout(context.Background())
	// 30*time.Second)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	return client, ctx, cancel, err
}

// close is a userdef func to close resources
// closes MongoDB connection and cancel context

func Close(client *mongo.Client, ctx context.Context,
	cancel context.CancelFunc,
) {
	defer cancel()
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
}

func FindUsers(category, severity string) ([]string, error) {
	MONGO_URI := os.Getenv("MONGODB_URI")
	client, ctx, cancelFunc, err := Connect(MONGO_URI)
	if err != nil {
		panic(err)
	}
	defer Close(client, ctx, cancelFunc)
	db := client.Database("AlertSimAndRemediation")
	collection := db.Collection("Users")
	filter := bson.M{
		"alert.categories": category,
		"alert.severities": severity,
	}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("error finding users: %w", err)
	}
	var results []string
	for cursor.Next(ctx) {
		var user map[string]interface{}
		err := cursor.Decode(&user)
		if err != nil {
			return nil, fmt.Errorf("error decoding user document: %w", err)
		}
		results = append(results, user["email"].(string))
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("error iterating cursor: %w", err)
	}

	return results, nil
}

func GetRules() ([]byte, error) {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	MONGO_URI := os.Getenv("MONGODB_URI")
	client, ctx, cancelFunc, err := Connect(MONGO_URI)
	if err != nil {
		panic(err)
	}
	defer Close(client, ctx, cancelFunc)

	db := client.Database("AlertSimAndRemediation")
	collection := db.Collection("Rules")
	curr, err := collection.Find(ctx, bson.M{}, options.Find().SetProjection(bson.D{{Key: "_id", Value: 0}}))
	if err != nil {
		return nil, err
	}
	var results []bson.M
	curr.All(ctx, &results)
	rules, err := json.Marshal(results)
	if err != nil {
		return nil, err
	}
	return rules, nil
}
