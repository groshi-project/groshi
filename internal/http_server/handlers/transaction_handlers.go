package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jieggii/groshi/internal/currency/currency_rates"
	"github.com/jieggii/groshi/internal/database"
	"github.com/jieggii/groshi/internal/http_server/handlers/util"
	"github.com/jieggii/groshi/internal/loggers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"math"
	"time"
)

type transactionCreateParams struct {
	Amount   int    `json:"amount" binding:"required"`
	Currency string `json:"currency" binding:"required,currency"`

	Description *string    `json:"description" binding:"omitempty,description"`
	Date        *time.Time `json:"date"`
}

func TransactionCreateHandler(c *gin.Context) {
	params := transactionCreateParams{}
	if ok := util.BindBody(c, &params); !ok {
		return
	}

	currentUser := c.MustGet("current_user").(*database.User)

	// use empty description as transaction description if description was not provided:
	if params.Description == nil {
		emptyDescription := ""
		params.Description = &emptyDescription
	}

	// use the current date as transaction date if date was not provided:
	if params.Date == nil {
		currentTime := time.Now()
		params.Date = &currentTime
	}

	transaction := database.Transaction{
		ID:   primitive.NewObjectID(),
		UUID: uuid.New().String(),

		OwnerID: currentUser.ID,

		Amount:   params.Amount,
		Currency: params.Currency,

		Description: *params.Description,
		Date:        *params.Date,

		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if _, err := database.TransactionsCol.InsertOne(database.Context, &transaction); err != nil {
		util.AbortWithStatusInternalServerError(c, err)
		return
	}

	util.ReturnSuccessfulResponse(c, gin.H{"uuid": transaction.UUID})
}

// TransactionReadOneHandler returns information about single transaction.
func TransactionReadOneHandler(c *gin.Context) {
	transactionUUID := c.Param("uuid")

	currentUser := c.MustGet("current_user").(*database.User)

	transaction := database.Transaction{}
	if err := database.TransactionsCol.FindOne(
		database.Context, bson.D{{"uuid", transactionUUID}},
	).Decode(&transaction); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			util.AbortWithStatusNotFound(c, APIErrorDetailTransactionNotFound)
		} else {
			util.AbortWithStatusInternalServerError(c, err)
		}
		return
	}
	if transaction.OwnerID != currentUser.ID {
		util.AbortWithStatusForbidden(c, APIErrorDetailTransactionDoesNotBelongToYou)
		return
	}

	util.ReturnSuccessfulResponse(c, transaction.JSON())
}

type transactionReadManyParams struct {
	Since time.Time `form:"since" binding:"required"`

	Before *time.Time `form:"before"`
}

// TransactionReadManyHandler returns all transactions for given period.
func TransactionReadManyHandler(c *gin.Context) {
	params := transactionReadManyParams{}
	if ok := util.BindQuery(c, &params); !ok {
		return
	}
	if params.Before == nil {
		currentDate := time.Now()
		params.Before = &currentDate
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
					{"$gte", params.Since},
					{"$lt", *params.Before},
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

	transactions := make([]gin.H, 0)
	for cursor.Next(database.Context) {
		transaction := database.Transaction{}
		if err := cursor.Decode(&transaction); err != nil {
			util.AbortWithStatusInternalServerError(c, err)
		}
		transactions = append(transactions, transaction.JSON())
	}

	util.ReturnSuccessfulResponse(c, transactions)
}

type transactionReadSummaryParams struct {
	Since    time.Time `form:"since" binding:"required"`
	Currency string    `form:"currency" binding:"required,currency"`

	Before *time.Time `form:"before"`
}

// TransactionReadSummary returns summary (count and sum of transaction)
// for given period in desired currency units.
func TransactionReadSummary(c *gin.Context) {
	params := transactionReadSummaryParams{}
	if ok := util.BindQuery(c, &params); !ok {
		return
	}
	currentUser := c.MustGet("current_user").(*database.User)

	if params.Before == nil { // set current date as end date if end date was not provided
		currentDate := time.Now()
		params.Before = &currentDate
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
					{"$gte", params.Since},
					{"$lte", *params.Before},
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
	}

	intTotal := int(math.Round(income - outcome))
	intIncome := int(math.Round(income))
	intOutcome := int(math.Round(outcome))

	util.ReturnSuccessfulResponse(c, gin.H{
		"currency":           params.Currency,
		"income":             intIncome,
		"outcome":            intOutcome,
		"total":              intTotal,
		"transactions_count": len(transactions),
	})
}

type transactionUpdateParams struct {
	NewAmount      *int       `json:"new_amount" binding:"omitempty,required_without:NewCurrency,NewDescription,NewDate"`
	NewCurrency    *string    `json:"new_currency" binding:"omitempty,currency,required_without:NewAmount,NewDate"`
	NewDescription *string    `json:"new_description" binding:"omitempty,description,required_without:NewAmount,NewCurrency,NewDate"`
	NewDate        *time.Time `json:"new_date" binding:"omitempty,required_without:NewAmount,NewCurrency,NewDescription"`
}

// TransactionUpdateHandler updates transaction.
func TransactionUpdateHandler(c *gin.Context) {
	params := transactionUpdateParams{}
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
			util.AbortWithStatusNotFound(c, APIErrorDetailTransactionNotFound)
		} else {
			util.AbortWithStatusInternalServerError(c, err)
		}
		return
	}

	if transaction.OwnerID != currentUser.ID {
		util.AbortWithStatusForbidden(c, APIErrorDetailTransactionDoesNotBelongToYou)
		return
	}

	var updateQueries bson.D
	if params.NewAmount != nil {
		updateQueries = append(updateQueries, bson.E{Key: "amount", Value: *params.NewAmount})
	}

	if params.NewCurrency != nil {
		updateQueries = append(updateQueries, bson.E{Key: "currency", Value: *params.NewCurrency})
	}

	if params.NewDescription != nil {
		updateQueries = append(updateQueries, bson.E{Key: "description", Value: *params.NewDescription})
	}

	if params.NewDate != nil {
		updateQueries = append(updateQueries, bson.E{Key: "date", Value: *params.NewDate})
	}

	updateQueries = append(updateQueries, bson.E{Key: "updated_at", Value: time.Now()})

	if _, err := database.TransactionsCol.UpdateOne(
		database.Context,
		bson.D{{"uuid", transaction.UUID}},
		bson.D{{"$set", updateQueries}},
	); err != nil {
		util.AbortWithStatusInternalServerError(c, err)
		return
	}

	util.ReturnSuccessfulResponse(c, gin.H{"uuid": transaction.UUID})
}

// TransactionDeleteHandler deletes transaction.
func TransactionDeleteHandler(c *gin.Context) {
	transactionUUID := c.Param("uuid")

	currentUser := c.MustGet("current_user").(*database.User)

	// fetch transaction:
	transaction := database.Transaction{}
	if err := database.TransactionsCol.FindOne(
		database.Context,
		bson.D{{"uuid", transactionUUID}},
	).Decode(&transaction); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			util.AbortWithStatusNotFound(c, APIErrorDetailTransactionNotFound)
		} else {
			util.AbortWithStatusInternalServerError(c, err)
		}
		return
	}

	// check if the current user owns transaction:
	if transaction.OwnerID != currentUser.ID {
		util.AbortWithStatusForbidden(c, APIErrorDetailTransactionDoesNotBelongToYou)
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

	util.ReturnSuccessfulResponse(c, gin.H{"uuid": transaction.UUID})
}
