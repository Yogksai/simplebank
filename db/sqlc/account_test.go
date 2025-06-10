package db

import (
	"context"
	"testing"

	"github.com/Yogksai/simplebank/util"
	"github.com/stretchr/testify/require"
)

func TestCreateAccount(t *testing.T) {
	CreateRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	account1 := CreateRandomAccount(t)
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)

	require.NoError(t, err, "Error getting account")
	require.NotEmpty(t, account2, "Account should not be empty")

	require.Equal(t, account1.ID, account2.ID, "Account IDs do not match")
	require.Equal(t, account1.Owner, account2.Owner, "Account owners do not match")
	require.Equal(t, account1.Balance, account2.Balance, "Account balances do not match")
	require.Equal(t, account1.Currency, account2.Currency, "Account currencies do not match")
}

func TestUpdateAccount(t *testing.T) {
	account1 := CreateRandomAccount(t)
	arg := UpdateAccountParams{
		ID:      account1.ID,
		Balance: util.RandomBalance(),
	}
	account2, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err, "Error getting account")
	require.NotEmpty(t, account2, "Account should not be empty")

	require.Equal(t, account1.ID, account2.ID, "Account IDs do not match")
	require.Equal(t, account1.Owner, account2.Owner, "Account owners do not match")
	require.NotEqual(t, arg.Balance, account1.Balance, "Account balances should not match")
	require.Equal(t, account1.Currency, account2.Currency, "Account currencies do not match")
}

func TestDeleteAccount(t *testing.T) {
	account1 := CreateRandomAccount(t)
	err := testQueries.DeleteAсcount(context.Background(), account1.ID)
	require.NoError(t, err)

	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.Error(t, err, "Error getting account after deletion")
	require.EqualError(t, err, "no rows in result set")
	require.Empty(t, account2, "Account should be empty after deletion")
}

func TestListAccounts(t *testing.T) {
	for i := 0; i < 2; i++ {
		CreateRandomAccount(t)
	}
	arg := ListAccountsParams{
		Limit:  2,
		Offset: 0,
	}
	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 2)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}

func CreateRandomAccount(t *testing.T) Account {
	user := CreateRandomUser(t)
	arg := CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomBalance(),
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

	return account
}
