package service

import (
	"github.com/julioisaac/daxxer-api/src/helpers/repository"
	entity2 "github.com/julioisaac/daxxer-api/src/wallet/currencies/entity"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

var (
	cryptoRepo repository.MockDBRepository
)

type CryptoServiceTestSuite struct {
	suite.Suite
	cryptoService CryptoService
}

func (suite *CryptoServiceTestSuite) SetupTest() {
	cryptoRepo = repository.MockDBRepository{}
	suite.cryptoService = NewCryptoCurrencyService(&cryptoRepo)
}

func (suite *CryptoServiceTestSuite) TestCryptoCurrencyWhenUpsert() {
	//given
	incomingCryptoCurrency := entity2.CryptoCurrency{ Id: "bitcoin", Symbol: "btc"}

	//when
	cryptoRepo.On("Upsert", mock.Anything, mock.Anything).Return(nil)

	//expected
	err := suite.cryptoService.Upsert(&incomingCryptoCurrency)
	cryptoRepo.AssertNumberOfCalls(suite.T(), "Upsert", 1)
	suite.Nil(err)
}

func (suite *CryptoServiceTestSuite) TestCryptoCurrencyWhenFindById() {
	//given
	incomingCryptoCurrencyId := "bitcoin"
	expectedCryptoCurrency := &entity2.CryptoCurrency{ Id: "bitcoin", Symbol: "btc"}

	//when
	cryptoRepo.On("FindOne", `{"id": "bitcoin"}`, &entity2.CryptoCurrency{}).Return(nil).Run(func(args mock.Arguments) {
		cryptoCurrency := args.Get(1).(*entity2.CryptoCurrency)
		cryptoCurrency.Id = "bitcoin"
		cryptoCurrency.Symbol = "btc"
	})

	//expected
	cryptoCurrency, _ := suite.cryptoService.FindById(incomingCryptoCurrencyId)
	cryptoRepo.AssertNumberOfCalls(suite.T(), "FindOne", 1)
	suite.Equal(expectedCryptoCurrency, cryptoCurrency)
}

func (suite *CryptoServiceTestSuite) TestCryptoCurrencyWhenFindAll() {
	//given
	var expectedCryptoCurrencies []interface{}
	expectedCryptoCurrencies = append(expectedCryptoCurrencies, &entity2.CryptoCurrency{ Id: "bitcoin", Symbol: "btc"})
	expectedCryptoCurrencies = append(expectedCryptoCurrencies, &entity2.CryptoCurrency{ Id: "ethereum", Symbol: "eth"})

	//when
	cryptoRepo.On("FindAll", mock.Anything, mock.Anything, mock.Anything, mock.Anything, new(entity2.CryptoCurrency)).Return(expectedCryptoCurrencies)

	//expected
	cryptoCurrencies := suite.cryptoService.FindAll()
	cryptoRepo.AssertNumberOfCalls(suite.T(), "FindAll", 1)
	suite.Equal(expectedCryptoCurrencies, cryptoCurrencies)
}

func (suite *CryptoServiceTestSuite) TestCryptoCurrencyWhenRemove() {
	//mock DeleteOne
	incomingCryptoCurrencyId := "bitcoin"

	//when
	cryptoRepo.On("DeleteOne", "id", "bitcoin").Return(int64(1), nil)

	//expected
	count, _ := suite.cryptoService.Remove(incomingCryptoCurrencyId)
	cryptoRepo.AssertNumberOfCalls(suite.T(), "DeleteOne", 1)
	suite.Equal(int64(1), count)
}

func TestCryptoServiceTestSuite(t *testing.T) {
	suite.Run(t, new(CryptoServiceTestSuite))
}