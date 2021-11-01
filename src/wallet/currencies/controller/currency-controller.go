package controller

import (
	"context"
	"encoding/json"
	"github.com/julioisaac/daxxer-api/internal/logs"
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

type currencyController struct {}

func NewCurrencyController() CurrencyHandler {
	return &currencyController{}
}

func (c *currencyController) Upsert(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var currency entity.Currency
	err := json.NewDecoder(request.Body).Decode(&currency)
	if err != nil {
		logs.Instance.Log.Error(context.Background(), "error trying decode currency upsert")
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(`{error: Error trying decode}`))
		return
	}
	err = currencyService.Validate(&currency)
	if err != nil {
		logs.Instance.Log.Error(context.Background(), "error invalid currency")
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(err.Error()))
		return
	}
	err = currencyService.Upsert(&currency)
	if err != nil {
		logs.Instance.Log.Error(context.Background(), "error trying upsert currency")
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(err.Error()))
		return
	}
	logs.Instance.Log.Debug(context.Background(), "currency: "+currency.Name+" successfully created")
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(currency)
}

func (c *currencyController) Delete(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	id := request.URL.Query().Get("id")
	currencies, err := currencyService.Remove(id)
	if err != nil {
		logs.Instance.Log.Error(context.Background(), "error trying remove currency")
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(err.Error()))
		return
	}
	logs.Instance.Log.Debug(context.Background(), "currency: "+id+" successfully deleted")
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(currencies)
}

func (c *currencyController) GetById(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	id := request.URL.Query().Get("id")
	currencies, err := currencyService.FindById(id)
	if err != nil {
		logs.Instance.Log.Error(context.Background(), "error trying find currency by id: "+id)
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(err.Error()))
		return
	}
	logs.Instance.Log.Debug(context.Background(), "currency: "+id+" successfully found")
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(currencies)
}

func (c *currencyController) GetAll(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	currencies := currencyService.FindAll()
	logs.Instance.Log.Debug(context.Background(), "currencies successfully found")
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(currencies)
}