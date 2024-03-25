package rule_engine

import "time"

// Implements RuleInput interface
// Holds all the input to rule engine parameters

type AlertInput struct {
	ID        string
	Category  string
	Source    string
	Origin    string
	Params    ParamInput
	CreatedAt time.Time
	Handled   bool
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
	Usage      uint
	PageFaults uint
	SwapUsge   uint
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
	Utilization uint
	Temperature uint
}

func (cpu *CPU) DataKey() string {
	return "CpuInput"
}

// Disk:
// Disk space usage (%)
// Read/write IOPS (Input/Output Operations Per Second)
// Throughput (MB/s)

type Disk struct {
	Usage       uint
	IOPs        uint
	ThroughtPut uint
}

// Network:
// interface traffic (bytes in/out per second)
// Packet loss (%)
// Latency (ms)

type Network struct {
	Traffic    uint
	PacketLoss uint
	Latency    uint
}

// Power:
// Battery level (%) (laptops)
// Power consumption (Watts)
// Remaining runtime on battery (minutes) (laptops)

type Power struct {
	BatteryLevel uint
	Consumption  uint
	Efficiency   uint
}

// Applications:
// No. of Processes running
// Max CPU usage by all Processes
// Max Memory usage by all Processes

type Applications struct {
	Processes   uint
	MaxCPUusage uint
	MaxMemUsage uint
}

type Security struct {
	LoginAttempts  uint
	FailedLogins   uint
	SuspectedFiles uint
	IDSEvents      uint
}
