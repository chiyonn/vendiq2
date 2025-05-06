package spapi

import (
	"context"
)

type Requester interface {
	SendRequest(ctx context.Context, e *Endpoint, params interface{}) ([]byte, error)
	Logger() Logger
}

type Logger interface {
	Printf(format string, v ...interface{})
}
