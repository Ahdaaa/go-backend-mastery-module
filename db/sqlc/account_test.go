package db

import (
	"context"
	"testing"

	"github.com/Ahdaaa/go-backend-mastery-module/tree/main/util"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Accounts {
	preBalance := util.RandomMoney()

	balance := pgtype.Numeric{}
	balance.Scan(preBalance)

	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  balance,
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, account)

	// Equal will make sure that the pre-inserted values
	// are the same with the post-inserted
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	// Check the ID and CreatedAt is automatically
	// generated by postgres
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
	return account
}
func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

// func TestGetAccount(t *testing.T) {
// 	account1 := createRandomAccount(t)
// 	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
// 	require.NoError(t, err)
// 	require.NotEmpty(t, account2)

// 	require.Equal(t, account1.ID, account2.Owner)
// 	require.Equal(t, account1.Owner, account2.Owner)
// 	require.Equal(t, account1.Balance, account2.Balance)
// 	require.Equal(t, account1.Currency, account2.Currency)
// }
