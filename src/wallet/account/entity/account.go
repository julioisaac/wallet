package entity

import (
	"errors"
	"github.com/julioisaac/daxxer-api/src/helpers/utils"
	"math"
)

type Account struct {
	UserName string   `json:"user"`
	Amounts  []Amount `json:"amounts,omitempty"`
}

func (account *Account) Deposit(amount Amount) (*Account, error) {
	if math.Signbit(amount.Value) || amount.Value == 0 {
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