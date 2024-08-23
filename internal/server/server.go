package server

import (
	"context"
	"github.com/nfangxu/mpc-go/internal/utils"
	"github.com/smallnest/rpcx/client"
	"time"
)

type option struct {
	retries  int
	interval time.Duration
}

func newOption(opts ...OptionFn) *option {
	opt := &option{
		retries:  10,
		interval: time.Second,
	}
	for _, fn := range opts {
		fn(opt)
	}
	return opt
}

type OptionFn func(*option)

func WithRetries(v int) OptionFn {
	return func(o *option) {
		o.retries = v
	}
}

func WithInterval(v time.Duration) OptionFn {
	return func(o *option) {
		o.interval = v
	}
}

type Server struct{}

type Empty struct{}

func (*Server) Ping(_ context.Context, _ *Empty, _ *Empty) error {
	return nil
}

func Ping(d client.ServiceDiscovery, opts ...OptionFn) error {
	opt := newOption(opts...)
	c := client.NewXClient("server", client.Failtry, client.RandomSelect, d, client.DefaultOption)
	defer c.Close()
	if err := utils.Try(func() error {
		return c.Call(context.Background(), "Ping", &Empty{}, &Empty{})
	}, opt.retries, opt.interval); err != nil {
		return err
	}
	return nil
}
