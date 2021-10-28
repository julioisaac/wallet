package utils

import (
	"github.com/julioisaac/daxxer-api/src/wallet/prices/entity"
	"reflect"
	"strings"
	"time"
)

type util struct {}

func Util() *util {
	return &util{}
}

func (u *util) ExtractAndJoinByField(sources []interface{}, field string, sep string) string {
	var extracted []string
	for _, source := range sources {
		v := reflect.ValueOf(source).Elem()
		extracted = append(extracted, v.FieldByName(field).String())
	}
	return strings.Join(extracted,sep)
}

func (u *util) BuildPrice(id string, currenciesPrices map[string]float64, exchangeBy string) entity.Price {
	return entity.Price{
		CryptoCurrency: id,
		Currencies: currenciesPrices,
		LastUpdate: time.Now(),
		ExchangeDataBy: exchangeBy,
	}
}