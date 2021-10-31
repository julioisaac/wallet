package utils

import (
	"encoding/json"
	"github.com/julioisaac/daxxer-api/src/wallet/prices/entity"
	"log"
	"strings"
	"time"
)

type util struct {}

func Util() *util {
	return &util{}
}

func (u *util) ExtractAndJoinByField(sources *[]interface{}, field string, sep string) string {
	var extracted []string
	for _, source := range *sources {
		m := make(map[string]interface{})
		src, _ := json.Marshal(&source)
		err := json.Unmarshal(src, &m)
		if err != nil {
			log.Fatal(err)
		}
		extracted = append(extracted, m[field].(string))
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