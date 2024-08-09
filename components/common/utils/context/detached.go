package context

import (
	"context"
	"time"
)

// You can use this method until golang implemented detach context function offcially.
// Related issue: https://github.com/golang/go/issues/40221

// Detached return no deadline and cancel context
func Detach(ctx context.Context) context.Context {
	return detached{ctx: ctx}
}

type detached struct {
	ctx context.Context
}

func (detached) Deadline() (time.Time, bool) {
	return time.Time{}, false
}

func (d detached) Done() <-chan struct{} {
	return nil
}

func (d detached) Err() error {
	return nil
}

func (d detached) Value(key interface{}) interface{} {
	return d.ctx.Value(key)
}

func NewCtx() context.Context {
	return Detach(context.Background())
}
