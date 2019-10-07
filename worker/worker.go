package worker

import (
	"context"
	"os"

	"github.com/dikaeinstein/job-runner/queue"
)

// Worker defines the interface for a worker
type Worker interface {
	Start(ctx context.Context, interrupt chan os.Signal)
	Stop(ctx context.Context, q queue.Queue)
}
