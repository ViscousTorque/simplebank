package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var txKey = struct{}{}

/*
Store provides all functions to execute all db queries and transactions

	Notes:
		Instead of inheritance, we have composition to extend a struct functionality
		The Queries struct doesn't support transactions, so it needs to be extended

ChatGPT update to course code:

	To make use of the pgx library instead of database/sql, you need to update the Store
	struct to include a pgx connection pool (*pgxpool.Pool) instead of *sql.DB. Additionally,
	you should ensure that the Queries struct is modified to work with pgx if necessary.
*/
type Store struct {
	*Queries
	db *pgxpool.Pool
}

// NewStore creates a new Store with the given pgxpool.Pool
func NewStore(pool *pgxpool.Pool) *Store {
	return &Store{
		Queries: New(pool),
		db:      pool,
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	q := New(tx)

	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return rbErr
		}
		return err
	}

	return tx.Commit(ctx)
}

/* TransferTXParams contains the input parameters for the transfer transactions
 */
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"amount"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

func (store *Store) TransferTx(ctx context.Context, args TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		txName := ctx.Value(txKey)
		fmt.Println(txName, "Create Transfer")
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams(args))
		if err != nil {
			return err
		}
		fmt.Println(txName, "Create Entry for FromAccount")
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: args.FromAccountID,
			Amount:    -args.Amount,
		})
		if err != nil {
			return err
		}
		fmt.Println(txName, "Create Entry for ToAccount")
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: args.ToAccountID,
			Amount:    args.Amount,
		})
		if err != nil {
			return err
		}

		fmt.Println(txName, "Transfer Money")
		if args.FromAccountID < args.ToAccountID {
			result.FromAccount, result.ToAccount, err = transferMoney(ctx, q, args.FromAccountID, -args.Amount, args.ToAccountID, args.Amount)
		} else {
			result.ToAccount, result.FromAccount, err = transferMoney(ctx, q, args.ToAccountID, args.Amount, args.FromAccountID, -args.Amount)
		}
		return err
	})

	return result, err
}

func transferMoney(ctx context.Context, q *Queries,
	fromAccountID int64, fromAccountAmount int64,
	toAccountID int64, toAccountAmount int64,
) (fromAccount Account, toAccount Account, err error) {
	txName := ctx.Value(txKey)
	fmt.Println(txName, "AddAccountBalance - fromAccount")
	fromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     fromAccountID,
		Amount: fromAccountAmount,
	})
	if err != nil {
		return
	}
	fmt.Println(txName, "AddAccountBalance - toAccount")
	toAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     toAccountID,
		Amount: toAccountAmount,
	})
	return
}
