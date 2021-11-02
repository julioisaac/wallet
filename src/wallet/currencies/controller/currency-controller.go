package controller

import (
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
		logs.Instance.Log.Error(request.Context(), "error trying decode currency upsert")
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(`{error: Error trying decode}`))
		return
	}
	err = currencyService.Validate(request.Context(), &currency)
	if err != nil {
		logs.Instance.Log.Error(request.Context(), "error invalid currency")
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(err.Error()))
		return
	}
	err = currencyService.Upsert(request.Context(), &currency)
	if err != nil {
		logs.Instance.Log.Error(request.Context(), "error trying upsert currency")
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(err.Error()))
		return
	}
	logs.Instance.Log.Debug(request.Context(), "currency: "+currency.Name+" successfully created")
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(currency)
}

func (c *currencyController) Delete(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	id := request.URL.Query().Get("id")
	currencies, err := currencyService.Remove(request.Context(), id)
	if err != nil {
		logs.Instance.Log.Error(request.Context(), "error trying remove currency")
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(err.Error()))
		return
	}
	logs.Instance.Log.Debug(request.Context(), "currency: "+id+" successfully deleted")
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(currencies)
}

func (c *currencyController) GetById(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	id := request.URL.Query().Get("id")
	currencies, err := currencyService.FindById(request.Context(), id)
	if err != nil {
		logs.Instance.Log.Error(request.Context(), "error trying find currency by id: "+id)
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(err.Error()))
		return
	}
	logs.Instance.Log.Debug(request.Context(), "currency: "+id+" successfully found")
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(currencies)
}

func (c *currencyController) GetAll(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	currencies := currencyService.FindAll(request.Context())
	logs.Instance.Log.Debug(request.Context(), "currencies successfully found")
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(currencies)
}