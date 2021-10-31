package pkg

import (
	"io"
	"net/http"
)

type healthCheck struct {}

func NewHealthCheck() *healthCheck {
	return &healthCheck{}
}

func (*healthCheck) IsAlive(response http.ResponseWriter, request *http.Request) {

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)

	io.WriteString(response, `{"alive": true}`)
}