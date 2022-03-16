package main

import (
	"os"

	_ "github.com/gin-gonic/gin"
)

func main() {
	servers := "localhost:9092"
	mode := "consumer"
	if v, ok := os.LookupEnv("mode"); ok {
		mode = v
	}
	if mode == "consumer" {
		RunConsumer(servers)
	} else if mode == "producer" {
		RunProducer(servers)
	}
}
