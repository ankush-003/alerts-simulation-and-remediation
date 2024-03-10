package main

import (
	"fmt"
	rule_engine "rule_engine_demo/engine"
	"time"

	"github.com/hyperjumptech/grule-rule-engine/logger"
)

// can be part of user serice and a separate directory
type Alert struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Description string    `json:"desc"`
	Severity    string    `json:"severity"`
	Origin      string    `json:"origin"`
	CreatedAt   time.Time `json:"createdAt"`
	Remedy      string    `json:"remedy"`
}

// can be moved to offer directory
type OfferService interface {
	getRemedyForAlert(user Alert) string
}

type OfferServiceClient struct {
	ruleEngineSvc *rule_engine.RuleEngineSvc
}

func NewOfferService(ruleEngineSvc *rule_engine.RuleEngineSvc) OfferService {
	return &OfferServiceClient{
		ruleEngineSvc: ruleEngineSvc,
	}
}

func (svc OfferServiceClient) getRemedyForAlert(user Alert) string {
	ruleContext := rule_engine.NewAlertContext()
	ruleContext.AlertInput = &rule_engine.InputData{
		ID:          user.ID,
		Name:        user.Name,
		Type:        user.Type,
		Description: user.Description,
		Severity:    user.Severity,
		Origin:      user.Origin,
		CreatedAt:   user.CreatedAt,
	}

	err := svc.ruleEngineSvc.Execute(ruleContext)
	if err != nil {
		logger.Log.Error("get user offer rule engine failed", err)
	}

	return ruleContext.AlertOutput.Remedy
}

func main() {
	ruleEngineSvc := rule_engine.NewRuleEngineSvc()
	offerSvc := NewOfferService(ruleEngineSvc)

	userA := Alert{
		ID:          "1",
		Name:        "Memory Alert",
		Type:        "Memory",
		Description: "Ram usage is above 80%",
		Severity:    "Critical",
		Origin:      "Computer FS",
		CreatedAt:   time.Now(),
	}

	fmt.Println("Remedy for Alert A: ", offerSvc.getRemedyForAlert(userA))

	userB := Alert{
		ID:          "2",
		Name:        "Memory Alert",
		Type:        "Memory",
		Description: "Ram usage is above 50%",
		Severity:    "Severe",
		Origin:      "Computer FS",
		CreatedAt:   time.Now(),
	}

	fmt.Println("Remedy for Alert B: ", offerSvc.getRemedyForAlert(userB))
}
