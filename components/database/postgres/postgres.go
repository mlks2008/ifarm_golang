package postgres

import (
	"github.com/go-pg/pg/extra/pgdebug/v10"
	"github.com/go-pg/pg/v10"
	"sync"
	"time"
)

type (
	Client struct {
		once sync.Once
		*pg.DB
		*Config
	}

	Config struct {
		Dsn             string
		ApplicationName string
		PoolSize        int
		ReadTimeout     time.Duration
		WriteTimeout    time.Duration
		MinIdleConns    int
		IdleTimeout     time.Duration
		Debug           bool
	}
)

func NewClient(config *Config) *Client {
	client := &Client{
		once:   sync.Once{},
		Config: config,
	}

	client.Init(config)

	return client
}

func (c *Client) Name() string {
	return c.Config.ApplicationName
}

func (c *Client) Init(config *Config) {
	if c.DB == nil {
		c.once.Do(func() {
			//opt, err := pg.ParseURL("postgresql://postgres:123456@127.0.0.1:5432/postgres?sslmode=disable")
			opt, err := pg.ParseURL(config.Dsn)
			if err != nil {
				panic(err)
			}

			if len(config.ApplicationName) > 0 {
				opt.ApplicationName = config.ApplicationName
			}
			if config.PoolSize != 0 {
				opt.PoolSize = config.PoolSize
			} else {
				opt.PoolSize = 100
			}
			if config.ReadTimeout != 0 {
				opt.ReadTimeout = config.ReadTimeout
			} else {
				opt.ReadTimeout = 3 * time.Second
			}
			if config.WriteTimeout != 0 {
				opt.WriteTimeout = config.WriteTimeout
			} else {
				opt.WriteTimeout = 3 * time.Second
			}
			if config.MinIdleConns != 0 {
				opt.MinIdleConns = config.MinIdleConns
			} else {
				opt.MinIdleConns = 20
			}
			if config.IdleTimeout != 0 {
				opt.IdleTimeout = config.IdleTimeout
			} else {
				opt.IdleTimeout = 3 * time.Second
			}

			connect := pg.Connect(opt)

			debugHook := pgdebug.NewDebugHook()
			debugHook.Verbose = config.Debug
			debugHook.EmptyLine = true
			connect.AddQueryHook(debugHook)

			c.DB = connect
		})
	}
}

func (c *Client) OnAfterInit() {
}

func (c *Client) OnBeforeStop() {
}

func (c *Client) OnStop() {
	c.DB.Close()
}
