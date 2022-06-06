package main

import (
	"fmt"
	"github.com/nats-io/stan.go"
	"time"
)

func main() {
	sc, err := stan.Connect("test-cluster", "master")
	if err != nil {
		fmt.Println(err)
	}
	defer sc.Close()
	// Simple Synchronous Publisher
	err = sc.Publish("test", []byte("Hello World")) // does not return until an ack has been received from NATS Streaming
	if err != nil {
		fmt.Println(err)
	}
}
