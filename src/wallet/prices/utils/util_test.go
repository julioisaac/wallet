package utils

import (
	"github.com/julioisaac/daxxer-api/src/wallet/account/entity"
	currencies2 "github.com/julioisaac/daxxer-api/src/wallet/currencies/entity"
	"github.com/stretchr/testify/assert"
	"testing"
)

type inputs struct {
	sources *[]interface{}
	field string
	sep string
}

func TestExtractAndJoinByField(t *testing.T) {

	var cryptoCurrencies []interface{}
	cryptoCurrencies = append(cryptoCurrencies, currencies2.CryptoCurrency{ Id: "bitcoin", Symbol: "btc"})
	cryptoCurrencies = append(cryptoCurrencies, currencies2.CryptoCurrency{ Id: "ethereum", Symbol: "eth"})
	var currencies []interface{}
	currencies = append(currencies, currencies2.Currency{ Id: "usd", Name: "dollar"})
	currencies = append(currencies, currencies2.Currency{ Id: "eur", Name: "euro"})
	var amounts []interface{}
	amounts = append(amounts, entity.Amount{ Id: "bitcoin", Currency: "btc", Value: 0.6 })
	amounts = append(amounts, entity.Amount{ Id: "ethereum", Currency: "eth", Value: 0.9 })

	tests := []struct {
		name string
		given inputs
		expected string
	}{
		{"Extract 'Symbol' value and join with Sep '-'", inputs{ &cryptoCurrencies, "Symbol", "-"}, "btc-eth"},
		{"Extract 'Name' value and join with Sep '@'", inputs{ &currencies, "Name", "@"}, "dollar@euro"},
		{"Extract 'Id' value and join with Sep ','", inputs{ &amounts, "id", ","}, "bitcoin,ethereum"},
	}

	asserts := assert.New(t)
	for _, test := range tests {
		asserts.Equal(test.expected, Util().ExtractAndJoinByField(test.given.sources, test.given.field, test.given.sep), test.name)
	}
}
