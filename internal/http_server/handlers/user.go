package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jieggii/groshi/internal/database"
	"github.com/jieggii/groshi/internal/http_server/handlers/util"
	"github.com/jieggii/groshi/internal/http_server/password_hashing"
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

// UserCreateHandler creates new user.
func UserCreateHandler(c *gin.Context) {
	params := userCreateParams{}
	if ok := util.BindBody(c, &params); !ok {
		return
	}

	// check if user already exists:
	err := database.UsersCol.FindOne(database.Context, bson.D{{"username", params.Username}}).Err()
	if err == nil {
		util.AbortWithConflictError(c, "user with such username already exists")
		return
	}
	if !errors.Is(err, mongo.ErrNoDocuments) {
		util.AbortWithInternalServerError(c, err)
		return
	}

	// hash user password:
	passwordHash, err := password_hashing.HashPassword(params.Password)
	if err != nil {
		util.AbortWithInternalServerError(c, err)
		return
	}

	// create and insert user:
	user := &database.User{
		ID: primitive.NewObjectID(),
		//UUID: database.GenerateUUID(),

		Username: params.Username,
		Password: passwordHash,

		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	_, err = database.UsersCol.InsertOne(database.Context, user)
	if err != nil {
		util.AbortWithInternalServerError(c, err)
		return
	}
	util.ReturnSuccessfulResponse(c, gin.H{"username": user.Username})
}

// UserReadHandler returns information about current user.
func UserReadHandler(c *gin.Context) {
	currentUser := c.MustGet("current_user").(*database.User)
	util.ReturnSuccessfulResponse(c, gin.H{"username": currentUser.Username})
}

type userUpdateParams struct {
	NewUsername *string `json:"new_username" binding:"omitempty"`
	NewPassword *string `json:"new_password" binding:"omitempty"`
}

// UserUpdateHandler updates current user credentials.
func UserUpdateHandler(c *gin.Context) {
	params := userUpdateParams{}
	if ok := util.BindBody(c, &params); !ok {
		return
	}

	currentUser := c.MustGet("current_user").(*database.User)

	// check if no update params were provided:
	if params.NewUsername == nil && params.NewPassword == nil {
		util.AbortWithBadRequest(
			c, "at least one of the following params is required: `new_username` or `new_password`",
		)
		return
	}

	// todo: update user using only one query to the database
	var updateQueries bson.D

	if params.NewUsername != nil {
		newUsername := *params.NewUsername

		// check if user already exists:
		err := database.UsersCol.FindOne(database.Context, bson.D{{"username", newUsername}}).Err()
		if err == nil {
			util.AbortWithErrorMessage(c, http.StatusConflict, "user with such username already exists")
			return
		}
		if !errors.Is(err, mongo.ErrNoDocuments) {
			util.AbortWithInternalServerError(c, err)
			return
		}
		updateQueries = append(updateQueries, bson.E{Key: "username", Value: newUsername})
	}

	if params.NewPassword != nil {
		newPassword := *params.NewPassword
		newPasswordHash, err := password_hashing.HashPassword(newPassword)
		if err != nil {
			util.AbortWithInternalServerError(c, err)
			return
		}
		updateQueries = append(updateQueries, bson.E{Key: "password", Value: newPasswordHash})
	}

	if _, err := database.UsersCol.UpdateOne(
		database.Context,
		bson.D{{"_id", currentUser.ID}},
		bson.D{{"$set", updateQueries}},
	); err != nil {
		util.AbortWithInternalServerError(c, err)
		return
	}
	util.ReturnSuccessfulResponse(c, gin.H{})
}

// UserDeleteHandler deletes current user.
func UserDeleteHandler(c *gin.Context) {
	currentUser := c.MustGet("current_user").(*database.User)
	if _, err := database.UsersCol.DeleteOne(
		database.Context,
		bson.D{{"_id", currentUser.ID}},
	); err != nil {
		util.AbortWithInternalServerError(c, err)
		return
	}

	util.ReturnSuccessfulResponse(c, gin.H{"username": currentUser.Username})
}
