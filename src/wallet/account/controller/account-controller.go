package controller

import (
	"encoding/json"
	"fmt"
	"github.com/julioisaac/daxxer-api/internal/logs"
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
	cryptoRepo repository.DBRepository  = mongodb.NewMongodbRepository("daxxer", "crypto_currencies")
	historyRepo repository.DBRepository = mongodb.NewMongodbRepository("daxxer", "history")
	pricesRepo  repository.DBRepository = mongodb.NewMongodbRepository("daxxer", "prices")
	accountService = service.NewAccountService(accountRepo, cryptoRepo, historyRepo, pricesRepo)
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
		logs.Instance.Log.Error(request.Context(), "error trying decode account create")
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(`{error: Error trying decode}`))
		return
	}
	err = accountService.Create(request.Context(), &account)
	if err != nil {
		logs.Instance.Log.Error(request.Context(), "error trying create new account: "+account.UserName)
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(err.Error()))
		return
	}
	logs.Instance.Log.Debug(request.Context(), "account create success user: "+account.UserName)
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(account)
}

func (*controller) Deposit(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var transaction entity.Transaction
	err := json.NewDecoder(request.Body).Decode(&transaction)
	if err != nil {
		logs.Instance.Log.Error(request.Context(), "error trying decode transaction deposit")
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(`{error: Error trying decode}`))
		return
	}
	err = accountService.Deposit(request.Context(), &transaction)
	if err != nil {
		logs.Instance.Log.Error(request.Context(), "error trying deposit account: "+transaction.UserName)
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(err.Error()))
		return
	}
	logs.Instance.Log.Debug(request.Context(), "deposit success user: "+transaction.UserName+" currency: "+transaction.Amount.Currency+" amount: "+fmt.Sprintf("%v", transaction.Amount.Value))
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(transaction)
}

func (*controller) Withdraw(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var transaction entity.Transaction
	err := json.NewDecoder(request.Body).Decode(&transaction)
	if err != nil {
		logs.Instance.Log.Error(request.Context(), "error trying decode transaction withdraw")
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(`{error: Error trying decode}`))
		return
	}
	err = accountService.Withdraw(request.Context(), &transaction)
	if err != nil {
		logs.Instance.Log.Error(request.Context(), "error trying withdraw account: "+transaction.UserName)
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(err.Error()))
		return
	}
	logs.Instance.Log.Debug(request.Context(), "withdraw success user: "+transaction.UserName+" currency: "+transaction.Amount.Currency+" amount: "+fmt.Sprintf("%v", transaction.Amount.Value))
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(transaction)
}

func (*controller) Balance(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	user := request.URL.Query().Get("user")
	amounts := accountService.Balance(request.Context(), user)
	if amounts == nil {
		logs.Instance.Log.Error(request.Context(), "no balance found for this user: "+user)
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte("no balance found for this user"))
		return
	}
	logs.Instance.Log.Debug(request.Context(), "balance success user: "+user)
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(amounts)
}


func (*controller) History(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	page, err := strconv.Atoi(request.URL.Query().Get("page"))
	if err != nil {
		logs.Instance.Log.Error(request.Context(), "history request - page must be int")
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte("Page must be int"))
		return
	}
	limit, err2 := strconv.Atoi(request.URL.Query().Get("limit"))
	if err2 != nil {
		logs.Instance.Log.Error(request.Context(), "history request - limit must be int")
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
		logs.Instance.Log.Error(request.Context(), "history request - bad startDate format")
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte("bad startDate format"))
		return
	}
	endDt, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		logs.Instance.Log.Error(request.Context(), "history request - bad endDate format")
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte("bad endDate format"))
		return
	}

	histories := accountService.Histories(request.Context(), user, page, limit, 1, startDt, endDt)
	logs.Instance.Log.Debug(request.Context(), "history request success user: "+user)
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(histories)
}

