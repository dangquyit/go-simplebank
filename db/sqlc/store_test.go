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
	n := 10
	amount := int64(10)
	errs := make(chan error)
	//results := make(chan TransferTxResult)
	for i := 0; i < n; i++ {
		fromAccountId := account1.ID
		toAccountId := account2.ID

		if i%2 == 1 {
			fromAccountId, toAccountId = toAccountId, fromAccountId
		}
		go func() {
			_, err := store.TransferTx(context.Background(), TransferParams{
				FromAccountId: fromAccountId,
				ToAccountId:   toAccountId,
				Amount:        amount,
			})
			errs <- err
			//results <- result
		}()
	}

	// check results
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		//result := <-results
		//require.NotEmpty(t, result)
		//
		//// check transfer
		//transfer := result.Transfer
		//require.Equal(t, account1.AccountNumber, transfer.FromAccountNumber)
		//require.Equal(t, account2.AccountNumber, transfer.ToAccountNumber)
		//require.Equal(t, amount, transfer.Amount)
		//require.NotZero(t, transfer.ID)
		//require.NotZero(t, transfer.CreatedAt)
		//require.NotZero(t, transfer.UpdatedAt)
		//
		//_, err = store.GetTransferById(context.Background(), transfer.ID)
		//require.NoError(t, err)
		//
		//// check entries
		//fromEntry := result.FromEntry
		//require.NotEmpty(t, fromEntry)
		//require.Equal(t, account1.AccountNumber, fromEntry.AccountNumber)
		//require.Equal(t, -amount, fromEntry.Amount)
		//require.NotZero(t, fromEntry.ID)
		//require.NotZero(t, fromEntry.CreatedAt)
		//require.NotZero(t, fromEntry.UpdatedAt)
		//_, err = store.GetEntryById(context.Background(), fromEntry.ID)
		//require.NoError(t, err)
		//
		//toEntry := result.ToEntry
		//require.NotEmpty(t, toEntry)
		//require.Equal(t, account2.AccountNumber, toEntry.AccountNumber)
		//require.Equal(t, amount, toEntry.Amount)
		//require.NotZero(t, toEntry.ID)
		//require.NotZero(t, toEntry.CreatedAt)
		//require.NotZero(t, toEntry.UpdatedAt)
		//
		//_, err = store.GetEntryById(context.Background(), toEntry.ID)
		//require.NoError(t, err)
		//
		//// check accounts
		//fromAccount := result.FromAccount
		//require.NotEmpty(t, fromAccount)
		//require.Equal(t, account1.ID, fromAccount.ID)
		//
		//toAccount := result.ToAccount
		//require.NotEmpty(t, toAccount)
		//require.Equal(t, account2.ID, toAccount.ID)
		//
		////check accounts balance
		//diff1 := account1.Balance - fromAccount.Balance
		//diff2 := toAccount.Balance - account2.Balance
		//require.Equal(t, diff1, diff2)
		//require.True(t, diff1 > 0)
		//require.True(t, diff1%amount == 0)
		//
		//k := int(diff1 / amount)
		//require.True(t, k >= 1 && k <= 5)
	}

	//check final updated balance
	updateAccount1, err := testQueries.GetAccountById(context.Background(), account1.ID)
	require.NoError(t, err)
	updateAccount2, err := testQueries.GetAccountById(context.Background(), account2.ID)
	require.NoError(t, err)
	//require.Equal(t, account1.Balance-int64(n)*amount, updateAccount1.Balance)
	//require.Equal(t, account2.Balance+int64(n)*amount, updateAccount2.Balance)
	require.Equal(t, account1.Balance, updateAccount1.Balance)
	require.Equal(t, account2.Balance, updateAccount2.Balance)
}
