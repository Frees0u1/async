package future

import (
	"context"
	"fmt"
	"github.com/Frees0u1/async/src/future"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestAwaitAll(t *testing.T) {
	future1 := future.NewFuture[int](func() (int, error) {
		time.Sleep(10 * time.Millisecond)
		return 1, nil
	})
	future2 := future.NewFuture[int](func() (int, error) {
		time.Sleep(20 * time.Millisecond)
		return 2, nil
	})
	future3 := future.NewFuture[int](func() (int, error) {
		time.Sleep(30 * time.Millisecond)
		return 3, nil
	})

	ctx := context.Background()
	futures := []future.Future[int]{future1, future2, future3}
	results, err := future.AwaitAll(ctx, nil, futures)
	assert.Nil(t, err)
	assert.Equal(t, 1, results[0])
	assert.Equal(t, 2, results[1])
	assert.Equal(t, 3, results[2])
}

func TestAwaitAllWhenError(t *testing.T) {
	future1 := future.NewFuture[int](func() (int, error) {
		time.Sleep(10 * time.Millisecond)
		return 1, nil
	})
	future2 := future.NewFuture[int](func() (int, error) {
		time.Sleep(20 * time.Millisecond)
		return 0, fmt.Errorf("error from 2")
	})
	future3 := future.NewFuture[int](func() (int, error) {
		time.Sleep(90 * time.Millisecond)
		return 3, nil
	})

	ctx := context.Background()
	futures := []future.Future[int]{future1, future2, future3}
	results, err := future.AwaitAll(ctx, nil, futures)
	assert.Errorf(t, err, "error from 2")
	assert.Equal(t, 0, results[0])
	assert.Equal(t, 0, results[1])
	assert.Equal(t, 0, results[2])
}

func TestFirst(t *testing.T) {
	future1 := future.NewFuture[int](func() (int, error) {
		time.Sleep(20 * time.Millisecond)
		return 1, nil
	})
	future2 := future.NewFuture[int](func() (int, error) {
		time.Sleep(10 * time.Millisecond)
		return 2, nil
	})
	future3 := future.NewFuture[int](func() (int, error) {
		time.Sleep(90 * time.Millisecond)
		return 3, nil
	})

	ctx := context.Background()
	futures := []future.Future[int]{future1, future2, future3}
	result, err := future.First(ctx, nil, futures)
	assert.Nil(t, err)
	assert.Equal(t, 2, result)
}

func TestFirstError(t *testing.T) {
	future1 := future.NewFuture[int](func() (int, error) {
		time.Sleep(20 * time.Millisecond)
		return 1, nil
	})
	future2 := future.NewFuture[int](func() (int, error) {
		time.Sleep(10 * time.Millisecond)
		return 0, fmt.Errorf("error from 2")
	})
	future3 := future.NewFuture[int](func() (int, error) {
		time.Sleep(90 * time.Millisecond)
		return 3, nil
	})

	ctx := context.Background()
	futures := []future.Future[int]{future1, future2, future3}
	result, err := future.First(ctx, nil, futures)
	assert.Errorf(t, err, "error from 2")
	assert.Equal(t, 0, result)
}
