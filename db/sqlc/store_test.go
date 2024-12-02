package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	testFromAccount := createRandomAccount(t)
	testToAccount := createRandomAccount(t)

	/* 	Run with 5 concurrent transactions, this is because if we are not careful with concurrency
	it can quickly become problematic
		Also we need to send the results from the go routines to the main testing routine using channels
	*/
	amount := int64(10) // Test below only works for +ve value here
	numberTestGoRoutines := 5

	errs := make(chan error)
	testResults := make(chan TransferTxResult)

	for i := 0; i < numberTestGoRoutines; i++ {
		txName := fmt.Sprintf("tx %d", i)
		go func() {
			ctx := context.WithValue(context.Background(), txKey, txName)
			result, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: testFromAccount.ID,
				ToAccountID:   testToAccount.ID,
				Amount:        amount,
			})

			errs <- err
			testResults <- result
		}()
	}

	// check results
	existed := make(map[int]bool) // TODO: need to check if this concurrent safe

	// Verify the results
	for i := 0; i < numberTestGoRoutines; i++ {
		actualErr := <-errs
		require.NoError(t, actualErr)

		actualResult := <-testResults
		require.NotEmpty(t, actualResult)

		// Verify Transfer
		actualTransfer := actualResult.Transfer
		require.Equal(t, testFromAccount.ID, actualTransfer.FromAccountID)
		require.Equal(t, testToAccount.ID, actualTransfer.ToAccountID)
		require.Equal(t, amount, actualTransfer.Amount)
		require.NotZero(t, actualTransfer.ID)
		require.NotZero(t, actualTransfer.CreatedAt)

		// To be sure the transaction exists in the DB, do at least a GetTransfer
		_, err := store.GetTransfer(context.Background(), actualTransfer.ID)
		require.NoError(t, err)

		// Verify From Entry
		actualFromEntry := actualResult.FromEntry
		require.NotEmpty(t, actualFromEntry)
		require.Equal(t, testFromAccount.ID, actualFromEntry.AccountID)
		require.Equal(t, -amount, actualFromEntry.Amount)
		require.NotZero(t, actualFromEntry.ID)
		require.NotZero(t, actualFromEntry.CreatedAt)

		// To be sure the Entry has been created in the db, do the minimum GetEntry
		_, err = store.GetEntry(context.Background(), actualFromEntry.ID)
		require.NoError(t, err)

		// Verify From Entry
		actualToEntry := actualResult.ToEntry
		require.NotEmpty(t, actualToEntry)
		require.Equal(t, testToAccount.ID, actualToEntry.AccountID)
		require.Equal(t, amount, actualToEntry.Amount)
		require.NotZero(t, actualToEntry.ID)
		require.NotZero(t, actualToEntry.CreatedAt)

		// To be sure the Entry has been created in the db, do the minimum GetEntry
		_, err = store.GetEntry(context.Background(), actualToEntry.ID)
		require.NoError(t, err)

		// check accounts
		actualFromAccount := actualResult.FromAccount
		require.NotEmpty(t, actualFromAccount)
		require.Equal(t, testFromAccount.ID, actualFromAccount.ID)

		actualToAccount := actualResult.ToAccount
		require.NotEmpty(t, actualToAccount)
		require.Equal(t, testToAccount.ID, actualToAccount.ID)

		// check balances
		expectedFromAccountNewBalance := testFromAccount.Balance - actualFromAccount.Balance
		expectedToAccountNewBalance := actualToAccount.Balance - testToAccount.Balance
		require.Equal(t, expectedFromAccountNewBalance, expectedToAccountNewBalance)
		require.True(t, expectedFromAccountNewBalance > 0)
		require.True(t, expectedFromAccountNewBalance%amount == 0) // 1 * amount, 2 * amount, 3 * amount, ..., n * amount

		// Check number of transactions is the same as the go routines
		k := int(expectedFromAccountNewBalance / amount)
		require.True(t, k >= 1 && k <= numberTestGoRoutines)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// // check the final updated balance
	txName := fmt.Sprintf("tx GetAccount %d", 1)
	ctx := context.WithValue(context.Background(), txKey, txName)
	updatedFromAccount, err := store.GetAccount(ctx, testFromAccount.ID)
	require.NoError(t, err)

	updateToAccount, err := store.GetAccount(context.Background(), testToAccount.ID)
	require.NoError(t, err)

	require.Equal(t, testFromAccount.Balance-int64(numberTestGoRoutines)*amount, updatedFromAccount.Balance)
	require.Equal(t, testToAccount.Balance+int64(numberTestGoRoutines)*amount, updateToAccount.Balance)
}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)

	testFromAccount := createRandomAccount(t)
	testToAccount := createRandomAccount(t)
	amount := int64(10) // Test below only works for +ve value here
	numberTestGoRoutines := 10

	errs := make(chan error)

	for i := 0; i < numberTestGoRoutines; i++ {
		testFromAccountId := testFromAccount.ID
		testToAccountId := testToAccount.ID

		if i%2 == 1 {
			testFromAccountId = testToAccount.ID
			testToAccountId = testFromAccount.ID
		}

		txName := fmt.Sprintf("tx %d", i)
		go func() {
			ctx := context.WithValue(context.Background(), txKey, txName)
			_, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: testFromAccountId,
				ToAccountID:   testToAccountId,
				Amount:        amount,
			})

			errs <- err
		}()
	}

	// Verify the results
	for i := 0; i < numberTestGoRoutines; i++ {
		actualErr := <-errs
		require.NoError(t, actualErr)
	}

	// // check the final updated balance
	txName := fmt.Sprintf("tx GetAccount %d", 1)
	ctx := context.WithValue(context.Background(), txKey, txName)
	updatedFromAccount, err := store.GetAccount(ctx, testFromAccount.ID)
	require.NoError(t, err)

	updateToAccount, err := store.GetAccount(context.Background(), testToAccount.ID)
	require.NoError(t, err)

	require.Equal(t, testFromAccount.Balance, updatedFromAccount.Balance)
	require.Equal(t, testToAccount.Balance, updateToAccount.Balance)
}
