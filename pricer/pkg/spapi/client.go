package spapi

import (
    "net/http"
    "context"
    "log"
)

type Client struct {
	httpClient *http.Client
	logger     *log.Logger
	// ...
}

func (c *Client) SendRequest(ctx context.Context, e *Endpoint, params interface{}) ([]byte, error) {
    return nil, nil
}

func (c *Client) Logger() Logger {
	return c.logger
}
