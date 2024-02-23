package store

import (
	"asmr/alerts"
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	Client *redis.Client
}

func NewRedisStore(addr string) *RedisStore {
	return &RedisStore{
		Client: redis.NewClient(&redis.Options{
			Addr: addr,
		}),
	}
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

	txn.Set(ctx, key, string(data), 0)

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

func (r *RedisStore) GetAlertsByNodeID(ctx context.Context, nodeID string) ([]alerts.Alerts, error) {
    keys, err := r.Client.Keys(ctx, "*").Result() // Get all keys in the Redis database
    if err != nil {
        return nil, fmt.Errorf("error getting keys: %s", err)
    }

    var alertsList []alerts.Alerts

    for _, key := range keys {
		

        data, err := r.Client.Get(ctx, key).Result()
        if err != nil {
            return nil, fmt.Errorf("error getting alert: %s", err)
        }

        var alert alerts.Alerts
        if err := json.Unmarshal([]byte(data), &alert); err != nil {
            return nil, fmt.Errorf("error unmarshalling alert: %s", err)
        }
		fmt.Println("data",data)

        // Check if the alert's nodeID matches the specified nodeID
        if alert.NodeID.String() == nodeID {
            alertsList = append(alertsList, alert)
        }
		
    }

    return alertsList, nil
}

