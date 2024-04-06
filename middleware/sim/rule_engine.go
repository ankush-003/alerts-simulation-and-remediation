package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"asmr/kafka"
	"asmr/rule_engine"

	"github.com/joho/godotenv"
)


type AlertContext struct {
	AlertInput  *rule_engine.AlertInput
	AlertOutput *rule_engine.AlertOutput
	AlertParam  *rule_engine.ParamInput
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

func NewAlert(alertInput *rule_engine.AlertInput, ruleEngineSvc *rule_engine.RuleEngineSvc) {
	defer wg.Done()

	alertContext := AlertContext{
		alertInput,
		&rule_engine.AlertOutput{Remedy: "Too be decided soon...", Severity: "NIL"},
		&alertInput.Params,
	}

	err := ruleEngineSvc.Execute(&alertContext)
	if err != nil {
		panic(err)
	}
	fmt.Println("Here")
	// Methods after parsing the alert
	fmt.Println("Alert -> ", alertInput)
	fmt.Println("Severity -> ", alertContext.AlertOutput.Severity)
	fmt.Println("Remedy -> ", alertContext.AlertOutput.Remedy)
	// Find the user associated with alertContext.AlertInput.source Node
	// Call Rest server notification handler
	// Call mail server

}

var wg sync.WaitGroup

func main() {
	ruleEngineSvc := rule_engine.NewRuleEngineSvc()

	alertA := rule_engine.AlertInput{
		ID:        "ID1",
		Category:  "Memory",
		Source:    "Hardware",
		Origin:    "NodeA",
		Params:    &rule_engine.Memory{Usage: 76, PageFaults: 30, SwapUsage: 2},
		CreatedAt: time.Now(),
		Handled:   false,
	}

	alertB := rule_engine.AlertInput{
		ID:        "ID2",
		Category:  "CPU",
		Source:    "Hardware",
		Origin:    "NodeA",
		Params:    &rule_engine.CPU{Utilization: 40, Temperature: 65},
		CreatedAt: time.Now(),
		Handled:   false,
	}

	wg.Add(1)
	wg.Add(1)

	go NewAlert(&alertA, ruleEngineSvc)
	go NewAlert(&alertB, ruleEngineSvc)

	wg.Wait()

}

func kafka_consumer() {

	// loading .env file
	err_load := godotenv.Load()
	if err_load != nil {
		log.Fatalf("Error loading .env file: %v\n", err_load)
	}

	logger := log.New(os.Stdout, "kafka-consumer: ", log.LstdFlags)

	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		broker = "localhost:9092"
		logger.Println("KAFKA_BROKER not set, using default %s\n", broker)
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

	alertsChan := make(chan rule_engine.AlertInput)
	doneChan := make(chan struct{})

	logger.Println("Consuming alerts !")

	go consumer.ConsumeAlerts("alerts", alertsChan, doneChan)
consumerLoop:

	for {
		select {
		case alert := <-alertsChan:
			logger.Printf("Received alert: alrtID: %s, NodeID: %s, Description: %s, Severity: %s, Source: %s, CreatedAt: %s\t", alert.ID.String(), alert.NodeID.String(), alert.Description, alert.Severity, alert.Source, alert.CreatedAt)
			logger.Printf("RuntimeMetrics: NumGoroutine: %d, CpuUsage: %f, RamUsage: %f\n\n", alert.RuntimeMetrics.NumGoroutine, alert.RuntimeMetrics.CpuUsage, alert.RuntimeMetrics.RamUsage)

		case <-doneChan:
			break consumerLoop
		}
	}
}
