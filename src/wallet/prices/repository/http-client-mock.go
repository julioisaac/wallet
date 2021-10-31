package repository

import (
	"github.com/stretchr/testify/mock"
	"net/http"
)

type MockClient struct {
	mock.Mock
}

func (c *MockClient) Do(req *http.Request) (*http.Response, error) {
	ret := c.Called(req)
	return ret.Get(0).(*http.Response), ret.Error(1)
}