package controller

import (
	"encoding/json"
	"github.com/julioisaac/daxxer-api/src/helpers/repository"
	"github.com/julioisaac/daxxer-api/src/helpers/repository/mongodb"
	"github.com/julioisaac/daxxer-api/src/wallet/currencies/entity"
	"github.com/julioisaac/daxxer-api/src/wallet/currencies/service"
	"net/http"
)

var (
	currencyRepo repository.DBRepository = mongodb.NewMongodbRepository("daxxer", "currencies")
	currencyService = service.NewCurrencyService(currencyRepo)
)

type controller struct {}

func NewCurrencyController() CurrencyHandler {
	return &controller{}
}

func (c *controller) Upsert(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var currency entity.Currency
	err := json.NewDecoder(request.Body).Decode(&currency)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(`{error: Error trying decode}`))
		return
	}
	err = currencyService.Validate(&currency)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(err.Error()))
		return
	}
	err = currencyService.Upsert(&currency)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(err.Error()))
		return
	}
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(currency)
}

func (c *controller) Delete(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	id := request.URL.Query().Get("id")
	currencies, err := currencyService.Remove(id)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(err.Error()))
		return
	}
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(currencies)
}

func (c *controller) GetById(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	id := request.URL.Query().Get("id")
	currencies, err := currencyService.FindById(id)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(err.Error()))
		return
	}
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(currencies)
}

func (c *controller) GetAll(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	currencies := currencyService.FindAll()
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(currencies)
}