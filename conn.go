package pool

import (
	"bufio"
	"context"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

var noDeadline = time.Time{}

type Conn struct {
	netConn   net.Conn
	writer    *bufio.Writer
	reader    *bufio.Reader
	Inited    bool      // 是否完成初始化
	pooled    bool      // 是否放进连接池
	createdAt time.Time // 创建时间
	usedAt    int64     // 使用时间
	ioLock    sync.Mutex
}

func NewConn(netConn net.Conn) *Conn {
	conn := &Conn{
		netConn:   netConn,
		writer:    bufio.NewWriter(netConn),
		reader:    bufio.NewReader(netConn),
		createdAt: time.Now(),
	}

	conn.SetUseAt(time.Now())
	return conn
}

func (c *Conn) UsedAt() time.Time {
	unix := atomic.LoadInt64(&c.usedAt)
	return time.Unix(unix, 0)
}

func (c *Conn) SetUseAt(now time.Time) {
	atomic.StoreInt64(&c.usedAt, now.Unix())
}

func (c *Conn) SetNetConn(netConn net.Conn) {
	c.netConn = netConn
	c.reader = bufio.NewReader(netConn)
	c.writer = bufio.NewWriter(netConn)
}

func (c *Conn) Write(b []byte) (int, error) {
	return c.netConn.Write(b)
}

func (c *Conn) RemoteAddr() net.Addr {
	if c.netConn != nil {
		return c.netConn.RemoteAddr()
	}
	return nil
}

func (c *Conn) WithReader(ctx context.Context, timeout time.Duration, fn func(rd net.Conn) error) error {
	if timeout != 0 {
		if err := c.netConn.SetReadDeadline(c.deadline(ctx, timeout)); err != nil {
			return err
		}
	}
	return fn(c.netConn)
}

func (c *Conn) WithWriter(ctx context.Context, timeout time.Duration, fn func(writer io.Writer) error) error {
	if timeout != 0 {
		if err := c.netConn.SetWriteDeadline(c.deadline(ctx, timeout)); err != nil {
			return err
		}
	}

	if c.writer.Buffered() > 0 {
		c.writer.Reset(c.netConn)
	}

	if err := fn(c.writer); err != nil {
		return err
	}

	return c.writer.Flush()
}

func (c *Conn) Close() error {
	return c.netConn.Close()
}

func (c *Conn) deadline(ctx context.Context, timeout time.Duration) time.Time {
	tm := time.Now()
	c.SetUseAt(tm)

	if timeout > 0 {
		tm.Add(timeout)
	}

	if ctx != nil {
		deadline, ok := ctx.Deadline()
		if ok {
			if timeout == 0 {
				return deadline
			}

			if deadline.Before(tm) {
				return deadline
			}
			return tm
		}
	}

	if timeout > 0 {
		return tm
	}
	return noDeadline
}
