package future

import (
	"context"
	"fmt"
	"github.com/Frees0u1/async/src/util"
	"sync"
	"time"
)

func Complete[T any](val T, err error) Future[T] {
	return NewFuture(func() (T, error) {
		return val, err
	})
}

func AwaitAll[T any](ctx context.Context, timeout *time.Duration, futures []Future[T]) ([]T, error) {
	size := len(futures)
	result := make([]T, size, size)
	emptyResult := make([]T, size, size)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	if size == 0 {
		return emptyResult, nil
	}

	errChan := make(chan error)
	allDone := make(chan interface{})
	var wg sync.WaitGroup
	wg.Add(size)
	for i, f := range futures {
		go func(i int, f Future[T]) {
			defer func() {
				e := util.RecoverAsError()
				if e != nil {
					errChan <- e
				}
				wg.Done()
			}()

			r, e := f.Await(ctx)
			if e != nil {
				errChan <- e
				return
			}
			result[i] = r
		}(i, f)
	}

	go func() {
		wg.Wait()
		close(allDone)
	}()

	var timeoutChan <-chan time.Time
	if timeout != nil {
		timeoutChan = time.After(*timeout)
	}

	select {
	case e := <-errChan:
		return emptyResult, e
	case <-allDone:
		return result, nil
	case <-timeoutChan:
		return emptyResult, fmt.Errorf("future timeout after %s", timeout)
	}
}

func First[T any](ctx context.Context, timeout *time.Duration, futures []Future[T]) (T, error) {
	var result T
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	errChan := make(chan error)
	resultChan := make(chan T)

	for i, f := range futures {
		go func(i int, f Future[T]) {
			defer func() {
				e := util.RecoverAsError()
				if e != nil {
					errChan <- e
				}
			}()

			r, e := f.Await(ctx)
			if e != nil {
				errChan <- e
				return
			}

			resultChan <- r
		}(i, f)
	}

	timeoutChan := make(<-chan time.Time)
	if timeout != nil {
		timeoutChan = time.After(*timeout)
	}

	select {
	case e := <-errChan:
		return result, e
	case result = <-resultChan:
		return result, nil
	case <-timeoutChan:
		return result, fmt.Errorf("future timeout after %s", timeout)
	}
}
