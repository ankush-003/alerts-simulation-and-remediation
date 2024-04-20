package alerts

import (
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

func SimulateHighCPULoad(metricsChan chan<- *RuntimeMetrics) {
    var wg sync.WaitGroup
    numGoroutines := rand.Intn(100) + 50 // Spawn between 50 and 150 goroutines

    for i := 0; i < numGoroutines; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            startTime := time.Now()
            for time.Since(startTime) < 10*time.Second {
                // Simulate high CPU load by doing a lot of math
				for j := 0; j < 1000; j++ {
					_ = j * j
				}
            }
        }()
    }

    // Call NewRuntimeMetrics before waiting for goroutines to complete
    metrics := NewRuntimeMetrics()
	metricsChan <- metrics

    wg.Wait()
	fmt.Println("High CPU load simulation complete")
}

func SimulateNetworkLoad(metricsChan chan<- *RuntimeMetrics) {
    client := &http.Client{}
    numRequests := rand.Intn(100) + 50 // Send between 50 and 150 requests

	var wg sync.WaitGroup
    for i := 0; i < numRequests; i++ {
		wg.Add(1)
        go func() {
			defer wg.Done()
            req, _ := http.NewRequest("GET", "https://google.com", nil)
            resp, _ := client.Do(req)
            resp.Body.Close()
        }()
    }

	// Call NewRuntimeMetrics before waiting for goroutines to complete
	metrics := NewRuntimeMetrics()
	metricsChan <- metrics

	wg.Wait()
	fmt.Println("Network load simulation complete")
}