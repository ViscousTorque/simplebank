package db

import (
	"context"
	"database/sql"
	"errors"
	"main/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	// We will need to re-use the TestCreateAccount so changing to ...

	/* Initial test demo
	We should consider creating random data so that we dont have to do manual test data creation
	arg := CreateAccountParams{
		Owner:    "Tom",
		Balance:  100,
		Currency: "USD",
	}
	*/
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	testAccount := createRandomAccount(t)

	getAccountResult, err := testQueries.GetAccount(context.Background(), testAccount.ID)
	require.NoError(t, err)
	require.NotEmpty(t, getAccountResult)

	require.Equal(t, testAccount.ID, getAccountResult.ID)
	require.Equal(t, testAccount.Owner, getAccountResult.Owner)
	require.Equal(t, testAccount.Balance, getAccountResult.Balance)
	require.Equal(t, testAccount.CreatedAt, getAccountResult.CreatedAt)
	require.Equal(t, testAccount.Currency, getAccountResult.Currency)
}

func TestUpdateAccount(t *testing.T) {
	testAccount := createRandomAccount(t)
	updateAccountArgs := UpdateAccountParams{
		ID:      testAccount.ID,
		Balance: util.RandomMoney(),
	}

	updateAccountResult, err := testQueries.UpdateAccount(context.Background(), updateAccountArgs)
	require.NoError(t, err)
	require.NotEmpty(t, updateAccountResult)

	require.Equal(t, testAccount.ID, updateAccountResult.ID)
	require.Equal(t, testAccount.Owner, updateAccountResult.Owner)
	require.Equal(t, updateAccountArgs.Balance, updateAccountResult.Balance)
	require.Equal(t, testAccount.CreatedAt, updateAccountResult.CreatedAt)
	require.Equal(t, testAccount.Currency, updateAccountResult.Currency)
}

func TestDeleteAccount(t *testing.T) {
	testAccount := createRandomAccount(t)

	err := testQueries.DeleteAccount(context.Background(), testAccount.ID)
	require.NoError(t, err)

	getAccountResult, err := testQueries.GetAccount(context.Background(), testAccount.ID)
	require.Error(t, err)
	require.True(t, errors.Is(err, sql.ErrNoRows), "expected sql.ErrNoRows error") //sql and pg errors are formatted slightly differently
	require.Empty(t, getAccountResult)
}

func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	args := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	listAccountsResult, err := testQueries.ListAccounts(context.Background(), args)
	require.NoError(t, err)
	require.Len(t, listAccountsResult, 5)

	for _, account := range listAccountsResult {
		require.NotEmpty(t, account)
	}
}
