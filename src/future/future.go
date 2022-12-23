package future

import (
	"context"
	"fmt"
	"github.com/Frees0u1/async/src/util"
	"sync"
	"time"
)

type Future[T any] interface {
	Await(ctx context.Context) (T, error)
	runAsync(fn func() (T, error))
	complete(T, error)
}

type futureResult[T any] struct {
	result T
	err    error
}

type FutureImpl[T any] struct {
	done        chan futureResult[T]
	timeout     *time.Duration
	timeoutChan <-chan time.Time
	doneOnce    sync.Once
	closeOnce   sync.Once
}

func NewFuture[T any](fn func() (T, error)) Future[T] {
	f := &FutureImpl[T]{
		done:        make(chan futureResult[T], 1),
		timeoutChan: nil,
	}
	f.runAsync(fn)
	return f
}

func NewFutureWithTimeout[T any](fn func() (T, error), timeoutMs int) Future[T] {
	timeout := time.Duration(timeoutMs) * time.Millisecond
	f := &FutureImpl[T]{
		done:        make(chan futureResult[T], 1),
		timeout:     &timeout,
		timeoutChan: time.After(time.Duration(timeoutMs) * time.Millisecond),
	}
	f.runAsync(fn)
	return f
}

func (f *FutureImpl[T]) Await(ctx context.Context) (T, error) {
	defer f.closeOnce.Do(func() {
		close(f.done)
	})
	var zeroR T

	go func() {
		select {
		case <-ctx.Done():
			f.complete(zeroR, ctx.Err())
		case <-f.timeoutChan:
			f.complete(zeroR, fmt.Errorf("timeout %v", f.timeout))
		}
	}()

	x := <-f.done
	return x.result, x.err
}

func (f *FutureImpl[T]) complete(result T, err error) {
	f.doneOnce.Do(func() {
		f.done <- futureResult[T]{
			result: result,
			err:    err,
		}
	})
}

func (f *FutureImpl[T]) runAsync(fn func() (T, error)) {
	go func() {
		defer func() {
			e := util.RecoverAsError()
			if e != nil {
				var zeroR T
				f.complete(zeroR, e)
			}
		}()

		r, e := fn()
		f.complete(r, e)
	}()
}
