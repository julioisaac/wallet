package service

import (
	"errors"
	"github.com/julioisaac/daxxer-api/src/helpers/repository"
	"github.com/julioisaac/daxxer-api/src/helpers/repository/mongodb"
	"github.com/julioisaac/daxxer-api/src/wallet/currencies/entity"
	"go.mongodb.org/mongo-driver/bson"
)

type CurrencyService interface {
	Validate(currency *entity.Currency) error
	Upsert(currency *entity.Currency) error
	FindById(id string) (*entity.Currency, error)
	FindAll() []interface{}
	Remove(id string) (int64, error)
}

type service struct {}

var (
	currencyRepo repository.DBRepository = mongodb.NewMongodbRepository("daxxer", "currencies")
)

func NewCurrencyService() CurrencyService {
	return &service{}
}

func (s *service) Validate(currency *entity.Currency) error {
	if currency.Id == "" {
		err := errors.New("Currency.Id must not be empty")
		return err
	}
	if currency.Name == "" {
		err := errors.New("Currency.Name must not be empty")
		return err
	}
	return nil
}

func (s *service) Upsert(currency *entity.Currency) error {
	selector := bson.M{"id": currency.Id}
	update := bson.M{
		"$set": currency,
	}
	err := currencyRepo.Upsert(selector, update)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) FindById(id string) (*entity.Currency, error) {
	var currency = entity.Currency{}
	var query = `{"id": "`+id+`"}`
	err := currencyRepo.FindOne(query, &currency)
	if err != nil {
		return nil, err
	}
	return &currency, err
}

func (s *service) FindAll() []interface{} {
	var query = bson.M{}
	currencies := currencyRepo.FindAll(0, 100, 1, query, new(entity.Currency))
	return currencies
}

func (s *service) Remove(id string)(int64, error)  {
	count, err := currencyRepo.DeleteOne("id", id)
	return count, err
}