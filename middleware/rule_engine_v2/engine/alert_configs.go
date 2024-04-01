package rule_engine

import "time"

// Implements RuleInput interface
// Holds all the input to rule engine parameters

type AlertInput struct {
	ID        string     `json:"id"`
	Category  string     `json:"category"`
	Source    string     `json:"source"`
	Origin    string     `json:"origin"`
	Params    ParamInput `json:"params"`
	CreatedAt time.Time  `json:"createdAt"`
	Handled   bool       `json:"handled"`
}

func (alert *AlertInput) DataKey() string {
	return "AlertInput"
}

type AlertOutput struct {
	Severity string
	Remedy   string
}

func (alert *AlertOutput) DataKey() string {
	return "AlertOutput"
}

// Memory:
// Usage (%)
// Swap Usage (%)
// Paging activity (page faults/second)

type Memory struct {
	Usage      uint `json:"usage"`
	PageFaults uint `json:"pageFaults"`
	SwapUsage  uint `json:"swapUsage"`
}

func (mem *Memory) DataKey() string {
	return "MemInput"
}

// CPU:
// Utilization (%)
// Temperature (°C)
// Core speeds (GHz)
// Context switches per second

type CPU struct {
	Utilization uint `json:"utilization"`
	Temperature uint `json:"temperature"`
}

func (cpu *CPU) DataKey() string {
	return "CpuInput"
}

// Disk:
// Disk space usage (%)
// Read/write IOPS (Input/Output Operations Per Second)
// Throughput (MB/s)

type Disk struct {
	Usage       uint `json:"usage"`
	IOPs        uint `json:"iops"`
	ThroughtPut uint `json:"throughPut"`
}

// Network:
// interface traffic (bytes in/out per second)
// Packet loss (%)
// Latency (ms)

type Network struct {
	Traffic    uint `json:"traffic"`
	PacketLoss uint `json:"packetLoss"`
	Latency    uint `json:"latency"`
}

// Power:
// Battery level (%) (laptops)
// Power consumption (Watts)
// Remaining runtime on battery (minutes) (laptops)

type Power struct {
	BatteryLevel uint `json:"batteryLevel"`
	Consumption  uint `json:"consumption"`
	Efficiency   uint `json:"efficiency"`
}

// Applications:
// No. of Processes running
// Max CPU usage by all Processes
// Max Memory usage by all Processes

type Applications struct {
	Processes   uint `json:"processes"`
	MaxCPUusage uint `json:"maxCPUusage"`
	MaxMemUsage uint `json:"maxMemUsage"`
}

type Security struct {
	LoginAttempts  uint `json:"loginAttempts"`
	FailedLogins   uint `json:"failedLogins"`
	SuspectedFiles uint `json:"suspectedFiles"`
	IDSEvents      uint `json:"idsEvents"`
}
