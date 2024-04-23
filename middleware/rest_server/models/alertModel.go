package models

import (
	"time"

	"github.com/google/uuid"
)

type Alerts struct {
	ID          uuid.UUID `json:"id"`
	NodeID      uuid.UUID `json:"node_id"`
	Description string    `json:"description"`
	Severity    string    `json:"severity"`
	Source      string    `json:"source"`
	CreatedAt   string    `json:"created_at"`
}

type AlertOutput struct {
	Severity string
	Remedy   string
}

type AlertConfig struct {
	ID          uuid.UUID `json:"id"`
	Description string    `json:"description"`
	Severity    string    `json:"severity"`
}

func NewAlert(alertConfig *AlertConfig, NodeID uuid.UUID, source string) *Alerts {
	return &Alerts{
		ID:          uuid.New(),
		NodeID:      NodeID,
		Description: alertConfig.Description,
		Severity:    alertConfig.Severity,
		Source:      source,
		CreatedAt:   time.Now().String(),
	}
}

func NewAlertConfig(description, severity string) *AlertConfig {
	return &AlertConfig{
		ID:          uuid.New(),
		Description: description,
		Severity:    severity,
	}
}

func NewAlertConfigWithID(description, severity string, id uuid.UUID) *AlertConfig {
	return &AlertConfig{
		ID:          id,
		Description: description,
		Severity:    severity,
	}
}