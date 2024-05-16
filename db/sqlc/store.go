package sqlc

import (
	"context"
	"fmt"
	"math/big"

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

// deadlock handling
var txKey = struct{}{} // empty object of a struct

// transferTx will perform money transfering
// it will create record, entries, and update balance within a single db transaction.
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		txName := ctx.Value(txKey)

		fmt.Println(txName, "create transfer")
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			SenderAccountID:   arg.FromAccountID,
			ReceiverAccountID: arg.ToAccountID,
			Amount:            arg.Amount,
		})

		if err != nil {
			return err
		}

		fmt.Println(txName, "create entry 1")
		decreaseAmount := util.CloneNegativeNumeric(arg.Amount)
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    decreaseAmount,
		})
		if err != nil {
			return err
		}

		fmt.Println(txName, "create entry 2")
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

		fmt.Println(txName, "get account 1")
		account1, err := q.GetAccountForUpdate(ctx, arg.FromAccountID)
		if err != nil {
			return err
		}

		newInt := new(big.Int)

		updateBalance1 := pgtype.Numeric{
			Int:              newInt.Sub(account1.Balance.Int, arg.Amount.Int),
			Exp:              -2,
			NaN:              false,
			Valid:            true,
			InfinityModifier: 0,
		}

		fmt.Println(txName, "update account 1")
		result.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID:      arg.FromAccountID,
			Balance: updateBalance1,
		})
		if err != nil {
			return err
		}

		fmt.Println(txName, "get account 2")
		account2, err := q.GetAccountForUpdate(ctx, arg.ToAccountID)
		if err != nil {
			return err
		}

		newInt2 := new(big.Int)

		updateBalance2 := pgtype.Numeric{
			Int:              newInt2.Add(account2.Balance.Int, arg.Amount.Int),
			Exp:              -2,
			NaN:              false,
			Valid:            true,
			InfinityModifier: 0,
		}

		fmt.Println(txName, "update account 2")
		result.ToAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID:      arg.ToAccountID,
			Balance: updateBalance2,
		})
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}
