package repository

import (
	"encoding/json"
	currencies "github.com/julioisaac/daxxer-api/src/wallet/currencies/entity"
	"github.com/julioisaac/daxxer-api/src/wallet/prices/entity"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type coinBaseRepo struct {
	Url	string
	Client http.Client
}

type CoinBasePrice struct {
	Data struct {
		Currency string
		Rates map[string]string
	}
}

func NewCoinBaseApiRepo(url string) ApiRepository {
	return &coinBaseRepo{
		Url: url,
		Client: http.Client {
			Timeout: time.Duration(5 * time.Second),
		},
	}
}

func (c *coinBaseRepo) GetPrices(cryptoCurrencies, currencies []interface{}) (*[]entity.Price, error) {

	var prices []entity.Price
	var currenciesIds = util.ExtractAndJoinByField(currencies, "Id", ",")
	var cryptos = buildCryptoCurrencies(cryptoCurrencies)

	for _, crypto := range cryptos {
		params := map[string]string{
			"symbol": crypto.Symbol,
		}

		resp, err := c.DoRequest(params)
		if err != nil {
			log.Fatalln(err)
		}

		coinBasePrices, err1 := buildCoinBasePriceByBody(resp.Body)
		if err1 != nil {
			log.Fatalln(err1)
		}

		currenciesPrices, err2 := extractCoinBasePrices(coinBasePrices, currenciesIds)
		if err2 != nil {
			log.Fatalln(err2)
		}
		price := util.BuildPrice(crypto.Id, currenciesPrices, "CoinBase")
		prices = append(prices, price)
	}

	return &prices, nil
}

func (c *coinBaseRepo) DoRequest(params map[string]string) (*http.Response, error) {
	req, err := http.NewRequest("GET", c.Url, nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Add("currency", params["symbol"])
	req.URL.RawQuery = q.Encode()

	resp, err1 := c.Client.Do(req)
	return resp, err1
}

func buildCryptoCurrencies(cryptoCurrencies []interface{}) []currencies.CryptoCurrency {
	var cryptos []currencies.CryptoCurrency
	for _, cryptoCurrency := range cryptoCurrencies {
		c := cryptoCurrency.(*currencies.CryptoCurrency)
		cryptos = append(cryptos, *c)
	}
	return cryptos
}

func buildCoinBasePriceByBody(body io.ReadCloser) (*CoinBasePrice, error) {
	coinBasePrices := CoinBasePrice{}
	b, _ := io.ReadAll(body)
	err := json.Unmarshal(b, &coinBasePrices)
	return &coinBasePrices, err
}

func extractCoinBasePrices(coinBasePrices *CoinBasePrice, currenciesIds string) (map[string]float64, error) {
	currenciesPrices := make(map[string]float64)
	for _, id := range strings.Split(currenciesIds, ",") {
		rate := coinBasePrices.Data.Rates[strings.ToUpper(id)]
		value, err := strconv.ParseFloat(rate, 64)
		if err != nil {
			return nil, err
		}
		currenciesPrices[id] = value
	}
	return currenciesPrices, nil
}