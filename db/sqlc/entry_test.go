package db

import (
	"context"
	"github.com/dangquyit/go-simplebank/util"
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
	"time"
)

func createRandomEntry(t *testing.T) Entry {
	accounts, _ := testQueries.ListAccount(context.Background(), ListAccountParams{
		Limit:  10,
		Offset: 0,
	})
	arg := CreateEntryParams{
		AccountNumber: accounts[rand.Intn(len(accounts))].AccountNumber,
		Amount:        util.RandomMoney(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.Equal(t, arg.Amount, entry.Amount)
	require.Equal(t, arg.AccountNumber, entry.AccountNumber)
	return entry
}

func TestCreateEntry(t *testing.T) {
	createRandomEntry(t)
}

func TestGetEntry(t *testing.T) {
	entry1 := createRandomEntry(t)
	entry2, err := testQueries.GetEntry(context.Background(), entry1.AccountNumber)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)
	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.AccountNumber, entry2.AccountNumber)
	require.WithinDuration(t, entry1.CreatedAt, entry2.CreatedAt, time.Second)
}

func TestListEntries(t *testing.T) {
	accounts, _ := testQueries.ListAccount(context.Background(), ListAccountParams{
		Limit:  10,
		Offset: 0,
	})
	arg := ListEntriesParams{
		AccountNumber: accounts[rand.Intn(len(accounts))].AccountNumber,
		Limit:         5,
		Offset:        0,
	}
	entries, err := testQueries.ListEntries(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, entries)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}
