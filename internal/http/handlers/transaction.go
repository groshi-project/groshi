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
	Amount   float64 `json:"amount" binding:"required"`
	Currency string  `json:"currency" binding:"required"`

	Description string    `json:"description" binding:"omitempty"`
	Date        time.Time `json:"date" binding:"omitempty"`
}

func TransactionCreateHandler(c *gin.Context) {
	params := transactionCreateParams{}
	if err := c.ShouldBindJSON(params); err != nil {
		util.AbortWithErrorMessage(
			c, http.StatusBadRequest, error_messages.ErrorInvalidRequestParams.Error(),
		)
		return
	}

	currentUser := c.MustGet("current_user").(*database.User)

	if params.Date.IsZero() {
		// use the current date as transaction date if the date was not provided
		params.Date = time.Now()
	}

	transaction := database.Transaction{
		ID:   primitive.NewObjectID(),
		UUID: database.GenerateUUID(),

		OwnerUUID: currentUser.UUID,

		Amount:   params.Amount,
		Currency: params.Currency,

		Description: params.Description,
		Date:        params.Date,
	}

	_, err := database.Transactions.InsertOne(database.Context, &transaction)
	if err != nil {
		util.AbortWithErrorMessage(
			c, http.StatusInternalServerError, err.Error(),
		)
		return
	}
	c.JSON(http.StatusOK, gin.H{"uuid": transaction.UUID})
}

type transactionReadParams struct {
	UUID string `json:"uuid" bind:"required"`
}

func TransactionReadHandler(c *gin.Context) {
	params := transactionReadParams{}
	if err := c.ShouldBindJSON(&params); err != nil {
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
			util.AbortWithErrorMessage(c, http.StatusNotFound, error_messages.TransactionNotFound.Error())
		} else {
			util.AbortWithErrorMessage(c, http.StatusInternalServerError, err.Error())
		}
		return
	}
	if transaction.OwnerUUID != currentUser.UUID {
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

	NewAmount   *float64 `json:"new_amount"`
	NewCurrency *string  `json:"new_currency"`

	NewDescription *string    `json:"new_description"`
	NewDate        *time.Time `json:"new_date"`
}

func TransactionUpdateHandler(c *gin.Context) {
	params := transactionUpdateParams{}
	if err := c.ShouldBind(&params); err != nil {
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
			util.AbortWithErrorMessage(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	if transaction.UUID != currentUser.UUID {
		util.AbortWithErrorMessage(
			c,
			http.StatusForbidden,
			error_messages.TransactionDoesNotBelongToYou.Error(),
		)
		return
	}

	//_, err = database.Users.UpdateOne(
	//	database.Context,
	//	bson.D{{"uuid", currentUser.UUID}},
	//	bson.D{{"$set", bson.D{{"username", newUsername}}}},
	//)
	//if err != nil {
	//	util.AbortWithErrorMessage(c, http.StatusInternalServerError, err.Error())
	//	return
	//}
}

func TransactionDeleteHandler(c *gin.Context) {

}
