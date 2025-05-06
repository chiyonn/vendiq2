package inventory

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/chiyonn/vendiq2/pricer/pkg/spapi"
)

type GetInventorySummaries struct {
	client   *spapi.Client
	endpoint *spapi.Endpoint
}

func (g *GetInventorySummaries) Run(ctx context.Context, params *GetInventorySummariesParams) (*GetInventorySummariesResponse, error) {
	body, err := g.client.SendRequest(ctx, g.endpoint, params)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	var response GetInventorySummariesResponse
	if err := json.Unmarshal(body, &response); err != nil {
		// 一応APIErrorもログに出す
		var apiErr spapi.APIError
		if json.Unmarshal(body, &apiErr) == nil {
			g.client.Logger.Printf("SPAPI error: code=%s, message=%s", apiErr.Code, apiErr.Message)
		}

		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}
