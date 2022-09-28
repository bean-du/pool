package pool

import (
	"context"
)

type Client struct {
	p    Pooler
	opts *Options
}

func NewClient(dialFunc DialFunc, log Logger, o ...Option) *Client {
	opts := &Options{Dialer: dialFunc}
	for _, fn := range o {
		fn(opts)
	}

	p := NewConnPool(opts, log)
	return &Client{
		p:    p,
		opts: opts,
	}
}

func (c *Client) Send(ctx context.Context, data []byte) error {
	var (
		conn *Conn
		err  error
	)
	if conn, err = c.p.Get(ctx); err != nil {
		return err
	}
	defer c.p.Put(ctx, conn)

	if c.opts.WriteFunc != nil {
		err = conn.WithWriter(ctx, 0, c.opts.WriteFunc(data))
	} else {
		_, err = conn.Write(data)
	}
	return err
}

func (c *Client) Close() error {
	return c.p.Close()
}

type PoolStats Stats

// PoolStats returns connection pool stats.
func (c *Client) PoolStats() *PoolStats {
	stats := c.p.Stats()
	return (*PoolStats)(stats)
}
