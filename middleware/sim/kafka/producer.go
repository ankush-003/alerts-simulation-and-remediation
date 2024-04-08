package kafka

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"

	"asmr/rule_engine"

	"github.com/IBM/sarama"
)

type Producer struct {
	Producer sarama.AsyncProducer
	Logger   *log.Logger
}

func NewProducer(brokers []string, config *sarama.Config, logger *log.Logger) (*Producer, error) {
	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		log.Printf("Error creating producer: %s\n", err)
		return &Producer{}, fmt.Errorf("Error creating producer: %s", err)
	}

	return &Producer{
		Producer: producer,
		Logger:   logger,
	}, nil
}

func (p *Producer) Close() {
	if err := p.Producer.Close(); err != nil {
		p.Logger.Printf("Error closing producer: %s\n", err)
	}
}

func (p *Producer) SendAlert(topic string, alert *rule_engine.AlertInput) error {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	alertJSON, err := json.Marshal(alert)
	if err != nil {
		p.Logger.Printf("Error marshalling alert: %s\n", err)
		return fmt.Errorf("Error marshalling alert: %s", err)
	}

	select {
	case p.Producer.Input() <- &sarama.ProducerMessage{Topic: topic, Value: sarama.StringEncoder(alertJSON)}:
		p.Logger.Printf("Produced alert: %s\n", alertJSON)
		return nil

	case err := <-p.Producer.Errors():
		p.Logger.Printf("Failed to send alert: %s\n", err)
		return fmt.Errorf("Failed to send alert: %s", err)

	case <-signals:
		p.Logger.Printf("Interrupted\n")
		return fmt.Errorf("Interrupted")
	}

	// return nil
}
