package entity

type Balances struct {
	User 		string				`json:"user"`
	ByCrypto	[]BalanceByCrypto	`json:"by-cryptos"`
	Total		map[string]float64	`json:"total,omitempty"`
}