package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"

	// "time"

	rule_engine "github.com/ankush-003/alerts-simulation-and-remediation/middleware/rule_engine/engine"
	"github.com/ankush-003/alerts-simulation-and-remediation/middleware/rule_engine/mailserver"
	"github.com/ankush-003/alerts-simulation-and-remediation/middleware/rule_engine/mongo"
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
	printStruct(alertInput, alertContext.AlertOutput)

	// Find the user associated with alertContext.AlertInput.source Node
	email, err := mongo.FindUser(alertInput.Origin)
	if err != nil {
		// panic(err)
		fmt.Println("Error in finding user")
	}

	// Call mail server
	if err = mailserver.SendEmail(*alertInput, *alertContext.AlertOutput, email); err != nil {
		fmt.Println(err)
	}

	// Call Rest server notification handler
	notifyRestServer(&alertContext)
}

func notifyRestServer(alertContext *AlertContext) {
	jsonData := map[string]interface{}{
		"ID":        alertContext.AlertInput.ID,
		"Category":  alertContext.AlertInput.Category,
		"Source":    alertContext.AlertInput.Source,
		"Origin":    alertContext.AlertInput.Origin,
		"CreatedAt": alertContext.AlertInput.CreatedAt,
		"Severity":  alertContext.AlertOutput.Severity,
		"Remedy":    alertContext.AlertOutput.Remedy,
	}

	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		panic(err)
	}

	host := os.Getenv("HOST")
	if host == "" {
		host = "localhost"
	}

	req, err := http.NewRequest("POST", "http://"+host+":8000/postRemedy", bytes.NewBuffer(jsonBytes))

	if err != nil {
		// panic(err)
		fmt.Println("Error in creating request")
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error in Connecting to Rest Server")
		return
		// panic(err)

	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// Convert the response body to a string and print it
	fmt.Println("Response Body:", string(body))
}

var wg sync.WaitGroup

func main() {

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

	// wg.Add(1)
	// wg.Add(1)
	wg.Add(1)

	// go NewAlert(&alertA, ruleEngineSvc)
	// go NewAlert(&alertB, ruleEngineSvc)
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
			printStruct(&alert, nil)
			wg.Add(1)
			NewAlert(&alert, ruleEngineSvc)
		case <-doneChan:
			break consumerLoop
		}
	}
}

func printStruct(alert *alerts.AlertInput, output *alerts.AlertOutput) {
	fmt.Println("------------ALERT--------------------")
	fmt.Println("ID: ", alert.ID)
	fmt.Println("Category: ", alert.Category)
	fmt.Println("Origin: ", alert.Origin)
	fmt.Println("Source: ", alert.Source)
	b, err := json.MarshalIndent(alert.Params, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Print("Params: ", string(b))
	if output != nil {
		fmt.Println("--------------OUTPUT-----------------")
		fmt.Println("Severity -> ", output.Severity)
		fmt.Println("Remedy -> ", output.Remedy)
	}
	fmt.Println("-------------------------------------")
}
