package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/jieggii/groshi/internal/database"
	"github.com/jieggii/groshi/internal/http/handlers/utils"
	"net/http"
	"time"
)

type transactionCreateParams struct {
	Amount   float64 `json:"amount" binding:"required"`
	Currency string  `json:"currency" binding:"required,currency"`

	BaseAmount  float64 `json:"base_amount"`
	Description *string `json:"description" binding:"omitempty,description"` // todo: think how to make it optional

	Date time.Time `json:"date" binding:""` // todo
}

// https://stackoverflow.com/questions/66432222/gin-validation-for-optional-pointer-to-be-uuid
func TransactionCreate(c *gin.Context) {
	var params transactionCreateParams
	if err := c.Bind(&params); err != nil {
		return
	}
	currentUser := c.MustGet("currentUser").(*database.User)

	transaction := database.Transaction{
		Date:        params.Date,
		OwnerId:     currentUser.ID,
		Amount:      params.Amount,
		Currency:    params.Currency,
		Description: params.Description,
	}
	_, err := database.DB.NewInsert().Model(&transaction).Exec(database.Ctx)
	if err != nil {
		utils.SendInternalServerErrorResponse(
			c,
			"could not insert new transaction",
			err,
		)
		return
	}
}

type transactionReadParams struct {
	UUID string `json:"uuid" bind:"required"`
}

func TransactionRead(c *gin.Context) {
	var params transactionReadParams
	if err := c.Bind(&params); err != nil {
		return
	}
	currentUser := c.MustGet("currentUser").(*database.User)

	var transaction database.Transaction
	err := database.SelectTransaction(params.UUID).Scan(database.Ctx, &transaction)
	if err != nil { // todo: transaction not found
		utils.SendInternalServerErrorResponse(
			c, "could not fetch transaction", err,
		)
		return
	}

	if transaction.Owner.ID != currentUser.ID { // todo test if this works
		utils.SendErrorResponse(c, http.StatusForbidden, "todo") // todo
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"uuid": transaction.UUID, // todo
		},
	)
}

type transactionUpdateParams struct {
}

func TransactionUpdate(c *gin.Context) {
	_ = c.Query("uuid")
}

func TransactionDelete(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(*database.User)

	uuid := c.Query("uuid")

	var transaction database.Transaction
	err := database.SelectTransaction(uuid).Scan(database.Ctx, &transaction)
	if err != nil {
		utils.SendErrorResponse(
			c, http.StatusNotFound, "transaction not found",
		)
		return
	}

	if transaction.Owner.ID != currentUser.ID {
		utils.SendErrorResponse(
			c, http.StatusForbidden, "", // todo
		)
		return
	}

	_, err = database.DB.NewDelete().Model(transaction).Exec(database.Ctx)
	if err != nil {
		utils.SendInternalServerErrorResponse(
			c, "could not delete user", err,
		)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
