package alerts

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type PromResponse struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric map[string]string `json:"metric"`
			Value  []interface{}     `json:"value"`
		} `json:"result"`
	} `json:"data"`
}

func fetchMetrics(url string) (float64, error) {
	resp, err := http.Get(url)

	if err != nil {
		return -1, fmt.Errorf("error fetching metrics: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return -1, fmt.Errorf("prometheus API returned non-200 status")
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return -1, fmt.Errorf("error reading response body: %v", err)
	}

	var promResp PromResponse
	err = json.Unmarshal(body, &promResp)
	if err != nil {
		return -1, fmt.Errorf("error unmarshalling response body: %v", err)
	}

	if promResp.Status != "success" {
		return -1, fmt.Errorf("prometheus API returned non-success status")
	}

	if len(promResp.Data.Result) == 0 {
		return 0, fmt.Errorf("no data returned")
	}

	valueTr, ok := promResp.Data.Result[0].Value[1].(string)
	if !ok {
		return -1, fmt.Errorf("unable to parse value")
	}

	usage, _ := strconv.ParseFloat(valueTr, 64)

	return usage, nil
}

type RuntimeMetrics struct {
	NumGoroutine uint64  `json:"num_goroutine"`
	CpuUsage     float64 `json:"cpu_usage"`
	RamUsage     float64 `json:"ram_usage"`
}

func (*RuntimeMetrics) DataKey() string {
	return "RuntimeMetrics"
}

func (rt *RuntimeMetrics) Unmarshal(paramsData map[string]interface{}) error {
	paramsBytes, err := json.Marshal(paramsData)
	if err != nil {
		return err
	}
	return json.Unmarshal(paramsBytes, rt)
}

func NewRuntimeMetrics() *RuntimeMetrics {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	cpu_url := "http://prometheus-sim:9090/api/v1/query?query=100%20-%20(avg(rate(node_cpu_seconds_total{mode%3D%22idle%22}[5m]))%20*%20100)"
	ram_url := "http://prometheus-sim:9090/api/v1/query?query=100%20-%20(node_memory_MemAvailable_bytes%20%2F%20node_memory_MemTotal_bytes%20*%20100)"
	CpuUsage, err := fetchMetrics(cpu_url)
	if err != nil {
		fmt.Println("Error:", err)
	}
	RamUsage, err := fetchMetrics(ram_url)
	if err != nil {
		fmt.Println("Error:", err)
	}
	return &RuntimeMetrics{
		NumGoroutine: uint64(runtime.NumGoroutine()),
		CpuUsage:     CpuUsage,
		RamUsage:     RamUsage,
	}
}

func NewAlertInput(alertConfig *AlertConfig, NodeID string, source string) *AlertInput {
	return &AlertInput{
		// ID:        uuid.New().String(),
		ID:        NodeID,
		// Category:  alertConfig.Description,
		Category: "RuntimeMetrics",
		Source:    source,
		Origin:    NodeID,
		Params:    NewRuntimeMetrics(),
		CreatedAt: time.Now().Format(time.DateTime),
		Handled:   false,
	}
}

// older version alert
type Alerts struct {
	ID             uuid.UUID       `json:"id"`
	NodeID         uuid.UUID       `json:"node_id"`
	Description    string          `json:"description"`
	Severity       string          `json:"severity"`
	Source         string          `json:"source"`
	CreatedAt      string          `json:"created_at"`
	RuntimeMetrics *RuntimeMetrics `json:"runtime_metrics"`
}

type AlertConfig struct {
	ID          uuid.UUID `json:"id"`
	Description string    `json:"description"`
	Severity    string    `json:"severity"`
}

func NewAlert(alertConfig *AlertConfig, NodeID uuid.UUID, source string) *Alerts {
	return &Alerts{
		ID:             uuid.New(),
		NodeID:         NodeID,
		Description:    alertConfig.Description,
		Severity:       alertConfig.Severity,
		Source:         source,
		CreatedAt:      time.Now().Format(time.DateTime),
		RuntimeMetrics: NewRuntimeMetrics(),
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
