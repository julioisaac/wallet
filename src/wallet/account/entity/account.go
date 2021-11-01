package entity

import (
	"context"
	"errors"
	"fmt"
	"github.com/julioisaac/daxxer-api/internal/logs"
	"github.com/julioisaac/daxxer-api/src/helpers/utils"
	"math"
)

type Account struct {
	UserName string   `json:"user"`
	Amounts  []Amount `json:"amounts,omitempty"`
}

func (account *Account) Deposit(amount Amount) (*Account, error) {
	if math.Signbit(amount.Value) || amount.Value == 0 {
		logs.Instance.Log.Warn(context.Background(), "value cannot be zero or negative")
		return nil, errors.New("value cannot be zero or negative")
	}
	if len(account.Amounts) > 0 {
		if account.contains(amount.Currency) {
			for i, am := range account.Amounts {
				if am.Currency == amount.Currency {
					account.Amounts[i].Value = utils.DecimalMaths().Sum(account.Amounts[i].Value, amount.Value)
				}
			}
		} else {
			account.Amounts = append(account.Amounts, amount)
		}
	} else {
		account.Amounts = append(account.Amounts, amount)
	}

	logs.Instance.Log.Debug(context.Background(), "successfully deposited into account: "+account.UserName)
	return account, nil
}

func (account *Account) Withdraw(amount Amount) (*Account, error) {
	if math.Signbit(amount.Value) || amount.Value == 0 {
		logs.Instance.Log.Warn(context.Background(), "value cannot be zero or negative")
		return nil, errors.New("value cannot be zero or negative")
	}
	if len(account.Amounts) > 0 {
		if account.contains(amount.Currency) {
			for i, am := range account.Amounts {
				if am.Currency == amount.Currency {
					if amount.Value > account.Amounts[i].Value {
						logs.Instance.Log.Warn(context.Background(), "insufficient "+amount.Currency+" funds")
						err := fmt.Errorf("insufficient %s funds", amount.Currency)
						return nil, err
					}
					account.Amounts[i].Value = utils.DecimalMaths().Sub(account.Amounts[i].Value, amount.Value)
				}
			}
		} else {
			logs.Instance.Log.Warn(context.Background(), "there is "+amount.Currency+" amount")
			return nil, fmt.Errorf("there is no %s amount", amount.Currency)
		}
	} else {
		logs.Instance.Log.Warn(context.Background(), "there is amounts")
		return nil, errors.New("there is no amounts")
	}

	logs.Instance.Log.Debug(context.Background(), "successfully withdrawn from the account: "+account.UserName)
	return account, nil
}

func (account *Account) contains(currency string) bool {
	for _, amount := range account.Amounts {
		if amount.Currency == currency {
			return true
		}
	}
	return false
}