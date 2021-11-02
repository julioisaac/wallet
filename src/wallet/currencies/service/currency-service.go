package service

import (
	"context"
	"errors"
	"github.com/julioisaac/daxxer-api/internal/logs"
	"github.com/julioisaac/daxxer-api/src/helpers/repository"
	"github.com/julioisaac/daxxer-api/src/wallet/currencies/entity"
	utils2 "github.com/julioisaac/daxxer-api/src/wallet/utils"
	"go.mongodb.org/mongo-driver/bson"
)

type CurrencyService interface {
	Validate(ctx context.Context, currency *entity.Currency) error
	Upsert(ctx context.Context, currency *entity.Currency) error
	FindById(ctx context.Context, id string) (*entity.Currency, error)
	FindAll(ctx context.Context, ) []interface{}
	Remove(ctx context.Context, id string) (int64, error)
}

type currencyService struct {
	currencyRepo repository.DBRepository
}

func NewCurrencyService(currencyRepo repository.DBRepository) CurrencyService {
	return &currencyService{currencyRepo}
}

func (s *currencyService) Validate(ctx context.Context, currency *entity.Currency) error {
	if currency.Id == "" {
		logs.Instance.Log.Warn(context.Background(), "Currency.Id must not be empty")
		err := errors.New("Currency.Id must not be empty")
		return err
	}
	if currency.Name == "" {
		logs.Instance.Log.Warn(context.Background(), "Currency.Name must not be empty")
		err := errors.New("Currency.Name must not be empty")
		return err
	}
	return nil
}

func (s *currencyService) Upsert(ctx context.Context, currency *entity.Currency) error {
	selector := bson.M{"id": currency.Id}
	update := bson.M{
		"$set": currency,
	}
	err := s.currencyRepo.Upsert(ctx, selector, update)
	if err != nil {
		return err
	}
	return nil
}

func (s *currencyService) FindById(ctx context.Context, id string) (*entity.Currency, error) {
	var currency = entity.Currency{}
	var query = utils2.QueryUtil().Build("id", id)

	err := s.currencyRepo.FindOne(ctx, query, &currency)
	if err != nil {
		return nil, err
	}
	return &currency, err
}

func (s *currencyService) FindAll(ctx context.Context, ) []interface{} {
	var query = bson.M{}
	currencies := s.currencyRepo.FindAll(ctx,0, 100, 1, query, new(entity.Currency))
	return currencies
}

func (s *currencyService) Remove(ctx context.Context, id string)(int64, error)  {
	count, err := s.currencyRepo.DeleteOne(ctx,"id", id)
	return count, err
}