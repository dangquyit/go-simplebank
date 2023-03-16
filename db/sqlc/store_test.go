package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStoreTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	// run n concurrent transfer transactions
	n := 5
	amount := int64(10)
	errs := make(chan error)
	results := make(chan TransferTxResult)
	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferParams{
				FromAccountNumber: account1.AccountNumber,
				ToAccountNumber:   account2.AccountNumber,
				Amount:            amount,
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

		// check transfer
		transfer := result.Transfer
		require.Equal(t, account1.AccountNumber, transfer.FromAccountNumber)
		require.Equal(t, account2.AccountNumber, transfer.ToAccountNumber)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)
		require.NotZero(t, transfer.UpdatedAt)

		_, err = store.GetTransferById(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.AccountNumber, fromEntry.AccountNumber)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)
		require.NotZero(t, fromEntry.UpdatedAt)

		_, err = store.GetEntryById(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.AccountNumber, toEntry.AccountNumber)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)
		require.NotZero(t, toEntry.UpdatedAt)

		_, err = store.GetEntryById(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// TODO: check accounts balance
	}
}