package store

import (
	"github.com/ankush-003/alerts-simulation-and-remediation/middleware/sim/alerts"

	"context"
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

	return &MongoStore{Client: client, Coll: collection}, nil
}

func (s *MongoStore) Close(ctx context.Context) error {
	return s.Client.Disconnect(ctx)
}

func (s *MongoStore) FetchAlertConfigs(ctx context.Context) ([]alerts.AlertConfig, error) {
	cur, err := s.Coll.Find(ctx, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("could not fetch alert configs: %v", err)
	}
	defer cur.Close(ctx)

	var alertConfigs []alerts.AlertConfig
	for cur.Next(ctx) {
		var alertConfig alerts.AlertConfig
		if err := cur.Decode(&alertConfig); err != nil {
			return nil, fmt.Errorf("cur.Decode: %v", err)
		}
		alertConfigs = append(alertConfigs, alertConfig)
	}
	if err := cur.Err(); err != nil {
		return nil, fmt.Errorf("error in cursor: %v", err)
	}

	return alertConfigs, nil
}

func (s *MongoStore) InsertAlertConfigs(ctx context.Context, alertConfigs []alerts.AlertConfig) error {
	for _, alertConfig := range alertConfigs {
		_, err := s.Coll.InsertOne(ctx, alertConfig)
		if err != nil {
			return fmt.Errorf("could not insert alert config: %v", err)
		}
	}

	return nil
}
