package db

import (
	"context"
	"github.com/dangquyit/go-simplebank/util"
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
)

func createRandomEntry(t *testing.T) Entry {
	accounts, _ := testQueries.ListAccount(context.Background(), ListAccountParams{
		Limit:  1000,
		Offset: 0,
	})
	arg := CreateEntryParams{
		AccountID: accounts[rand.Intn(len(accounts))].ID,
		Amount:    util.RandomMoney(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.Equal(t, arg.Amount, entry.Amount)
	require.Equal(t, arg.AccountID, entry.AccountID)
	return entry
}

func TestCreateEntry(t *testing.T) {
	createRandomEntry(t)
}

func TestGetEntry(t *testing.T) {
	entry1 := createRandomEntry(t)
	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
}

func TestListEntries(t *testing.T) {
	accounts, _ := testQueries.ListAccount(context.Background(), ListAccountParams{
		Limit:  1000,
		Offset: 0,
	})
	arg := ListEntriesParams{
		AccountID: accounts[rand.Intn(len(accounts)-1)].ID,
		Limit:     5,
		Offset:    0,
	}
	entries, err := testQueries.ListEntries(context.Background(), arg)

	require.NoError(t, err)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}
