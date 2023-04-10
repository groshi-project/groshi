package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/jieggii/groshi/internal/database"
	"github.com/jieggii/groshi/internal/passhash"
)

type userCreateParams struct {
	Username     string `json:"username" binding:"required,username"`
	Password     string `json:"password" binding:"required,password"`
	BaseCurrency string `json:"base_currency" binding:"required,currency"`
}

func UserCreate(c *gin.Context) {
	var params userCreateParams
	var err error

	if err := c.Bind(&params); err != nil {
		return
	}

	userExists, err := database.SelectUser(params.Username).Exists(database.Ctx)

	if err != nil {
		// todo
	}
	if userExists {
		return // todo
	}

	passwordHash, err := passhash.HashPassword(params.Password)
	if err != nil {
		return // todo
	}

	user := database.User{
		Username:     params.Username,
		Password:     passwordHash,
		BaseCurrency: params.BaseCurrency,
	}

	_, err = database.DB.NewInsert().Model(&user).Exec(database.Ctx)
	if err != nil {
		// todo
	}

}
