package main

import (
	"log"

	"github.com/dikaeinstein/job-runner/queue"
)

// newProductConsumer is the callback that handles the `new.product` event
func newProductConsumer(message queue.Message) error {
	log.Printf("Received message: %v, %v, %v\n", message.ID, message.Name, message.Payload)

	return nil // successfully processed message
}
