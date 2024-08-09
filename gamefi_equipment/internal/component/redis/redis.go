package redis

import (
	"components/database/redis"
	"errors"
	"gamefi_equipment/internal/conf"
	"github.com/go-kratos/kratos/v2/log"
)

type Redis struct {
	RedisCli *redis.Client
}

func NewRedis(c *conf.Data, logger log.Logger) (*Redis, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}

	if c.Redis == nil {
		return nil, nil, errors.New("redis config is nil")
	}

	var rediscli = redis.NewClient(&redis.Config{
		Addr:         c.Redis.Addr,
		Auth:         c.Redis.Auth,
		DialTimeout:  c.Redis.DialTimeout.AsDuration(),
		ReadTimeout:  c.Redis.ReadTimeout.AsDuration(),
		WriteTimeout: c.Redis.WriteTimeout.AsDuration(),
		Active:       int(c.Redis.Active),
		Idle:         int(c.Redis.Idle),
		IdleTimeout:  c.Redis.IdleTimeout.AsDuration(),
		SlowLog:      c.Redis.SlowLog.AsDuration(),
	})

	log.NewHelper(logger).Debug("redis client success")

	return &Redis{RedisCli: rediscli}, cleanup, nil
}
