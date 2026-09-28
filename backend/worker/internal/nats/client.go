package nats

import (
	"context"
	"encoding/json"
	"log"

	"github.com/NirajDonga/pingpong/backend/worker/internal/worker"
	natsgo "github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

const (
	CheckJobsSubject    = "check.jobs"
	CheckResultsSubject = "check.results"
	WorkersQueue        = "workers"
)

type Client struct {
	conn *natsgo.Conn
	js   jetstream.JetStream
}

func NewClient(url string) (*Client, error) {
	conn, err := natsgo.Connect(url)
	if err != nil {
		return nil, err
	}

	js, err := jetstream.New(conn)
	if err != nil {
		conn.Close()
		return nil, err
	}

	_, err = js.CreateStream(context.Background(), jetstream.StreamConfig{
		Name:     "RESULTS",
		Subjects: []string{CheckResultsSubject},
	})
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &Client{conn: conn, js: js}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) SubscribeCheckJobs(ctx context.Context, handler func(context.Context, worker.CheckJob) error) (jetstream.ConsumeContext, error) {
	cons, err := c.js.CreateOrUpdateConsumer(ctx, "JOBS", jetstream.ConsumerConfig{
		Durable:       WorkersQueue,
		DeliverPolicy: jetstream.DeliverAllPolicy,
		AckPolicy:     jetstream.AckExplicitPolicy,
	})
	if err != nil {
		return nil, err
	}

	return cons.Consume(func(msg jetstream.Msg) {
		var job worker.CheckJob
		if err := json.Unmarshal(msg.Data(), &job); err != nil {
			log.Printf("failed to decode check job: %v", err)
			msg.Term()
			return
		}

		if err := handler(ctx, job); err != nil {
			log.Printf("failed to process job: %v", err)
			msg.Nak()
			return
		}

		msg.Ack()
	})
}

func (c *Client) PublishCheckResult(ctx context.Context, result worker.CheckResult) error {
	data, err := json.Marshal(result)
	if err != nil {
		return err
	}

	_, err = c.js.Publish(ctx, CheckResultsSubject, data)
	return err
}
