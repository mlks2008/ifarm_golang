package redis

import (
	"context"
	"testing"
)

func TestConn(t *testing.T) {
	client := NewClient(&Config{
		Addr:         "10.40.10.15:6379",
		Auth:         "game123456",
		DialTimeout:  10,
		ReadTimeout:  10,
		WriteTimeout: 10,
		Active:       10,
		Idle:         10,
		IdleTimeout:  10,
	})

	val := client.Get(context.Background(), "GF01:platform:user:526:rune:31:speed:lock").Val()

	t.Logf("val:%s", val)
}
