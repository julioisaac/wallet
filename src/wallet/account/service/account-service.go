package service

import (
	"context"
	"errors"
	"github.com/julioisaac/daxxer-api/internal/logs"
	"github.com/julioisaac/daxxer-api/src/helpers/repository"
	"github.com/julioisaac/daxxer-api/src/helpers/utils"
	"github.com/julioisaac/daxxer-api/src/wallet/account/entity"
	entity3 "github.com/julioisaac/daxxer-api/src/wallet/currencies/entity"
	entity2 "github.com/julioisaac/daxxer-api/src/wallet/prices/entity"
	utils2 "github.com/julioisaac/daxxer-api/src/wallet/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
	"time"
)

type AccountService interface {
	Create(ctx context.Context, account *entity.Account) error
	Deposit(ctx context.Context, transaction *entity.Transaction) error
	Withdraw(ctx context.Context, transaction *entity.Transaction) error
	Histories(ctx context.Context, user string, page, limit int, sort int, startDate time.Time, endDate time.Time) []interface{}
	Balance(ctx context.Context, username string) *entity.Balances
}

type accountService struct {
	accountRepo        repository.DBRepository
	cryptoCurrencyRepo repository.DBRepository
	historyRepo        repository.DBRepository
	pricesRepo         repository.DBRepository
}

func NewAccountService(accountRepo, cryptoCurrencyRepo, historyRepo, pricesRepo repository.DBRepository) AccountService {
	return &accountService{accountRepo, cryptoCurrencyRepo, historyRepo, pricesRepo}
}

func (s *accountService) Create(ctx context.Context, account *entity.Account) error {
	var storedAccount entity.Account
	var query = utils2.QueryUtil().Build("username", account.UserName)
	_ = s.accountRepo.FindOne(ctx, query, &storedAccount)
	if storedAccount.UserName != "" {
		logs.Instance.Log.Warn(ctx, "the account: '"+account.UserName+"' already exists")
		err := errors.New("the account already exists")
		return err
	}
	return s.accountRepo.Insert(ctx, account)
}

func (s *accountService) Deposit(ctx context.Context, transaction *entity.Transaction) error {
	var cryptoCurrency = entity3.CryptoCurrency{}
	var queryCrypto = utils2.QueryUtil().Build("Symbol", transaction.Amount.Currency)
	err := s.cryptoCurrencyRepo.FindOne(ctx, queryCrypto, &cryptoCurrency)
	if err != nil {
		logs.Instance.Log.Warn(ctx, transaction.Amount.Currency+" is not supported yet", zap.Error(err))
		return errors.New(transaction.Amount.Currency + " is not supported yet")
	}

	var account = entity.Account{}
	transaction.Type = "deposit"
	var query = utils2.QueryUtil().Build("username", transaction.UserName)
	err1 := s.accountRepo.FindOne(ctx, query, &account)
	if err1 != nil {
		logs.Instance.Log.Warn(ctx, "account: '"+transaction.UserName+"'  not found", zap.Error(err1))
		return errors.New("account not found")
	}
	err2 := s.execTransaction(ctx, account.Deposit, transaction)
	if err2 != nil {
		logs.Instance.Log.Error(ctx, "error trying to execute deposit to account"+transaction.UserName, zap.Error(err2))
		return err2
	}
	return nil
}

func (s *accountService) Withdraw(ctx context.Context, transaction *entity.Transaction) error {
	var account = entity.Account{}
	transaction.Type = "withdraw"
	var query = utils2.QueryUtil().Build("username", transaction.UserName)
	err := s.accountRepo.FindOne(ctx, query, &account)
	if err != nil {
		logs.Instance.Log.Warn(ctx, "account: '"+transaction.UserName+"'  not found", zap.Error(err))
		return errors.New("account not found")
	}
	err = s.execTransaction(ctx, account.Withdraw, transaction)
	if err != nil {
		logs.Instance.Log.Error(ctx, "error trying to execute withdraw from account"+transaction.UserName, zap.Error(err))
		return err
	}

	return nil
}

func (s *accountService) Histories(ctx context.Context, user string, page, limit int, sort int, startDate time.Time, endDate time.Time) []interface{} {
	var endDt = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, endDate.Nanosecond(), endDate.Location())
	var query = bson.M{"eventtime": bson.M{"$gte": startDate, "$lte": endDt}, "username": user}
	var histories = s.historyRepo.FindAll(ctx, page, limit, sort, query, new(entity.History))
	logs.Instance.Log.Debug(ctx, "returning histories from account: "+user)
	return histories
}

func (s *accountService) Balance(ctx context.Context, username string) *entity.Balances {
	var account = entity.Account{}
	var query = utils2.QueryUtil().Build("username", username)
	err := s.accountRepo.FindOne(ctx, query, &account)
	if err != nil {
		logs.Instance.Log.Warn(ctx, "account: '"+username+"'  not found", zap.Error(err))
		return nil
	}

	var byCryptos []entity.BalanceByCrypto
	for _, amount := range account.Amounts {
		var price = entity2.Price{}
		var queryPrice = utils2.QueryUtil().Build("cryptocurrency", amount.Id)
		s.pricesRepo.FindOne(ctx, queryPrice, &price)
		totalByCurrency := calcTotalByCurrency(price.Currencies, amount.Value)
		byCrypto := buildBalanceByCrypto(price, amount, totalByCurrency)
		byCryptos = append(byCryptos, byCrypto)
	}
	total := calcTotal(byCryptos)

	return buildBalances(username, byCryptos, total)
}

func (s *accountService) execTransaction(ctx context.Context, operation func(context.Context, entity.Amount) (*entity.Account, error), transaction *entity.Transaction) error {
	account, err := operation(ctx, transaction.Amount)
	if err != nil {
		return err
	}
	selector := bson.M{"username": transaction.UserName}
	update := bson.M{
		"$set": account,
	}
	err1 := s.accountRepo.Upsert(ctx, selector, update)
	if err1 != nil {
		return err1
	}
	err2 := s.saveHistory(ctx, transaction)
	if err2 != nil {
		return err2
	}
	return nil
}

func (s *accountService) saveHistory(ctx context.Context, transaction *entity.Transaction) error {
	err := s.historyRepo.Insert(ctx, buildHistory(transaction))
	if err != nil {
		logs.Instance.Log.Error(ctx, "error saving history to account: "+transaction.UserName, zap.Error(err))
		return err
	}
	return nil
}

func buildHistory(transaction *entity.Transaction) *entity.History {
	return &entity.History{
		UserName:  transaction.UserName,
		Type:      transaction.Type,
		Amount:    transaction.Amount,
		EventTime: time.Now(),
	}
}

func calcTotal(byCryptos []entity.BalanceByCrypto) map[string]float64 {
	var total = make(map[string]float64)
	for _, crypto := range byCryptos {
		for currency, value := range crypto.TotalByCurrency {
			if val, ok := total[currency]; ok {
				total[currency] = utils.DecimalMaths().Sum(value, val)
			} else {
				total[currency] = value
			}
		}
	}
	return total
}

func calcTotalByCurrency(currencies map[string]float64, amountValue float64) map[string]float64 {
	var totalByCurrency = make(map[string]float64)
	for currency, value := range currencies {
		totalByCurrency[currency] = utils.DecimalMaths().Mul(value, amountValue)
	}
	return totalByCurrency
}

func buildBalanceByCrypto(price entity2.Price, amount entity.Amount, totalByCurrency map[string]float64) entity.BalanceByCrypto {
	return entity.BalanceByCrypto{
		Currency:        amount.Currency,
		Amount:          amount.Value,
		ExchangeDataBy:  price.ExchangeDataBy,
		Prices:          price.Currencies,
		TimeOfRate:      price.LastUpdate,
		TotalByCurrency: totalByCurrency,
	}
}

func buildBalances(user string, byCryptos []entity.BalanceByCrypto, total map[string]float64) *entity.Balances {
	return &entity.Balances{
		User:     user,
		ByCrypto: byCryptos,
		Total:    total,
	}
}
