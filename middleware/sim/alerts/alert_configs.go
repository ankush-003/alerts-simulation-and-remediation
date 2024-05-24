package alerts

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

// Implements RuleInput interface
// Holds all the input to rule engine parameters

type ParamInput interface {
	DataKey() string
	GenerateMetrics()
	Unmarshal(paramsData map[string]interface{}) error
}

type AlertInput struct {
	ID        string     `json:"id"`
	Category  string     `json:"category"`
	Source    string     `json:"source"`
	Origin    string     `json:"origin"`
	Params    ParamInput `json:"params"`
	CreatedAt string     `json:"createdAt"`
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
	alert.CreatedAt, _ = data["createdAt"].(string)

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
	case "RuntimeMetrics":
		var rt RuntimeMetrics
		if err := rt.Unmarshal(paramsData); err != nil {
			return err
		}
		alert.Params = &rt
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
	On         bool `json:"On"`
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

func (mem *Memory) DataKey() string {
	return "MemInput"
}

func (cpu *CPU) DataKey() string {
	return "CpuInput"
}

func (disk *Disk) DataKey() string {
	return "DiskInput"
}

func (net *Network) DataKey() string {
	return "NetworkInput"
}

func (pow *Power) DataKey() string {
	return "PowerInput"
}

func (apps *Applications) DataKey() string {
	return "ApplicationsInput"
}

func (sec *Security) DataKey() string {
	return "SecurityInput"
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

func (m *Memory) GenerateMetrics() {
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	m.Usage = uint(rand.Intn(90)) + 10 // Assuming usage percentage
	m.PageFaults = uint(rand.Intn(100))
	m.SwapUsage = uint(rand.Intn(50))
}

func (c *CPU) GenerateMetrics() {
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	c.Utilization = uint(rand.Intn(101))     // Utilization can range from 0 to 100%
	c.Temperature = uint(rand.Intn(50)) + 30 // Temperature in Celsius
}

func (d *Disk) GenerateMetrics() {
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	d.Usage = uint(rand.Intn(90)) + 10   // Disk usage percentage
	d.IOPs = uint(rand.Intn(10000))      // Assuming IOPs range
	d.ThroughtPut = uint(rand.Intn(100)) // Throughput in MB/s
}

func (n *Network) GenerateMetrics() {
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	n.On = rand.Intn(2) == 1
	n.Traffic = uint(rand.Intn(500))            // Traffic in Mbytes per second
	n.PacketLoss = uint(rand.Intn(4000)) + 1000 // Packet loss percentage
	n.Latency = uint(rand.Intn(100))            // Latency in milliseconds
}

func (p *Power) GenerateMetrics() {
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	p.BatteryLevel = uint(rand.Intn(101)) // Battery level percentage
	p.Consumption = uint(rand.Intn(100))  // Power consumption in Watts
	p.Efficiency = uint(rand.Intn(100))   // Remaining runtime in minutes
}

func (a *Applications) GenerateMetrics() {
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	a.Processes = uint(rand.Intn(10000)) // Number of processes running (0-100)
	a.MaxCPUusage = uint(rand.Intn(101)) // Max CPU usage by all processes (0-100)
	a.MaxMemUsage = uint(rand.Intn(101)) // Max memory usage by all processes (0-100)
}

func (s *Security) GenerateMetrics() {
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	s.LoginAttempts = uint(rand.Intn(101))  // Number of login attempts (0-100)
	s.FailedLogins = uint(rand.Intn(30))    // Number of failed logins (0-100)
	s.SuspectedFiles = uint(rand.Intn(101)) // Number of suspected files (0-100)
	s.IDSEvents = uint(rand.Intn(101))      // Number of IDS events (0-100)
	println(s)
}

type RuntimeMetrics struct {
	NumGoroutine uint64  `json:"num_goroutine"`
	CpuUsage     float64 `json:"cpu_usage"`
	RamUsage     float64 `json:"ram_usage"`
}

func (*RuntimeMetrics) DataKey() string {
	return "RuntimeMetricsInput"
}

func (rt *RuntimeMetrics) Unmarshal(paramsData map[string]interface{}) error {
	paramsBytes, err := json.Marshal(paramsData)
	if err != nil {
		return err
	}
	return json.Unmarshal(paramsBytes, rt)
}

func (rt *RuntimeMetrics) GenerateMetrics() {
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomSimulator := rand.Intn(2)

	if randomSimulator == 0 {
		SimulateHighCPULoad(rt)
	} else {
		SimulateNetworkLoad(rt)
	}

	fmt.Println("Simulated metrics: ", rt)
}
