package metrics

import "net/http"

type Metrics struct {
	Tracer interface{}
	Middleware interface{}
	Client *http.Client
}