package service

import (
	"errors"
	"github.com/julioisaac/daxxer-api/src/helpers/repository"
	"github.com/julioisaac/daxxer-api/src/wallet/currencies/entity"
	utils2 "github.com/julioisaac/daxxer-api/src/wallet/utils"
	"go.mongodb.org/mongo-driver/bson"
)

type CurrencyService interface {
	Validate(currency *entity.Currency) error
	Upsert(currency *entity.Currency) error
	FindById(id string) (*entity.Currency, error)
	FindAll() []interface{}
	Remove(id string) (int64, error)
}

type currencyService struct {
	currencyRepo repository.DBRepository
}

func NewCurrencyService(currencyRepo repository.DBRepository) CurrencyService {
	return &currencyService{currencyRepo}
}

func (s *currencyService) Validate(currency *entity.Currency) error {
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

func (s *currencyService) Upsert(currency *entity.Currency) error {
	selector := bson.M{"id": currency.Id}
	update := bson.M{
		"$set": currency,
	}
	err := s.currencyRepo.Upsert(selector, update)
	if err != nil {
		return err
	}
	return nil
}

func (s *currencyService) FindById(id string) (*entity.Currency, error) {
	var currency = entity.Currency{}
	var query = utils2.QueryUtil().Build("id", id)

	err := s.currencyRepo.FindOne(query, &currency)
	if err != nil {
		return nil, err
	}
	return &currency, err
}

func (s *currencyService) FindAll() []interface{} {
	var query = bson.M{}
	currencies := s.currencyRepo.FindAll(0, 100, 1, query, new(entity.Currency))
	return currencies
}

func (s *currencyService) Remove(id string)(int64, error)  {
	count, err := s.currencyRepo.DeleteOne("id", id)
	return count, err
}