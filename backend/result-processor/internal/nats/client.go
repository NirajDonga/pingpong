package nats

import (
	"context"
	"encoding/json"
	"log"

	"github.com/NirajDonga/pingpong/backend/result-processor/internal/processor"
	natsgo "github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

const (
	CheckResultsSubject = "check.results"
	ProcessorQueue      = "processor"
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

	return &Client{conn: conn, js: js}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) SubscribeCheckResults(ctx context.Context, handler func(context.Context, processor.CheckResult) error) (jetstream.ConsumeContext, error) {
	cons, err := c.js.CreateOrUpdateConsumer(ctx, "RESULTS", jetstream.ConsumerConfig{
		Durable:       ProcessorQueue,
		DeliverPolicy: jetstream.DeliverAllPolicy,
		AckPolicy:     jetstream.AckExplicitPolicy,
	})
	if err != nil {
		return nil, err
	}

	return cons.Consume(func(msg jetstream.Msg) {
		var checkResult processor.CheckResult
		if err := json.Unmarshal(msg.Data(), &checkResult); err != nil {
			log.Printf("failed to decode check result: %v", err)
			msg.Term()
			return
		}

		if err := handler(ctx, checkResult); err != nil {
			log.Printf("failed to process check result: %v", err)
			msg.Nak()
			return
		}

		msg.Ack()
	})
}
