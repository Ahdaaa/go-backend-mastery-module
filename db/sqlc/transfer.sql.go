// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: transfer.sql

package sqlc

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createTransfer = `-- name: CreateTransfer :one
INSERT INTO transfers (
  sender_account_id, 
  receiver_account_id,
  amount
) VALUES (
  $1, $2, $3
) RETURNING id, sender_account_id, receiver_account_id, amount, created_at
`

type CreateTransferParams struct {
	SenderAccountID   int64          `json:"sender_account_id"`
	ReceiverAccountID int64          `json:"receiver_account_id"`
	Amount            pgtype.Numeric `json:"amount"`
}

func (q *Queries) CreateTransfer(ctx context.Context, arg CreateTransferParams) (Transfers, error) {
	row := q.db.QueryRow(ctx, createTransfer, arg.SenderAccountID, arg.ReceiverAccountID, arg.Amount)
	var i Transfers
	err := row.Scan(
		&i.ID,
		&i.SenderAccountID,
		&i.ReceiverAccountID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

const getTransfer = `-- name: GetTransfer :one
SELECT id, sender_account_id, receiver_account_id, amount, created_at FROM transfers
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetTransfer(ctx context.Context, id int64) (Transfers, error) {
	row := q.db.QueryRow(ctx, getTransfer, id)
	var i Transfers
	err := row.Scan(
		&i.ID,
		&i.SenderAccountID,
		&i.ReceiverAccountID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

const listTransfers = `-- name: ListTransfers :many
SELECT id, sender_account_id, receiver_account_id, amount, created_at FROM transfers
ORDER BY id
LIMIT $1
OFFSET $2
`

type ListTransfersParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListTransfers(ctx context.Context, arg ListTransfersParams) ([]Transfers, error) {
	rows, err := q.db.Query(ctx, listTransfers, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Transfers
	for rows.Next() {
		var i Transfers
		if err := rows.Scan(
			&i.ID,
			&i.SenderAccountID,
			&i.ReceiverAccountID,
			&i.Amount,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
