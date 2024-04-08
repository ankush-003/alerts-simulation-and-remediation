package rule_engine

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
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
	alert.CreatedAt = data["createdAt"].(time.Time)
	switch paramsType {
	case "Memory":
		var memory Memory
		if err := unmarshalParam(memory, paramsData); err != nil {
			return err
		}
		alert.Params = &memory
	case "CPU":
		var cpu CPU
		if err := unmarshalParam(cpu, paramsData); err != nil {
			return err
		}
		alert.Params = &cpu
	case "Disk":
		var disk Disk
		if err := unmarshalParam(disk, paramsData); err != nil {
			return err
		}
		alert.Params = &disk
	case "Network":
		var network Network
		if err := unmarshalParam(network, paramsData); err != nil {
			return err
		}
		alert.Params = &network
	case "Power":
		var power Power
		if err := unmarshalParam(power, paramsData); err != nil {
			return err
		}
		alert.Params = &power
	case "Applications":
		var apps Applications
		if err := unmarshalParam(apps, paramsData); err != nil {
			return err
		}
		alert.Params = &apps
	case "Security":
		var security Security
		if err := unmarshalParam(security, paramsData); err != nil {
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

func unmarshalParam(ptr interface{}, paramsData map[string]interface{}) error {
	structType := reflect.TypeOf(ptr).Elem()
	fmt.Println(structType)
	paramsField, found := structType.FieldByName("Params")
	if !found {
		return errors.New("struct has no field named Params")
	}
	fmt.Println(paramsField.Name)
	expectedType := paramsField.Type.Elem()

	paramsValue := reflect.New(expectedType).Interface()

	paramsBytes, err := json.Marshal(paramsData)
	if err != nil {
		return err
	}
	err = json.Unmarshal(paramsBytes, paramsValue)
	if err != nil {
		return err
	}

	// Set the value of the Params field in the original struct
	reflect.ValueOf(ptr).Elem().FieldByName("Params").Set(reflect.ValueOf(paramsValue))
	return nil
}
