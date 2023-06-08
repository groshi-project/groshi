package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jieggii/groshi/internal/database"
	"github.com/jieggii/groshi/internal/hashing"
	"github.com/jieggii/groshi/internal/http/handlers/utils"
	"net/http"
)

type userCreateParams struct {
	Username     string `json:"username" binding:"required,username"`
	Password     string `json:"password" binding:"required,password"`
	BaseCurrency string `json:"base_currency" binding:"required,currency"`
}

func UserCreate(c *gin.Context) {
	var params userCreateParams
	if err := c.Bind(&params); err != nil {
		return
	}

	userExists, err := database.SelectUser(params.Username).Exists(database.Ctx)

	if err != nil {
		utils.SendInternalServerErrorResponse(
			c, "could not check if user exists", err,
		)
		return
	}
	if userExists {
		utils.SendErrorResponse(
			c, http.StatusConflict, "user with such username already exists",
		)
		return
	}

	passwordHash, err := hashing.HashPassword(params.Password)
	if err != nil {
		utils.SendInternalServerErrorResponse(
			c, "could not hash password", err,
		)
		return
	}

	user := database.User{
		Username:     params.Username,
		Password:     passwordHash,
		BaseCurrency: params.BaseCurrency,
	}

	_, err = database.DB.NewInsert().Model(&user).Exec(database.Ctx)
	if err != nil {
		utils.SendInternalServerErrorResponse(
			c, "could not insert new user to the database", err,
		)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

func UserRead(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(*database.User)
	c.JSON(http.StatusOK, gin.H{
		"username":      currentUser.Username,
		"base_currency": currentUser.BaseCurrency,
	})

}

type userUpdateParams struct {
	NewUsername *string `json:"new_username" binding:"omitempty,username"`
	NewPassword *string `json:"new_password" binding:"omitempty,password"`
}

func UserUpdate(c *gin.Context) {
	var params userUpdateParams
	if err := c.Bind(&params); err != nil {
		return
	}
	currentUser := c.MustGet("currentUser").(*database.User)
	fmt.Println(params, currentUser.ID)

	if params.NewUsername == nil && params.NewPassword == nil {
		utils.SendErrorResponse(
			c, http.StatusBadRequest, "at least one of these parameters is required: username, new_password",
		)
		return
	}

	if params.NewUsername != nil {
		exists, err := database.SelectUser(*params.NewUsername).Exists(database.Ctx)
		if err != nil {
			utils.SendInternalServerErrorResponse(
				c, "could not check if user exists", err,
			)
			return
		}
		if exists {
			utils.SendErrorResponse(
				c, http.StatusConflict, "user with such username already exists",
			)
			return
		}
	}

	if params.NewPassword != nil {
		passwordHash, err := hashing.HashPassword(*params.NewPassword)
		if err != nil {
			utils.SendInternalServerErrorResponse(
				c, "could not hash password", err,
			)
			return
		}
		currentUser.Password = passwordHash
	}

	_, err := database.DB.NewUpdate().Model(currentUser).WherePK().Exec(database.Ctx)
	if err != nil {
		utils.SendInternalServerErrorResponse(
			c, "could not update user", err,
		)
		return
	}

}

func UserDelete(c *gin.Context) {
	// todo: also delete transactions owned by user
	currentUser := c.MustGet("currentUser").(*database.User)

	_, err := database.DB.NewDelete().Model(currentUser).Exec(database.Ctx)
	if err != nil {
		utils.SendInternalServerErrorResponse(
			c, "could not delete user from the database", err,
		)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
