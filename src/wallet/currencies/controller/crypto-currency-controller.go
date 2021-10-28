package controller

import (
	"encoding/json"
	"github.com/julioisaac/daxxer-api/src/wallet/currencies/entity"
	"github.com/julioisaac/daxxer-api/src/wallet/currencies/service"
	"net/http"
)

var (
	cryptoService = service.NewCryptoCurrencyService()
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
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{error: Error trying decode}`))
		return
	}
	err = cryptoService.Validate(&crypto)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(err.Error()))
		return
	}
	err = cryptoService.Upsert(&crypto)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
		return
	}
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(crypto)
}

func (c *cryptoController) Delete(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	id := request.URL.Query().Get("id")
	currencies, err := cryptoService.Remove(id)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
		return
	}
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(currencies)
}

func (c *cryptoController) GetById(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	id := request.URL.Query().Get("id")
	currencies, err := cryptoService.FindById(id)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
		return
	}
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(currencies)
}

func (c *cryptoController) GetAll(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	cryptocurrencies := cryptoService.FindAll()
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(cryptocurrencies)
}