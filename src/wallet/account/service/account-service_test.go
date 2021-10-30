package service

import (
	"github.com/julioisaac/daxxer-api/src/helpers/repository"
	"github.com/julioisaac/daxxer-api/src/wallet/account/entity"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
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

func (suite *AccountServiceTestSuite) SetupSuite() {
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
			"ethereum",
			"eth",
			0.7,
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
	mockAccountRepo.On("FindOne", `{"username": "ravi"}`, &entity.Account{}).Return(nil)

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
			"bitcoin",
			"btc",
			0.5,
		},
	}

	//when
	mockAccountRepo.On("FindOne", `{"username": "ravi"}`, &entity.Account{}).Return(nil).Run(func(args mock.Arguments) {
		account := args.Get(1).(*entity.Account)
		account.UserName = "julio"
		account.Amounts = []entity.Amount{
			{
				"bitcoin",
				"btc",
				0.6,
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

func TestAccountServiceTestSuite(t *testing.T) {
	suite.Run(t, new(AccountServiceTestSuite))
}