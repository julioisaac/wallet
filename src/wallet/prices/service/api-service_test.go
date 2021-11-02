package service

import (
	"bytes"
	"context"
	"errors"
	"github.com/julioisaac/daxxer-api/internal/logs"
	"github.com/julioisaac/daxxer-api/src/helpers/repository"
	entity2 "github.com/julioisaac/daxxer-api/src/wallet/currencies/entity"
	repository2 "github.com/julioisaac/daxxer-api/src/wallet/prices/repository"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"io/ioutil"
	"net/http"
	"testing"
)

var (
	mockCryptoRepo repository.MockDBRepository
	mockCurrencyRepo repository.MockDBRepository
	mockPricesRepo repository.MockDBRepository
	mockClient repository2.MockClient
	mockCoinGeckoResponse string
	mockCoinBaseBtcResponse string
	mockCoinBaseEthResponse string
	cryptoCurrencies, currencies []interface{}
)

type ApiServiceTestSuite struct {
	suite.Suite
	apiService ApiService
}

func (suite *ApiServiceTestSuite) SetupSuite() {
	logs.NewZapLogger().Init()
	responseCoinGecko, _ := ioutil.ReadFile("./mock/coingecko-response.json")
	responseCoinBaseBtc, _ := ioutil.ReadFile("./mock/coinbase-response-btc.json")
	responseCoinBaseEth, _ := ioutil.ReadFile("./mock/coinbase-response-eth.json")

	mockCoinGeckoResponse = string(responseCoinGecko)
	mockCoinBaseBtcResponse = string(responseCoinBaseBtc)
	mockCoinBaseEthResponse = string(responseCoinBaseEth)

	cryptoCurrencies = append(cryptoCurrencies, entity2.CryptoCurrency{ Id: "bitcoin", Symbol: "btc"})
	cryptoCurrencies = append(cryptoCurrencies, entity2.CryptoCurrency{ Id: "ethereum", Symbol: "eth"})
	currencies = append(currencies, entity2.Currency{ Id: "usd", Name: "dollar"})
	currencies = append(currencies, entity2.Currency{ Id: "eur", Name: "euro"})
}

func (suite *ApiServiceTestSuite) SetupTest() {
	mockCryptoRepo = repository.MockDBRepository{}
	mockCurrencyRepo = repository.MockDBRepository{}
	mockPricesRepo = repository.MockDBRepository{}
	mockClient = repository2.MockClient{}
}

func (suite *ApiServiceTestSuite) TestNoCurrenciesToUpdate() {
	//given
	ctx := context.TODO()
	var cryptoCurrenciesEmpty, currenciesEmpty []interface{}
	coinGeckoApiRepo := repository2.NewCoinGeckoApiRepo("http://test", &mockClient)
	suite.apiService = NewApiService(&mockCryptoRepo, &mockCurrencyRepo, &mockPricesRepo, coinGeckoApiRepo)

	//when
	mockCryptoRepo.On("FindAll", ctx, 0, 100, 1, bson.M{}, new(entity2.CryptoCurrency)).Return(cryptoCurrenciesEmpty)
	mockCurrencyRepo.On("FindAll", ctx, 0, 100, 1, bson.M{}, new(entity2.Currency)).Return(currenciesEmpty)

	//expected
	err := suite.apiService.Update(ctx)
	mockCryptoRepo.AssertNumberOfCalls(suite.T(), "FindAll", 1)
	mockCurrencyRepo.AssertNumberOfCalls(suite.T(), "FindAll", 1)
	suite.Error(err, "no currencies to update")
}

func (suite *ApiServiceTestSuite) TestErrorWhenGetPrices() {
	//given
	ctx := context.TODO()
	var cryptoCurrenciesError, currenciesError []interface{}
	cryptoCurrenciesError = append(cryptoCurrenciesError, entity2.Currency{})
	currenciesError = append(currenciesError, entity2.CryptoCurrency{})
	coinGeckoApiRepo := repository2.NewCoinGeckoApiRepo("http://test", &mockClient)
	suite.apiService = NewApiService(&mockCryptoRepo, &mockCurrencyRepo, &mockPricesRepo, coinGeckoApiRepo)

	//when
	mockCryptoRepo.On("FindAll", ctx, 0, 100, 1, bson.M{}, new(entity2.CryptoCurrency)).Return(cryptoCurrenciesError)
	mockCurrencyRepo.On("FindAll", ctx, 0, 100, 1, bson.M{}, new(entity2.Currency)).Return(currenciesError)
	mockClient.On("Do", mock.Anything).Return(&http.Response{}, errors.New("error trying to get prices"))

	//expected
	err := suite.apiService.Update(ctx)
	mockCryptoRepo.AssertNumberOfCalls(suite.T(), "FindAll", 1)
	mockCurrencyRepo.AssertNumberOfCalls(suite.T(), "FindAll", 1)
	suite.Error(err, "error trying to get prices")
}

func (suite *ApiServiceTestSuite) TestSuccessCoinGeckoToPrice() {
	//given
	ctx := context.TODO()
	coinGeckoApiRepo := repository2.NewCoinGeckoApiRepo("http://test", &mockClient)
	suite.apiService = NewApiService(&mockCryptoRepo, &mockCurrencyRepo, &mockPricesRepo, coinGeckoApiRepo)

	//when
	mockCryptoRepo.On("FindAll", ctx, 0, 100, 1, bson.M{}, new(entity2.CryptoCurrency)).Return(cryptoCurrencies)
	mockCurrencyRepo.On("FindAll", ctx, 0, 100, 1, bson.M{}, new(entity2.Currency)).Return(currencies)
	mockPricesRepo.On("Upsert", ctx, mock.Anything, mock.Anything).Return(nil)
	response := &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(mockCoinGeckoResponse))}
	mockClient.On("Do", mock.Anything).Return(response, nil)

	//expected
	err := suite.apiService.Update(ctx)
	mockClient.AssertNumberOfCalls(suite.T(), "Do", 1)
	mockPricesRepo.AssertNumberOfCalls(suite.T(), "Upsert", 2)
	suite.Nil(err)
}

func (suite *ApiServiceTestSuite) TestSuccessCoinBaseToPrice() {
	//given
	ctx := context.TODO()
	coinGeckoApiRepo := repository2.NewCoinBaseApiRepo("http://test", &mockClient)
	suite.apiService = NewApiService(&mockCryptoRepo, &mockCurrencyRepo, &mockPricesRepo, coinGeckoApiRepo)

	//when
	mockCryptoRepo.On("FindAll", ctx, 0, 100, 1, bson.M{}, new(entity2.CryptoCurrency)).Return(cryptoCurrencies)
	mockCurrencyRepo.On("FindAll", ctx, 0, 100, 1, bson.M{}, new(entity2.Currency)).Return(currencies)
	mockPricesRepo.On("Upsert", ctx, mock.Anything, mock.Anything).Return(nil)

	requestBtc, _ := http.NewRequest("GET", "http://test", nil)
	qBtc := requestBtc.URL.Query()
	qBtc.Add("currency", "btc")
	requestBtc.URL.RawQuery = qBtc.Encode()
	responseBtc := &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(mockCoinBaseBtcResponse))}
	mockClient.On("Do", requestBtc).Return(responseBtc, nil)

	requestEth, _ := http.NewRequest("GET", "http://test", nil)
	qEth := requestEth.URL.Query()
	qEth.Add("currency", "eth")
	requestEth.URL.RawQuery = qEth.Encode()
	responseEth := &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(mockCoinBaseEthResponse))}
	mockClient.On("Do", requestEth).Return(responseEth, nil)

	//expected
	err := suite.apiService.Update(ctx)
	mockClient.AssertNumberOfCalls(suite.T(), "Do", 2)
	mockPricesRepo.AssertNumberOfCalls(suite.T(), "Upsert", 2)
	suite.Nil(err)
}

func TestApiServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ApiServiceTestSuite))
}