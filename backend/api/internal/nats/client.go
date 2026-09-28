package nats

import (
	natsgo "github.com/nats-io/nats.go"
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

func (c *Client) Subscribe(subject string, handler func(data []byte)) (*natsgo.Subscription, error) {
	return c.conn.Subscribe(subject, func(msg *natsgo.Msg) {
		handler(msg.Data)
	})
}
