package nats

import (
	"context"
	"encoding/json"

	"github.com/NirajDonga/pingpong/backend/scheduler/internal/scheduler"
	natsgo "github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

const CheckJobsSubject = "check.jobs"

type Publisher struct {
	conn *natsgo.Conn
	js   jetstream.JetStream
}

func NewPublisher(url string) (*Publisher, error) {
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
		Name:     "JOBS",
		Subjects: []string{CheckJobsSubject},
	})
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &Publisher{conn: conn, js: js}, nil
}

func (p *Publisher) Close() {
	p.conn.Close()
}

func (p *Publisher) PublishCheckJob(ctx context.Context, job scheduler.CheckJob) error {
	data, err := json.Marshal(job)
	if err != nil {
		return err
	}

	_, err = p.js.Publish(ctx, CheckJobsSubject, data)
	return err
}
