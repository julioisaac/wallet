package service

import (
	"context"
	"errors"
	"github.com/julioisaac/daxxer-api/internal/logs"
	"github.com/julioisaac/daxxer-api/src/helpers/repository"
	"github.com/julioisaac/daxxer-api/src/wallet/account/entity"
	entity3 "github.com/julioisaac/daxxer-api/src/wallet/currencies/entity"
	entity2 "github.com/julioisaac/daxxer-api/src/wallet/prices/entity"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

var (
	mockAccountRepo        repository.MockDBRepository
	mockCryptoCurrencyRepo repository.MockDBRepository
	mockHistoryRepo        repository.MockDBRepository
	mockPricesRepo         repository.MockDBRepository
)

type AccountServiceTestSuite struct {
	suite.Suite
	accountService AccountService
}

func (suite *AccountServiceTestSuite) SetupTest() {
	logs.NewZapLogger().Init()
	mockAccountRepo = repository.MockDBRepository{}
	mockCryptoCurrencyRepo = repository.MockDBRepository{}
	mockHistoryRepo = repository.MockDBRepository{}
	mockPricesRepo = repository.MockDBRepository{}
	suite.accountService = NewAccountService(&mockAccountRepo, &mockCryptoCurrencyRepo, &mockHistoryRepo, &mockPricesRepo)
}

func (suite *AccountServiceTestSuite) TestAccountExistsWhenCreate() {
	//given
	ctx := context.Background()
	incomingAccount := entity.Account{UserName: "julio"}

	//when
	mockAccountRepo.On("FindOne", ctx, `{"username": "julio"}`, &entity.Account{}).Return(nil).Run(func(args mock.Arguments) {
		account := args.Get(2).(*entity.Account)
		account.UserName = "julio"
	})

	//expected
	err := suite.accountService.Create(ctx, &incomingAccount)
	mockAccountRepo.AssertNumberOfCalls(suite.T(), "FindOne", 1)
	suite.Error(err, "the account already exists")
}

func (suite *AccountServiceTestSuite) TestSuccessWhenCreate() {
	//given
	ctx := context.TODO()
	incomingAccount := entity.Account{UserName: "gabi"}

	//when
	mockAccountRepo.On("FindOne", ctx, `{"username": "gabi"}`, &entity.Account{}).Return(nil)
	mockAccountRepo.On("Insert", ctx, &entity.Account{UserName: "gabi"}).Return(nil)

	//expected
	err := suite.accountService.Create(ctx, &incomingAccount)
	mockAccountRepo.AssertNumberOfCalls(suite.T(), "FindOne", 1)
	mockAccountRepo.AssertNumberOfCalls(suite.T(), "Insert", 1)
	suite.Nil(err)
}

func (suite *AccountServiceTestSuite) TestCryptoNotFoundWhenDeposit() {
	//given
	ctx := context.TODO()
	incomingTransaction := entity.Transaction{
		UserName: "gabi",
		Amount: entity.Amount{
			Id:       "cardano",
			Currency: "ada",
			Value:    0.7,
		},
	}

	//when
	mockCryptoCurrencyRepo.On("FindOne", ctx, `{"Symbol": "ada"}`, &entity3.CryptoCurrency{}).Return(errors.New("ada is not supported yet"))

	//expected
	err := suite.accountService.Deposit(ctx, &incomingTransaction)
	mockCryptoCurrencyRepo.AssertNumberOfCalls(suite.T(), "FindOne", 1)
	mockAccountRepo.AssertNumberOfCalls(suite.T(), "FindOne", 0)
	suite.Error(err, "ada is not supported yet")
}

func (suite *AccountServiceTestSuite) TestAccountNotFoundWhenDeposit() {
	//given
	ctx := context.TODO()
	incomingTransaction := entity.Transaction{
		UserName: "gabi",
		Amount: entity.Amount{
			Id:       "ethereum",
			Currency: "eth",
			Value:    0.7,
		},
	}

	//when
	mockCryptoCurrencyRepo.On("FindOne", ctx, `{"Symbol": "eth"}`, &entity3.CryptoCurrency{}).Return(nil).Run(func(args mock.Arguments) {
		cryptoCurrency := args.Get(2).(*entity3.CryptoCurrency)
		cryptoCurrency.Symbol = "eth"
		cryptoCurrency.Id = "ethereum"
	})
	mockAccountRepo.On("FindOne", ctx, `{"username": "gabi"}`, &entity.Account{}).Return(errors.New("account not found"))

	//expected
	err := suite.accountService.Deposit(ctx, &incomingTransaction)
	mockAccountRepo.AssertNumberOfCalls(suite.T(), "FindOne", 1)
	suite.Error(err, "account not found")
}

func (suite *AccountServiceTestSuite) TestSuccessWhenDeposit() {
	//given
	ctx := context.TODO()
	incomingTransaction := &entity.Transaction{
		UserName: "gabi",
		Amount: entity.Amount{
			Id:       "ethereum",
			Currency: "eth",
			Value:    0.7,
		},
	}

	//when
	mockCryptoCurrencyRepo.On("FindOne", ctx, `{"Symbol": "eth"}`, &entity3.CryptoCurrency{}).Return(nil).Run(func(args mock.Arguments) {
		cryptoCurrency := args.Get(2).(*entity3.CryptoCurrency)
		cryptoCurrency.Symbol = "eth"
		cryptoCurrency.Id = "ethereum"
	})
	mockAccountRepo.On("FindOne", ctx, `{"username": "gabi"}`, &entity.Account{}).Return(nil)
	mockAccountRepo.On("Upsert", ctx, mock.Anything, mock.Anything).Return(nil)
	mockHistoryRepo.On("Insert", ctx, mock.Anything).Return(nil)

	//expected
	err := suite.accountService.Deposit(ctx, incomingTransaction)
	mockAccountRepo.AssertNumberOfCalls(suite.T(), "Upsert", 1)
	mockHistoryRepo.AssertNumberOfCalls(suite.T(), "Insert", 1)
	suite.Nil(err)
}

func (suite *AccountServiceTestSuite) TestAccountNotFoundWhenWithdraw() {
	//given
	ctx := context.TODO()
	incomingTransaction := entity.Transaction{UserName: "ravi"}

	//when
	mockAccountRepo.On("FindOne", ctx, `{"username": "ravi"}`, &entity.Account{}).Return(errors.New("account not found"))

	//expected
	err := suite.accountService.Withdraw(ctx, &incomingTransaction)
	mockAccountRepo.AssertNumberOfCalls(suite.T(), "FindOne", 1)
	suite.Error(err, "account not found")
}

func (suite *AccountServiceTestSuite) TestSuccessWhenWithdraw() {
	//given
	ctx := context.TODO()
	incomingTransaction := &entity.Transaction{
		UserName: "ravi",
		Amount: entity.Amount{
			Id:       "bitcoin",
			Currency: "btc",
			Value:    0.5,
		},
	}

	//when
	mockAccountRepo.On("FindOne", ctx, `{"username": "ravi"}`, &entity.Account{}).Return(nil).Run(func(args mock.Arguments) {
		account := args.Get(2).(*entity.Account)
		account.UserName = "julio"
		account.Amounts = []entity.Amount{
			{
				Id:       "bitcoin",
				Currency: "btc",
				Value:    0.6,
			},
		}
	})
	mockAccountRepo.On("Upsert", ctx, mock.Anything, mock.Anything).Return(nil)
	mockHistoryRepo.On("Insert", ctx, mock.Anything).Return(nil)

	//expected
	err := suite.accountService.Withdraw(ctx, incomingTransaction)
	mockAccountRepo.AssertNumberOfCalls(suite.T(), "Upsert", 1)
	mockHistoryRepo.AssertNumberOfCalls(suite.T(), "Insert", 1)
	suite.Nil(err)
}

func (suite *AccountServiceTestSuite) TestSuccessBalance() {
	//given
	ctx := context.TODO()
	userName := "ravi"
	expectedBalances := entity.Balances{
		User: "ravi",
		ByCrypto: []entity.BalanceByCrypto{
			{
				Currency: "btc",
				Amount:   0.6,
				Prices: map[string]float64{
					"eur": 5266, "usd": 61227.755,
				},
				TimeOfRate:     time.Date(2020, time.October, 31, 6, 30, 01, 0, time.UTC),
				ExchangeDataBy: "CoinGecko",
				TotalByCurrency: map[string]float64{
					"eur": 3159.6, "usd": 36736.653,
				},
			},
		},
		Total: map[string]float64{
			"eur": 3159.6, "usd": 36736.653,
		},
	}

	//when
	mockAccountRepo.On("FindOne", ctx, `{"username": "ravi"}`, &entity.Account{}).Return(nil).Run(func(args mock.Arguments) {
		account := args.Get(2).(*entity.Account)
		account.UserName = "ravi"
		account.Amounts = []entity.Amount{
			{
				Id:       "bitcoin",
				Currency: "btc",
				Value:    0.6,
			},
		}
	})
	mockPricesRepo.On("FindOne", ctx, `{"cryptocurrency": "bitcoin"}`, &entity2.Price{}).Return(nil).Run(func(args mock.Arguments) {
		price := args.Get(2).(*entity2.Price)
		price.CryptoCurrency = "bitcoin"
		price.Currencies = map[string]float64{
			"eur": 5266, "usd": 61227.755,
		}
		price.ExchangeDataBy = "CoinGecko"
		price.LastUpdate = time.Date(2020, time.October, 31, 6, 30, 01, 0, time.UTC)
	})
	mockAccountRepo.On("Upsert", ctx, mock.Anything, mock.Anything).Return(nil)
	mockHistoryRepo.On("Insert", ctx, mock.Anything).Return(nil)

	//expected
	balances := suite.accountService.Balance(ctx, userName)
	mockAccountRepo.AssertNumberOfCalls(suite.T(), "FindOne", 1)
	mockPricesRepo.AssertNumberOfCalls(suite.T(), "FindOne", 1)
	suite.Equal(&expectedBalances, balances)
}

func (suite *AccountServiceTestSuite) TestErrorBalance() {
	//given
	ctx := context.TODO()
	userName := "ravi"

	//when
	mockAccountRepo.On("FindOne", ctx, `{"username": "ravi"}`, &entity.Account{}).Return(errors.New("not found"))

	//expected
	balances := suite.accountService.Balance(ctx, userName)
	mockAccountRepo.AssertNumberOfCalls(suite.T(), "FindOne", 1)
	mockPricesRepo.AssertNumberOfCalls(suite.T(), "FindOne", 0)
	suite.Nil(balances)
}

func (suite *AccountServiceTestSuite) TestSuccessHistories() {
	//given
	ctx := context.TODO()
	userName := "ravi"
	page := 0
	limit := 100
	order := 1
	startDt := time.Date(2020, time.October, 30, 6, 30, 01, 0, time.UTC)
	endDt := time.Date(2020, time.October, 31, 6, 30, 01, 0, time.UTC)

	var expectedHistories []interface{}
	expectedHistories = append(expectedHistories, entity.History{
		UserName: "ravi",
		Type:     "deposit",
		Amount: entity.Amount{
			Id:       "bitcoin",
			Currency: "btc",
			Value:    0.6,
		},
		EventTime: time.Date(2020, time.October, 31, 6, 30, 01, 0, time.UTC),
	})
	expectedHistories = append(expectedHistories, entity.History{
		UserName: "ravi",
		Type:     "withdraw",
		Amount: entity.Amount{
			Id:       "bitcoin",
			Currency: "btc",
			Value:    0.3,
		},
		EventTime: time.Date(2020, time.October, 31, 8, 30, 01, 0, time.UTC),
	})

	//when
	mockHistoryRepo.On("FindAll", ctx, page, limit, order, mock.Anything, new(entity.History)).Return(expectedHistories)

	//expected
	histories := suite.accountService.Histories(ctx, userName, page, limit, order, startDt, endDt)
	mockHistoryRepo.AssertNumberOfCalls(suite.T(), "FindAll", 1)
	suite.NotNil(histories)
}

func (suite *AccountServiceTestSuite) TestEmptyHistories() {
	//given
	ctx := context.TODO()
	userName := "ravi"
	page := 0
	limit := 100
	order := 1
	startDt := time.Date(2020, time.October, 30, 6, 30, 01, 0, time.UTC)
	endDt := time.Date(2020, time.October, 31, 6, 30, 01, 0, time.UTC)

	var expectedHistories []interface{}

	//when
	mockHistoryRepo.On("FindAll", ctx, page, limit, order, mock.Anything, new(entity.History)).Return(expectedHistories)

	//expected
	histories := suite.accountService.Histories(ctx, userName, page, limit, order, startDt, endDt)
	mockHistoryRepo.AssertNumberOfCalls(suite.T(), "FindAll", 1)
	suite.Empty(histories)
}

func TestAccountServiceTestSuite(t *testing.T) {
	suite.Run(t, new(AccountServiceTestSuite))
}
