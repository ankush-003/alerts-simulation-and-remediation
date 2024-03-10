package rule_engine

import "time"

type AlertContext struct {
	AlertInput  *InputData
	AlertOutput *OutputData
}

func (ac *AlertContext) RuleName() string {
	return "Alerts"
}

func (ac *AlertContext) RuleInput() RuleInput {
	return ac.AlertInput
}

func (ac *AlertContext) RuleOutput() RuleOutput {
	return ac.AlertOutput
}

// User data attributes
type InputData struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Description string    `json:"desc"`
	Severity    string    `json:"severity"`
	Origin      string    `json:"origin"`
	CreatedAt   time.Time `json:"createdAt"`
}

func (id *InputData) DataKey() string {
	return "InputData"
}

// Offer output object
type OutputData struct {
	Remedy string `josn:"remedy"`
}

func (od *OutputData) DataKey() string {
	return "OutputData"
}

func NewAlertContext() *AlertContext {
	return &AlertContext{
		AlertInput:  &InputData{},
		AlertOutput: &OutputData{},
	}
}
