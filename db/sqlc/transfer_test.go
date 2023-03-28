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
	var fromAccountId int64
	var toAccountId int64
	for {
		fromAccountId = listAccount[rand.Intn(n)].ID
		toAccountId = listAccount[rand.Intn(n)].ID
		if fromAccountId != toAccountId {
			break
		}
	}
	transfer, err := testQueries.CreateTransfer(context.Background(), CreateTransferParams{
		FromAccountID: fromAccountId,
		ToAccountID:   toAccountId,
		Amount:        util.RandomMoney(),
	})

	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, fromAccountId, transfer.FromAccountID)
	require.Equal(t, toAccountId, transfer.ToAccountID)

	require.NotEmpty(t, transfer.Amount)
	return transfer
}

func TestCreateTransfer(t *testing.T) {
	createRandomTransfer(t)
}

func TestListTransfers(t *testing.T) {
	transfer := createRandomTransfer(t)

	transfers, err := testQueries.ListTransfers(context.Background(), ListTransfersParams{
		FromAccountID: transfer.FromAccountID,
		ToAccountID:   transfer.ToAccountID,
		Limit:         5,
		Offset:        0,
	})

	require.NoError(t, err)
	require.NotEmpty(t, transfers)

	for _, i := range transfers {
		require.NotEmpty(t, i)
	}
}
