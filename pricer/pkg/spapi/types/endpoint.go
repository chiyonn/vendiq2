package types

import "golang.org/x/time/rate"

type Endpoint struct {
	Method string
	Path   string
	Rate   rate.Limit
	Burst  int
}

type ErrorList struct {
	code    string
	message *string
	details *string
}

type APIError struct {
	Code    string
	Message string
}
