package models

import (
	"fmt"
	"time"
)

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

func (t *Transaction) DisplayAmount() string {
	return fmt.Sprintf("%.2f", float64(t.Amount)/100.0)
}
