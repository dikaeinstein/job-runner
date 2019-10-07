package queue

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"time"

	"github.com/adjust/rmq"
)

// redisQueue implements the Queue interface for a redis based message queue
type redisQueue struct {
	queue    rmq.Queue
	name     string
	callback Callback
}

var seriaLNumberLimit = new(big.Int).Lsh(big.NewInt(1), 128)

// NewRedisQueue creates a new redisQueue with the given redisClient and queueName
func NewRedisQueue(conn rmq.Connection, queueName string) (Queue, error) {
	taskQueue := conn.OpenQueue(queueName)
	return &redisQueue{queue: taskQueue, name: queueName}, nil
}

// Add a new message with the given payload to the queue
func (r *redisQueue) Add(messageName string, payload []byte) error {
	message := Message{Name: messageName, Payload: string(payload)}
	return r.AddMessage(message)
}

// AddMessage to the queue, generating a unique ID for the message before dispatch
func (r *redisQueue) AddMessage(message Message) error {
	serialNumber, err := rand.Int(rand.Reader, seriaLNumberLimit)
	if err != nil {
		return err
	}
	message.ID = strconv.Itoa(time.Now().Nanosecond()) + serialNumber.String()

	payloadBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	log.Printf("Add event to Redis %s queue: %s\n", message.Name, string(payloadBytes))
	if !r.queue.PublishBytes(payloadBytes) {
		return fmt.Errorf("unable to add message to %s queue", message.Name)
	}

	return nil
}

// StartConsuming starts consuming message from the queue
func (r *redisQueue) StartConsuming(size int, pollInterval time.Duration,
	numOfConsumers int, callback Callback) {
	r.callback = callback
	r.queue.StartConsuming(size, pollInterval)
	for i := 0; i < numOfConsumers; i++ {
		r.queue.AddConsumer(fmt.Sprintf("job_runner_%s_consumer_%d", r.name, i), r)
	}
}

// StopConsuming stops message from the queue and returns a channel.
// If you want to wait until all consumers are idle you can wait on the returned channel
func (r *redisQueue) StopConsuming() <-chan struct{} {
	return r.queue.StopConsuming()
}

// Consume is the internal callback for the message queue
func (r *redisQueue) Consume(delivery rmq.Delivery) {
	log.Println("Got event from the queue:", delivery.Payload())

	message := Message{}

	err := json.Unmarshal([]byte(delivery.Payload()), &message)
	if err != nil {
		log.Println("Error consuming event, unable to deserialize event")
		delivery.Reject()
		return
	}

	if err := r.callback(message); err != nil {
		delivery.Reject()
		return
	}

	delivery.Ack()
}
