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
		// using if handle deadlock when account1 transfer to account2 and account2 transfer to account1
		if arg.FromAccountNumber < arg.ToAccountNumber {
			result.FromAccount, result.ToAccount, err = updateMoney(ctx, q,
				arg.FromAccountNumber, -arg.Amount,
				arg.ToAccountNumber, arg.Amount)
		} else {
			result.ToAccount, result.FromAccount, err = updateMoney(ctx, q,
				arg.ToAccountNumber, arg.Amount,
				arg.FromAccountNumber, -arg.Amount)
		}
		return nil
	})

	return result, err
}

func updateMoney(
	ctx context.Context,
	q *Queries,
	fromAccountNumber, amount1, toAccountNumber, amount2 int64,
) (account1, account2 Account, err error) {
	account1, err = q.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
		Amount:        amount1,
		AccountNumber: fromAccountNumber,
	})

	if err != nil {
		return
	}
	account2, err = q.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
		Amount:        amount2,
		AccountNumber: toAccountNumber,
	})
	return
}
