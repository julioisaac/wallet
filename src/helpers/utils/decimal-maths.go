package utils

import "github.com/shopspring/decimal"

type calc struct {}

func DecimalMaths() *calc {
	return &calc{}
}

func (*calc) Sum(amount, incoming float64) float64 {
	return decimal.NewFromFloat(amount).Add(decimal.NewFromFloat(incoming)).InexactFloat64()
}

func (*calc) Sub(amount, incoming float64) float64  {
	return decimal.NewFromFloat(amount).Sub(decimal.NewFromFloat(incoming)).InexactFloat64()
}

func (*calc) Mul(amount, incoming float64) float64 {
	return decimal.NewFromFloat(amount).Mul(decimal.NewFromFloat(incoming)).InexactFloat64()
}