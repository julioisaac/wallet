package service

import (
	"errors"
	"github.com/julioisaac/daxxer-api/src/helpers/repository"
	"github.com/julioisaac/daxxer-api/src/wallet/account/entity"
	entity2 "github.com/julioisaac/daxxer-api/src/wallet/prices/entity"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

var (
	mockAccountRepo repository.MockDBRepository
	mockHistoryRepo repository.MockDBRepository
	mockPricesRepo repository.MockDBRepository
)

type AccountServiceTestSuite struct {
	suite.Suite
	accountService AccountService
}

func (suite *AccountServiceTestSuite) SetupTest() {
	mockAccountRepo = repository.MockDBRepository{}
	mockHistoryRepo = repository.MockDBRepository{}
	mockPricesRepo = repository.MockDBRepository{}
	suite.accountService = NewAccountService(&mockAccountRepo, &mockHistoryRepo, &mockPricesRepo)
}

func (suite *AccountServiceTestSuite) TestAccountExistsWhenCreate() {
	//given
	incomingAccount := entity.Account{UserName: "julio"}

	//when
	mockAccountRepo.On("FindOne", `{"username": "julio"}`, &entity.Account{}).Return(nil).Run(func(args mock.Arguments) {
		account := args.Get(1).(*entity.Account)
		account.UserName = "julio"
	})

	//expected
	err := suite.accountService.Create(&incomingAccount)
	mockAccountRepo.AssertNumberOfCalls(suite.T(), "FindOne", 1)
	suite.Error(err, "the account already exists")
}


func (suite *AccountServiceTestSuite) TestSuccessWhenCreate() {
	//given
	incomingAccount := entity.Account{UserName: "gabi"}

	//when
	mockAccountRepo.On("FindOne", `{"username": "gabi"}`, &entity.Account{}).Return(nil)
	mockAccountRepo.On("Insert",  &entity.Account{UserName: "gabi"}).Return(nil)

	//expected
	err := suite.accountService.Create(&incomingAccount)
	mockAccountRepo.AssertNumberOfCalls(suite.T(), "FindOne", 1)
	mockAccountRepo.AssertNumberOfCalls(suite.T(), "Insert", 1)
	suite.Nil(err)
}

func (suite *AccountServiceTestSuite) TestAccountNotFoundWhenDeposit() {
	//given
	incomingTransaction := entity.Transaction{UserName: "gabi"}

	//when
	mockAccountRepo.On("FindOne", `{"username": "gabi"}`, &entity.Account{}).Return(nil)

	//expected
	err := suite.accountService.Deposit(&incomingTransaction)
	mockAccountRepo.AssertNumberOfCalls(suite.T(), "FindOne", 1)
	suite.Error(err, "account not found")
}

func (suite *AccountServiceTestSuite) TestSuccessWhenDeposit() {
	//given
	incomingTransaction := &entity.Transaction{
		UserName: "gabi",
		Amount: entity.Amount{
			Id:       "ethereum",
			Currency: "eth",
			Value:    0.7,
		},
	}

	//when
	mockAccountRepo.On("FindOne", `{"username": "gabi"}`, &entity.Account{}).Return(nil)
	mockAccountRepo.On("Upsert", mock.Anything, mock.Anything).Return(nil)
	mockHistoryRepo.On("Insert", mock.Anything).Return(nil)

	//expected
	err := suite.accountService.Deposit(incomingTransaction)
	mockAccountRepo.AssertNumberOfCalls(suite.T(), "Upsert", 1)
	mockHistoryRepo.AssertNumberOfCalls(suite.T(), "Insert", 1)
	suite.Nil(err)
}

func (suite *AccountServiceTestSuite) TestAccountNotFoundWhenWithdraw() {
	//given
	incomingTransaction := entity.Transaction{UserName: "ravi"}

	//when
	mockAccountRepo.On("FindOne", `{"username": "ravi"}`, &entity.Account{}).Return(errors.New("account not found"))

	//expected
	err := suite.accountService.Withdraw(&incomingTransaction)
	mockAccountRepo.AssertNumberOfCalls(suite.T(), "FindOne", 1)
	suite.Error(err, "account not found")
}

func (suite *AccountServiceTestSuite) TestSuccessWhenWithdraw() {
	//given
	incomingTransaction := &entity.Transaction{
		UserName: "ravi",
		Amount: entity.Amount{
			Id:       "bitcoin",
			Currency: "btc",
			Value:    0.5,
		},
	}

	//when
	mockAccountRepo.On("FindOne", `{"username": "ravi"}`, &entity.Account{}).Return(nil).Run(func(args mock.Arguments) {
		account := args.Get(1).(*entity.Account)
		account.UserName = "julio"
		account.Amounts = []entity.Amount{
			{
				Id:       "bitcoin",
				Currency: "btc",
				Value:    0.6,
			},
		}
	})
	mockAccountRepo.On("Upsert", mock.Anything, mock.Anything).Return(nil)
	mockHistoryRepo.On("Insert", mock.Anything).Return(nil)

	//expected
	err := suite.accountService.Withdraw(incomingTransaction)
	mockAccountRepo.AssertNumberOfCalls(suite.T(), "Upsert", 1)
	mockHistoryRepo.AssertNumberOfCalls(suite.T(), "Insert", 1)
	suite.Nil(err)
}

func (suite *AccountServiceTestSuite) TestSuccessBalance() {
	//given
	userName := "ravi"
	expectedBalances := entity.Balances{
		User: "ravi",
		ByCrypto: []entity.BalanceByCrypto{
			{
				Currency: "btc",
				Amount: 0.6,
				Prices: map[string]float64{
					"eur": 5266, "usd": 61227.755,
				},
				TimeOfRate: time.Date(2020, time.October, 31, 6, 30, 01, 0, time.UTC),
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
	mockAccountRepo.On("FindOne", `{"username": "ravi"}`, &entity.Account{}).Return(nil).Run(func(args mock.Arguments) {
		account := args.Get(1).(*entity.Account)
		account.UserName = "ravi"
		account.Amounts = []entity.Amount{
			{
				Id:       "bitcoin",
				Currency: "btc",
				Value:    0.6,
			},
		}
	})
	mockPricesRepo.On("FindOne", `{"cryptocurrency": "bitcoin"}`, &entity2.Price{}).Return(nil).Run(func(args mock.Arguments) {
		price := args.Get(1).(*entity2.Price)
		price.CryptoCurrency = "bitcoin"
		price.Currencies = map[string]float64{
			"eur": 5266, "usd": 61227.755,
		}
		price.ExchangeDataBy = "CoinGecko"
		price.LastUpdate = time.Date(2020, time.October, 31, 6, 30, 01, 0, time.UTC)
	})
	mockAccountRepo.On("Upsert", mock.Anything, mock.Anything).Return(nil)
	mockHistoryRepo.On("Insert", mock.Anything).Return(nil)

	//expected
	balances := suite.accountService.Balance(userName)
	mockAccountRepo.AssertNumberOfCalls(suite.T(), "FindOne", 1)
	mockPricesRepo.AssertNumberOfCalls(suite.T(), "FindOne", 1)
	suite.Equal(&expectedBalances, balances)
}

func (suite *AccountServiceTestSuite) TestErrorBalance() {
	//given
	userName := "ravi"

	//when
	mockAccountRepo.On("FindOne", `{"username": "ravi"}`, &entity.Account{}).Return(errors.New("not found"))

	//expected
	balances := suite.accountService.Balance(userName)
	mockAccountRepo.AssertNumberOfCalls(suite.T(), "FindOne", 1)
	mockPricesRepo.AssertNumberOfCalls(suite.T(), "FindOne", 0)
	suite.Nil(balances)
}


func (suite *AccountServiceTestSuite) TestSuccessHistories() {
	//given
	userName := "ravi"
	page := 0
	limit := 100
	order := 1
	startDt := time.Date(2020, time.October, 30, 6, 30, 01, 0, time.UTC)
	endDt := time.Date(2020, time.October, 31, 6, 30, 01, 0, time.UTC)

	var expectedHistories []interface{}
	expectedHistories = append(expectedHistories, entity.History{
		UserName: "ravi",
		Type: "deposit",
		Amount: entity.Amount{
			Id:       "bitcoin",
			Currency: "btc",
			Value:    0.6,
		},
		EventTime: time.Date(2020, time.October, 31, 6, 30, 01, 0, time.UTC),
	})
	expectedHistories = append(expectedHistories, entity.History{
		UserName: "ravi",
		Type: "withdraw",
		Amount: entity.Amount{
			Id:       "bitcoin",
			Currency: "btc",
			Value:    0.3,
		},
		EventTime: time.Date(2020, time.October, 31, 8, 30, 01, 0, time.UTC),
	})

	//when
	mockHistoryRepo.On("FindAll", page, limit, order, mock.Anything, new(entity.History)).Return(expectedHistories)

	//expected
	histories := suite.accountService.Histories(userName, page, limit, order, startDt, endDt)
	mockHistoryRepo.AssertNumberOfCalls(suite.T(), "FindAll", 1)
	suite.NotNil(histories)
}

func (suite *AccountServiceTestSuite) TestEmptyHistories() {
	//given
	userName := "ravi"
	page := 0
	limit := 100
	order := 1
	startDt := time.Date(2020, time.October, 30, 6, 30, 01, 0, time.UTC)
	endDt := time.Date(2020, time.October, 31, 6, 30, 01, 0, time.UTC)

	var expectedHistories []interface{}

	//when
	mockHistoryRepo.On("FindAll", page, limit, order, mock.Anything, new(entity.History)).Return(expectedHistories)

	//expected
	histories := suite.accountService.Histories(userName, page, limit, order, startDt, endDt)
	mockHistoryRepo.AssertNumberOfCalls(suite.T(), "FindAll", 1)
	suite.Empty(histories)
}

func TestAccountServiceTestSuite(t *testing.T) {
	suite.Run(t, new(AccountServiceTestSuite))
}