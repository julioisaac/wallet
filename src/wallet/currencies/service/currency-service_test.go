package service

import (
	"github.com/julioisaac/daxxer-api/src/helpers/repository"
	entity2 "github.com/julioisaac/daxxer-api/src/wallet/currencies/entity"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

var (
	currencyRepo repository.MockDBRepository
)

type CurrencyServiceTestSuite struct {
	suite.Suite
	currencyService CurrencyService
}

func (suite *CurrencyServiceTestSuite) SetupTest() {
	currencyRepo = repository.MockDBRepository{}
	suite.currencyService = NewCurrencyService(&currencyRepo)
}

func (suite *CurrencyServiceTestSuite) TestCurrencyCurrencyWhenUpsert() {
	//given
	incomingCurrency := entity2.Currency{ Id: "usd", Name: "dollar"}

	//when
	currencyRepo.On("Upsert", mock.Anything, mock.Anything).Return(nil)

	//expected
	err := suite.currencyService.Upsert(&incomingCurrency)
	currencyRepo.AssertNumberOfCalls(suite.T(), "Upsert", 1)
	suite.Nil(err)
}

func (suite *CurrencyServiceTestSuite) TestCurrencyWhenFindById() {
	//given
	incomingCurrencyId := "usd"
	expectedCurrency := &entity2.Currency{ Id: "usd", Name: "dollar"}

	//when
	currencyRepo.On("FindOne", `{"id": "usd"}`, &entity2.Currency{}).Return(nil).Run(func(args mock.Arguments) {
		currency := args.Get(1).(*entity2.Currency)
		currency.Id = "usd"
		currency.Name = "dollar"
	})

	//expected
	currency, _ := suite.currencyService.FindById(incomingCurrencyId)
	currencyRepo.AssertNumberOfCalls(suite.T(), "FindOne", 1)
	suite.Equal(expectedCurrency, currency)
}

func (suite *CurrencyServiceTestSuite) TestCurrencyWhenFindAll() {
	//given
	var expectedCurrencies []interface{}
	expectedCurrencies = append(expectedCurrencies, &entity2.Currency{ Id: "usd", Name: "dollar"})
	expectedCurrencies = append(expectedCurrencies, &entity2.Currency{ Id: "eur", Name: "euro"})

	//when
	currencyRepo.On("FindAll", mock.Anything, mock.Anything, mock.Anything, mock.Anything, new(entity2.Currency)).Return(expectedCurrencies)

	//expected
	currencies := suite.currencyService.FindAll()
	currencyRepo.AssertNumberOfCalls(suite.T(), "FindAll", 1)
	suite.Equal(expectedCurrencies, currencies)
}

func (suite *CurrencyServiceTestSuite) TestCurrencyWhenRemove() {
	//given
	incomingCurrencyId := "usd"

	//when
	currencyRepo.On("DeleteOne", "id", "usd").Return(int64(1), nil)

	//expected
	count, _ := suite.currencyService.Remove(incomingCurrencyId)
	currencyRepo.AssertNumberOfCalls(suite.T(), "DeleteOne", 1)
	suite.Equal(int64(1), count)
}

func TestCurrencyServiceTestSuite(t *testing.T) {
	suite.Run(t, new(CurrencyServiceTestSuite))
}