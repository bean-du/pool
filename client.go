package pool

import (
	"context"
)

type Client struct {
	p    Pooler
	opts *Options
}

func NewClient(dialFunc DialFunc, o ...Option) *Client {
	opts := &Options{Dialer: dialFunc}
	for _, fn := range o {
		fn(opts)
	}

	p := NewConnPool(opts)
	return &Client{
		p:    p,
		opts: opts,
	}
}

func (c *Client) Send(ctx context.Context, data []byte) error {
	conn, err := c.p.Get(ctx)
	if err != nil {
		return err
	}
	return conn.WithWriter(ctx, 0, WsWriter(data))
}

func (c *Client) Close() error {
	return c.p.Close()
}
