package sqlc

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/Ahdaaa/go-backend-mastery-module/tree/main/util"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println(">> before:", account1.Balance, account2.Balance)

	// run n concurrent transfer tx (co-routines)
	n := 5
	preAmount := "10.00"
	amount := pgtype.Numeric{}
	amount.Scan(preAmount)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		// handle deadlock
		txName := fmt.Sprintf("tx %d", i+1)
		go func() { // go routines
			ctx := context.WithValue(context.Background(), txKey, txName)
			result, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}

	// check results
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		negAmount := util.CloneNegativeNumeric(amount)

		//check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.SenderAccountID)
		require.Equal(t, account2.ID, transfer.ReceiverAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, negAmount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// get account results
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		// check accounts balance
		fmt.Println(">> tx:", fromAccount.Balance, toAccount.Balance)
		diff1 := new(big.Int)
		diff1.Sub(account1.Balance.Int, fromAccount.Balance.Int)

		diff2 := new(big.Int)
		diff2.Sub(toAccount.Balance.Int, account2.Balance.Int)

		modResult := new(big.Int)
		modResult.Mod(diff1, amount.Int)

		require.Equal(t, diff1, diff2)                     // they should be equal
		require.True(t, diff1.Cmp(big.NewInt(0)) > 0)      // they should be positive
		require.True(t, modResult.Cmp(big.NewInt(0)) == 0) // they should be divisible to the amount, because account balance will be decreased by 1*amount, 2*amount, till n*amount.

	}

	// check the final updated balances
	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	mulResult := new(big.Int)
	mulResult.Mul(amount.Int, big.NewInt(int64(n))) // n * amount
	subResult := new(big.Int)
	subResult.Sub(account1.Balance.Int, mulResult)
	sumResult := new(big.Int)
	sumResult.Add(account2.Balance.Int, mulResult)

	fmt.Println(">> after:", updatedAccount1.Balance, updatedAccount2.Balance)
	require.Equal(t, subResult, updatedAccount1.Balance.Int)
	require.Equal(t, sumResult, updatedAccount2.Balance.Int)

}
