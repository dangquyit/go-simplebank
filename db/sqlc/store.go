package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferParams) (TransferTxResult, error)
}

type SQLStore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		Queries: New(db),
		db:      db,
	}
}

func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)

	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v rb err: %v", err, rbErr)
		}
	}

	return tx.Commit()
}

type TransferParams struct {
	FromAccountId int64 `json:"from_account_id"`
	ToAccountId   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

func (store *SQLStore) TransferTx(ctx context.Context, arg TransferParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		// transfer
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountId,
			ToAccountID:   arg.ToAccountId,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		// entry
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountId,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountId,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		//update balance
		// using if handle deadlock when account1 transfer to account2 and account2 transfer to account1
		if arg.FromAccountId < arg.ToAccountId {
			result.FromAccount, result.ToAccount, err = updateMoney(ctx, q,
				arg.FromAccountId, -arg.Amount,
				arg.ToAccountId, arg.Amount)
		} else {
			result.ToAccount, result.FromAccount, err = updateMoney(ctx, q,
				arg.ToAccountId, arg.Amount,
				arg.FromAccountId, -arg.Amount)
		}
		return err
	})

	return result, err
}

func updateMoney(
	ctx context.Context,
	q *Queries,
	fromAccountId, amount1, toAccountId, amount2 int64,
) (account1, account2 Account, err error) {
	account1, err = q.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
		Amount: amount1,
		ID:     fromAccountId,
	})

	if err != nil {
		return
	}
	account2, err = q.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
		Amount: amount2,
		ID:     toAccountId,
	})
	return
}
