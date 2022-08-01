package pool

import (
	"context"
	"net"
	"time"
)

type DialFunc func(context.Context) (net.Conn, error)

type ReadFunc func(conn net.Conn) error

type ReceiveHandler func(context.Context, []byte)

type KeepAliveFunc func(conn net.Conn)

type Option func(o *Options)

type Options struct {
	Dialer         DialFunc
	OnClose        func(*Conn) error
	ReceiveHandler ReceiveHandler
	ReadFunc       ReadFunc
	Keepalive      KeepAliveFunc

	PoolFIFO           bool
	PoolSize           int
	MinIdleConns       int
	MaxConnAge         time.Duration
	PoolTimeout        time.Duration
	IdleTimeout        time.Duration
	IdleCheckFrequency time.Duration
}

func WithReadFunc(fn ReadFunc) Option {
	return func(o *Options) {
		o.ReadFunc = fn
	}
}

func WithKeepAlive(fn KeepAliveFunc) Option {
	return func(o *Options) {
		o.Keepalive = fn
	}
}

func WithReceiveHandle(fn ReceiveHandler) Option {
	return func(o *Options) {
		o.ReceiveHandler = fn
	}
}

func WithPoolFIFO(b bool) Option {
	return func(o *Options) {
		o.PoolFIFO = b
	}
}

func WithPoolSize(i int) Option {
	return func(o *Options) {
		o.PoolSize = i
	}
}

func WithMinIdleConns(i int) Option {
	return func(o *Options) {
		o.MinIdleConns = i
	}
}
func WithMaxConnAge(d time.Duration) Option {
	return func(o *Options) {
		o.MaxConnAge = d
	}
}

func WithPoolTimeout(d time.Duration) Option {
	return func(o *Options) {
		o.PoolTimeout = d
	}
}

func WithIdleTimeout(d time.Duration) Option {
	return func(o *Options) {
		o.IdleTimeout = d
	}
}

func WithIdleCheckFrequency(d time.Duration) Option {
	return func(o *Options) {
		o.IdleCheckFrequency = d
	}
}
