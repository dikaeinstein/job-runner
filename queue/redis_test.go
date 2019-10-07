package queue

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/adjust/rmq"
)

func setupRedisQueue(queueName string) (Queue, rmq.TestConnection, error) {
	testConn := rmq.NewTestConnection()
	r, err := NewRedisQueue(testConn, queueName)
	return r, testConn, err
}

func TestProducer(t *testing.T) {
	queueName := "test.add"
	r, testConn, err := setupRedisQueue(queueName)
	if err != nil {
		t.Fatal(err)
	}

	taskPayload := "test add"
	err = r.Add(queueName, []byte(taskPayload))
	if err != nil {
		t.Fatal(err)
	}

	expected := Message{Name: queueName, Payload: taskPayload}
	got := Message{}
	json.Unmarshal([]byte(testConn.GetDelivery(queueName, 0)), &got)

	if expected.Payload != got.Payload {
		t.Fatalf("expected: %s, got: %s", expected, got)
	}
}

func TestConsumer(t *testing.T) {
	queueName := "test.startConsuming"
	r, _, err := setupRedisQueue(queueName)
	if err != nil {
		t.Fatal(err)
	}

	m := Message{Name: queueName, Payload: "test start consuming"}
	b, _ := json.Marshal(m)
	delivery := rmq.NewTestDeliveryString(string(b))

	r.StartConsuming(1, 2*time.Millisecond, 1, func(m Message) error {
		return nil
	})
	rq, _ := r.(*redisQueue)
	rq.Consume(delivery)

	if delivery.State != rmq.Acked {
		t.Fatalf("expected: %s, got: %s", rmq.Acked, delivery.State)
	}

	got := Message{}
	json.Unmarshal([]byte(delivery.Payload()), &got)
	if got.Payload != m.Payload {
		t.Fatalf("expected: %s, got: %s", m, got)
	}
}
