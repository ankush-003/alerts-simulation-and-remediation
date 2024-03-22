package store

import (
	"context"
	"encoding/json"
	"asmr/alerts"
	"github.com/redis/go-redis/v9"
	"fmt"
)

type RedisStore struct {
	Client *redis.Client
}

func NewRedisStore(ctx context.Context, addr string) (*RedisStore, error) {

	opt, _ := redis.ParseURL(addr)	

	client := redis.NewClient(opt)

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("error connecting to redis: %s", err)
	}
	
	return &RedisStore{
		Client: client,
	}, nil
}

func (r *RedisStore) Close() error {
	return r.Client.Close()
}

func (r *RedisStore) StoreAlert(ctx context.Context, alert *alerts.Alerts) error {
	key := alert.ID.String()
	data, err := json.Marshal(alert)
	if err != nil {
		return fmt.Errorf("error marshalling alert: %s", err)
	}

	txn := r.Client.TxPipeline()

	// txn.Set(ctx, key, string(data), 0)
	txn.SetNX(ctx, key, string(data), 0)

	if _, err := txn.Exec(ctx); err != nil {
		txn.Discard()
		return fmt.Errorf("error storing alert: %s", err)
	}

	fmt.Printf("Stored alert: %s\n", string(data))
	return nil
}

func (r *RedisStore) StoreAlertConfig(ctx context.Context, alertConfig *alerts.AlertConfig) error {
	key := alertConfig.ID.String()
	data, err := json.Marshal(alertConfig)
	if err != nil {
		return fmt.Errorf("error marshalling alert: %s", err)
	}

	txn := r.Client.TxPipeline()

	txn.Set(ctx, key, string(data), 0)

	if _, err := txn.Exec(ctx); err != nil {
		txn.Discard()
		return fmt.Errorf("error storing alert: %s", err)
	}

	fmt.Printf("Stored alert: %s\n", string(data))
	return nil
}

func (r *RedisStore) GetAlertConfigByID(ctx context.Context, id string) (alerts.AlertConfig, error) {
	data, err := r.Client.Get(ctx, id).Result()
	if err != nil {
		return alerts.AlertConfig{}, fmt.Errorf("error getting alert: %s", err)
	}

	var alertConfig alerts.AlertConfig
	if err := json.Unmarshal([]byte(data), &alertConfig); err != nil {
		return alerts.AlertConfig{}, fmt.Errorf("error unmarshalling alert: %s", err)
	}

	return alertConfig, nil

}


func (r *RedisStore) GetAlertByID(ctx context.Context, id string) (alerts.Alerts, error) {
	data, err := r.Client.Get(ctx, id).Result()
	if err != nil {
		return alerts.Alerts{}, fmt.Errorf("error getting alert: %s", err)
	}

	var alert alerts.Alerts
	if err := json.Unmarshal([]byte(data), &alert); err != nil {
		return alerts.Alerts{}, fmt.Errorf("error unmarshalling alert: %s", err)
	}

	return alert, nil
}

func (r *RedisStore) GetRandomAlert(ctx context.Context) (alerts.Alerts, error) {
	key, err := r.Client.RandomKey(ctx).Result()
	if err != nil {
		return alerts.Alerts{}, fmt.Errorf("error getting random key: %s", err)
	}

	data, err := r.Client.Get(ctx, key).Result()
	if err != nil {
		return alerts.Alerts{}, fmt.Errorf("error getting alert: %s", err)
	}

	var alert alerts.Alerts
	if err := json.Unmarshal([]byte(data), &alert); err != nil {
		return alerts.Alerts{}, fmt.Errorf("error unmarshalling alert: %s", err)
	}

	return alert, nil
}

func (r *RedisStore) GetRandomAlertConfig(ctx context.Context) (alerts.AlertConfig, error) {
	key, err := r.Client.RandomKey(ctx).Result()
	if err != nil {
		return alerts.AlertConfig{}, fmt.Errorf("error getting random key: %s", err)
	}

	data, err := r.Client.Get(ctx, key).Result()
	if err != nil {
		return alerts.AlertConfig{}, fmt.Errorf("error getting alert: %s", err)
	}

	var alertConfig alerts.AlertConfig
	if err := json.Unmarshal([]byte(data), &alertConfig); err != nil {
		return alerts.AlertConfig{}, fmt.Errorf("error unmarshalling alert: %s", err)
	}

	return alertConfig, nil
}

func (r *RedisStore) PublishAlerts(ctx context.Context, alert *alerts.Alerts) error {
	alertBytes, err := json.Marshal(alert)
    if err != nil {
        return err
    }

    // Create a Redis Streams entry with the alert data
    entry := &redis.XAddArgs{
        Stream: "alerts",
        Values: map[string]interface{}{
            "alert": alertBytes,
        },
    }

    // Publish the entry to the Redis Stream
    result, err := r.Client.XAdd(ctx, entry).Result()
    if err != nil {
        return err
    }

    fmt.Printf("Published alert to Redis Stream: %s\n", result)
    return nil
}

func (r *RedisStore) ConsumeAlerts(ctx context.Context, alertsChan chan<- alerts.Alerts, doneChan chan struct{}, stream string, groupName string) {
    status, err := r.Client.XGroupCreate(ctx, stream, groupName, "0").Result()
    if err != nil {
        fmt.Printf("Error creating group: %s\n", err)
		// close(doneChan)
        // return
    } else {
        fmt.Printf("Group created: %s\n", status)
    }

    for {
        select {
        case <-doneChan:
            return
        default:
            streamData, err := r.Client.XReadGroup(ctx, &redis.XReadGroupArgs{
                Streams: []string{stream, ">"},
                Group:   groupName,
                Count:   1,
                NoAck:   true,
            }).Result()
            if err != nil {
                fmt.Printf("Error reading from stream: %s\n", err)
                continue
            }

            for _, stream := range streamData {
                for _, message := range stream.Messages {
                    var alert alerts.Alerts
                    err := json.Unmarshal([]byte(message.Values["alert"].(string)), &alert)
                    if err != nil {
                        fmt.Printf("Error unmarshaling alert: %s\n", err)
                        continue
                    }

                    alertsChan <- alert
                }
            }
        }
    }
}