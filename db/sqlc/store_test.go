package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTranferTx(t *testing.T) {
	store := NewStore(testPool)

	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)

	amount := int64(10)
	n := 4

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			ctx := context.Background()
			result, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: int64(account1.ID),
				ToAccountID:   int64(account2.ID),
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}
	//check results
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err, "Error in transfer transaction")

		result := <-results
		require.NotEmpty(t, result, "Transfer transaction result should not be empty")

		//check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer, "Transfer should not be empty")
		require.Equal(t, account1.ID, transfer.FromAccountID, "FromAccountID should match")
		require.Equal(t, account2.ID, transfer.ToAccountID, "ToAccountID should match")
		require.Equal(t, amount, transfer.Amount, "Transfer amount should match")
		require.NotZero(t, transfer.ID, "Transfer ID should not be zero")
		require.NotZero(t, transfer.CreatedAt, "Transfer created_at should not be zero")

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err, "Error getting transfer after transaction")

		//check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry, "From entry should not be empty")
		require.Equal(t, account1.ID, fromEntry.AccountID, "From entry AccountID should match")
		require.Equal(t, -amount, fromEntry.Amount, "From entry amount should be negative")
		require.NotZero(t, fromEntry.ID, "From entry ID should not be zero")
		require.NotZero(t, fromEntry.CreatedAt, "From entry created_at should not be zero")
		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err, "Error getting from entry after transaction")

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry, "To entry should not be empty")
		require.Equal(t, account2.ID, toEntry.AccountID, "To entry AccountID should match")
		require.Equal(t, amount, toEntry.Amount, "To entry amount should match")
		require.NotZero(t, toEntry.ID, "To entry ID should not be zero")
		require.NotZero(t, toEntry.CreatedAt, "To entry created_at should not be zero")
		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err, "Error getting to entry after transaction")

		//check accounts
		fromAccount := result.AccountFrom
		require.NotEmpty(t, fromAccount, "From account should not be empty")
		require.Equal(t, account1.ID, fromAccount.ID, "From account ID should match")
		toAccount := result.AccountTo
		require.NotEmpty(t, toAccount, "To account should not be empty")
		require.Equal(t, account2.ID, toAccount.ID, "To account ID should match")
		//check balances
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0, "Balance difference should be positive")
		require.True(t, diff1%amount == 0, "Balance difference should be a multiple of the transfer amount")
	}
	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	expectedFromBalance := account1.Balance - int64(n)*amount
	expectedToBalance := account2.Balance + int64(n)*amount

	require.Equal(t, expectedFromBalance, updatedAccount1.Balance)
	require.Equal(t, expectedToBalance, updatedAccount2.Balance)
}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testPool)

	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)

	amount := int64(10)
	n := 10

	errs := make(chan error)

	for i := 0; i < n; i++ {
		fromAccountID := account1.ID
		toAccountID := account2.ID
		if i%2 == 1 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}

		go func() {
			ctx := context.Background()
			_, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})
			errs <- err
		}()
	}
	//check results
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err, "Error in transfer transaction")
	}

	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	require.Equal(t, account1.Balance, updatedAccount1.Balance)
	require.Equal(t, account2.Balance, updatedAccount2.Balance)
}
