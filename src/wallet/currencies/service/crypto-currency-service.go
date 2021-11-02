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

type CryptoService interface {
	Validate(ctx context.Context, currency *entity.CryptoCurrency) error
	Upsert(ctx context.Context, currency *entity.CryptoCurrency) error
	FindById(ctx context.Context, id string) (*entity.CryptoCurrency, error)
	FindAll(ctx context.Context, ) []interface{}
	Remove(ctx context.Context, id string) (int64, error)
}

type cryptoService struct {
	cryptoRepo repository.DBRepository
}

func NewCryptoCurrencyService(cryptoRepo repository.DBRepository) CryptoService {
	return &cryptoService{cryptoRepo}
}

func (s *cryptoService) Validate(ctx context.Context, crypto *entity.CryptoCurrency) error {
	if crypto.Id == "" {
		logs.Instance.Log.Warn(ctx,"Crypto.Id must not be empty")
		err := errors.New("Crypto.Id must not be empty")
		return err
	}
	if crypto.Symbol == "" {
		logs.Instance.Log.Warn(ctx,"Crypto.Symbol must not be empty")
		err := errors.New("Crypto.Symbol must not be empty")
		return err
	}
	return nil
}

func (s *cryptoService) Upsert(ctx context.Context, crypto *entity.CryptoCurrency) error {
	selector := bson.M{"symbol": crypto.Symbol}
	update := bson.M{
		"$set": crypto,
	}
	err := s.cryptoRepo.Upsert(ctx, selector, update)
	if err != nil {
		return err
	}
	return nil
}

func (s *cryptoService) FindById(ctx context.Context, id string) (*entity.CryptoCurrency, error) {
	var crypto = entity.CryptoCurrency{}
	var query = utils2.QueryUtil().Build("id", id)

	err := s.cryptoRepo.FindOne(ctx,query, &crypto)
	if err != nil {
		return nil, err
	}
	return &crypto, err
}

func (s *cryptoService) FindAll(ctx context.Context, ) []interface{} {
	var query = bson.M{}
	cryptoCurrencies := s.cryptoRepo.FindAll(ctx,0, 100, 1, query, new(entity.CryptoCurrency))
	return cryptoCurrencies
}

func (s *cryptoService) Remove(ctx context.Context, id string)(int64, error)  {
	count, err := s.cryptoRepo.DeleteOne(ctx,"id", id)
	return count, err
}