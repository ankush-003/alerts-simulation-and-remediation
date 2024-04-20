package store

import (
	"context"
	"log"
	// "encoding/json"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoStore struct {
	Client *mongo.Client
	Coll   *mongo.Collection
}

func NewMongoStore(ctx context.Context, uri, db, coll string) (*MongoStore, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("mongo.Connect: %v", err)
	}

	collection := client.Database(db).Collection(coll)

	log.Println("Connected to MongoDB")

	return &MongoStore{Client: client, Coll: collection}, nil
}

func (s *MongoStore) Close(ctx context.Context) error {
	return s.Client.Disconnect(ctx)
}

func (s *MongoStore) GetNodeId(ctx context.Context) (string, func(), error) {
	// get a random node id to be assigned to the current simulator
	var result bson.M
	
	if err := s.Coll.FindOne(ctx, bson.M{"available": true}).Decode(&result); err != nil {
		return "", nil, fmt.Errorf("error finding available node: %v", err)
	}
	
	nodeId := result["node"].(string)

	_, err := s.Coll.UpdateOne(ctx, bson.M{"node": nodeId}, bson.M{"$set": bson.M{"available": false}})

	if err != nil {
		return "", nil, fmt.Errorf("error updating node availability: %v", err)
	}

	closeFunc := func() {
		_, err := s.Coll.UpdateOne(ctx, bson.M{"node": nodeId}, bson.M{"$set": bson.M{"available": true}})
		if err != nil {
			fmt.Printf("error updating node availability: %v\n", err)
		}
	}

	return nodeId, closeFunc, nil	
}