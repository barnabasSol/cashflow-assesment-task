package models

import "time"

const (
	TransactionPending = iota
	TransactionSucceeded
	TransactionFailed = -1
)

type Transaction struct {
	ID        string    `json:"id"`
	Amount    int64     `json:"amount"`
	Currency  string    `json:"currency"`
	Ref       string    `json:"ref"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
