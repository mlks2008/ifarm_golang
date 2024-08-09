package redis2

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/shopspring/decimal"
	"time"
)

type RedisCli struct {
	client *redis.Client
}

func NewRedisCli(host, password string, db int) *RedisCli {
	rdb := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: password,
		DB:       db,
	})

	return &RedisCli{client: rdb}
}

func (this *RedisCli) SetEX(key string, val interface{}, expiration time.Duration) {
	this.client.SetEX(context.Background(), key, val, expiration)
}

func (this *RedisCli) GetDecimal(key string) (decimal.Decimal, error) {
	val, err := this.client.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return decimal.Zero, nil
	} else if err != nil {
		return decimal.Zero, err
	} else {
		return decimal.NewFromString(val)
	}
}

func (this *RedisCli) GetString(key string) (string, error) {
	val, err := this.client.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		return "", err
	} else {
		return val, nil
	}
}
