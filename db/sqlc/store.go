package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		Queries: New(db),
		db:      db,
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
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
	FromAccountNumber int64 `json:"from_account_number"`
	ToAccountNumber   int64 `json:"to_account_number"`
	Amount            int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

func (store *Store) TransferTx(ctx context.Context, arg TransferParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		// transfer
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountNumber: arg.FromAccountNumber,
			ToAccountNumber:   arg.ToAccountNumber,
			Amount:            arg.Amount,
		})
		if err != nil {
			return err
		}

		// entry
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountNumber: arg.FromAccountNumber,
			Amount:        -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountNumber: arg.ToAccountNumber,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		//update balance
		result.FromAccount, err = q.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
			AccountNumber: arg.FromAccountNumber,
			Amount:        -arg.Amount,
		})

		result.ToAccount, err = q.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
			AccountNumber: arg.ToAccountNumber,
			Amount:        arg.Amount,
		})

		return nil
	})

	return result, err
}
