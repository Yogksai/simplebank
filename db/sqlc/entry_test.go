package db

import (
	"context"
	"testing"

	"github.com/Yogksai/backend-projects/simplebank/util"
	"github.com/stretchr/testify/require"
)

func CreateRandomEntry(t *testing.T) Entry {
	randomAccount := CreateRandomAccount(t)
	arg := CreateEntryParams{
		AccountID: randomAccount.ID, // Assuming account IDs are in this range
		Amount:    randomAccount.Balance,
	}
	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}

func TestCreateEntry(t *testing.T) {
	CreateRandomEntry(t)
}

func TestGetEntry(t *testing.T) {
	entry1 := CreateRandomEntry(t)
	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)
	require.NoError(t, err, "Error getting entry")
	require.NotEmpty(t, entry2, "Entry should not be empty")

	require.Equal(t, entry1.ID, entry2.ID, "Entry IDs should match")
	require.Equal(t, entry1.AccountID, entry2.AccountID, "Account IDs should match")
	require.Equal(t, entry1.Amount, entry2.Amount, "Amounts should match")
}
func TestDeleteEntry(t *testing.T) {
	entry := CreateRandomEntry(t)
	err := testQueries.DeleteEntry(context.Background(), entry.ID)
	require.NoError(t, err, "Error deleting entry")

	entry2, err := testQueries.GetEntry(context.Background(), entry.ID)
	require.Error(t, err, "Expected error when getting deleted entry")
	require.Empty(t, entry2, "Entry should be empty after deletion")
}

func TestUpdateEntry(t *testing.T) {
	entry1 := CreateRandomEntry(t)
	arg := UpdateEntryParams{
		ID:     entry1.ID,
		Amount: util.RandomBalance(),
	}
	entry2, err := testQueries.UpdateEntry(context.Background(), arg)
	require.NoError(t, err, "Error updating entry")
	require.NotEmpty(t, entry2, "Updated entry should not be empty")

	require.Equal(t, entry1.AccountID, entry2.AccountID, "Account IDs should match after update")
	require.Equal(t, arg.Amount, entry2.Amount, "Amounts should match after update")
}

func TestListEntry(t *testing.T) {
	// Create multiple entries
	for i := 0; i < 5; i++ {
		CreateRandomEntry(t)
	}

	// List entries
	entries, err := testQueries.ListEntries(context.Background(), ListEntriesParams{
		Limit:  5,
		Offset: 0,
	})
	require.NoError(t, err, "Error listing entries")
	require.Len(t, entries, 5, "Expected 5 entries")

	for _, entry := range entries {
		require.NotEmpty(t, entry.ID, "Entry ID should not be empty")
		require.NotEmpty(t, entry.AccountID, "Account ID should not be empty")
		require.NotEmpty(t, entry.Amount, "Amount should not be empty")
	}
}
