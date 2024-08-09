package sqlite

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"sync"
)

type (
	Client struct {
		once sync.Once
		*gorm.DB
		*Config
	}

	Config struct {
		Dsn string
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

func (c *Client) Init(config *Config) {
	if c.DB == nil {
		c.once.Do(func() {
			db, err := gorm.Open(sqlite.Open(config.Dsn), &gorm.Config{})
			if err != nil {
				panic("failed to connect database")
			}
			c.DB = db
		})
	}
}
