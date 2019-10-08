package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/adjust/rmq"
	nc "github.com/dikaeinstein/job-runner/nats"
	"github.com/dikaeinstein/job-runner/server/handler"
)

func main() {
	PORT := os.Getenv("PORT")
	natsName := "job_runner NATS publisher"
	natsClient, err := nc.ConnectNATS(os.Getenv("NATS"), natsName)
	if err != nil {
		log.Fatal(err)
	}
	defer natsClient.Close()

	connection := rmq.OpenConnection("handler", "tcp", os.Getenv("REDIS"), 1)

	jobs := handler.NewJobHandler(natsClient)
	http.Handle("/jobs", handler.NewValidationHandler(natsClient, jobs))
	http.Handle("/rmqstats", handler.NewRMQStatsHandler(connection))

	log.Printf("server listening on :%s", PORT)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", PORT), nil))
}
