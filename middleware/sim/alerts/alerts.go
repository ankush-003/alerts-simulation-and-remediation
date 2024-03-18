package alerts

import (
	"github.com/google/uuid"
	"runtime"
	"time"
)

type RuntimeMetrics struct {
	NumGoroutine           uint64 `json:"num_goroutine"`
	AllocatedMemBytes      uint64 `json:"allocated_mem_bytes"`
	TotalAllocatedMemBytes uint64 `json:"total_allocated_mem_bytes"`
	SysMemBytes            uint64 `json:"sys_mem_bytes"`
}

func NewRuntimeMetrics() *RuntimeMetrics {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return &RuntimeMetrics{
		NumGoroutine:           uint64(runtime.NumGoroutine()),
		AllocatedMemBytes:      memStats.Alloc,
		TotalAllocatedMemBytes: memStats.TotalAlloc,
		SysMemBytes:            memStats.Sys,
	}
}

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
		ID:          uuid.New(),
		NodeID:      NodeID,
		Description: alertConfig.Description,
		Severity:    alertConfig.Severity,
		Source:      source,
		CreatedAt:   time.Now().Format(time.DateTime),
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
