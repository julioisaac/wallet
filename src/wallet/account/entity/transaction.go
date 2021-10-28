package entity

type Transaction struct {
	Type 	 string `json:"type"`
	UserName string `json:"user"`
	Amount   Amount `json:"amount"`
}