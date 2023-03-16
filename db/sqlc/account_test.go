package db

import (
	"context"
	"github.com/dangquyit/go-simplebank/util"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func createRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		AccountNumber: util.RandomAccountNumber(),
		Owner:         util.RandomOwner(),
		UserName:      util.RandomUserName(),
		Email:         util.RandomEmail(),
		Password:      util.RandomPassword(),
		PhoneNumber:   util.RandomAccountNumber(),
		Balance:       util.RandomMoney(),
		Currency:      util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.AccountNumber, account.AccountNumber)
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.UserName, account.UserName)
	require.Equal(t, arg.Email, account.Email)
	require.Equal(t, arg.Password, account.Password)
	require.Equal(t, arg.PhoneNumber, account.PhoneNumber)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
	require.NotZero(t, account.UpdatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	account2, err := testQueries.GetAccount(context.Background(), account1.AccountNumber)

	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.AccountNumber, account2.AccountNumber)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.UserName, account2.UserName)
	require.Equal(t, account1.Email, account2.Email)
	require.Equal(t, account1.Password, account2.Password)
	require.Equal(t, account1.PhoneNumber, account2.PhoneNumber)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
	require.WithinDuration(t, account1.UpdatedAt, account2.UpdatedAt, time.Second)

}

func TestUpdateAccount(t *testing.T) {
	account1 := createRandomAccount(t)

	arg := UpdateAccountParams{
		AccountNumber: account1.AccountNumber,
		Balance:       util.RandomMoney(),
	}

	account2, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account2)
	require.Equal(t, account1.AccountNumber, account2.AccountNumber)
	require.Equal(t, arg.Balance, account2.Balance)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), account1.AccountNumber)
	require.NoError(t, err)

	account2, err := testQueries.GetAccount(context.Background(), account1.AccountNumber)
	require.Error(t, err)
	require.Empty(t, account2)
}

func TestListAccount(t *testing.T) {
	arg := ListAccountParams{
		Limit:  1000,
		Offset: 0,
	}
	listAccount, err := testQueries.ListAccount(context.Background(), arg)

	require.NoError(t, err)
	for _, account := range listAccount {
		require.NotEmpty(t, account)
	}
}
