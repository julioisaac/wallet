package entity

type Account struct {
	UserName string   `json:"user"`
	Amounts  []Amount `json:"amounts"`
}
