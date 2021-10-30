package service

import (
	"errors"
	"github.com/julioisaac/daxxer-api/src/helpers/repository"
	"github.com/julioisaac/daxxer-api/src/helpers/utils"
	"github.com/julioisaac/daxxer-api/src/wallet/account/entity"
	entity2 "github.com/julioisaac/daxxer-api/src/wallet/prices/entity"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

type AccountService interface {
	Create(account *entity.Account) error
	Deposit(transaction *entity.Transaction) error
	Withdraw(transaction *entity.Transaction) error
	Histories(user string, page, limit int, sort int, startDate time.Time, endDate time.Time) []interface{}
	Balance(username string) *entity.Balances
}

type service struct {
	accountRepo repository.DBRepository
	historyRepo repository.DBRepository
	pricesRepo  repository.DBRepository
}

func NewAccountService(accountRepo, historyRepo, pricesRepo repository.DBRepository) AccountService  {
	return &service{accountRepo, historyRepo, pricesRepo}
}

func (s *service) Create(account *entity.Account) error {
	var storedAccount entity.Account
	var query = `{"username": "`+account.UserName+`"}`
	_ = s.accountRepo.FindOne(query, &storedAccount)
	if storedAccount.UserName != "" {
		err := errors.New("the account already exists")
		return err
	}
	return s.accountRepo.Insert(account)
}

func (s *service) Deposit(transaction *entity.Transaction) error {
	var account = entity.Account{}
	transaction.Type = "deposit"
	var query = `{"username": "`+transaction.UserName+`"}`
	err := s.accountRepo.FindOne(query, &account)
	if err != nil {
		return errors.New("account not found")
	}
	err1 := s.execTransaction(account.Deposit, transaction)
	if err1 != nil {
		return err1
	}
	return nil
}

func (s *service) Withdraw(transaction *entity.Transaction) error {
	var account = entity.Account{}
	transaction.Type = "withdraw"
	var query = `{"username": "`+transaction.UserName+`"}`
	err := s.accountRepo.FindOne(query, &account)
	if err != nil {
		return errors.New("account not found")
	}
	err1 := s.execTransaction(account.Withdraw, transaction)
	if err1 != nil {
		return err1
	}

	return nil
}

func (s *service) Histories(user string, page, limit int, sort int, startDate time.Time, endDate time.Time) []interface{} {
	var endDt = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, endDate.Nanosecond(), endDate.Location())
	var query = bson.M{ "eventtime": bson.M{ "$gte": startDate, "$lte": endDt }, "username": user}
	var histories = s.historyRepo.FindAll(page, limit, sort, query, new(entity.History))
	return histories
}

func (s *service) Balance(username string) *entity.Balances {
	var account = entity.Account{}
	var query = `{"username": "`+username+`"}`
	err := s.accountRepo.FindOne(query, &account)
	if err != nil {
		return nil
	}

	var byCryptos []entity.BalanceByCrypto
	for _, amount := range account.Amounts {
		var price = entity2.Price{}
		var queryPrice = `{"cryptocurrency": "`+amount.Id+`"}`
		s.pricesRepo.FindOne(queryPrice, &price)
		totalByCurrency := calcTotalByCurrency(price.Currencies, amount.Value)
		byCrypto := buildBalanceByCrypto(price, amount, totalByCurrency)
		byCryptos = append(byCryptos, byCrypto)
	}
	total := calcTotal(byCryptos)

	return buildBalances(username, byCryptos ,total)
}

func (s *service) execTransaction(operation func(entity.Amount) (*entity.Account, error), transaction *entity.Transaction) error {
	account, err := operation(transaction.Amount)
	if err != nil {
		return err
	}
	selector := bson.M{"username": transaction.UserName}
	update := bson.M{
		"$set": account,
	}
	err1 := s.accountRepo.Upsert(selector, update)
	if err1 != nil {
		return err1
	}
	err2 := s.saveHistory(transaction)
	if err2 != nil {
		return err2
	}
	return nil
}

func (s *service) saveHistory(transaction *entity.Transaction) error {
	err := s.historyRepo.Insert(buildHistory(transaction))
	if err != nil {
		return err
	}
	return nil
}

func buildHistory(transaction *entity.Transaction) entity.History {
	return entity.History{
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
		Currency: amount.Currency,
		Amount: amount.Value,
		ExchangeDataBy: price.ExchangeDataBy,
		Prices: price.Currencies,
		TimeOfRate: price.LastUpdate,
		TotalByCurrency: totalByCurrency,
	}
}

func buildBalances(user string, byCryptos []entity.BalanceByCrypto, total map[string]float64) *entity.Balances{
	return &entity.Balances{
		User: user,
		ByCrypto: byCryptos,
		Total: total,
	}

}