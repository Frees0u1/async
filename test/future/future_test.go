package future

import (
	"context"
	"fmt"
	"github.com/Frees0u1/async/src/future"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAwait(t *testing.T) {
	future1 := future.NewFuture[int](func() (int, error) {
		time.Sleep(10 * time.Millisecond)
		return 1, nil
	})

	future2HasError := future.NewFuture[float64](func() (float64, error) {
		time.Sleep(20 * time.Millisecond)
		return 0, fmt.Errorf("error")
	})

	ctx := context.Background()
	res1, e := future1.Await(ctx, nil)
	assert.Equal(t, 1, res1)
	assert.Nil(t, e)

	res2, e := future2HasError.Await(ctx, nil)
	assert.Error(t, e)
	assert.True(t, math.Abs(res2-0) < 0.001)
}
