package alerts

type Alerts struct {
	Description string `json:"description"`
	Severity    string `json:"severity"`
	Source      string `json:"source"`
}

func NewAlert(description, severity, source string) *Alerts {
	return &Alerts{
		Description: description,
		Severity:    severity,
		Source:      source,
	}
}
