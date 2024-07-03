package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	rule_engine "github.com/ankush-003/alerts-simulation-and-remediation/middleware/rule_engine/engine"
	"github.com/ankush-003/alerts-simulation-and-remediation/middleware/rule_engine/mailserver"
	"github.com/ankush-003/alerts-simulation-and-remediation/middleware/rule_engine/mongo"
	"github.com/ankush-003/alerts-simulation-and-remediation/middleware/sim/alerts"
	"github.com/ankush-003/alerts-simulation-and-remediation/middleware/sim/kafka"
	"github.com/joho/godotenv"
)

func NewAlert(alertInput *alerts.AlertInput, ruleEngineSvc *rule_engine.RuleEngineSvc, cache rule_engine.Cache) {
	defer wg.Done()

	alertContext := rule_engine.AlertContext{
		AlertInput:  alertInput,
		AlertOutput: &alerts.AlertOutput{Remedy: "Too be decided soon...", Severity: "NIL"},
		AlertParam:  &alertInput.Params,
	}

	err := ruleEngineSvc.Execute(alertContext)
	if err != nil {
		panic(err)
	}

	// Methods after parsing the alert
	if alertContext.AlertOutput.Severity == "Safe" {
		fmt.Println("Safe event")
		return
	} else if !rule_engine.CacheChecker(alertInput.Category, alertInput.Source, alertContext.AlertOutput, &cache) {
		rule_engine.PrintStruct(alertInput, alertContext.AlertOutput)
		fmt.Println("Not alert")
		return
	}
	rule_engine.PrintStruct(alertInput, alertContext.AlertOutput)

	// Find the user associated with alertContext.AlertInput.source Node
	emails, err := mongo.FindUsers(alertInput.Category, alertContext.AlertOutput.Severity)
	if err != nil {
		// panic(err)
		fmt.Println("Error in finding user")
	}

	// Call mail server
	for _, email := range emails {
		if err = mailserver.SendEmail(*alertInput, *alertContext.AlertOutput, email); err != nil {
			fmt.Println(err)
		}
	}

	// Call Rest server notification handler
	alertContext.NotifyRestServer()
}

var wg sync.WaitGroup

func main() {

	ctx := context.Background()

	ruleEngineSvc := rule_engine.NewRuleEngineSvc()

	// alertA := alerts.AlertInput{
	// 	ID:        "ID1",
	// 	Category:  "Memory",
	// 	Source:    "Hardware",
	// 	Origin:    "NodeB",
	// 	Params:    &alerts.Memory{Usage: 76, PageFaults: 30, SwapUsage: 2},
	// 	CreatedAt: time.Now().Format(time.DateTime),
	// 	Handled:   false,
	// }

	// alertB := alerts.AlertInput{
	// 	ID:        "ID2",
	// 	Category:  "CPU",
	// 	Source:    "Hardware",
	// 	Origin:    "NodeA",
	// 	Params:    &alerts.CPU{Utilization: 40, Temperature: 65},
	// 	CreatedAt: time.Now().Format(time.DateTime),
	// 	Handled:   false,
	// }
	
	kafka_consumer(ctx, ruleEngineSvc)
}

func kafka_consumer(ctx context.Context, ruleEngineSvc *rule_engine.RuleEngineSvc) {
	// loading .env file
	err_load := godotenv.Load()
	if err_load != nil {
		log.Println("Error loading .env file")
	}

	logger := log.New(os.Stdout, "[kafka-consumer] ", log.LstdFlags)

	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		broker = "localhost:9092"
		logger.Println("KAFKA_BROKER not set, using default ", broker)
	}
	brokers := []string{broker}
	username := os.Getenv("KAFKA_USERNAME")
	password := os.Getenv("KAFKA_PASSWORD")

	config := kafka.NewConfig(username, password)
	groupID := "rule-engine-consumer"
	topics := []string{"alerts"}

	consumerGroup, err := kafka.NewConsumerGroup(brokers, groupID, topics, config, logger)
	if err != nil {
		logger.Fatalf("Error creating consumer: %s\n", err)
	}

	defer consumerGroup.Close(logger)

	alertsChan := make(chan alerts.AlertInput)
	doneChan := make(chan struct{})
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT)

	logger.Println("Consuming alerts !")

	go consumerGroup.Consume(ctx, topics, alertsChan, doneChan)
	cache := rule_engine.Cache{}
	cache = cache.New()

consumerLoop:
	for {
		select {
		case alert := <-alertsChan:
			rule_engine.PrintStruct(&alert, nil)
			wg.Add(1)
			NewAlert(&alert, ruleEngineSvc, cache)
		case <-signalChan:
			logger.Println("Interrupted")
			close(doneChan)
			break consumerLoop
		}
	}
}
