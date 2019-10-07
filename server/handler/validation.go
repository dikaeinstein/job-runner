package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/dikaeinstein/job-runner/server/entity"
	"github.com/nats-io/nats.go"
	"gopkg.in/go-playground/validator.v9"
)

type validation struct {
	natsClient *nats.Conn
	next       http.Handler
}

type job string

// NewValidationHandler creates a new handler with the given nats.io connection
// and next handler
func NewValidationHandler(nc *nats.Conn, next http.Handler) http.Handler {
	return &validation{nc, next}
}

func (v *validation) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var jobRequest entity.JobRequest

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&jobRequest); err != nil {
		http.Error(w, "error decoding request", http.StatusBadRequest)
		return
	}

	if err := validateJobRequest(jobRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	c := context.WithValue(context.Background(), job("jobRequest"), jobRequest)
	r = r.WithContext(c)

	v.next.ServeHTTP(w, r)
}

var supportedJobTypes = []string{"new.product"}

var validate = validator.New()

// validateJobRequest checks if job type is supported
func validateJobRequest(j entity.JobRequest) error {
	if err := validate.Struct(j); err != nil {
		return err
	}

	for _, t := range supportedJobTypes {
		if t == j.Name {
			return nil
		}
	}

	return errors.New("unsupported job type")
}
