package alerts

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strconv"
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

func NewRuntimeMetrics() *RuntimeMetrics {
	prom_host := os.Getenv("PROMETHEUS_HOST")
	if prom_host == "" {
		prom_host = "prometheus-kube-prometheus-prometheus.default"
	}
	
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	// cpu_url := "http://prometheus-kube-prometheus-prometheus.default:9090/api/v1/query?query=100%20-%20(avg(rate(node_cpu_seconds_total{mode%3D%22idle%22}[5m]))%20*%20100)"
	prom_cpu_url := fmt.Sprintf("http://%s:9090/api/v1/query?query=100%%20-%%20(avg(rate(node_cpu_seconds_total{mode%%3D%%22idle%%22}[5m]))%%20*%%20100)", prom_host)
	// ram_url := "http://prometheus-kube-prometheus-prometheus.default:9090/api/v1/query?query=100%20-%20(node_memory_MemAvailable_bytes%20%2F%20node_memory_MemTotal_bytes%20*%20100)"
	prom_ram_url := fmt.Sprintf("http://%s:9090/api/v1/query?query=100%%20-%%20(node_memory_MemAvailable_bytes%%20%%2F%%20node_memory_MemTotal_bytes%%20*%%20100)", prom_host)
	CpuUsage, err := fetchMetrics(prom_cpu_url)
	if err != nil {
		fmt.Println("Error:", err)
	}
	RamUsage, err := fetchMetrics(prom_ram_url)
	if err != nil {
		fmt.Println("Error:", err)
	}
	return &RuntimeMetrics{
		NumGoroutine: uint64(runtime.NumGoroutine()),
		CpuUsage:     CpuUsage,
		RamUsage:     RamUsage,
	}
}