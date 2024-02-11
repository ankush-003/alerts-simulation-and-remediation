package kafka

import (
	"asmr/alerts"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"log"
	"os"
	"os/signal"
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

func (c *Consumer) ConsumeAlerts(topic string, alertsChan chan alerts.Alerts, doneChan chan struct{}) {
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
			var alert alerts.Alerts
			err := json.Unmarshal(msg.Value, &alert)
			if err != nil {
				c.Logger.Printf("Error unmarshalling alert: %s\n", err)
				continue
			}
			// c.Logger.Printf("Consumed alert: %v\n", alert)
			alertsChan <- alert
		case <-signals:
			c.Logger.Printf("Interrupted\n")
			break consumerLoop
		}
	}
}
