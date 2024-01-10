package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/Ahdaaa/go-backend-mastery-module/tree/main/util"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func TestCreateAccount(t *testing.T) {

	preBalance := util.RandomMoney()
	convert := fmt.Sprintf("%v", preBalance)
	balance := pgtype.Numeric{}
	balance.Scan(convert)

	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  balance,
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

}
