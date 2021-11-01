package service

import (
	"context"
	"errors"
	"github.com/julioisaac/daxxer-api/internal/logs"
	"github.com/julioisaac/daxxer-api/src/helpers/repository"
	entity2 "github.com/julioisaac/daxxer-api/src/wallet/prices/entity"
	"go.mongodb.org/mongo-driver/bson"
)

type PricesService interface {
	GetAll() (*[]interface{}, error)
}

type pricesService struct {
	pricesRepo repository.DBRepository
}

func NewPricesService(pricesRepo repository.DBRepository) PricesService {
	return &pricesService{pricesRepo}
}

func (s *pricesService) GetAll() (*[]interface{}, error) {
	currencies := s.pricesRepo.FindAll(0, 100, 1, bson.M{}, new(entity2.Price))
	if currencies == nil {
		logs.Instance.Log.Warn(context.Background(), "there is no prices")
		return nil, errors.New("there is no prices")
	}
	return &currencies, nil
}