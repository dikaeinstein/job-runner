package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/adjust/rmq"
	nc "github.com/dikaeinstein/job-runner/nats"
	"github.com/dikaeinstein/job-runner/server/handler"
	"github.com/nats-io/nats.go"
)

func main() {
	const PORT = 8912
	natsClient, err := nc.ConnectNATS(nats.DefaultURL, "job_runner NATS publisher")
	if err != nil {
		log.Fatal(err)
	}
	defer natsClient.Close()

	connection := rmq.OpenConnection("handler", "tcp", "localhost:6379", 1)

	jobs := handler.NewJobHandler(natsClient)
	http.Handle("/jobs", handler.NewValidationHandler(natsClient, jobs))
	http.Handle("/rmqstats", handler.NewRMQStatsHandler(connection))

	log.Printf("server listening on :%d", PORT)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil))
}
