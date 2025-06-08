package db

import (
	"context"
	"testing"

	"github.com/Yogksai/backend-projects/simplebank/util"
	"github.com/stretchr/testify/require"
)

func CreateRandomTransfer(t *testing.T) Transfer {
	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)
	arg := CreateTransferParams{
		FromAccountID: int64(account1.ID),   // Replace with a valid account ID
		ToAccountID:   int64(account2.ID),   // Replace with a valid account ID
		Amount:        util.RandomBalance(), // Replace with a valid amount
	}
	transfer, _ := testQueries.CreateTransfer(context.Background(), arg)
	return transfer
}

func TestCreateTransfer(t *testing.T) {
	CreateRandomTransfer(t)
}

func TestGetTransfer(t *testing.T) {
	transfer1 := CreateRandomTransfer(t)
	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)
	require.NoError(t, err, "Error getting transfer")
	require.NotEmpty(t, transfer2, "Transfer should not be empty")
	require.Equal(t, transfer1.ID, transfer2.ID, "Transfer IDs do not match")
	require.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID, "FromAccountIDs do not match")
	require.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID, "ToAccountIDs do not match")
	require.Equal(t, transfer1.Amount, transfer2.Amount, "Transfer amounts do not match")
	require.Equal(t, transfer1.CreatedAt, transfer2.CreatedAt, "Transfer created_at timestamps do not match")
}
func TestUpdateTransfer(t *testing.T) {
	transfer1 := CreateRandomTransfer(t)
	arg := UpdateTransferParams{
		ID:     transfer1.ID,
		Amount: util.RandomBalance(),
	}
	transfer2, err := testQueries.UpdateTransfer(context.Background(), arg)
	require.NoError(t, err, "Error updating transfer")
	require.NotEmpty(t, transfer2, "Transfer should not be empty")

	require.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID, "FromAccountIDs do not match")
	require.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID, "ToAccountIDs do not match")
	require.NotEqual(t, transfer1.Amount, transfer2.Amount, "Transfer amounts should not match")
}

func TestDeleteTransfer(t *testing.T) {
	transfer := CreateRandomTransfer(t)
	err := testQueries.DeleteTransfer(context.Background(), transfer.ID)
	require.NoError(t, err, "Error deleting transfer")

	transfer2, err := testQueries.GetTransfer(context.Background(), transfer.ID)
	require.Error(t, err, "Expected error when getting deleted transfer")
	require.Empty(t, transfer2, "Transfer should be empty after deletion")
}

func TestListTransfer(t *testing.T) {
	arg := ListTransfersParams{
		Limit:  2,
		Offset: 0,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err, "Error listing transfers")
	require.NotEmpty(t, transfers, "Transfers should not be empty")

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer.ID, "Transfer ID should not be empty")
	}
}
