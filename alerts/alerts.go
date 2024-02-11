package alerts

import (
	"github.com/google/uuid"
	"time"
)

type Alerts struct {
	ID          uuid.UUID `json:"id"`
	Description string    `json:"description"`
	Severity    string    `json:"severity"`
	Source      string    `json:"source"`
	CreatedAt   string    `json:"created_at"`
}

func NewAlert(description, severity, source string) *Alerts {
	return &Alerts{
		ID:          uuid.New(),
		Description: description,
		Severity:    severity,
		Source:      source,
		CreatedAt:   time.Now().String(),
	}
}
