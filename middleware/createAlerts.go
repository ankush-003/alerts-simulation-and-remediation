package main

import (
	"asmr/alerts"
	"asmr/store"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
)

func main() {
	redis := store.NewRedisStore("localhost:6379")

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	//predefined alert configs
	descriptions := []string{"High CPU usage", "Low disk space", "Network down", "Service unavailable"}
	severities := []string{"critical", "warning", "info"}

	ctx := context.Background()
	// push predefined alert configs to redis
	for _, desc := range descriptions {
		for _, sev := range severities {
			alertConfig := alerts.NewAlertConfig(desc, sev)
			err := redis.StoreAlertConfig(ctx, alertConfig)
			if err != nil {
				log.Printf("Error storing alert config: %s\n", err)
			}
		}
	}
	fmt.Println("Alert configs stored")
}