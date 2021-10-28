package service

import (
	"errors"
	"github.com/julioisaac/daxxer-api/src/helpers/repository"
	"github.com/julioisaac/daxxer-api/src/helpers/repository/mongodb"
	"github.com/julioisaac/daxxer-api/src/wallet/currencies/entity"
	"go.mongodb.org/mongo-driver/bson"
)

type CryptoService interface {
	Validate(currency *entity.CryptoCurrency) error
	Upsert(currency *entity.CryptoCurrency) error
	FindById(id string) (*entity.CryptoCurrency, error)
	FindAll() []interface{}
	Remove(id string) (int64, error)
}

type cryptoService struct {}

var (
	cryptoRepo repository.DBRepository = mongodb.NewMongodbRepository("daxxer", "cryptoCurrencies")
)

func NewCryptoCurrencyService() CryptoService {
	return &cryptoService{}
}

func (s *cryptoService) Validate(crypto *entity.CryptoCurrency) error {
	if crypto.Id == "" {
		err := errors.New("Crypto.Id must not be empty")
		return err
	}
	if crypto.Symbol == "" {
		err := errors.New("Crypto.Symbol must not be empty")
		return err
	}
	return nil
}

func (s *cryptoService) Upsert(crypto *entity.CryptoCurrency) error {
	selector := bson.M{"symbol": crypto.Symbol}
	update := bson.M{
		"$set": crypto,
	}
	err := cryptoRepo.Upsert(selector, update)
	if err != nil {
		return err
	}
	return nil
}

func (s *cryptoService) FindById(id string) (*entity.CryptoCurrency, error) {
	var crypto = entity.CryptoCurrency{}
	var query = `{"id": "`+id+`"}`
	err := cryptoRepo.FindOne(query, &crypto)
	if err != nil {
		return nil, err
	}
	return &crypto, err
}

func (s *cryptoService) FindAll() []interface{} {
	var query = bson.M{}
	cryptoCurrencies := cryptoRepo.FindAll(0, 100, 1, query, new(entity.CryptoCurrency))
	return cryptoCurrencies
}

func (s *cryptoService) Remove(id string)(int64, error)  {
	count, err := cryptoRepo.DeleteOne("id", id)
	return count, err
}