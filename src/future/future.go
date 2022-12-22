package future

import (
	"async/src/constant"
	"async/src/util"
	"context"
	"fmt"
	"time"
)

type Future[T any] interface {
	Await(ctx context.Context, timeout *time.Duration) (T, error)
}

type FutureImpl[T any] struct {
	done   chan interface{}
	result T
	err    error
}

func NewFuture[T any](fn func() (T, error)) Future[T] {
	f := &FutureImpl[T]{
		done: make(chan interface{}),
	}
	go func() {
		defer func() {
			close(f.done)
			e := util.RecoverAsError()

			if e != nil {
				var zeroR T
				f.result, f.err = zeroR, e
			}
		}()

		f.result, f.err = fn()
	}()

	return f
}

func (f *FutureImpl[T]) Await(ctx context.Context, timeout *time.Duration) (T, error) {
	var timer *time.Timer
	if timeout == nil {
		timer = time.NewTimer(time.Duration(constant.DefaultTimeoutInSeconds) * time.Second)
	} else {
		timer = time.NewTimer(*timeout)
	}

	var zeroR T
	select {
	case <-ctx.Done():
		return zeroR, ctx.Err()
	case <-f.done:
		return f.result, f.err
	case <-timer.C:
		return zeroR, fmt.Errorf("Future timeout after %s", timeout)
	}
}
