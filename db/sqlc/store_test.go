package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStore(t *testing.T) {
	store := NewStore(testDb)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	n := 5
	var amount int64 = 10

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        int64(amount),
			})

			errs <- err
			results <- *result
		}()
	}

	var existed = make(map[int]bool)

	// check result
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		//check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, transfer.FromAccountID, account1.ID)
		require.Equal(t, transfer.ToAccountID, account2.ID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check from account
		FromEntry := result.FromEntry
		require.NotEmpty(t, FromEntry)
		require.Equal(t, FromEntry.AccountID, account1.ID)
		require.Equal(t, -amount, FromEntry.Amount)
		require.NotZero(t, FromEntry.ID)
		require.NotZero(t, FromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), FromEntry.ID)
		require.NoError(t, err)

		// check from account
		ToEntry := result.ToEntry
		require.NotEmpty(t, ToEntry)
		require.Equal(t, ToEntry.AccountID, account2.ID)
		require.Equal(t, amount, ToEntry.Amount)
		require.NotZero(t, ToEntry.ID)
		require.NotZero(t, ToEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), ToEntry.ID)
		require.NoError(t, err)

		// TODO: check account balance
		FromAccount := result.FromAccount
		require.NotEmpty(t, FromAccount)
		require.Equal(t, FromAccount.ID, account1.ID)

		ToAccount := result.ToAccount
		require.NotEmpty(t, ToAccount)
		require.Equal(t, ToAccount.ID, account2.ID)

		diff1 := account1.Balance - FromAccount.Balance
		diff2 := ToAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	require.Equal(t, account1.Balance-int64(n)*amount, updatedAccount1.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updatedAccount2.Balance)
}

func TestStoreDeadlock(t *testing.T) {
	store := NewStore(testDb)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	n := 10
	var amount int64 = 10

	errs := make(chan error)

	for i := 0; i < n; i++ {
		fromAccountID := account1.ID
		toAccountID := account2.ID

		if i%2 != 0 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}

		go func() {
			_, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        int64(amount),
			})

			errs <- err
		}()
	}

	// check result
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	require.Equal(t, account1.Balance, updatedAccount1.Balance)
	require.Equal(t, account2.Balance, updatedAccount2.Balance)
}
