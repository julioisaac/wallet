package service

import (
	"errors"
	"github.com/julioisaac/daxxer-api/src/helpers/repository"
	"github.com/julioisaac/daxxer-api/src/helpers/repository/mongodb"
	"github.com/julioisaac/daxxer-api/src/wallet/account/entity"
	"go.mongodb.org/mongo-driver/bson"
)

type service struct {}

var (
	accountRepo repository.DBRepository = mongodb.NewMongodbRepository("daxxer", "account")
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
	return nil
}