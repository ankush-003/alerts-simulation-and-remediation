package main

import (
	"fmt"
	rule_engine "rule_engine/engine"
	"sync"
	"time"
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
		Params:    &rule_engine.Memory{Usage: 76, PageFaults: 30, SwapUsge: 2},
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
