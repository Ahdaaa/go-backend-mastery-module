package sqlc

import (
	"context"
	"fmt"

	"github.com/Ahdaaa/go-backend-mastery-module/tree/main/util"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

// This will provide all functions to execute db queries and transactions
// Queries are handled by the Queries struct inside db.go
// But it cant handle transactions so we will extend it in store.
type Store struct {
	*Queries // instead of inheritance, we'll use this
	db       *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function in db transaction
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit(ctx)
}

type TransferTxParams struct {
	FromAccountID int64          `json:"sender_account_id"`
	ToAccountID   int64          `json:"receiver_account_id"`
	Amount        pgtype.Numeric `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfers `json:"transfer"`
	FromAccount Accounts  `json:"from_account"`
	ToAccount   Accounts  `json:"to_account"`
	FromEntry   Entries   `json:"from_entry"`
	ToEntry     Entries   `json:"to_entry"`
}

// transferTx will perform money transfering
// it will create record, entries, and update balance within a single db transaction.
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			SenderAccountID:   arg.FromAccountID,
			ReceiverAccountID: arg.ToAccountID,
			Amount:            arg.Amount,
		})

		if err != nil {
			return err
		}

		// setNegative := fmt.Sprintf("-%.2f", arg.Amount)
		decreaseAmount := util.CloneNegativeNumeric(arg.Amount)

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    decreaseAmount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		// Below is for Updating Account Balance
		// Will be a bit complicated, to prevent deadlock, etc.
		// TODO!

		return nil
	})

	return result, err
}
