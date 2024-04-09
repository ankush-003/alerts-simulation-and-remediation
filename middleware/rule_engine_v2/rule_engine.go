package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	rule_engine "github.com/ankush-003/alerts-simulation-and-remediation/middleware/rule_engine_v2/engine"
	"github.com/ankush-003/alerts-simulation-and-remediation/middleware/sim/alerts"
	"github.com/ankush-003/alerts-simulation-and-remediation/middleware/sim/kafka"
	"github.com/joho/godotenv"
)

type AlertContext struct {
	AlertInput  *alerts.AlertInput
	AlertOutput *alerts.AlertOutput
	AlertParam  *alerts.ParamInput
}

func (alertContext *AlertContext) RuleInput() rule_engine.RuleInput {
	return alertContext.AlertInput
}

func (alertContext *AlertContext) RuleOutput() rule_engine.RuleOutput {
	return alertContext.AlertOutput
}

func (alertContext *AlertContext) ParamInput() rule_engine.ParamInput {
	return *alertContext.AlertParam
}

func NewAlert(alertInput *alerts.AlertInput, ruleEngineSvc *rule_engine.RuleEngineSvc) {
	defer wg.Done()

	alertContext := AlertContext{
		alertInput,
		&alerts.AlertOutput{Remedy: "Too be decided soon...", Severity: "NIL"},
		&alertInput.Params,
	}

	err := ruleEngineSvc.Execute(&alertContext)
	if err != nil {
		panic(err)
	}
	// Methods after parsing the alert
	printStruct(*alertInput)
	fmt.Println("Severity -> ", alertContext.AlertOutput.Severity)
	fmt.Println("Remedy -> ", alertContext.AlertOutput.Remedy)
	
	// Find the user associated with alertContext.AlertInput.source Node
	// Call Rest server notification handler
	// Call mail server

}

var wg sync.WaitGroup

func main() {

	ruleEngineSvc := rule_engine.NewRuleEngineSvc()

	alertA := alerts.AlertInput{
		ID:        "ID1",
		Category:  "Memory",
		Source:    "Hardware",
		Origin:    "NodeA",
		Params:    &alerts.Memory{Usage: 76, PageFaults: 30, SwapUsage: 2},
		CreatedAt: time.Now(),
		Handled:   false,
	}

	alertB := alerts.AlertInput{
		ID:        "ID2",
		Category:  "CPU",
		Source:    "Hardware",
		Origin:    "NodeA",
		Params:    &alerts.CPU{Utilization: 40, Temperature: 65},
		CreatedAt: time.Now(),
		Handled:   false,
	}

	wg.Add(1)
	wg.Add(1)
	wg.Add(1)

	go NewAlert(&alertA, ruleEngineSvc)
	go NewAlert(&alertB, ruleEngineSvc)
	go kafka_consumer(ruleEngineSvc)

	wg.Wait()

}

func kafka_consumer(ruleEngineSvc *rule_engine.RuleEngineSvc) {
	defer wg.Done()
	// loading .env file
	err_load := godotenv.Load()
	if err_load != nil {
		log.Fatalf("Error loading .env file: %v\n", err_load)
	}

	logger := log.New(os.Stdout, "kafka-consumer: ", log.LstdFlags)

	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		broker = "localhost:9092"
		logger.Println("KAFKA_BROKER not set, using default ", broker)
	}
	brokers := []string{broker}
	username := os.Getenv("KAFKA_USERNAME")
	password := os.Getenv("KAFKA_PASSWORD")

	config := kafka.NewConfig(username, password)

	consumer, err := kafka.NewConsumer(brokers, config, logger)
	if err != nil {
		logger.Fatalf("Error creating consumer: %s\n", err)
	}

	defer consumer.Close()

	alertsChan := make(chan alerts.AlertInput)
	doneChan := make(chan struct{})

	logger.Println("Consuming alerts !")

	go consumer.ConsumeAlerts("alerts", alertsChan, doneChan)

consumerLoop:
	for {
		select {
		case alert := <-alertsChan:
			printStruct(alert)
			wg.Add(1)
			NewAlert(&alert, ruleEngineSvc)
		case <-doneChan:
			break consumerLoop
		}
	}
}

func printStruct(alert alerts.AlertInput) {
	fmt.Println("ID: ", alert.ID)
	fmt.Println("Category: ", alert.Category)
	fmt.Println("Origin: ", alert.Origin)
	fmt.Println("Source: ", alert.Source)
	fmt.Println("Params: ", alert.Params)
}
