package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/julioisaac/daxxer-api/internal/logs"
	"github.com/julioisaac/daxxer-api/src/helpers/repository"
	"github.com/julioisaac/daxxer-api/src/wallet/currencies/entity"
	api "github.com/julioisaac/daxxer-api/src/wallet/prices/repository"
	errors2 "github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

type ApiService interface {
	Update(ctx context.Context) error
}

type apiService struct {
	cryptoRepo repository.DBRepository
	currencyRepo repository.DBRepository
	pricesRepo repository.DBRepository
	apiRepo  api.ApiRepository
}

func NewApiService(cryptoRepo, currencyRepo, pricesRepo repository.DBRepository, apiRepo api.ApiRepository) ApiService {
	return &apiService{cryptoRepo, currencyRepo, pricesRepo, apiRepo}
}

func (s *apiService) Update(ctx context.Context) error {
	cryptoCurrencies := s.cryptoRepo.FindAll(ctx,0, 100, 1, bson.M{}, new(entity.CryptoCurrency))
	currencies := s.currencyRepo.FindAll(ctx,0, 100, 1, bson.M{}, new(entity.Currency))
	if cryptoCurrencies == nil || len(cryptoCurrencies) == 0 || currencies == nil || len(currencies) == 0 {
		logs.Instance.Log.Info(ctx, "no currencies to update")
		return errors.New("no currencies to update")
	}
	fmt.Printf("update prices %s", time.Now())
	prices, err := s.apiRepo.GetPrices(&cryptoCurrencies, &currencies)
	if err != nil {
		return errors2.Wrap(err, "error trying to get prices")
	}

	for _, price := range *prices {
		selector := bson.M{"cryptocurrency": price.CryptoCurrency}
		update := bson.M{
			"$set": price,
		}
		err2 := s.pricesRepo.Upsert(ctx, selector, update)
		if err2 != nil {
			priceStr, _ := json.Marshal(price)
			logs.Instance.Log.Error(ctx, "error trying to upsert price "+string(priceStr)+" from "+price.ExchangeDataBy)
			return errors2.Wrap(err, "error trying to upsert price "+string(priceStr)+" from "+price.ExchangeDataBy)
		}
	}
	return nil
}
