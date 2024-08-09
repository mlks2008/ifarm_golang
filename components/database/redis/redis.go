package redis

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/gomodule/redigo/redis"
	"golang.org/x/sync/singleflight"
	"sync"
	"time"
)

type (
	Client struct {
		once sync.Once
		pool *redis.Pool
		conf *Config
		g    singleflight.Group

		cmdable
	}

	Config struct {
		Addr         string
		Auth         string
		DialTimeout  time.Duration
		ReadTimeout  time.Duration
		WriteTimeout time.Duration
		Active       int
		Idle         int
		IdleTimeout  time.Duration
		SlowLog      time.Duration
	}
)

func NewClient(conf *Config) *Client {
	client := &Client{
		once: sync.Once{},
		conf: conf,
	}

	if conf.DialTimeout == 0 {
		conf.DialTimeout = 1
	}
	if conf.ReadTimeout == 0 {
		conf.ReadTimeout = 1
	}
	if conf.WriteTimeout == 0 {
		conf.WriteTimeout = 1
	}
	if conf.Active == 0 {
		conf.Active = 100
	}
	if conf.Idle == 0 {
		conf.Idle = 10
	}
	if conf.IdleTimeout == 0 {
		conf.IdleTimeout = 10
	}

	client.Init(conf)

	return client
}

func (c *Client) Name() string {
	return "redis"
}

func (c *Client) Init(conf *Config) {
	if c.pool == nil {
		c.once.Do(func() {
			c.pool = &redis.Pool{
				MaxIdle:     conf.Idle,
				MaxActive:   conf.Active,
				IdleTimeout: conf.IdleTimeout,
				Dial: func() (redis.Conn, error) {
					c, err := redis.Dial("tcp", conf.Addr,
						redis.DialClientName(c.Name()),
						redis.DialConnectTimeout(c.conf.DialTimeout),
						redis.DialReadTimeout(c.conf.ReadTimeout),
						redis.DialWriteTimeout(c.conf.WriteTimeout),
						redis.DialPassword(conf.Auth),
					)
					if err != nil {
						log.Infof("redis dial failed:%s,err:%v", conf.Addr, err)
						return nil, err
					}
					return c, err
				},
				TestOnBorrow: func(c redis.Conn, t time.Time) error {
					_, err := c.Do("PING")
					return err
				},
			}
			c.cmdable = c.process
		})
	}
}

func (c *Client) OnAfterInit() {
}

func (c *Client) OnBeforeStop() {
}

func (c *Client) OnStop() {
	_ = c.pool.Close()
}

func (c *Client) GetConn() redis.Conn {
	return c.pool.Get()
}

func (c *Client) Do(ctx context.Context, args ...interface{}) *Cmd {
	cmd := NewCmd(ctx, args...)
	_ = c.process(ctx, cmd)
	return cmd
}

func (c *Client) process(ctx context.Context, cmd Cmder) error {
	if cmd.Err() != nil {
		return cmd.Err()
	}

	conn := c.pool.Get()
	defer conn.Close()
	cmd.setReply(conn.Do(cmd.Name(), cmd.Args()[1:]...))
	return cmd.Err()
}

func (c *Client) processPipeline(ctx context.Context, cmds []Cmder) (err error) {
	// 验证参数构建问题
	for _, cmd := range cmds {
		if err = cmd.Err(); err != nil {
			return
		}
	}

	conn := c.pool.Get()
	defer conn.Close()

	for _, cmd := range cmds {
		err = conn.Send(cmd.Name(), cmd.Args()[1:]...)
		if err != nil {
			return
		}
	}

	err = conn.Flush()
	if err != nil {
		return
	}

	for _, cmd := range cmds {
		cmd.setReply(conn.Receive())
	}
	return
}

func (c *Client) Pipelined(ctx context.Context, fn func(Pipeliner) error) ([]Cmder, error) {
	return c.Pipeline().Pipelined(ctx, fn)
}

func (c *Client) Pipeline() Pipeliner {
	pipe := Pipeline{
		exec: c.processPipeline,
	}
	pipe.init()
	return &pipe
}
