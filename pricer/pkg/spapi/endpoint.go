package spapi

type Endpoint struct {
	Method string
	Path   string
	Rate   float64
	Burst  int
}
