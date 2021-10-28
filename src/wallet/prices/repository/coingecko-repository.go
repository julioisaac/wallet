package repository

import (
	"encoding/json"
	"github.com/julioisaac/daxxer-api/src/wallet/prices/entity"
	"github.com/julioisaac/daxxer-api/src/wallet/prices/utils"
	"io"
	"log"
	"net/http"
	"time"
)

var (
	util = utils.Util()
)

type repo struct {
	Url	string
	Client http.Client
}

type GeckoSimplePrice map[string]map[string]float64

func NewCoinGeckoApiRepo(url string) ApiRepository {
	return &repo{
		Url: url,
		Client: http.Client {
			Timeout: time.Duration(5 * time.Second),
		},
	}
}

func (r *repo) GetPrices(cryptoCurrencies, currencies []interface{}) (*[]entity.Price, error) {

	cryptoIds := util.ExtractAndJoinByField(cryptoCurrencies, "Id", ",")
	vsCurrencies := util.ExtractAndJoinByField(currencies, "Id", ",")

	params := map[string]string{
		"ids": cryptoIds,
		"currencies": vsCurrencies,
	}
	resp, err := r.DoRequest(params)
	if err != nil {
		log.Fatalln(err)
	}

	geckoPrices, err1 := buildCoinGeckoPriceByBody(resp.Body)
	if err1 != nil {
		return nil, err1
	}
	var prices []entity.Price
	for id, currencies := range *geckoPrices {
		price := util.BuildPrice(id, currencies, "CoinGecko")
		prices = append(prices, price)
	}

	return &prices, err
}

func (r *repo) DoRequest(params map[string]string) (*http.Response, error) {
	req, err := http.NewRequest("GET", r.Url, nil)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	q := req.URL.Query()
	q.Add("ids", params["ids"])
	q.Add("vs_currencies", params["currencies"])
	req.URL.RawQuery = q.Encode()
	resp, err2 := r.Client.Do(req)
	return resp, err2
}

func buildCoinGeckoPriceByBody(body io.ReadCloser) (*GeckoSimplePrice, error) {
	geckoPrices := GeckoSimplePrice{}
	b, _ := io.ReadAll(body)
	err := json.Unmarshal(b, &geckoPrices)
	return &geckoPrices, err
}