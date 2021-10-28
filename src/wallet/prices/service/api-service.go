package service

import (
	"fmt"
	"github.com/julioisaac/daxxer-api/src/helpers/repository"
	"github.com/julioisaac/daxxer-api/src/helpers/repository/mongodb"
	"github.com/julioisaac/daxxer-api/src/wallet/currencies/entity"
	api "github.com/julioisaac/daxxer-api/src/wallet/prices/repository"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

var (
	pricesRepo   repository.DBRepository = mongodb.NewMongodbRepository("daxxer", "prices")
	cryptoRepo   repository.DBRepository = mongodb.NewMongodbRepository("daxxer", "cryptoCurrencies")
	currencyRepo repository.DBRepository = mongodb.NewMongodbRepository("daxxer", "currencies")
	apiRepo      api.ApiRepository       = api.NewCoinGeckoApiRepo("https://api.coingecko.com/api/v3/simple/price")
)

func NewApiService() *service {
	return &service{}
}

type service struct {}

func (s service) Update() {
	cryptoCurrencies := cryptoRepo.FindAll(0, 100, 1, bson.M{}, new(entity.CryptoCurrency))
	currencies := currencyRepo.FindAll(0, 100, 1, bson.M{}, new(entity.Currency))
	if cryptoCurrencies == nil || len(cryptoCurrencies) == 0 || currencies == nil || len(currencies) == 0 {
		fmt.Printf("no currencies to update\n")
		return
	}
	fmt.Printf("update prices %s\n", time.Now())
	prices, err := apiRepo.GetPrices(cryptoCurrencies, currencies)
	if err != nil {
		return
	}

	for _, price := range *prices {
		selector := bson.M{"cryptocurrency": price.CryptoCurrency}
		update := bson.M{
			"$set": price,
		}
		err2 := pricesRepo.Upsert(selector, update)
		if err2 != nil {
			return
		}
	}
}
