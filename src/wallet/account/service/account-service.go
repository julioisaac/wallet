package service

import (
	"errors"
	"github.com/julioisaac/daxxer-api/src/helpers/repository"
	"github.com/julioisaac/daxxer-api/src/helpers/repository/mongodb"
	"github.com/julioisaac/daxxer-api/src/wallet/account/entity"
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