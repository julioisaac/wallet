package service

import (
	"errors"
	"github.com/julioisaac/daxxer-api/src/helpers/repository"
	"github.com/julioisaac/daxxer-api/src/helpers/repository/mongodb"
	"github.com/julioisaac/daxxer-api/src/wallet/account/entity"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

type service struct {}

var (
	accountRepo repository.DBRepository = mongodb.NewMongodbRepository("daxxer", "account")
	historyRepo repository.DBRepository = mongodb.NewMongodbRepository("daxxer", "history")
)

func NewAccountService() *service  {
	return &service{}
}

func (s *service) Create(account *entity.Account) error {
	var storedAccount entity.Account
	var query = `{"username": "`+account.UserName+`"}`
	_ = accountRepo.FindOne(query, &storedAccount)
	if storedAccount.UserName != "" {
		err := errors.New("the account already exists")
		return err
	}
	return accountRepo.Insert(account)
}

func (s *service) Deposit(transaction *entity.Transaction) error {
	var account = entity.Account{}
	transaction.Type = "deposit"
	var query = `{"username": "`+transaction.UserName+`"}`
	err := accountRepo.FindOne(query, &account)
	if err != nil {
		return errors.New("account not found")
	}
	err1 := execTransaction(account.Deposit, transaction)
	if err1 != nil {
		return err1
	}
	return nil
}

func (s *service) Withdraw(transaction *entity.Transaction) error {
	var account = entity.Account{}
	transaction.Type = "withdraw"
	var query = `{"username": "`+transaction.UserName+`"}`
	err := accountRepo.FindOne(query, &account)
	if err != nil {
		return errors.New("account not found")
	}
	err1 := execTransaction(account.Withdraw, transaction)
	if err1 != nil {
		return err1
	}

	return nil
}

func (s *service) Histories(user string, page, limit int, sort int, startDate time.Time, endDate time.Time) []interface{} {
	var endDt = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, endDate.Nanosecond(), endDate.Location())
	var query = bson.M{ "eventtime": bson.M{ "$gte": startDate, "$lte": endDt }, "username": user}
	var histories = historyRepo.FindAll(page, limit, sort, query, new(entity.History))
	return histories
}

func execTransaction(operation func(entity.Amount) (*entity.Account, error), transaction *entity.Transaction) error {
	account, err := operation(transaction.Amount)
	if err != nil {
		return err
	}
	selector := bson.M{"username": transaction.UserName}
	update := bson.M{
		"$set": account,
	}
	err1 := accountRepo.Upsert(selector, update)
	if err1 != nil {
		return err1
	}
	err2 := saveHistory(transaction)
	if err2 != nil {
		return err2
	}
	return nil
}

func saveHistory(transaction *entity.Transaction) error {
	err := historyRepo.Insert(buildHistory(transaction))
	if err != nil {
		return err
	}
	return nil
}

func buildHistory(transaction *entity.Transaction) entity.History {
	return entity.History{
		UserName:  transaction.UserName,
		Type:      transaction.Type,
		Amount:    transaction.Amount,
		EventTime: time.Now(),
	}
}