package db

import "context"

/* TransferTXParams contains the input parameters for the transfer transactions
 */
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

func (store *SQLStore) TransferTx(ctx context.Context, args TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams(args))
		if err != nil {
			return err
		}
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: args.FromAccountID,
			Amount:    -args.Amount,
		})
		if err != nil {
			return err
		}
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: args.ToAccountID,
			Amount:    args.Amount,
		})
		if err != nil {
			return err
		}
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
	fromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     fromAccountID,
		Amount: fromAccountAmount,
	})
	if err != nil {
		return
	}
	toAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     toAccountID,
		Amount: toAccountAmount,
	})
	return
}
