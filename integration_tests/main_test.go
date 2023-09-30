package integration_tests

import (
	"fmt"
	groshi "github.com/groshi-project/go-groshi"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"testing"
	"time"
)

const GroshiSocketEnvVarName = "GROSHI_TEST_SOCKET"

var GroshiSocket string

func TestMain(m *testing.M) {
	GroshiSocket = os.Getenv(GroshiSocketEnvVarName)
	if GroshiSocket == "" {
		panic(
			fmt.Sprintf("environmental variable %v is not set", GroshiSocketEnvVarName),
		)
	}

	os.Exit(m.Run())
}

func TestUserCreate(t *testing.T) {
	client := NewPureGroshiClient(GroshiSocket)

	// create a new user:
	username, password := GenerateCredentials()
	user, err := client.UserCreate(username, password)
	assert.NoError(t, err)

	assert.Equal(t, username, user.Username)

	// try to create user with the same username:
	user, err = client.UserCreate(username, password)
	if assert.Error(t, err) {
		if assert.IsType(t, groshi.GroshiAPIError{}, err) {
			assert.Equal(t, http.StatusConflict, err.(groshi.GroshiAPIError).HTTPStatusCode)
		}
	}
	assert.Nil(t, user)

	// try to create user with  username and password
	user, err = client.UserCreate("", "")
	if assert.Error(t, err) {
		if assert.IsType(t, groshi.GroshiAPIError{}, err) {
			assert.Equal(t, http.StatusBadRequest, err.(groshi.GroshiAPIError).HTTPStatusCode)
		}
	}
	assert.Nil(t, user)
}

func TestAuth(t *testing.T) {
	username, password, client := NewGroshiClientWithUser(GroshiSocket)

	credentials, err := client.AuthLogin(username, password)
	assert.NoError(t, err)

	assert.NotEmpty(t, credentials.Token)
	assert.NotEmpty(t, credentials.ExpiresAt)
}

func TestAuthRefresh(t *testing.T) {
	_, _, client := NewAuthorizedGroshiClientWithUser(GroshiSocket)
	newCredentials, err := client.AuthRefresh()
	assert.NoError(t, err)

	assert.NotEmpty(t, newCredentials.Token)
	assert.NotEmpty(t, newCredentials.ExpiresAt)
}

func TestUserRead(t *testing.T) {
	username, _, client := NewAuthorizedGroshiClientWithUser(GroshiSocket)
	user, err := client.UserRead()
	assert.NoError(t, err)

	assert.Equal(t, username, user.Username)
}

func TestUserUpdate(t *testing.T) {
	_, _, client := NewAuthorizedGroshiClientWithUser(GroshiSocket)

	// update the current user:
	newUsername, newPassword := GenerateCredentials()
	user, err := client.UserUpdate(&newUsername, &newPassword)
	assert.NoError(t, err)

	assert.Equal(t, newUsername, user.Username)

	// read the current user:
	readUser, err := client.UserRead()
	assert.NoError(t, err)

	assert.Equal(t, user.Username, readUser.Username)
}

func TestUserDelete(t *testing.T) {
	username, _, client := NewAuthorizedGroshiClientWithUser(GroshiSocket)
	user, err := client.UserDelete()
	assert.NoError(t, err)

	assert.Equal(t, username, user.Username)

	// read the current deleted user
	readUser, err := client.UserRead()
	if assert.Error(t, err) {
		if assert.IsType(t, groshi.GroshiAPIError{}, err) {
			assert.Equal(t, http.StatusUnauthorized, err.(groshi.GroshiAPIError).HTTPStatusCode)
		}
	}
	assert.Nil(t, readUser)
}

func TestTransactionsCreate(t *testing.T) {
	_, _, client := NewAuthorizedGroshiClientWithUser(GroshiSocket)

	// create a new transaction:
	amount := 500
	currency := "USD"
	description := "Hello world"
	date := time.Now()

	transaction, err := client.TransactionsCreate(
		amount,
		currency,
		&description,
		&date,
	)
	assert.NoError(t, err)

	assert.Equal(t, amount, transaction.Amount)
	assert.Equal(t, currency, transaction.Currency)
	assert.Equal(t, description, transaction.Description)
	assert.Equal(
		t, date.In(time.UTC).Format(time.RFC3339), transaction.Timestamp.Format(time.RFC3339),
	)

	// fetch the transaction:
	uuid := transaction.UUID
	transaction, err = client.TransactionsReadOne(uuid, nil)
	assert.NoError(t, err)

	assert.Equal(t, uuid, transaction.UUID)
	assert.Equal(t, amount, transaction.Amount)
	assert.Equal(t, currency, transaction.Currency)
	assert.Equal(t, description, transaction.Description)
	assert.Equal(
		t, date.In(time.UTC).Format(time.RFC3339), transaction.Timestamp.Format(time.RFC3339),
	)
}

func TestTransactionsReadOne(t *testing.T) {
	_, _, client := NewAuthorizedGroshiClientWithUser(GroshiSocket)

	// create a new transaction:
	amount := 1000
	currency := "USD"
	description := "Description of the transaction"
	timestamp := time.Now()

	transaction, err := client.TransactionsCreate(
		amount, currency, &description, &timestamp,
	)
	assert.NoError(t, err)
	assert.NotEmpty(t, transaction.UUID)

	// read the transaction without currency conversion:
	readTransaction, err := client.TransactionsReadOne(transaction.UUID, nil)
	assert.NoError(t, err)

	assert.Equal(t, transaction.UUID, readTransaction.UUID)
	assert.Equal(t, transaction.Amount, readTransaction.Amount)
	assert.Equal(t, transaction.Currency, readTransaction.Currency)
	assert.Equal(t, transaction.Description, readTransaction.Description)
	assert.Equal(t, transaction.Timestamp, readTransaction.Timestamp)

	// read the transaction with currency conversion to the same currency:
	readTransaction, err = client.TransactionsReadOne(transaction.UUID, &currency)
	assert.NoError(t, err)

	assert.Equal(t, transaction.UUID, readTransaction.UUID)
	assert.Equal(t, transaction.Amount, readTransaction.Amount)
	assert.Equal(t, transaction.Currency, readTransaction.Currency)
	assert.Equal(t, transaction.Description, readTransaction.Description)
	assert.Equal(t, transaction.Timestamp, readTransaction.Timestamp)

	// read the transaction with currency conversion to EUR:
	currency = "EUR"
	readTransaction, err = client.TransactionsReadOne(transaction.UUID, &currency)
	assert.NoError(t, err)

	assert.Equal(t, transaction.UUID, readTransaction.UUID)
	assert.NotEmpty(t, readTransaction.Amount)
	assert.Equal(t, currency, readTransaction.Currency)
	assert.Equal(t, transaction.Description, readTransaction.Description)
	assert.Equal(t, transaction.Timestamp, readTransaction.Timestamp)
}

func TestTransactionsReadMany(t *testing.T) {
	_, _, client := NewAuthorizedGroshiClientWithUser(GroshiSocket)

	// create 10 test transactions:
	transactionsCount := 10
	for i := 0; i < transactionsCount; i++ {
		amount := -200
		currency := "USD"
		description := "Transaction description"
		timestamp := time.Now()
		_, err := client.TransactionsCreate(
			amount, currency, &description, &timestamp,
		)
		assert.NoError(t, err)
	}

	// read transactions without currency conversion:
	startTime := time.Now().Add(-time.Hour) // an hour ago
	readTransactions, err := client.TransactionsReadMany(startTime, nil, nil)
	assert.NoError(t, err)
	assert.Len(t, readTransactions, transactionsCount)

	// read transactions with currency conversions to the same currency:
	currency := "USD"
	readTransactions, err = client.TransactionsReadMany(startTime, nil, &currency)
	assert.NoError(t, err)
	assert.Len(t, readTransactions, transactionsCount)

	// read transactions with currency conversions to EUR:
	currency = "EUR"
	readTransactions, err = client.TransactionsReadMany(startTime, nil, &currency)
	assert.NoError(t, err)
	assert.Len(t, readTransactions, transactionsCount)

}

func TestTransactionsSummary(t *testing.T) {
	_, _, client := NewAuthorizedGroshiClientWithUser(GroshiSocket)

	currency := "USD"
	// create 2 transactions:
	_, err := client.TransactionsCreate(150, currency, nil, nil) // $1.50 income
	assert.NoError(t, err)

	_, err = client.TransactionsCreate(-100, currency, nil, nil) // $1 outcome
	assert.NoError(t, err)

	startTime := time.Now().Add(-time.Hour) // an hour ago

	// fetch summary without currency conversions
	summary, err := client.TransactionsReadSummary(currency, startTime, nil)
	assert.NoError(t, err)

	assert.Equal(t, 2, summary.TransactionsCount)
	assert.Equal(t, 150, summary.Income)  // $1.50 income
	assert.Equal(t, 100, summary.Outcome) // $1 outcome
	assert.Equal(t, 50, summary.Total)    // $1.50 - $1 = $0.50
	assert.Equal(t, "USD", currency)

	// fetch summary with currency conversions
	summary, err = client.TransactionsReadSummary("EUR", startTime, nil)
	assert.NoError(t, err)
	assert.Equal(t, 2, summary.TransactionsCount)
	assert.NotEmpty(t, summary.Income)
	assert.NotEmpty(t, summary.Outcome)
	assert.NotEmpty(t, summary.Total)
	assert.NotEmpty(t, "EUR")
}

func TestTransactionsUpdate(t *testing.T) {
	_, _, client := NewAuthorizedGroshiClientWithUser(GroshiSocket)

	amount := -350
	currency := "EUR"
	description := "Payed bills for goddamn drum lessons"
	timestamp := time.Now()

	transaction, err := client.TransactionsCreate(
		amount,
		currency,
		&description,
		&timestamp,
	)
	assert.NoError(t, err)
	assert.NotEmpty(t, transaction)

	// update without providing anything new (error is to be expected)
	_, err = client.TransactionsUpdate(transaction.UUID, nil, nil, nil, nil)
	if assert.Error(t, err) {
		if assert.IsType(t, groshi.GroshiAPIError{}, err) {
			assert.Equal(t, http.StatusBadRequest, err.(groshi.GroshiAPIError).HTTPStatusCode)
		}
	}

	newAmount1 := -600
	newCurrency := "USD"
	newDescription := "Payed bills for my car insurance (I think I cannot afford to own a car)"
	newTimestamp := time.Date(2023, 1, 1, 15, 30, 1, 0, time.Local)

	// update amount only:
	updatedTransaction1, err := client.TransactionsUpdate(
		transaction.UUID,
		&newAmount1,
		nil,
		nil,
		nil,
	)
	assert.NoError(t, err)
	if assert.NotNil(t, updatedTransaction1) {
		assert.Equal(t, newAmount1, updatedTransaction1.Amount)
	}

	// update all fields:
	newAmount2 := -999
	updatedTransaction2, err := client.TransactionsUpdate(
		transaction.UUID,
		&newAmount2,
		&newCurrency,
		&newDescription,
		&newTimestamp,
	)

	assert.NoError(t, err)
	if assert.NotNil(t, updatedTransaction2) {
		assert.Equal(t, newAmount2, updatedTransaction2.Amount)
		assert.Equal(t, newCurrency, updatedTransaction2.Currency)
		assert.Equal(t, newDescription, updatedTransaction2.Description)
		assert.Equal(t, newTimestamp.In(time.UTC), updatedTransaction2.Timestamp)
	}
}

func TestTransactionsDelete(t *testing.T) {
	_, _, client := NewAuthorizedGroshiClientWithUser(GroshiSocket)

	// create a new transaction:
	transaction, err := client.TransactionsCreate(100, "USD", nil, nil)
	assert.NoError(t, err)

	// delete the transaction:
	deletedTransaction, err := client.TransactionsDelete(transaction.UUID)
	assert.NoError(t, err)

	assert.Equal(t, transaction.UUID, deletedTransaction.UUID)

	// try to read the deleted transaction:
	_, err = client.TransactionsReadOne(deletedTransaction.UUID, nil)
	if assert.Error(t, err) {
		if assert.IsType(t, groshi.GroshiAPIError{}, err) {
			assert.Equal(t, http.StatusNotFound, err.(groshi.GroshiAPIError).HTTPStatusCode)
		}
	}

}

func TestCurrenciesRead(t *testing.T) {
	client := NewPureGroshiClient(GroshiSocket)

	// fetch slice of supported currencies:
	currencies, err := client.CurrenciesRead()
	assert.NoError(t, err)
	assert.NotEmpty(t, currencies)
}
