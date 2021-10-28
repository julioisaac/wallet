package repository

import (
	"github.com/julioisaac/daxxer-api/src/wallet/prices/entity"
	"net/http"
)

type ApiRepository interface {
	GetPrices([]interface{}, []interface{}) (*[]entity.Price, error)
	DoRequest(params map[string]string) (*http.Response, error)
}