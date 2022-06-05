package main

import (
	"fmt"
	"github.com/nats-io/stan.go"
	"time"
)

func main() {
	sc, err := stan.Connect("test-cluster", "client12")
	if err != nil {
		fmt.Println(err)
	}
	// Simple Synchronous Publisher
	err = sc.Publish("foo", []byte("Hello World")) // does not return until an ack has been received from NATS Streaming

	if err != nil {
		fmt.Println(err)
	}

	timer, err := time.ParseDuration("30s")
	// Simple Async Subscriber
	sub, err := sc.Subscribe("foo", func(m *stan.Msg) {
		fmt.Printf("Received a message: %s\n", string(m.Data))
	}, stan.StartAtTimeDelta(timer))
	if err != nil {
		fmt.Println(err)
	}
	// Unsubscribe
	sub.Unsubscribe()

	// Close connection
	sc.Close()
}
