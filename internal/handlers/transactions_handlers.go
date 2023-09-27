package handlers

import (
	"errors"
	"math"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/groshi-project/groshi/internal/currency/currency_rates"
	"github.com/groshi-project/groshi/internal/database"
	"github.com/groshi-project/groshi/internal/handlers/util"
	"github.com/groshi-project/groshi/internal/loggers"
	"github.com/groshi-project/groshi/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const errorDescriptionTransactionNotFound = "transaction was not found"
const errorDescriptionTransactionForbidden = "you have no right to access to this transaction"

type transactionsCreateParams struct {
	Amount   int    `json:"amount" binding:"required"`
	Currency string `json:"currency" binding:"required,currency"`

	Description *string    `json:"description" binding:"omitempty,description"`
	Timestamp   *time.Time `json:"timestamp"`
}

// TransactionsCreateHandler creates a new transaction.
//
//	@summary		create a new transaction
//	@description	Creates a new transaction owned by current user.
//	@tags			transactions
//	@accept			json
//	@produce		json
//	@param			amount		body		integer				true	"Negative or positive amount of transaction in minor units."
//	@param			currency	body		string				true	"Currency code of transaction in ISO-4217 format."
//	@param			description	body		string				false	"Description of transaction."
//	@param			timestamp	body		string				false	"Timestamp of transaction in RFC-3339 format."
//	@success		200			{object}	models.Transaction	"Object of newly created transaction is returned."
//	@router			/transactions [post]
func TransactionsCreateHandler(c *gin.Context) {
	params := transactionsCreateParams{}
	if ok := util.BindBody(c, &params); !ok {
		return
	}

	currentUser := c.MustGet("current_user").(*database.User)

	// use empty description as transaction description if description was not provided:
	if params.Description == nil {
		emptyDescription := ""
		params.Description = &emptyDescription
	}

	// use the current time as transaction time if time was not provided:
	if params.Timestamp == nil {
		currentTime := time.Now()
		params.Timestamp = &currentTime
	}

	transaction := database.Transaction{
		ID:   primitive.NewObjectID(),
		UUID: uuid.New().String(),

		OwnerID: currentUser.ID,

		Amount:   params.Amount,
		Currency: params.Currency,

		Description: *params.Description,
		Timestamp:   *params.Timestamp,

		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if _, err := database.TransactionsCol.InsertOne(database.Context, &transaction); err != nil {
		util.AbortWithStatusInternalServerError(c, err)
		return
	}

	util.ReturnSuccessfulResponse(c, transaction.APIModel())
}

type transactionsReadOneParams struct {
	Currency string `form:"currency" binding:"optional_currency"`
}

// TransactionsReadOneHandler returns information about single transaction.
// Optionally converts its amount to the given currency units.
//
//	@summary		fetch one transaction
//	@description	Returns information about one transaction.
//	@tags			transactions
//	@accept			json
//	@produce		json
//	@param			uuid		path		string				true	"UUID of transaction."
//	@param			currency	query		string				false	"The currency to convert amount to"
//	@success		200			{object}	models.Transaction	"Transaction object is returned."
//	@failure		404			{object}	models.Error		"Transaction was not found."
//	@failure		403			{object}	models.Error		"You have no right to read this transaction."
//	@router			/transactions/{uuid} [get]
func TransactionsReadOneHandler(c *gin.Context) {
	params := transactionsReadOneParams{}
	if ok := util.BindQuery(c, &params); !ok {
		return
	}

	transactionUUID := c.Param("uuid")

	currentUser := c.MustGet("current_user").(*database.User)

	transaction := database.Transaction{}
	if err := database.TransactionsCol.FindOne(
		database.Context, bson.D{{"uuid", transactionUUID}},
	).Decode(&transaction); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			util.AbortWithStatusNotFound(c, errorDescriptionTransactionNotFound)
		} else {
			util.AbortWithStatusInternalServerError(c, err)
		}
		return
	}
	if transaction.OwnerID != currentUser.ID {
		util.AbortWithStatusForbidden(c, errorDescriptionTransactionForbidden)
		return
	}

	// convert amount to the new currency if needed:
	if params.Currency != "" {
		if params.Currency != transaction.Currency {
			newAmount, err := currency_rates.Convert(transaction.Currency, params.Currency, float64(transaction.Amount))
			if err != nil {
				util.AbortWithStatusInternalServerError(c, err)
			}
			transaction.Amount = int(math.Round(newAmount))
			transaction.Currency = params.Currency
		}
	}

	util.ReturnSuccessfulResponse(c, transaction.APIModel())
}

type transactionsReadManyParams struct {
	StartTime time.Time `form:"start_time" binding:"required,nonzero_time"`

	EndTime  *time.Time `form:"end_time"`
	Currency string     `form:"currency" binding:"optional_currency"`
}

// TransactionsReadManyHandler returns all transactions for time given period.
// Optionally converts all amounts to given currency units.
//
//	@summary		fetch many transactions
//	@description	Returns array of transactions for given time period.
//	@tags			transactions
//	@accept			json
//	@produce		json
//	@param			start_time	query		string					true	"Beginning of the time period in RFC-3339 format."
//	@param			end_time	query		string					false	"End of the time period in RFC-3339 format (current time is used by default if no value provided)."
//	@param			currency	query		string					false	"The currency to convert amount to."
//	@success		200			{object}	[]models.Transaction	"Array of transaction objects is returned."
//	@router			/transactions [get]
func TransactionsReadManyHandler(c *gin.Context) {
	params := transactionsReadManyParams{}
	if ok := util.BindQuery(c, &params); !ok {
		return
	}
	if params.EndTime == nil {
		currentTime := time.Now()
		params.EndTime = &currentTime
	}

	currentUser := c.MustGet("current_user").(*database.User)

	cursor, err := database.TransactionsCol.Find(
		database.Context,
		bson.D{
			{
				"owner_id", currentUser.ID,
			},
			{
				"created_at",
				bson.D{
					{"$gte", params.StartTime},
					{"$lt", *params.EndTime},
				},
			}},
	)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) { // unexpected error happened
		util.AbortWithStatusInternalServerError(c, err)
		return
	}

	defer func() {
		if err := cursor.Close(database.Context); err != nil {
			loggers.Error.Printf("could not close cursor: %v", err)
		}
	}()

	transactions := make([]*models.Transaction, 0)
	for cursor.Next(database.Context) {
		transaction := database.Transaction{}
		if err := cursor.Decode(&transaction); err != nil {
			util.AbortWithStatusInternalServerError(c, err)
		}

		// convert amount to the new currency if needed:
		if params.Currency != "" {
			if params.Currency != transaction.Currency {
				newAmount, err := currency_rates.Convert(transaction.Currency, params.Currency, float64(transaction.Amount))
				if err != nil {
					util.AbortWithStatusInternalServerError(c, err)
					return
				}
				transaction.Amount = int(math.Round(newAmount))
				transaction.Currency = params.Currency
			}
		}

		transactions = append(transactions, transaction.APIModel())
	}

	util.ReturnSuccessfulResponse(c, transactions)
}

type transactionsReadSummaryParams struct {
	StartTime time.Time `form:"start_time" binding:"required,nonzero_time"`
	Currency  string    `form:"currency" binding:"required,currency"`

	EndTime *time.Time `form:"end_time"`
}

// TransactionsReadSummary returns summary (count and sum of transaction)
// for given time period in desired currency units.
//
//	@summary		fetch summary of transactions for given time period
//	@description	Returns summary of transactions for given time period in desired currency units.
//	@tags			transactions
//	@accept			json
//	@produce		json
//	@param			start_time	query		string			true	"Beginning of the time period in RFC-3339 format."
//	@param			currency	query		string			true	"Desired currency of sum of transactions in ISO-4217 format."
//	@param			end_time	query		string			false	"End of the time period in RFC-3339 format (current time is used by default if no value provided)."
//	@success		200			{object}	models.Summary	"Summary object is returned."
//	@router			/transactions/summary [get]
func TransactionsReadSummary(c *gin.Context) {
	params := transactionsReadSummaryParams{}
	if ok := util.BindQuery(c, &params); !ok {
		return
	}
	currentUser := c.MustGet("current_user").(*database.User)

	if params.EndTime == nil { // set current time as end time if end time was not provided
		currentTime := time.Now()
		params.EndTime = &currentTime
	}

	cursor, err := database.TransactionsCol.Find(
		database.Context,
		bson.D{
			{
				"owner_id", currentUser.ID,
			},
			{
				"created_at",
				bson.D{
					{"$gte", params.StartTime},
					{"$lte", *params.EndTime},
				},
			},
		},
	)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) { // if unexpected error happened
		util.AbortWithStatusInternalServerError(c, err)
		return
	}

	defer func() {
		if err := cursor.Close(database.Context); err != nil {
			loggers.Error.Printf("could not close cursor: %v", err)
		}
	}()

	var transactions []*database.Transaction
	for cursor.Next(database.Context) {
		transaction := database.Transaction{}
		if err := cursor.Decode(&transaction); err != nil {
			util.AbortWithStatusInternalServerError(c, err)
			return
		}
		transactions = append(transactions, &transaction)
	}

	// sum of all transaction amounts (in the same `params.Currency` units)
	//
	// important notice: only final sum must be rounded,
	// not intermediate amounts.
	transactionsCount := 0
	income := 0.0
	outcome := 0.0

	for _, transaction := range transactions {
		var transactionAmountInSameUnits float64         // transaction amount in the `params.Currency` units
		transactionAmount := float64(transaction.Amount) // transaction amount in the original currency units

		if transaction.Currency == params.Currency {
			transactionAmountInSameUnits = transactionAmount
		} else {
			transactionAmountInSameUnits, err = currency_rates.Convert(
				transaction.Currency, params.Currency, transactionAmount,
			)
			if err != nil {
				util.AbortWithStatusInternalServerError(c, err)
				return
			}
		}
		if transactionAmountInSameUnits >= 0 {
			income += transactionAmountInSameUnits
		} else {
			outcome += -transactionAmountInSameUnits
		}
		transactionsCount++
	}

	intTotal := int(math.Round(income - outcome))
	intIncome := int(math.Round(income))
	intOutcome := int(math.Round(outcome))

	util.ReturnSuccessfulResponse(c, &models.Summary{
		Currency:          params.Currency,
		Income:            intIncome,
		Outcome:           intOutcome,
		Total:             intTotal,
		TransactionsCount: transactionsCount,
	})
}

type transactionsUpdateParams struct {
	NewAmount      *int       `json:"new_amount" binding:"omitempty,required_without:NewCurrency,NewDescription,NewTimestamp"`
	NewCurrency    *string    `json:"new_currency" binding:"omitempty,optional_currency,required_without:NewAmount,NewTimestamp"`
	NewDescription *string    `json:"new_description" binding:"omitempty,description,required_without:NewAmount,NewCurrency,NewTimestamp"`
	NewTimestamp   *time.Time `json:"new_timestamp" binding:"omitempty,required_without:NewAmount,NewCurrency,NewDescription"`
}

// TransactionsUpdateHandler updates transaction.
//
//	@summary		update transaction
//	@description	Updates transaction.
//	@tags			transactions
//	@accept			json
//	@produce		json
//	@param			uuid			path		string				true	"UUID of transaction."
//	@param			new_amount		body		integer				false	"New negative or positive amount of transaction in minor units."
//	@param			new_currency	body		string				false	"New currency of transaction in ISO-4217 format."
//	@param			new_description	body		string				false	"New description of transaction."
//	@param			new_timestamp	body		string				false	"New timestamp of transaction in RFC-3339 format."
//	@success		200				{object}	models.Transaction	"Updated transaction object is returned."
//	@failure		404				{object}	models.Error		"Transaction was not found."
//	@failure		403				{object}	models.Error		"You have no right to update the transaction."
//	@router			/transactions/{uuid} [put]
func TransactionsUpdateHandler(c *gin.Context) {
	params := transactionsUpdateParams{}
	if ok := util.BindBody(c, &params); !ok {
		return
	}

	transactionUUID := c.Param("uuid")

	currentUser := c.MustGet("current_user").(*database.User)

	transaction := database.Transaction{}
	if err := database.TransactionsCol.FindOne(
		database.Context, bson.D{{"uuid", transactionUUID}},
	).Decode(&transaction); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			util.AbortWithStatusNotFound(c, errorDescriptionTransactionNotFound)
		} else {
			util.AbortWithStatusInternalServerError(c, err)
		}
		return
	}

	if transaction.OwnerID != currentUser.ID {
		util.AbortWithStatusForbidden(c, errorDescriptionTransactionForbidden)
		return
	}

	var updateQueries bson.D
	if params.NewAmount != nil {
		updateQueries = append(updateQueries, bson.E{Key: "amount", Value: *params.NewAmount})
		transaction.Amount = *params.NewAmount
	}

	if params.NewCurrency != nil {
		updateQueries = append(updateQueries, bson.E{Key: "currency", Value: *params.NewCurrency})
		transaction.Currency = *params.NewCurrency
	}

	if params.NewDescription != nil {
		updateQueries = append(updateQueries, bson.E{Key: "description", Value: *params.NewDescription})
		transaction.Description = *params.NewDescription
	}

	if params.NewTimestamp != nil {
		updateQueries = append(updateQueries, bson.E{Key: "timestamp", Value: *params.NewTimestamp})
		transaction.Timestamp = *params.NewTimestamp
	}

	currentTime := time.Now()
	updateQueries = append(updateQueries, bson.E{Key: "updated_at", Value: currentTime})
	transaction.UpdatedAt = currentTime

	if _, err := database.TransactionsCol.UpdateOne(
		database.Context,
		bson.D{{"uuid", transaction.UUID}},
		bson.D{{"$set", updateQueries}},
	); err != nil {
		util.AbortWithStatusInternalServerError(c, err)
		return
	}

	util.ReturnSuccessfulResponse(c, transaction.APIModel())
}

// TransactionsDeleteHandler deletes transaction.
//
//	@summary		delete transaction
//	@description	Deletes transaction.
//	@tags			transactions
//	@accept			json
//	@produce		json
//	@param			uuid	path		string				true	"UUID of transaction."
//	@success		200		{object}	models.Transaction	"Deleted transaction object is returned."
//	@failure		404		{object}	models.Error		"Transaction was not found."
//	@failure		403		{object}	models.Error		"You have no right to delete the transaction."
//	@router			/transactions/{uuid} [delete]
func TransactionsDeleteHandler(c *gin.Context) {
	transactionUUID := c.Param("uuid")

	currentUser := c.MustGet("current_user").(*database.User)

	// fetch transaction:
	transaction := database.Transaction{}
	if err := database.TransactionsCol.FindOne(
		database.Context,
		bson.D{{"uuid", transactionUUID}},
	).Decode(&transaction); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			util.AbortWithStatusNotFound(c, errorDescriptionTransactionNotFound)
		} else {
			util.AbortWithStatusInternalServerError(c, err)
		}
		return
	}

	// check if the current user owns transaction:
	if transaction.OwnerID != currentUser.ID {
		util.AbortWithStatusForbidden(c, errorDescriptionTransactionForbidden)
		return
	}

	// delete transaction:
	if _, err := database.TransactionsCol.DeleteOne(
		database.Context,
		bson.D{{"uuid", transaction.UUID}},
	); err != nil {
		util.AbortWithStatusInternalServerError(c, err)
		return
	}

	util.ReturnSuccessfulResponse(c, transaction.APIModel())
}
