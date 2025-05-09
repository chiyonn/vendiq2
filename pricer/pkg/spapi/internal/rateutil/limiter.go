package rateutil

import (
	"context"
	"fmt"
	"sync"

	"golang.org/x/time/rate"

	"github.com/chiyonn/vendiq2/pricer/pkg/spapi/client"
	"github.com/chiyonn/vendiq2/pricer/pkg/spapi/types"
)

// 内部キャッシュ用
var limiterMap sync.Map // key: *types.Endpoint, value: *rate.Limiter

func getLimiter(e *types.Endpoint) *rate.Limiter {
	actual, _ := limiterMap.LoadOrStore(e, rate.NewLimiter(e.Rate, e.Burst))
	return actual.(*rate.Limiter)
}

// Endpointとレスポンス型を受け取る関数を返す
func NewEndpointLimitedFunc[T any](e *types.Endpoint) func(
	ctx context.Context,
	c *client.Client,
	fn func(ctx context.Context, c *client.Client, e *types.Endpoint) (T, error),
) (T, error) {
	limiter := getLimiter(e)

	return func(ctx context.Context, c *client.Client, fn func(ctx context.Context, c *client.Client, e *types.Endpoint) (T, error)) (T, error) {
		var zero T
		if err := limiter.Wait(ctx); err != nil {
			return zero, fmt.Errorf("rate limit exceeded: %w", err)
		}
		return fn(ctx, c, e)
	}
}

func Wait(ctx context.Context, e *types.Endpoint) error {
	key := fmt.Sprintf("%s:%s", e.Method, e.Path)
	limiterIface, _ := limiterMap.LoadOrStore(key, rate.NewLimiter(e.Rate, e.Burst))
	limiter := limiterIface.(*rate.Limiter)
	return limiter.Wait(ctx)
}

