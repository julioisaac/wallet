package repository

import (
	entity2 "github.com/julioisaac/daxxer-api/src/wallet/prices/entity"
	"github.com/stretchr/testify/mock"
	"net/http"
)

type MockApiRepository struct {
	mock.Mock
	Client MockClient
}

func (m *MockApiRepository) GetPrices(cryptoCurrencies, currencies []interface{}) (*[]entity2.Price, error) {
	ret := m.Called(cryptoCurrencies, currencies)
	return ret.Get(0).(*[]entity2.Price), ret.Error(1)
}

func (m *MockApiRepository) DoRequest(params map[string]string) (*http.Response, error) {
	ret := m.Called(params)
	return ret.Get(0).(*http.Response), ret.Error(1)
}