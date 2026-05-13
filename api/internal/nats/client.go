package nats

import (
	"encoding/json"
	"log"

	"github.com/NirajDonga/pingpong/api/internal/result"
	natsgo "github.com/nats-io/nats.go"
)

const (
	CheckResultsSubject = "check.results"
	APIQueue            = "api"
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

func (c *Client) SubscribeCheckResults(handler func(result.CheckResult)) (*natsgo.Subscription, error) {
	return c.conn.QueueSubscribe(CheckResultsSubject, APIQueue, func(msg *natsgo.Msg) {
		var checkResult result.CheckResult
		if err := json.Unmarshal(msg.Data, &checkResult); err != nil {
			log.Printf("failed to decode check result: %v", err)
			return
		}

		handler(checkResult)
	})
}
