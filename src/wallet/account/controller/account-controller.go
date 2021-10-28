package controller

import (
	"encoding/json"
	"github.com/julioisaac/daxxer-api/src/wallet/account/entity"
	"github.com/julioisaac/daxxer-api/src/wallet/account/service"
	"net/http"
)

var (
	accountService = service.NewAccountService()
)

type controller struct {}

func NewAccountController() *controller {
	return &controller{}
}

func (*controller) Create(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var account entity.Account
	err := json.NewDecoder(request.Body).Decode(&account)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{error: Error trying decode}`))
		return
	}
	err = accountService.Create(&account)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
		return
	}
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(account)
}

func (*controller) Deposit(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var transaction *entity.Transaction
	err := json.NewDecoder(request.Body).Decode(&transaction)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{error: Error trying decode}`))
		return
	}
	err = accountService.Deposit(transaction)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
		return
	}
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(transaction)
}