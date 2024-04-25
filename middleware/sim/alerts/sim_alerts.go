package alerts

import (
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"
    "os"
    "os/signal"
)

func SimulateHighCPULoad(rt *RuntimeMetrics) {
    signalChan := make(chan os.Signal, 1)
    signal.Notify(signalChan, os.Interrupt)

    fmt.Println("Simulating high CPU load")
    var wg sync.WaitGroup
    numGoroutines := rand.Intn(100) + 50 // Spawn between 50 and 150 goroutines

// simulatorLoop:
    for i := 0; i < numGoroutines; i++ {
        select {
        case <-signalChan:
            fmt.Println("Received signal to stop")
            return
        default:
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
    }

    // Call NewRuntimeMetrics before waiting for goroutines to complete
    metrics := NewRuntimeMetrics()
	rt.NumGoroutine = metrics.NumGoroutine
    rt.CpuUsage = metrics.CpuUsage
    rt.RamUsage = metrics.RamUsage

    // Wait for all goroutines to complete
    wg.Wait()
	fmt.Println("High CPU load simulation complete")
}

func SimulateNetworkLoad(rt *RuntimeMetrics) {
    fmt.Println("Simulating network load")

    signalChan := make(chan os.Signal, 1)
    signal.Notify(signalChan, os.Interrupt)

    client := &http.Client{}
    numRequests := rand.Intn(50) + 50 // Send between 50 and 100 requests

	var wg sync.WaitGroup

// simulatorLoop:    
    for i := 0; i < numRequests; i++ {
		select {
        case <-signalChan:
            fmt.Println("Received signal to stop")
            return
        default:
            wg.Add(1)
            go func() {
                defer wg.Done()
                req, _ := http.NewRequest("GET", "https://google.com", nil)
                resp, err := client.Do(req)
                if err != nil {
                    fmt.Println("Error sending request: ", err)
                    return
                }
                defer resp.Body.Close()
                fmt.Println("Request sent")
            }()
        }
    }

	// Call NewRuntimeMetrics before waiting for goroutines to complete
	metrics := NewRuntimeMetrics()
    rt.NumGoroutine = metrics.NumGoroutine
    rt.CpuUsage = metrics.CpuUsage
    rt.RamUsage = metrics.RamUsage

	wg.Wait()
	fmt.Println("Network load simulation complete")
}