package main

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"
)

type Event struct {
	Criticality  int    `json:"criticality"`
	Timestamp    string `json:"timestamp"`
	EventMessage string `json:"eventMessage"`
}

func consumer(ctx context.Context, id int, topic string) error {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return err
	}
	defer nc.Close()

	// Subscribe to NATS subject using SubscribeSync
	sub, err := nc.SubscribeSync(topic)
	if err != nil {
		return err
	}

loop:
	for {
		select {
		case <-ctx.Done():
			log.Printf("Shut down consumer %d...", id)
			break loop
		default:
			msg, err := sub.NextMsg(2 * time.Second)
			if err != nil {
				if err == nats.ErrTimeout {
					// No message received within the timeout, continue the loop
					continue
				}
				log.Printf("ConsumerID: %d Error fetching message from NATS: %v", id, err)
				continue
			}

			// Process the message
			var event Event
			if err := json.Unmarshal(msg.Data, &event); err != nil {
				log.Printf("ConsumerID: %d Error unmarshalling event: %v", id, err)
				continue
			}

			// Log the received event
			log.Printf("ConsumerID: %d Received event: %+v", id, event)
		}
	}
	return nil
}

func producer(ctx context.Context, topic string) error {

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return err
	}
	defer nc.Close()
loop:
	for {
		select {
		case <-ctx.Done():
			log.Println("Shut down producer...")
			break loop
		default:
			event := Event{
				Criticality:  rand.Intn(10) + 1, // Criticality: 1-10
				Timestamp:    time.Now().Format(time.RFC3339),
				EventMessage: "Random security event",
			}

			data, err := json.Marshal(event)
			if err != nil {
				log.Printf("Error marshalling event: %v", err)
				continue
			}

			err = nc.Publish(topic, data)
			if err != nil {
				log.Printf("Error publishing to NATS: %v", err)
			} else {
				log.Printf("Published event: %s", data)
			}

			time.Sleep(2 * time.Second) // Generate event every 2 seconds
		}
	}
	return nil
}

func main() {
	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, os.Interrupt, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup
	wg.Add(6)
	go func() {
		<-osSignal
		cancel()
		wg.Done()
	}()

	go func() {
		err := producer(ctx, "time.eu.uk.london")
		if err != nil {
			log.Printf("Producer error: %v", err)
		}
		wg.Done()
	}()

	go func() {
		err := consumer(ctx, 1, "time.eu.uk.london")
		if err != nil {
			log.Printf("Consumer error: %v", err)
		}
		wg.Done()
	}()

	go func() {
		err := consumer(ctx, 2, "time.eu.uk.*")
		if err != nil {
			log.Printf("Consumer error: %v", err)
		}
		wg.Done()
	}()

	// This will not receive any message
	go func() {
		err := consumer(ctx, 3, "time.eu.*")
		if err != nil {
			log.Printf("Consumer error: %v", err)
		}
		wg.Done()
	}()

	go func() {
		err := consumer(ctx, 4, "time.eu.>")
		if err != nil {
			log.Printf("Consumer error: %v", err)
		}
		wg.Done()
	}()

	wg.Wait()
}
