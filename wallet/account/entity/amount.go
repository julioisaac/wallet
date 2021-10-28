package entity

type Amount struct {
	Id       string  `json:"id"`
	Currency string  `json:"currency"`
	Value    float64 `json:"value"`
}
