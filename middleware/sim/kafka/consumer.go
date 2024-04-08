package kafka

import (
	rule_engine "github.com/ankush-003/alerts-simulation-and-remediation/middleware/rule_engine_v2/engine"

	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/IBM/sarama"
)

type Consumer struct {
	Consumer sarama.Consumer
	Logger   *log.Logger
}

func NewConsumer(brokers []string, config *sarama.Config, logger *log.Logger) (*Consumer, error) {
	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Printf("Error creating consumer: %s\n", err)
		return &Consumer{}, fmt.Errorf("Error creating consumer: %s", err)
	}

	return &Consumer{
		Consumer: consumer,
		Logger:   logger,
	}, nil
}

func (c *Consumer) Close() {
	if err := c.Consumer.Close(); err != nil {
		c.Logger.Printf("Error closing consumer: %s\n", err)
	}
}

func (c *Consumer) ConsumeAlerts(topic string, alertsChan chan rule_engine.AlertInput, doneChan chan struct{}) {
	partitionConsumer, err := c.Consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)

	if err != nil {
		c.Logger.Printf("Error creating partition consumer: %s\n", err)
		close(doneChan)
		return
	}

	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			c.Logger.Printf("Error closing partition consumer: %s\n", err)
		}
		close(doneChan)
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

consumerLoop:
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			var data map[string]interface{}
			err := json.Unmarshal(msg.Value, &data)
			if err != nil {
				c.Logger.Printf("Error unmarshalling alert: %s\n", err)
				continue
			}
			fmt.Println(data)
			paramsData := data["params"].(map[string]interface{})
			paramsType := data["category"].(string)

			var parsedAlert rule_engine.AlertInput
			parsedAlert.Category = data["category"].(string)
			parsedAlert.ID = data["id"].(string)
			parsedAlert.Source = data["source"].(string)
			parsedAlert.Handled = data["handled"].(bool)
			switch paramsType {
			case "Memory":
				var memory rule_engine.Memory
				paramsBytes, err := json.Marshal(paramsData) // Convert map to JSON bytes
				if err != nil {
					c.Logger.Printf("Error unmarshalling alert: %s\n", err)
					continue
				}
				err = json.Unmarshal(paramsBytes, &memory)
				if err != nil {
					c.Logger.Printf("Error unmarshalling alert: %s\n", err)
					continue
				}
				parsedAlert.Params = &memory
			case "CPU":
			default:
				c.Logger.Printf("Error unmarshalling alert: %s\n", err)
			}
			// err = json.Unmarshal(msg.Value, &alert)
			// c.Logger.Printf("Consumed alert: %v\n", parsedAlert)
			alertsChan <- parsedAlert
		case <-signals:
			c.Logger.Printf("Interrupted\n")
			break consumerLoop
		}
	}
}
