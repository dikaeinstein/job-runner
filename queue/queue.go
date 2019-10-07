package queue

import "time"

// Message represent messages stored on the queue
type Message struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Payload string `json:"payload"`
}

// Callback represents the callback registered
type Callback func(Message) error

// Queue defines the interface for a message queue
type Queue interface {
	Add(messageName string, payload []byte) error
	AddMessage(message Message) error
	StartConsuming(size int, pollInterval time.Duration,
		numOfConsumers int, callback Callback)
	StopConsuming() <-chan struct{}
}
