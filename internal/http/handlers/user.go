package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/jieggii/groshi/internal/loggers"
)

type userCreateParams struct {
	Username     string `json:"username" binding:"required"`
	Password     string `json:"password" binding:"required"`
	BaseCurrency string `json:"base_currency" binding:"required"`
}

func UserCreate(c *gin.Context) {
	var params userCreateParams
	if err := c.Bind(&params); err != nil {
		return
	}
	loggers.Info.Println(params)
}
