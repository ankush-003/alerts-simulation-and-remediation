package rule_engine

import (
	"encoding/json"
	"errors"
	"time"
)

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

func (alert *AlertInput) Unmarshal(obj []byte) error {
	var data map[string]interface{}
	err := json.Unmarshal(obj, &data)
	if err != nil {
		return err
	}
	paramsData := data["params"].(map[string]interface{})
	paramsType := data["category"].(string)
	alert.Category = data["category"].(string)
	alert.ID = data["id"].(string)
	alert.Source = data["source"].(string)
	alert.Handled = data["handled"].(bool)
	alert.Origin = data["origin"].(string)
	alert.CreatedAt, _ = time.Parse(time.RFC3339, data["createdAt"].(string))
	switch paramsType {
	case "Memory":
		var memory Memory
		if err := memory.Unmarshal(paramsData); err != nil {
			return err
		}
		alert.Params = &memory
	case "CPU":
		var cpu CPU
		if err := cpu.Unmarshal(paramsData); err != nil {
			return err
		}
		alert.Params = &cpu
	case "Disk":
		var disk Disk
		if err := disk.Unmarshal(paramsData); err != nil {
			return err
		}
		alert.Params = &disk
	case "Network":
		var network Network
		if err := network.Unmarshal(paramsData); err != nil {
			return err
		}
		alert.Params = &network
	case "Power":
		var power Power
		if err := power.Unmarshal(paramsData); err != nil {
			return err
		}
		alert.Params = &power
	case "Applications":
		var apps Applications
		if err := apps.Unmarshal(paramsData); err != nil {
			return err
		}
		alert.Params = &apps
	case "Security":
		var security Security
		if err := security.Unmarshal(paramsData); err != nil {
			return err
		}
		alert.Params = &security
	default:
		return errors.New("WRONG PARAM INPUT TYPE")
	}
	return nil
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

// CPU:
// Utilization (%)
// Temperature (°C)
// Core speeds (GHz)
// Context switches per second

type CPU struct {
	Utilization uint `json:"utilization"`
	Temperature uint `json:"temperature"`
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

func (cpu *CPU) DataKey() string {
	return "CpuInput"
}

func (mem *Memory) DataKey() string {
	return "MemInput"
}

func (disk *Disk) DataKey() string {
	return "DiskInput"
}

func (net *Network) DataKey() string {
	return "NetInput"
}

func (pow *Power) DataKey() string {
	return "PowerInput"
}

func (apps *Applications) DataKey() string {
	return "AppInput"
}

func (sec *Security) DataKey() string {
	return "SecInput"
}

func (mem *Memory) Unmarshal(paramsData map[string]interface{}) error {
	paramsBytes, err := json.Marshal(paramsData) // Convert map to JSON bytes
	if err != nil {
		return err
	}
	return json.Unmarshal(paramsBytes, mem)

}

func (cpu *CPU) Unmarshal(paramsData map[string]interface{}) error {
	paramsBytes, err := json.Marshal(paramsData) // Convert map to JSON bytes
	if err != nil {
		return err
	}
	return json.Unmarshal(paramsBytes, cpu)
}

func (disk *Disk) Unmarshal(paramsData map[string]interface{}) error {
	paramsBytes, err := json.Marshal(paramsData)
	if err != nil {
		return err
	}
	return json.Unmarshal(paramsBytes, disk)
}

// Network Unmarshal
func (network *Network) Unmarshal(paramsData map[string]interface{}) error {
	paramsBytes, err := json.Marshal(paramsData)
	if err != nil {
		return err
	}
	return json.Unmarshal(paramsBytes, network)
}

// Power Unmarshal
func (power *Power) Unmarshal(paramsData map[string]interface{}) error {
	paramsBytes, err := json.Marshal(paramsData)
	if err != nil {
		return err
	}
	return json.Unmarshal(paramsBytes, power)
}

// Applications Unmarshal
func (applications *Applications) Unmarshal(paramsData map[string]interface{}) error {
	paramsBytes, err := json.Marshal(paramsData)
	if err != nil {
		return err
	}
	return json.Unmarshal(paramsBytes, applications)
}

// Security Unmarshal
func (security *Security) Unmarshal(paramsData map[string]interface{}) error {
	paramsBytes, err := json.Marshal(paramsData)
	if err != nil {
		return err
	}
	return json.Unmarshal(paramsBytes, security)
}