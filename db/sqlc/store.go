package db

import (
	"context"

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
type SQLStore struct {
	*Queries
	db *pgxpool.Pool
}

type Store interface {
	Querier
	TransferTx(ctx context.Context, args TransferTxParams) (TransferTxResult, error)
	CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error)
	VerifyEmailTx(ctx context.Context, arg VerifyEmailTxParams) (VerifyEmailTxResult, error)
}

// NewStore creates a new Store with the given pgxpool.Pool
func NewStore(pool *pgxpool.Pool) Store {
	return &SQLStore{
		Queries: New(pool),
		db:      pool,
	}
}
