package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	os.Setenv("KAFKA_BROKER", "localhost:9092")
	os.Setenv("REDIS_ADDR", "localhost:6379")

}