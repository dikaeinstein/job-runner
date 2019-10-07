package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/adjust/rmq"
	"github.com/dikaeinstein/job-runner/queue"
	"github.com/dikaeinstein/job-runner/worker"
	"github.com/nats-io/nats.go"
)

// newProductWorker represents a new product worker
type newProductWorker struct {
	natsClient     *nats.Conn
	conn           rmq.Connection
	numOfConsumers int
	name           string
}

// NewProductWorker creates a newProduct worker instance and starts the
// given `num` of consumers
func NewProductWorker(natsClient *nats.Conn, conn rmq.Connection, numOfConsumers int) worker.Worker {
	return &newProductWorker{natsClient, conn, numOfConsumers, Subject}
}

// Start subscribes to the subject and starts the product worker
func (w *newProductWorker) Start(ctx context.Context, interrupt chan os.Signal) {
	log.Printf("Starting [%s] worker...", w.name)
	// Setup new.product worker Redis queue
	q := setupNewProductWorkerQueue(w.conn, Subject, w.numOfConsumers)

	// Subscribe to NATS new.product subject
	sub, err := w.natsClient.QueueSubscribe(Subject,
		"worker", makeNatsMsgHandler(q))
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Started [%s] worker", w.name)
	log.Printf("%s worker listening on [%s]", Subject, Subject)
	w.natsClient.Flush()

	<-interrupt
	fmt.Println()
	fmt.Println("interrupt: quit via SIGINT (Ctrl+C)")

	// Call Drain on the subscription. It unsubscribes but
	// wait for all pending messages to be processed.
	if err := sub.Drain(); err != nil {
		log.Fatal(err)
	}

	// Stop the worker
	w.Stop(ctx, q)
}

// Stop the worker
func (w *newProductWorker) Stop(ctx context.Context, q queue.Queue) {
	// Create a deadline to wait for
	d := time.Now().Add(5 * time.Second)
	ctx, cancel := context.WithDeadline(ctx, d)
	defer cancel()

	done := make(chan struct{})
	go func() {
		log.Printf("Stopping [%s] worker", Subject)
		log.Println("Cleaning up resources...")
	Loop:
		for {
			select {
			// Stop consuming from queue
			case <-q.StopConsuming():
				log.Printf("all messages of [%s] Redis queue consumed", Subject)
				break Loop
			case <-ctx.Done():
				log.Println(ctx.Err())
				break Loop
			default:
			}
		}
		close(done)
	}()

	<-done
	log.Printf("Stopped [%s] worker", Subject)
}

func makeNatsMsgHandler(q queue.Queue) nats.MsgHandler {
	return func(m *nats.Msg) {
		log.Printf("Received message: %v, %#v", m.Subject, string(m.Data))

		if err := q.Add(m.Subject, m.Data); err != nil {
			log.Printf("Failed to add message to %s Redis queue", m.Subject)
			return
		}
	}
}

func setupNewProductWorkerQueue(conn rmq.Connection, queueName string,
	numOfConsumers int) queue.Queue {
	log.Printf("Starting Redis [%s] queue...", queueName)

	q, err := queue.NewRedisQueue(conn, queueName)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Started Redis [%s] queue", queueName)

	log.Printf("Starting [%d] consumers...", numOfConsumers)
	q.StartConsuming(10, 100*time.Millisecond,
		numOfConsumers, newProductConsumer)
	log.Printf("Started [%d] consumers", numOfConsumers)

	return q
}
