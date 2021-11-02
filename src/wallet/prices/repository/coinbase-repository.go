package repository

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/julioisaac/daxxer-api/internal/logs"
	currencies "github.com/julioisaac/daxxer-api/src/wallet/currencies/entity"
	"github.com/julioisaac/daxxer-api/src/wallet/prices/entity"
	errors2 "github.com/pkg/errors"
	"go.uber.org/zap"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type coinBaseRepo struct {
	Url	string
	Client HttpClient
}

type CoinBasePrice struct {
	Data struct {
		Currency string
		Rates map[string]string
	}
}

func NewCoinBaseApiRepo(client HttpClient) ApiRepository {
	return &coinBaseRepo{
		Url: os.Getenv("API_ENDPOINT"),
		Client: client,
	}
}

func (c *coinBaseRepo) GetPrices(cryptoCurrencies, currencies *[]interface{}) (*[]entity.Price, error) {

	var prices []entity.Price
	var currenciesIds = util.ExtractAndJoinByField(currencies, "Id", ",")
	var cryptos = buildCryptoCurrencies(cryptoCurrencies)

	if currenciesIds == "" || cryptos == nil {
		logs.Instance.Log.Warn(context.Background(), "error trying to extract currencies")
		return nil, errors.New("error trying to extract currencies")
	}

	for _, crypto := range cryptos {
		params := map[string]string{
			"symbol": crypto.Symbol,
		}

		resp, err := c.DoRequest(params)
		if err != nil {
			return nil, err
		}

		coinBasePrices, err1 := buildCoinBasePriceByBody(resp.Body)
		if err1 != nil {
			logs.Instance.Log.Error(context.Background(), "error trying to build CoinBasePrice by Body", zap.Error(err1))
			return nil, err1
		}

		currenciesPrices, err2 := extractCoinBasePrices(coinBasePrices, currenciesIds)
		if err2 != nil {
			return nil, err2
		}
		price := util.BuildPrice(crypto.Id, currenciesPrices, "CoinBase")
		prices = append(prices, price)
	}

	return &prices, nil
}

func (c *coinBaseRepo) DoRequest(params map[string]string) (*http.Response, error) {
	req, err := http.NewRequest("GET", c.Url, nil)
	if err != nil {
		logs.Instance.Log.Error(context.Background(), "error trying to create new request "+c.Url, zap.Error(err))
		return nil, errors2.Wrap(err, "error trying to create new request "+c.Url)
	}
	q := req.URL.Query()
	q.Add("currency", params["symbol"])
	req.URL.RawQuery = q.Encode()

	resp, err1 := c.Client.Do(req)
	if err1 != nil {
		logs.Instance.Log.Error(context.Background(), "error trying to do request "+c.Url, zap.Error(err1))
		return nil, errors2.Wrap(err1, "error trying to do request "+c.Url)
	}
	return resp, err1
}

func buildCryptoCurrencies(cryptoCurrencies *[]interface{}) []currencies.CryptoCurrency {
	var cryptos []currencies.CryptoCurrency
	for _, cryptoCurrency := range *cryptoCurrencies {
		c := cryptoCurrency.(currencies.CryptoCurrency)
		cryptos = append(cryptos, c)
	}
	return cryptos
}

func buildCoinBasePriceByBody(body io.ReadCloser) (*CoinBasePrice, error) {
	coinBasePrices := CoinBasePrice{}
	b, _ := io.ReadAll(body)
	err := json.Unmarshal(b, &coinBasePrices)
	return &coinBasePrices, errors2.Wrap(err, "error trying to bind coinBasePrices from body")
}

func extractCoinBasePrices(coinBasePrices *CoinBasePrice, currenciesIds string) (map[string]float64, error) {
	currenciesPrices := make(map[string]float64)
	for _, id := range strings.Split(currenciesIds, ",") {
		rate := coinBasePrices.Data.Rates[strings.ToUpper(id)]
		value, err := strconv.ParseFloat(rate, 64)
		if err != nil {
			logs.Instance.Log.Error(context.Background(), "error trying to parse rate from CoinBase", zap.Error(err))
			return nil, errors2.Wrap(err, "error trying to parse rate from CoinBase")
		}
		currenciesPrices[id] = value
	}
	return currenciesPrices, nil
}