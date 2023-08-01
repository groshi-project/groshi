package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jieggii/groshi/internal/database"
	"github.com/jieggii/groshi/internal/http/error_messages"
	utils "github.com/jieggii/groshi/internal/http/handlers/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

type userCreateParams struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func UserCreateHandler(c *gin.Context) {
	params := userCreateParams{}
	if err := c.ShouldBindJSON(&params); err != nil {
		utils.AbortWithErrorMessage(c, http.StatusBadRequest, error_messages.ErrorInvalidRequestParams.Error())
		return
	}
	err := database.Users.FindOne(database.Context, bson.D{{"username", params.Username}}).Err()
	if err == nil {
		utils.AbortWithErrorMessage(c, http.StatusConflict, "user with such username already exists")
		return
	}
	if !errors.Is(err, mongo.ErrNoDocuments) {
		utils.AbortWithErrorMessage(
			c, http.StatusInternalServerError, err.Error(),
		)
		return
	}

	passwordHashBytes, err := bcrypt.GenerateFromPassword([]byte(params.Password), 12) // todo const
	if err != nil {
		utils.AbortWithErrorMessage(c, http.StatusInternalServerError, err.Error())
	}
	passwordHash := string(passwordHashBytes)

	if err != nil {
		utils.AbortWithErrorMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	user := &database.User{
		UUID:      database.GenerateUUID(),
		Username:  params.Username,
		Password:  passwordHash,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	_, err = database.Users.InsertOne(database.Context, user)
	if err != nil {
		utils.AbortWithErrorMessage(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"username": user.Username})
}

func UserReadHandler(c *gin.Context) {

}

type userUpdateParams struct {
	NewUsername *string `json:"new_username" binding:"omitempty"`
	NewPassword *string `json:"new_password" binding:"omitempty"`
}

func UserUpdateHandler(c *gin.Context) {

}

func UserDeleteHandler(c *gin.Context) {

}
