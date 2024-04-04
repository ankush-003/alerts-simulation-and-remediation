package mongo

import (
	"context"
	"encoding/json"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// connect is a userdefined function returns mongo.Client,
// context.Context, context.CancelFunc and error.
// context.Context is used to set a deadline for process
// context.CancelFunc is used to cancel context and its resources

func Connect(uri string) (*mongo.Client, context.Context,
	context.CancelFunc, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	// ctx, cancel := context.WithTimeout(context.Background())
	// 30*time.Second)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	return client, ctx, cancel, err
}

// close is a userdef func to close resources
// closes MongoDB connection and cancel context

func Close(client *mongo.Client, ctx context.Context,
	cancel context.CancelFunc) {

	defer cancel()
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
}

func GetRules() ([]byte, error) {
	client, ctx, cancelFunc, err := Connect("mongodb://127.0.0.1:27017")
	if err != nil {
		panic(err)
	}
	defer Close(client, ctx, cancelFunc)

	db := client.Database("hpe_cty")
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
