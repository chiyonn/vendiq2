package inventory

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/chiyonn/vendiq2/pricer/pkg/spapi/client"
	"github.com/chiyonn/vendiq2/pricer/pkg/spapi/internal/rateutil"
	"github.com/chiyonn/vendiq2/pricer/pkg/spapi/types"
	"golang.org/x/time/rate"
)

var getInventorySummaries = rateutil.NewEndpointLimitedFunc[*GetInventorySummariesResponse](&types.Endpoint{
	Method: "GET",
	Path:   "/fba/inventory/v1/summaries",
	Rate:   rate.Every(2 * 1e9), // 2秒に1回
	Burst:  2,
})

func GetInventorySummaries(ctx context.Context, c *client.Client, params *GetInventorySummariesParams) (*GetInventorySummariesResponse, error) {
	return getInventorySummaries(ctx, c, func(ctx context.Context, c *client.Client, e *types.Endpoint) (*GetInventorySummariesResponse, error) {
		if err := rateutil.Wait(ctx, e); err != nil {
			return nil, fmt.Errorf("rate limit exceeded: %w", err)
		}

		body, err := c.SendRequest(ctx, e, params)
		if err != nil {
			return nil, fmt.Errorf("failed to send request: %w", err)
		}

		var response GetInventorySummariesResponse
		if err := json.Unmarshal(body, &response); err != nil {
			var apiErr types.APIError
			if json.Unmarshal(body, &apiErr) == nil {
				c.Logger.Printf("SPAPI error: code=%s, message=%s", apiErr.Code, apiErr.Message)
			}
			return nil, fmt.Errorf("failed to decode response: %w", err)
		}

		return &response, nil
	})
}

