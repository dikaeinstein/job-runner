package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/adjust/rmq"
	nc "github.com/dikaeinstein/job-runner/nats"
	"github.com/go-redis/redis"
)

func main() {
	numOfConsumers := flag.Int("consumers", 1, "number of worker consumers")
	flag.Parse()
	natsName := fmt.Sprintf("job_runner %s NATS subscriber", Subject)
	natsClient, err := nc.ConnectNATS(os.Getenv("NATS"), natsName)
	if err != nil {
		log.Fatal(err)
	}
	defer natsClient.Close()

	redisClient := setupRedisClient()
	defer redisClient.Close()
	conn := rmq.OpenConnectionWithRedisClient("job_runner", redisClient)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	var wg sync.WaitGroup
	wg.Add(1)

	w := NewProductWorker(natsClient, conn, *numOfConsumers)
	go func() {
		defer wg.Done()
		w.Start(ctx, interrupt)
	}()

	wg.Wait()
	fmt.Println("exiting main")
	os.Exit(0)
}

func setupRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Network:      "tcp",
		Addr:         os.Getenv("REDIS"),
		DB:           1,
		MaxRetries:   5,
		DialTimeout:  time.Second * 15,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	})
}
