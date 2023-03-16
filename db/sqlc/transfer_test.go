package db

import (
	"context"
	"github.com/dangquyit/go-simplebank/util"
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
)

func createRandomTransfer(t *testing.T) Transfer {
	listAccount, _ := testQueries.ListAccount(context.Background(), ListAccountParams{
		Limit:  100,
		Offset: 0,
	})
	n := len(listAccount)
	var fromAccountNumber int64
	var toAccountNumber int64
	for {
		fromAccountNumber = listAccount[rand.Intn(n)].AccountNumber
		toAccountNumber = listAccount[rand.Intn(n)].AccountNumber
		if fromAccountNumber != toAccountNumber {
			break
		}
	}
	transfer, err := testQueries.CreateTransfer(context.Background(), CreateTransferParams{
		FromAccountNumber: fromAccountNumber,
		ToAccountNumber:   toAccountNumber,
		Amount:            util.RandomMoney(),
	})

	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, fromAccountNumber, transfer.FromAccountNumber)
	require.Equal(t, toAccountNumber, transfer.ToAccountNumber)

	require.NotEmpty(t, transfer.Amount)
	return transfer
}

func TestCreateTransfer(t *testing.T) {
	createRandomTransfer(t)
}

func TestListTransfers(t *testing.T) {
	transfer := createRandomTransfer(t)

	transfers, err := testQueries.ListTransfers(context.Background(), ListTransfersParams{
		FromAccountNumber: transfer.FromAccountNumber,
		ToAccountNumber:   transfer.ToAccountNumber,
		Limit:             5,
		Offset:            0,
	})

	require.NoError(t, err)
	require.NotEmpty(t, transfers)

	for _, i := range transfers {
		require.NotEmpty(t, i)
	}
}
