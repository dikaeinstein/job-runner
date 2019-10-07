package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/adjust/rmq"
)

type rmqStatsHandler struct {
	connection rmq.Connection
}

// NewRMQStatsHandler creates a new stats handler with the given connection
func NewRMQStatsHandler(connection rmq.Connection) http.Handler {
	return &rmqStatsHandler{connection: connection}
}

func (rmq *rmqStatsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	layout := r.FormValue("layout")
	refresh := r.FormValue("refresh")

	queues := rmq.connection.GetOpenQueues()
	stats := rmq.connection.CollectStats(queues)
	log.Printf("queue stats\n%s", stats)
	fmt.Fprint(w, stats.GetHtml(layout, refresh))
}
