package kafka

import (
	"context"
	"errors"
	"github.com/IBM/sarama"
	"github.com/ankush-003/alerts-simulation-and-remediation/middleware/sim/alerts"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type ConsumerGroupHandler struct {
	lock      sync.Mutex
	ready     chan bool
	alertChan chan alerts.AlertInput
	logger    *log.Logger
}

type ConsumerGroup struct {
	client  sarama.ConsumerGroup
	handler ConsumerGroupHandler
}

func NewConsumerGroup(brokers []string, groupID string, topics []string, config *sarama.Config, logger *log.Logger) (*ConsumerGroup, error) {
	sarama.Logger = log.New(os.Stdout, "[kafka] ", log.LstdFlags)

	config.Consumer.Offsets.AutoCommit.Enable = false

	client, err := sarama.NewConsumerGroup(brokers, groupID, config)

	if err != nil {
		logger.Panic("Error creating consumer group: ", err)
		return nil, err
	}

	cg := &ConsumerGroup{
		client: client,
		handler: ConsumerGroupHandler{
			ready:  make(chan bool),
			logger: logger,
			lock:   sync.Mutex{},
		},
	}

	return cg, nil
}

func (cg *ConsumerGroup) Consume(ctx context.Context, topics []string, alertsChan chan alerts.AlertInput, doneChan chan struct{}) {

	cur_ctx, cancel := context.WithCancel(ctx)

	cg.handler.alertChan = alertsChan

	keepRunning := true
	consumptionIsPaused := false

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			if err := cg.client.Consume(cur_ctx, topics, &cg.handler); err != nil {
				if errors.Is(err, sarama.ErrClosedConsumerGroup) {
					return
				}
				log.Panicf("Error from consumer: %v", err)
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
			cg.handler.ready = make(chan bool)
		}
	}()

	<-cg.handler.ready // Await till the consumer has been set up
	cg.handler.logger.Println("Sarama consumer up and running!...")

	sigusr1 := make(chan os.Signal, 1)
	signal.Notify(sigusr1, syscall.SIGUSR1)

	for keepRunning {
		select {
		case <-ctx.Done():
			cg.handler.logger.Println("terminating: via context")
			keepRunning = false
		case <-sigusr1:
			toggleConsumptionFlow(cg.client, &consumptionIsPaused)

		case <-doneChan:
			cg.handler.logger.Println("terminating: via doneChan")
		}
	}

	cancel()
	wg.Wait()

	close(doneChan)
}

func (cg *ConsumerGroup) Close(logger *log.Logger) {
	if err := cg.client.Close(); err != nil {
		logger.Printf("Error closing consumer group: %s\n", err)
	}
	logger.Println("Consumer group closed")
}

func toggleConsumptionFlow(client sarama.ConsumerGroup, isPaused *bool) {
	if *isPaused {
		client.ResumeAll()
		log.Println("Resuming consumption")
	} else {
		client.PauseAll()
		log.Println("Pausing consumption")
	}

	*isPaused = !*isPaused
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *ConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(consumer.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *ConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
// Once the Messages() channel is closed, the Handler must finish its processing
// loop and exit.
func (consumer *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/IBM/sarama/blob/main/consumer_group.go#L27-L29
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				consumer.logger.Println("Claim messages channel closed")
				return nil
			}
			consumer.lock.Lock()

			consumer.SendAlert(message.Value)

			consumer.lock.Unlock()
			// consumer.logger.Printf("Consumed message: topic=%s, partition=%d, offset=%d, key=%s, value=%s\n", message.Topic, message.Partition, message.Offset, message.Key, string(message.Value))
			session.MarkMessage(message, "")

		case <-session.Context().Done():
			return nil
		}
	}
}

func (consumer *ConsumerGroupHandler) SendAlert(value []byte) {
	// consumer.logger.Printf("Consumed message: value=%s\n", string(value))
	var parsedAlert alerts.AlertInput

	if err := parsedAlert.Unmarshal(value); err != nil {
		consumer.logger.Println("ERR: ", err)
	}

	consumer.alertChan <- parsedAlert
	consumer.logger.Printf("Received alert: %s\n", parsedAlert.ID)

}
