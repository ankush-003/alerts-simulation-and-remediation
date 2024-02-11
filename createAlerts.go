package main

import (
	"asmr/alerts"
	"log"
	"context"
	"os"
	"os/signal"
	"asmr/store"
	"time"
)

func main() {
	ticker := time.NewTicker(10 * time.Second)

	redis := store.NewRedisStore("localhost:6379")

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	ctx := context.Background()
	for {
		select {
		case <-ticker.C:
			alert := alerts.NewAlert("This is a test alert", "critical", "simulator-test")
			if err := redis.StoreAlert(ctx, alert); err != nil {
				log.Printf("Error saving alert: %s\n", err)
			}
			
		case <-signalChan:
			log.Printf("Interrupted\n")
			return
		}
	}
}