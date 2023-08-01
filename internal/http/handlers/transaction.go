package handlers

import (
	"github.com/gin-gonic/gin"
	"time"
)

// https://stackoverflow.com/questions/66432222/gin-validation-for-optional-pointer-to-be-uuid
type transactionCreateParams struct {
	Amount   float64 `json:"amount" binding:"required"`
	Currency string  `json:"currency" binding:"required"`

	Description *string   `json:"description" binding:"omitempty"` // todo: make optional
	Date        time.Time `json:"date" binding:""`                 // todo: make optional
}

func TransactionCreateHandler(c *gin.Context) {

}

type transactionReadParams struct {
	UUID string `json:"uuid" bind:"required"`
}

func TransactionReadHandler(c *gin.Context) {

}

type transactionUpdateParams struct {
}

func TransactionUpdateHandler(c *gin.Context) {

}

func TransactionDeleteHandler(c *gin.Context) {

}
