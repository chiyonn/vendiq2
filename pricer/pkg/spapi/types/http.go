package types

import (
	"net/http"
	"net/url"
)

type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

type Queryable interface {
	Stringfy() url.Values
}


