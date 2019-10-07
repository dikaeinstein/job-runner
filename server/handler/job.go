package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/dikaeinstein/job-runner/server/entity"
	"github.com/nats-io/nats.go"
)

type jobHandler struct {
	natsClient *nats.Conn
}

// NewJobHandler creates a new handler with the given queue
func NewJobHandler(natsClient *nats.Conn) http.Handler {
	return &jobHandler{natsClient}
}

func (j *jobHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	jobRequest := r.Context().Value(job("jobRequest")).(entity.JobRequest)
	data, _ := json.Marshal(jobRequest.Payload)

	err := j.natsClient.Publish(jobRequest.Name, data)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	j.natsClient.Flush()
	if err := j.natsClient.LastError(); err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Published [%s] : '%s'\n", jobRequest.Name, string(data))
	}

	w.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.Encode(entity.JobResponse{Message: "job received"})
}
