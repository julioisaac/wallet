package controller

import (
	"encoding/json"
	"github.com/julioisaac/daxxer-api/src/helpers/repository"
	"github.com/julioisaac/daxxer-api/src/helpers/repository/mongodb"
	"github.com/julioisaac/daxxer-api/src/wallet/account/entity"
	"github.com/julioisaac/daxxer-api/src/wallet/account/service"
	"net/http"
	"strconv"
	"time"
)

var (
	accountRepo repository.DBRepository = mongodb.NewMongodbRepository("daxxer", "account")
	historyRepo repository.DBRepository = mongodb.NewMongodbRepository("daxxer", "history")
	pricesRepo  repository.DBRepository = mongodb.NewMongodbRepository("daxxer", "prices")
	accountService = service.NewAccountService(accountRepo, historyRepo, pricesRepo)
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
	var transaction entity.Transaction
	err := json.NewDecoder(request.Body).Decode(&transaction)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{error: Error trying decode}`))
		return
	}
	err = accountService.Deposit(&transaction)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
		return
	}
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(transaction)
}

func (*controller) Withdraw(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var transaction entity.Transaction
	err := json.NewDecoder(request.Body).Decode(&transaction)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{error: Error trying decode}`))
		return
	}
	err = accountService.Withdraw(&transaction)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
		return
	}
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(transaction)
}

func (*controller) Balance(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	user := request.URL.Query().Get("user")
	amounts := accountService.Balance(user)
	if amounts == nil {
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte("no balance found for this user"))
		return
	}
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(amounts)
}


func (*controller) History(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	page, err := strconv.Atoi(request.URL.Query().Get("page"))
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte("Page must be int"))
		return
	}
	limit, err2 := strconv.Atoi(request.URL.Query().Get("limit"))
	if err2 != nil {
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte("Limit must be int"))
		return
	}
	if limit > 100 {
		limit = 100
	}
	user := request.URL.Query().Get("user")
	startDate := request.URL.Query().Get("startDate")
	endDate := request.URL.Query().Get("endDate")

	startDt, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte("Bad startDate format"))
		return
	}
	endDt, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte("Bad endDate format"))
		return
	}

	histories := accountService.Histories(user, page, limit, 1, startDt, endDt)
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(histories)
}

