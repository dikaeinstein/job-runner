package nats

import (
	"log"
	"os"
	"time"

	"github.com/nats-io/nats.go"
)

// ConnectNATS connects the client to the nats-server with given url and name.
// The name is used as the NATS name option so pick wisely. good names can be
// `Example Subscriber` or `Example Publisher`
func ConnectNATS(url, name string) (*nats.Conn, error) {
	opts := []nats.Option{nats.Name(name)}
	opts = setupConnOptions(opts)

	log.Println("Connecting to NATS server")
	nc, err := nats.Connect(nats.DefaultURL, opts...)
	if err != nil {
		return nil, err
	}

	return nc, nil
}

func setupConnOptions(opts []nats.Option) []nats.Option {
	totalWait := 10 * time.Minute
	reconnectDelay := time.Second

	opts = append(opts, nats.ReconnectWait(reconnectDelay))
	opts = append(opts, nats.MaxReconnects(int(totalWait/reconnectDelay)))
	opts = append(opts, nats.DisconnectHandler(func(nc *nats.Conn) {
		log.Printf("Disconnected: will attempt reconnects for %.0fm", totalWait.Minutes())
	}))
	opts = append(opts, nats.ReconnectHandler(func(nc *nats.Conn) {
		log.Printf("Reconnected [%s]", nc.ConnectedUrl())
	}))
	opts = append(opts, nats.ClosedHandler(func(nc *nats.Conn) {
		if err := nc.LastError(); err != nil {
			log.Fatalf("Exiting %v", err)
		} else {
			log.Println("Exited NATS")
			os.Exit(0)
		}
	}))

	return opts
}
