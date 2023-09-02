package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/groshi-project/groshi/internal/database"
	"github.com/groshi-project/groshi/internal/handlers/util"
	"github.com/groshi-project/groshi/internal/models"
	"github.com/groshi-project/groshi/internal/password_hashing"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type userCreateParams struct {
	Username string `json:"username" binding:"required,username" example:"pipka500"`
	Password string `json:"password" binding:"required,password"`
}

// UserCreateHandler creates new user.
//
//	@summary		create new user
//	@description	Creates a new user.
//	@tags			user
//	@accept			json
//	@produce		json
//	@param			username	body		string			true	"Username for new user."
//	@param			password	body		string			true	"Password for new user."
//	@success		200			{object}	models.User		"Username of the newly created user is returned."
//	@failure		409			{object}	models.Error	"User with such username already exists."
//	@router			/user [post]
func UserCreateHandler(c *gin.Context) {
	params := userCreateParams{}
	if ok := util.BindBody(c, &params); !ok {
		return
	}

	// check if user already exists:
	err := database.UsersCol.FindOne(database.Context, bson.D{{"username", params.Username}}).Err()
	if err == nil {
		util.AbortWithStatusConflict(c, "user with such username already exists")
		return
	}
	if !errors.Is(err, mongo.ErrNoDocuments) {
		util.AbortWithStatusInternalServerError(c, err)
		return
	}

	// hash user password:
	passwordHash, err := password_hashing.HashPassword(params.Password)
	if err != nil {
		util.AbortWithStatusInternalServerError(c, err)
		return
	}

	// create and insert user:
	user := &database.User{
		ID: primitive.NewObjectID(),

		Username: params.Username,
		Password: passwordHash,

		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	_, err = database.UsersCol.InsertOne(database.Context, user)
	if err != nil {
		util.AbortWithStatusInternalServerError(c, err)
		return
	}
	util.ReturnSuccessfulResponse(c, &models.User{
		Username: user.Username,
	})
}

// UserReadHandler returns information about current user.
//
//	@summary		get information about current user
//	@description	Returns information about current user.
//	@tags			user
//	@accept			json
//	@produce		json
//	@success		200	{object}	models.User	"Information about current user is returned."
//	@router			/user [get]
func UserReadHandler(c *gin.Context) {
	currentUser := c.MustGet("current_user").(*database.User)
	util.ReturnSuccessfulResponse(c, &models.User{
		Username: currentUser.Username,
	})
}

type userUpdateParams struct {
	NewUsername *string `json:"new_username" binding:"required_without=NewPassword,omitempty"`
	NewPassword *string `json:"new_password" binding:"required_without=NewUsername,omitempty"`
}

// UserUpdateHandler updates current user credentials.
//
//	@summary		update current user
//	@description	Updates username and/or password of current user.
//	@tags			user
//	@accept			json
//	@produce		json
//	@success		200	{object}	models.User	"Information about current user is returned."
//	@router			/user [put]
func UserUpdateHandler(c *gin.Context) {
	params := userUpdateParams{}
	if ok := util.BindBody(c, &params); !ok {
		return
	}

	currentUser := c.MustGet("current_user").(*database.User)

	var updateQueries bson.D
	if params.NewUsername != nil {
		newUsername := *params.NewUsername

		// check if user already exists:
		err := database.UsersCol.FindOne(database.Context, bson.D{{"username", newUsername}}).Err()
		if err == nil {
			util.AbortWithStatusConflict(c, "user with such username already exists")
			return
		}
		if !errors.Is(err, mongo.ErrNoDocuments) {
			util.AbortWithStatusInternalServerError(c, err)
			return
		}
		updateQueries = append(updateQueries, bson.E{Key: "username", Value: newUsername})
	}

	if params.NewPassword != nil {
		newPassword := *params.NewPassword
		newPasswordHash, err := password_hashing.HashPassword(newPassword)
		if err != nil {
			util.AbortWithStatusInternalServerError(c, err)
			return
		}
		updateQueries = append(updateQueries, bson.E{Key: "password", Value: newPasswordHash})
	}

	if _, err := database.UsersCol.UpdateOne(
		database.Context,
		bson.D{{"_id", currentUser.ID}},
		bson.D{{"$set", updateQueries}},
	); err != nil {
		util.AbortWithStatusInternalServerError(c, err)
		return
	}

	var displayUsername string
	if params.NewUsername != nil {
		displayUsername = *params.NewUsername
	} else {
		displayUsername = currentUser.Username
	}
	util.ReturnSuccessfulResponse(c, &models.User{
		Username: displayUsername,
	})
}

// UserDeleteHandler deletes current user.
//
//	@summary		delete current user
//	@description	Deletes current user.
//	@tags			user
//	@accept			json
//	@produce		json
//	@success		200	{object}	models.User	"Object of deleted user is returned."
//	@router			/user [delete]
func UserDeleteHandler(c *gin.Context) {
	currentUser := c.MustGet("current_user").(*database.User)
	if _, err := database.UsersCol.DeleteOne(
		database.Context,
		bson.D{{"_id", currentUser.ID}},
	); err != nil {
		util.AbortWithStatusInternalServerError(c, err)
		return
	}

	util.ReturnSuccessfulResponse(c, &models.User{
		Username: currentUser.Username,
	})
}
