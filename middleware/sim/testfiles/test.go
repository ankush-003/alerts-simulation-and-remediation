package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

	valueTr, ok := promResp.Data.Result[0].Value[1].(string)
	if !ok {
		return -1, fmt.Errorf("unable to parse value")
	}

	usage, _ := strconv.ParseFloat(valueTr, 64)

	return usage, nil
}

func main() {
	url := "http://localhost:9090/api/v1/query?query=100%20-%20(avg(rate(node_cpu_seconds_total{mode%3D%22idle%22}[5m]))%20*%20100)"
	url2 := "http://localhost:9090/api/v1/query?query=100%20-%20(node_memory_MemAvailable_bytes%20%2F%20node_memory_MemTotal_bytes%20*%20100)"

	usage, err := fetchMetrics(url)
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println("CPU Usage:", usage)

	usage2, err := fetchMetrics(url2)
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println("RAM Usage:", usage2)
}
