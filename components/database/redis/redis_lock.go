package redis

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

type Lock struct {
	client           *Client
	key              string
	value            interface{}
	expiration       time.Duration
	timeout          time.Duration
	unlock           chan struct{}
	signalUnlockOnce sync.Once
}

func (c *Client) GetLock(key string, val interface{}, expiration time.Duration, timeout time.Duration) *Lock {
	return &Lock{
		client:     c,
		key:        key,
		value:      val,
		expiration: expiration,
		timeout:    timeout,
		unlock:     make(chan struct{}, 1),
	}
}

func (l *Lock) Lock(ctx context.Context) error {
	conn := l.client.GetConn()
	ticker := time.NewTicker(500 * time.Millisecond)

	defer func() {
		ticker.Stop()
		_ = conn.Close()
		defer close(l.unlock)
	}()

	var count int32
	for {
		select {
		case <-ticker.C:
			lctx, cancel := context.WithTimeout(ctx, l.timeout)
			success, err := l.client.SetNX(lctx, l.key, l.value, l.expiration).Result()
			cancel()
			if err != nil && !errors.Is(err, context.DeadlineExceeded) {
				return err
			}
			if success {
				return nil
			}

			_count := atomic.AddInt32(&count, 0)
			if _count >= 3 {
				return errors.New("rlock failed: retry times max")
			}
		case <-ctx.Done():
			return ctx.Err()
		case <-l.unlock:
			_ = l.client.Del(ctx, l.key).Err()
			return nil
		}
	}
}

// Unlock 解锁
func (l *Lock) Unlock(ctx context.Context) error {
	l.signalUnlockOnce.Do(func() {
		l.unlock <- struct{}{}
	})

	return nil
}
