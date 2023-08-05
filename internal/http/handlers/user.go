package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jieggii/groshi/internal/database"
	"github.com/jieggii/groshi/internal/http/error_messages"
	"github.com/jieggii/groshi/internal/http/handlers/util"
	"github.com/jieggii/groshi/internal/http/password_hashing"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
		util.AbortWithErrorMessage(c, http.StatusBadRequest, error_messages.ErrorInvalidRequestParams.Error())
		return
	}
	err := database.Users.FindOne(database.Context, bson.D{{"username", params.Username}}).Err()
	if err == nil {
		util.AbortWithErrorMessage(c, http.StatusConflict, "user with such username already exists")
		return
	}
	if !errors.Is(err, mongo.ErrNoDocuments) {
		util.AbortWithErrorMessage(
			c, http.StatusInternalServerError, err.Error(),
		)
		return
	}

	passwordHash, err := password_hashing.HashPassword(params.Password)
	if err != nil {
		util.AbortWithErrorMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	if err != nil {
		util.AbortWithErrorMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	user := &database.User{
		ID:   primitive.NewObjectID(),
		UUID: database.GenerateUUID(),

		Username: params.Username,
		Password: passwordHash,

		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	_, err = database.Users.InsertOne(database.Context, user)
	if err != nil {
		util.AbortWithErrorMessage(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"username": user.Username})
}

func UserReadHandler(c *gin.Context) {
	currentUser := c.MustGet("current_user").(*database.User)
	c.JSON(
		http.StatusOK, gin.H{"username": currentUser.Username},
	)
}

type userUpdateParams struct {
	NewUsername *string `json:"new_username" binding:"omitempty"`
	NewPassword *string `json:"new_password" binding:"omitempty"`
}

func UserUpdateHandler(c *gin.Context) {
	currentUser := c.MustGet("current_user").(*database.User)
	params := userUpdateParams{}
	if err := c.ShouldBind(&params); err != nil {
		util.AbortWithErrorMessage(
			c, http.StatusBadRequest, error_messages.ErrorInvalidRequestParams.Error(),
		)
		return
	}

	if params.NewUsername == nil && params.NewPassword == nil {
		util.AbortWithErrorMessage(c, http.StatusBadRequest, "TODO")
		return
	}

	if params.NewUsername != nil {
		newUsername := *params.NewUsername
		// check if user already exists
		err := database.Users.FindOne(database.Context, bson.D{{"username", newUsername}}).Err()
		if err == nil {
			util.AbortWithErrorMessage(c, http.StatusConflict, "user with such username already exists")
			return
		}
		if !errors.Is(err, mongo.ErrNoDocuments) {
			util.AbortWithErrorMessage(c, http.StatusInternalServerError, err.Error())
			return
		}

		_, err = database.Users.UpdateOne(
			database.Context,
			bson.D{{"uuid", currentUser.UUID}},
			bson.D{{"$set", bson.D{{"username", newUsername}}}},
		)
		if err != nil {
			util.AbortWithErrorMessage(c, http.StatusInternalServerError, err.Error())
			return
		}
	}

	if params.NewPassword != nil {
		newPassword := *params.NewPassword
		newPasswordHash, err := password_hashing.HashPassword(newPassword)
		if err != nil {
			util.AbortWithErrorMessage(c, http.StatusInternalServerError, err.Error())
			return
		}
		_, err = database.Users.UpdateOne(
			database.Context,
			bson.D{{"uuid", currentUser.UUID}},
			bson.D{{"$set", bson.D{{"password", newPasswordHash}}}},
		)
	}

	c.JSON(http.StatusOK, gin.H{})
}

func UserDeleteHandler(c *gin.Context) {
	currentUser := c.MustGet("current_user").(*database.User)
	_, err := database.Users.DeleteOne(
		database.Context,
		bson.D{{"uuid", currentUser.UUID}},
	)
	if err != nil {
		util.AbortWithErrorMessage(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(
		http.StatusOK, gin.H{"username": currentUser.Username},
	)
}
