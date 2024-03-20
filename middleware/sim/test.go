package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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

func main() {
	url := "http://localhost:9090/api/v1/query?query=100 - (avg(rate(node_cpu_seconds_total{mode=\"idle\"}[5m])) * 100)"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var promResp PromResponse
	err = json.Unmarshal(body, &promResp)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if promResp.Status != "success" {
		fmt.Println("Error: Prometheus API returned non-success status")
		return
	}

	if len(promResp.Data.Result) == 0 {
		fmt.Println("No results found")
		return
	}

	// Assuming there is only one result
	value := promResp.Data.Result[0].Value[1].(float64)
	fmt.Printf("CPU Usage: %.2f%%\n", value)
}
