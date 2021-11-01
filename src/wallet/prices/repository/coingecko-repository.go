package repository

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/julioisaac/daxxer-api/internal/logs"
	"github.com/julioisaac/daxxer-api/src/wallet/prices/entity"
	"github.com/julioisaac/daxxer-api/src/wallet/prices/utils"
	errors2 "github.com/pkg/errors"
	"go.uber.org/zap"
	"io"
	"net/http"
)

var (
	util = utils.Util()
)

type coinGeckoRepo struct {
	Url	string
	Client HttpClient
}

type GeckoSimplePrice map[string]map[string]float64

func NewCoinGeckoApiRepo(url string, client HttpClient) ApiRepository {
	return &coinGeckoRepo{
		Url: url,
		Client: client,
	}
}

func (r *coinGeckoRepo) GetPrices(cryptoCurrencies, currencies *[]interface{}) (*[]entity.Price, error) {

	cryptoIds := util.ExtractAndJoinByField(cryptoCurrencies, "Id", ",")
	vsCurrencies := util.ExtractAndJoinByField(currencies, "Id", ",")

	if cryptoIds == "" || vsCurrencies == "" {
		logs.Instance.Log.Warn(context.Background(), "error trying to extract currencies")
		return nil, errors.New("error trying to extract currencies")
	}

	params := map[string]string{
		"ids": cryptoIds,
		"currencies": vsCurrencies,
	}
	resp, err := r.DoRequest(params)
	if err != nil {
		return nil, err
	}

	geckoPrices, err1 := buildCoinGeckoPriceByBody(resp.Body)
	if err1 != nil {
		logs.Instance.Log.Error(context.Background(), "error trying to build CoinGeckoPrice by Body", zap.Error(err1))
		return nil, err1
	}
	var prices []entity.Price
	for id, currencies := range *geckoPrices {
		price := util.BuildPrice(id, currencies, "CoinGecko")
		prices = append(prices, price)
	}

	return &prices, err
}

func (r *coinGeckoRepo) DoRequest(params map[string]string) (*http.Response, error) {
	req, err := http.NewRequest("GET", r.Url, nil)
	if err != nil {
		logs.Instance.Log.Error(context.Background(), "error trying to create new request "+r.Url, zap.Error(err))
		return nil, errors2.Wrap(err, "error trying to create new request "+r.Url)
	}

	q := req.URL.Query()
	q.Add("ids", params["ids"])
	q.Add("vs_currencies", params["currencies"])
	req.URL.RawQuery = q.Encode()

	resp, err1 := r.Client.Do(req)
	if err1 != nil {
		logs.Instance.Log.Error(context.Background(), "error trying to do request "+r.Url, zap.Error(err1))
		return nil, errors2.Wrap(err1, "error trying to do request "+r.Url)
	}
	return resp, nil
}

func buildCoinGeckoPriceByBody(body io.ReadCloser) (*GeckoSimplePrice, error) {
	geckoPrices := GeckoSimplePrice{}
	b, _ := io.ReadAll(body)
	err := json.Unmarshal(b, &geckoPrices)
	if err != nil {
		return nil, errors2.Wrap(err, "error trying to bind coinGeckoPrice from body")
	}
	return &geckoPrices, nil
}