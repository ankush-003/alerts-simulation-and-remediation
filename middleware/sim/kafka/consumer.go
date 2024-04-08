package kafka

import (
	"asmr/rule_engine"

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
			var parsedAlert rule_engine.AlertInput
			if err := parsedAlert.Unmarshal(msg.Value); err != nil {
				c.Logger.Println("ERR: ", err)
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