package store

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ankush-003/alerts-simulation-and-remediation/middleware/sim/alerts"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
	"time"
)

type RedisStore struct {
	Client *redis.Client
	Logger *log.Logger
}

func NewRedisStore(ctx context.Context, addr string) (*RedisStore, error) {

	opt, _ := redis.ParseURL(addr)
	// add dial timeout
	opt.DialTimeout = 10 * time.Second

	client := redis.NewClient(opt)

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("error connecting to redis: %s", err)
	}

	return &RedisStore{
		Client: client,
		Logger: log.New(os.Stdout, "[redis] ", log.LstdFlags),
	}, nil
}

func (r *RedisStore) Close() error {
	return r.Client.Close()
}

// Redis Stream functions
func (r *RedisStore) PublishData(ctx context.Context, data map[string]interface{}, stream string) error {
	entry := &redis.XAddArgs{
		Stream: stream,
		Values: data,
	}

	txn := r.Client.TxPipeline()

	txn.XAdd(ctx, entry)

	result, err := txn.Exec(ctx)
	if err != nil {
		txn.Discard()
		r.Logger.Printf("Error publishing data to Redis Stream: %s\n", err)
		return fmt.Errorf("error publishing data to Redis Stream: %s", err)
	}

	r.Logger.Printf("Published data to Redis Stream: %s\n", result)
	return nil
}

func (r *RedisStore) PublishAlerts(ctx context.Context, alert *alerts.AlertInput) error {
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

	r.Logger.Printf("Published alert to Redis Stream: %s\n", result)
	return nil
}

func (r *RedisStore) PublishAlertInputs(ctx context.Context, alert *alerts.AlertInput, stream string) error {
	alertBytes, err := json.Marshal(alert)
	if err != nil {
		return err
	}

	// Create a Redis Streams entry with the alert data
	entry := &redis.XAddArgs{
		Stream: stream,
		Values: map[string]interface{}{
			"alert": alertBytes,
		},
	}

	// Publish the entry to the Redis Stream
	result, err := r.Client.XAdd(ctx, entry).Result()
	if err != nil {
		return err
	}

	r.Logger.Printf("Published alert to Redis Stream: %s\n", result)
	return nil
}

func (r *RedisStore) ConsumeData(ctx context.Context, stream string, dataChan chan<- map[string]interface{}, doneChan chan struct{}) {
	lastID := "$"
	for {
		select {
		case <-doneChan:
			return
		default:
			streamData, err := r.Client.XRead(ctx, &redis.XReadArgs{
				Streams: []string{stream, lastID},
				Count:   1,
				Block:   0,
			}).Result()
			if err != nil {
				r.Logger.Printf("Error reading from stream: %s\n", err)
				continue
			}

			for _, stream := range streamData {
				for _, message := range stream.Messages {
					r.Logger.Println("received message: %s\n", message.ID)
					for k, v := range message.Values {
						r.Logger.Printf("Key: %s, Value: %s\n\t", k, v)
					}
					dataChan <- message.Values
					lastID = message.ID
				}
			}
		}
	}
}

func (r *RedisStore) ConsumeDataGroup(ctx context.Context, stream string, dataChan chan<- map[string]interface{}, doneChan chan struct{}, groupName string) {
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
				r.Logger.Printf("Error reading from stream: %s\n", err)
				continue
			}

			for _, stream := range streamData {
				for _, message := range stream.Messages {
					dataChan <- message.Values
				}
			}
		}
	}
}

func (r *RedisStore) ConsumeAlertInputs(ctx context.Context, alertsChan chan<- alerts.AlertInput, doneChan chan struct{}, stream string) {
	lastID := "$"
	for {
		select {
		case <-doneChan:
			return
		default:
			streamData, err := r.Client.XRead(ctx, &redis.XReadArgs{
				Streams: []string{stream, lastID},
				Count:   1,
				Block:   0,
			}).Result()
			if err != nil {
				fmt.Printf("Error reading from stream: %s\n", err)
				continue
			}

			for _, stream := range streamData {
				for _, message := range stream.Messages {
					var alert alerts.AlertInput
					err := json.Unmarshal([]byte(message.Values["alert"].(string)), &alert)
					if err != nil {
						fmt.Printf("Error unmarshaling alert: %s\n", err)
						continue
					}

					alertsChan <- alert
					lastID = message.ID
				}
			}
		}
	}
}

func (r *RedisStore) StoreHeartBeat(ctx context.Context, NodeID string, metrics *alerts.RuntimeMetrics, logger *log.Logger) error {
	values := map[string]interface{}{
		"nodeID":       NodeID,
		"numGoroutine": metrics.NumGoroutine,
		"cpuUsage":     metrics.CpuUsage,
		"ramUsage":     metrics.RamUsage,
		"status":       "UP",
	}

	txn := r.Client.TxPipeline()

	txn.HSet(ctx, NodeID, values)

	if _, err := txn.Exec(ctx); err != nil {
		txn.Discard()
		return fmt.Errorf("error storing heartbeat: %s", err)
	}

	logger.Printf("Published heartbeat to Redis Stream: %s\n %s\n", NodeID, values)

	return nil
}

func (r *RedisStore) StreamHeartBeat(ctx context.Context, NodeID string, metrics *alerts.RuntimeMetrics, logger *log.Logger) error {
	values := map[string]interface{}{
		"nodeID":       NodeID,
		"numGoroutine": metrics.NumGoroutine,
		"cpuUsage":     metrics.CpuUsage,
		"ramUsage":     metrics.RamUsage,
		"status":       "UP",
	}

	if err := r.PublishData(ctx, values, "nodeHeartbeats"); err != nil {
		return fmt.Errorf("error streaming heartbeat: %s", err)
	}

	logger.Printf("Published heartbeat to Redis Stream: %s\n %s\n", NodeID, values)
	return nil
}

func (r *RedisStore) KillHeartBeat(ctx context.Context, NodeID string, logger *log.Logger) error {
	values := map[string]interface{}{
		"status": "DOWN",
	}

	txn := r.Client.TxPipeline()

	txn.HSet(ctx, NodeID, values)

	if _, err := txn.Exec(ctx); err != nil {
		txn.Discard()
		return fmt.Errorf("error killing heartbeat: %s", err)
	}

	logger.Printf("Killed heartbeat to Redis Stream: %s\n %s\n", NodeID, values)

	return nil
}
