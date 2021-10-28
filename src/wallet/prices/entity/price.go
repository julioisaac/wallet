package entity

import "time"

type Price struct {
	CryptoCurrency 	string
	Currencies    	map[string]float64
	ExchangeDataBy	string
	LastUpdate		time.Time
}