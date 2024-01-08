// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package db

import (
	"time"
)

type Accounts struct {
	ID        int64     `json:"id"`
	Owner     string    `json:"owner"`
	Balance   string    `json:"balance"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
}

type Entries struct {
	ID        int64 `json:"id"`
	AccountID int64 `json:"account_id"`
	// can be negative/positive
	Amount    string    `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

type Transfers struct {
	ID                int64 `json:"id"`
	SenderAccountID   int64 `json:"sender_account_id"`
	ReceiverAccountID int64 `json:"receiver_account_id"`
	// must be positive
	Amount    string    `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}
