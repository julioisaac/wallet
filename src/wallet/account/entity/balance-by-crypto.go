package entity

import "time"

type BalanceByCrypto struct {
	Currency 		string  			`json:"currency"`
	Amount   		float64 			`json:"amount"`
	Prices	 		map[string]float64 	`json:"prices,omitempty"`
	TimeOfRate		time.Time			`json:"time-of-rate,omitempty"`
	ExchangeDataBy	string				`json:"exchange-data-by,omitempty"`
	TotalByCurrency	map[string]float64 	`json:"total-by-currency,omitempty"`
}