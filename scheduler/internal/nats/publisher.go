package nats

import (
	"encoding/json"

	"github.com/NirajDonga/pingpong/scheduler/internal/scheduler"
	natsgo "github.com/nats-io/nats.go"
)

const CheckJobsSubject = "check.jobs"

type Publisher struct {
	conn *natsgo.Conn
}

func NewPublisher(url string) (*Publisher, error) {
	conn, err := natsgo.Connect(url)
	if err != nil {
		return nil, err
	}

	return &Publisher{conn: conn}, nil
}

func (p *Publisher) Close() {
	p.conn.Close()
}

func (p *Publisher) PublishCheckJob(job scheduler.CheckJob) error {
	data, err := json.Marshal(job)
	if err != nil {
		return err
	}

	return p.conn.Publish(CheckJobsSubject, data)
}
