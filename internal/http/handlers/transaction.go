package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jieggii/groshi/internal/database"
	"github.com/jieggii/groshi/internal/http/error_messages"
	"github.com/jieggii/groshi/internal/http/handlers/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
)

// https://stackoverflow.com/questions/66432222/gin-validation-for-optional-pointer-to-be-uuid
type transactionCreateParams struct {
	Amount   int    `json:"amount" binding:"required"`
	Currency string `json:"currency" binding:"required"`

	Description *string    `json:"description"`
	Date        *time.Time `json:"date"`
}

func TransactionCreateHandler(c *gin.Context) {
	params := transactionCreateParams{}
	if ok := util.BindParams(c, &params); !ok {
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
		UUID: database.GenerateUUID(),

		OwnerID: currentUser.ID,

		Amount:   params.Amount,
		Currency: params.Currency,

		Description: *params.Description,
		Date:        *params.Date,

		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if _, err := database.Transactions.InsertOne(database.Context, &transaction); err != nil {
		util.AbortWithInternalServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"uuid": transaction.UUID})
}

type transactionReadParams struct {
	UUID string `json:"uuid" bind:"required"`
}

func TransactionReadHandler(c *gin.Context) {
	params := transactionReadParams{}
	if ok := util.BindParams(c, &params); !ok {
		return
	}

	currentUser := c.MustGet("currentUser").(*database.User)

	transaction := database.Transaction{}
	if err := database.Transactions.FindOne(
		database.Context, bson.D{{"uuid", params.UUID}},
	).Decode(&transaction); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			util.AbortWithErrorMessage(c, http.StatusNotFound, error_messages.TransactionNotFound.Error())
		} else {
			util.AbortWithInternalServerError(c, err)
		}
		return
	}
	if transaction.OwnerID != currentUser.ID {
		util.AbortWithErrorMessage(
			c,
			http.StatusForbidden,
			error_messages.TransactionDoesNotBelongToYou.Error(),
		)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"uuid":        transaction.UUID,
		"amount":      transaction.Amount,
		"currency":    transaction.Currency,
		"description": transaction.Description,
		"date":        transaction.Date,
	})
}

type transactionUpdateParams struct {
	UUID string `json:"uuid" binding:"required"`

	NewAmount   *int    `json:"new_amount"`
	NewCurrency *string `json:"new_currency"`

	NewDescription *string    `json:"new_description"`
	NewDate        *time.Time `json:"new_date"`
}

func TransactionUpdateHandler(c *gin.Context) {
	params := transactionUpdateParams{}
	if ok := util.BindParams(c, &params); !ok {
		return
	}

	if params.NewAmount == nil &&
		params.NewCurrency == nil &&
		params.NewDescription == nil &&
		params.NewDate == nil {
		util.AbortWithErrorMessage(
			c, http.StatusBadRequest, error_messages.ErrorInvalidRequestParams.Error(),
		)
		return
	}

	currentUser := c.MustGet("currentUser").(*database.User)

	transaction := database.Transaction{}
	if err := database.Transactions.FindOne(
		database.Context, bson.D{{"uuid", params.UUID}},
	).Decode(&transaction); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			util.AbortWithErrorMessage(
				c, http.StatusNotFound, error_messages.TransactionNotFound.Error(),
			)
		} else {
			util.AbortWithInternalServerError(c, err)
		}
		return
	}

	if transaction.OwnerID != currentUser.ID {
		util.AbortWithErrorMessage(
			c,
			http.StatusForbidden,
			error_messages.TransactionDoesNotBelongToYou.Error(),
		)
		return
	}

	// todo: update transaction using only one query to the database
	var updateQueries bson.D

	if params.NewAmount != nil {
		updateQueries = append(updateQueries, bson.E{Key: "amount", Value: *params.NewAmount})
	}

	if params.NewCurrency != nil {
		// todo: validate currency code
		updateQueries = append(updateQueries, bson.E{Key: "currency", Value: *params.NewCurrency})
	}

	if params.NewDescription != nil {
		// todo: validate description
		updateQueries = append(updateQueries, bson.E{Key: "description", Value: *params.NewDescription})
	}

	if params.NewDate != nil {
		updateQueries = append(updateQueries, bson.E{Key: "date", Value: *params.NewDate})
	}

	updateQueries = append(updateQueries, bson.E{Key: "updated_at", Value: time.Now()})

	if _, err := database.Transactions.UpdateOne(
		database.Context,
		bson.D{{"uuid", transaction.UUID}},
		bson.D{{"$set", updateQueries}},
	); err != nil {
		util.AbortWithInternalServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"uuid": transaction.UUID})
}

type transactionDeleteParams struct {
	UUID string `json:"uuid" binding:"required"`
}

func TransactionDeleteHandler(c *gin.Context) {
	params := transactionDeleteParams{}
	if ok := util.BindParams(c, &params); !ok {
		return
	}
	currentUser := c.MustGet("currentUser").(*database.User)

	// fetch transaction:
	transaction := database.Transaction{}
	if err := database.Transactions.FindOne(
		database.Context,
		bson.D{{"uuid", params.UUID}},
	).Decode(&transaction); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			util.AbortWithErrorMessage(
				c, http.StatusNotFound, error_messages.TransactionNotFound.Error(),
			)
		} else {
			util.AbortWithInternalServerError(c, err)
		}
		return
	}

	// check if the current user owns transaction:
	if transaction.OwnerID != currentUser.ID {
		util.AbortWithErrorMessage(
			c, http.StatusNotFound, error_messages.TransactionDoesNotBelongToYou.Error(),
		)
		return
	}

	// delete transaction:
	if _, err := database.Transactions.DeleteOne(
		database.Context,
		bson.D{{"uuid", transaction.UUID}},
	); err != nil {
		util.AbortWithInternalServerError(c, err)
		return
	}

	c.JSON(
		http.StatusOK, gin.H{"uuid": transaction.UUID},
	)
}
