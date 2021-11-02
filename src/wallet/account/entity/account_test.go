package entity

import (
	"context"
	"errors"
	"github.com/julioisaac/daxxer-api/internal/logs"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSuccessWhenDeposit(t *testing.T) {
	logs.NewZapLogger().Init()

	tests := []struct {
		name  string
		given *Account
		when  Amount
		expected *Account
		wantErr bool
	}{
		{
			name: "Expected balance equal 0.7",
			given: &Account{
				UserName: "julio",
				Amounts: []Amount{
					{
						"ethereum",
						"eth",
						0.5,
					},
				},
			},
			when: Amount{
				"ethereum",
				"eth",
				0.2,
			},
			expected: &Account{
				UserName: "julio",
				Amounts: []Amount{
					{
						"ethereum",
						"eth",
						0.7,
					},
				},
			},
		},
		{
			name: "Expected balance equal 1.01",
			given: &Account{
				UserName: "julio",
				Amounts: []Amount{
					{
						"ethereum",
						"eth",
						0.5,
					},
				},
			},
			when: Amount{
				"ethereum",
				"eth",
				0.51,
			},
			expected: &Account{
				UserName: "julio",
				Amounts: []Amount{
					{
						"ethereum",
						"eth",
						1.01,
					},
				},
			},
		},
		{
			name: "Expected added the new currency",
			given: &Account{
				UserName: "julio",
				Amounts: []Amount{
					{
						"bitcoin",
						"btc",
						0.2,
					},
				},
			},
			when: Amount{
				"ethereum",
				"eth",
				1.2,
			},
			expected: &Account{
				UserName: "julio",
				Amounts: []Amount{
					{
						"bitcoin",
						"btc",
						0.2,
					},
					{
						"ethereum",
						"eth",
						1.2,
					},
				},
			},
		},
		{
			name: "Expected amounts creation",
			given: &Account{
				UserName: "julio",
			},
			when: Amount{
				"ethereum",
				"eth",
				0.6,
			},
			expected: &Account{
				UserName: "julio",
				Amounts: []Amount{
					{
						"ethereum",
						"eth",
						0.6,
					},
				},
			},
		},
	}
	asserts := assert.New(t)
	for _, test := range tests {
		account, _ := test.given.Deposit(context.Background(), test.when)
		asserts.Equal(test.expected, account, test.name)
	}
}

func TestErrorWhenDeposit(t *testing.T) {

	test := struct {
		name  string
		given *Account
		when  Amount
		expected error
		wantErr bool
	}{
		name: "Expected error cannot be negative",
		given: &Account{
			UserName: "julio",
			Amounts: []Amount{
				{
					"ethereum",
					"eth",
					0.5,
				},
			},
		},
		when: Amount{
			"ethereum",
			"eth",
			-1,
		},
		expected: errors.New("value cannot be zero or negative"),
	}
	_, err := test.given.Deposit(context.Background(), test.when)
	assert.Equal(t, test.expected, err, test.name)
}

func TestSuccessWhenWithdraw(t *testing.T) {

	tests := []struct {
		name  string
		given *Account
		when  Amount
		expected *Account
		wantErr bool
	}{
		{
			name: "Expected balance equal 0.2",
			given: &Account{
				UserName: "julio",
				Amounts: []Amount{
					{
						"ethereum",
						"eth",
						0.7,
					},
				},
			},
			when: Amount{
				"ethereum",
				"eth",
				0.5,
			},
			expected: &Account{
				UserName: "julio",
				Amounts: []Amount{
					{
						"ethereum",
						"eth",
						0.2,
					},
				},
			},
		},
		{
			name: "Expected balance equal 0",
			given: &Account{
				UserName: "julio",
				Amounts: []Amount{
					{
						"ethereum",
						"eth",
						0.9,
					},
				},
			},
			when: Amount{
				"ethereum",
				"eth",
				0.9,
			},
			expected: &Account{
				UserName: "julio",
				Amounts: []Amount{
					{
						"ethereum",
						"eth",
						0,
					},
				},
			},
		},
	}
	asserts := assert.New(t)
	for _, test := range tests {
		account, _ := test.given.Withdraw(context.Background(), test.when)
		asserts.Equal(test.expected, account, test.name)
	}
}

func TestErrorWhenWithdraw(t *testing.T) {

	tests := []struct {
		name  string
		given *Account
		when  Amount
		expected error
		wantErr bool
	}{
		{
			name: "Expected error cannot be negative",
			given: &Account{
				UserName: "julio",
				Amounts: []Amount{
					{
						"ethereum",
						"eth",
						0.7,
					},
				},
			},
			when: Amount{
				"ethereum",
				"eth",
				-4,
			},
			expected: errors.New("value cannot be zero or negative"),
		},
		{
			name: "Expected error insufficient eth funds",
			given: &Account{
				UserName: "julio",
				Amounts: []Amount{
					{
						"ethereum",
						"eth",
						0.9,
					},
				},
			},
			when: Amount{
				"ethereum",
				"eth",
				1,
			},
			expected: errors.New("insufficient eth funds"),
		},
		{
			name: "Expected error there is no eth amount",
			given: &Account{
				UserName: "julio",
				Amounts: []Amount{
					{
						"bitcoin",
						"btc",
						0.04,
					},
				},
			},
			when: Amount{
				"ethereum",
				"eth",
				0.3,
			},
			expected: errors.New("there is no eth amount"),
		},
		{
			name: "Expected error there is no amounts",
			given: &Account{
				UserName: "julio",
			},
			when: Amount{
				"bitcoin",
				"btc",
				0.1,
			},
			expected: errors.New("there is no amounts"),
		},
	}
	asserts := assert.New(t)
	for _, test := range tests {
		_, err := test.given.Withdraw(context.Background(), test.when)
		asserts.Equal(test.expected, err, test.name)
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		name  string
		given *Account
		when  string
		expected bool
	}{
		{
			name: "Expected false account does not contain received currency",
			given: &Account{
				UserName: "julio",
				Amounts: []Amount{
					{
						"ethereum",
						"eth",
						0.5,
					},
				},
			},
			when: "btc",
			expected: false,
		},
		{
			name: "Expected true account contain received currency",
			given: &Account{
				UserName: "julio",
				Amounts: []Amount{
					{
						"bitcoin",
						"btc",
						0.5,
					},
				},
			},
			when: "btc",
			expected: true,
		},
	}
	asserts := assert.New(t)
	for _, test := range tests {
		asserts.Equal(test.expected, test.given.contains(test.when), test.name)
	}
}
