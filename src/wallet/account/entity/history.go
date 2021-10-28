package entity

import "time"

type History struct {
	UserName  string    `json:"user"`
	Type      string    `json:"type"`
	Amount    Amount    `json:"amount"`
	EventTime time.Time `json:"event_time"`
}