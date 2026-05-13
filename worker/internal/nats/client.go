package nats

import (
	"encoding/json"
	"log"

	"github.com/NirajDonga/pingpong/worker/internal/worker"
	natsgo "github.com/nats-io/nats.go"
)

const (
	CheckJobsSubject    = "check.jobs"
	CheckResultsSubject = "check.results"
	WorkersQueue        = "workers"
)

type Client struct {
	conn *natsgo.Conn
}

func NewClient(url string) (*Client, error) {
	conn, err := natsgo.Connect(url)
	if err != nil {
		return nil, err
	}

	return &Client{conn: conn}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) SubscribeCheckJobs(handler func(worker.CheckJob)) (*natsgo.Subscription, error) {
	return c.conn.QueueSubscribe(CheckJobsSubject, WorkersQueue, func(msg *natsgo.Msg) {
		var job worker.CheckJob
		if err := json.Unmarshal(msg.Data, &job); err != nil {
			log.Printf("failed to decode check job: %v", err)
			return
		}

		handler(job)
	})
}

func (c *Client) PublishCheckResult(result worker.CheckResult) error {
	data, err := json.Marshal(result)
	if err != nil {
		return err
	}

	return c.conn.Publish(CheckResultsSubject, data)
}
