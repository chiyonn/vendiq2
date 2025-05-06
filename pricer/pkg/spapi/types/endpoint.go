package types

type Endpoint struct {
	Method string
	Path   string
	Rate   float64
	Burst  int
}

type ErrorList struct {
	code    string
	message *string
	details *string
}
