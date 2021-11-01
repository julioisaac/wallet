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
	cryptoRepo repository.DBRepository = mongodb.NewMongodbRepository("daxxer", "cryptoCurrencies")
	cryptoService = service.NewCryptoCurrencyService(cryptoRepo)
)

type cryptoController struct {}

func NewCryptoCurrencyController() CurrencyHandler {
	return &cryptoController{}
}

func (c *cryptoController) Upsert(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var crypto entity.CryptoCurrency
	err := json.NewDecoder(request.Body).Decode(&crypto)
	if err != nil {
		logs.Instance.Log.Error(context.Background(), "error trying decode crypto currency upsert")
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(`{error: Error trying decode}`))
		return
	}
	err = cryptoService.Validate(&crypto)
	if err != nil {
		logs.Instance.Log.Error(context.Background(), "error invalid crypto currency")
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(err.Error()))
		return
	}
	err = cryptoService.Upsert(&crypto)
	if err != nil {
		logs.Instance.Log.Error(context.Background(), "error trying upsert crypto currency")
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(err.Error()))
		return
	}
	logs.Instance.Log.Debug(context.Background(), "crypto currency: "+crypto.Symbol+" successfully created")
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(crypto)
}

func (c *cryptoController) Delete(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	id := request.URL.Query().Get("id")
	currencies, err := cryptoService.Remove(id)
	if err != nil {
		logs.Instance.Log.Error(context.Background(), "error trying remove crypto currency")
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(err.Error()))
		return
	}
	logs.Instance.Log.Debug(context.Background(), "crypto currency: "+id+" successfully deleted")
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(currencies)
}

func (c *cryptoController) GetById(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	id := request.URL.Query().Get("id")
	currencies, err := cryptoService.FindById(id)
	if err != nil {
		logs.Instance.Log.Error(context.Background(), "error trying find crypto currency by id: "+id)
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(err.Error()))
		return
	}
	logs.Instance.Log.Debug(context.Background(), "crypto currency: "+id+" successfully found")
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(currencies)
}

func (c *cryptoController) GetAll(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	cryptocurrencies := cryptoService.FindAll()
	logs.Instance.Log.Debug(context.Background(), "crypto currencies successfully found")
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(cryptocurrencies)
}