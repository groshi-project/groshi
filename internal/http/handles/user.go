package handles

import (
	"errors"
	"fmt"
	"github.com/jieggii/groshi/internal/database"
	"github.com/jieggii/groshi/internal/http/ghttp"
	"github.com/jieggii/groshi/internal/http/ghttp/schema"
	"github.com/jieggii/groshi/internal/http/handles/datatypes"
	"github.com/jieggii/groshi/internal/http/jwt"
	"github.com/jieggii/groshi/internal/passhash"
	"regexp"
)

const usernameRegex = ".+"
const minPasswordLen = 8
const maxPasswordLen = 128

type userCreateRequest struct {
	Username     string             `json:"username"`
	Password     string             `json:"password"`
	BaseCurrency datatypes.Currency `json:"base_currency"`
}

func (p *userCreateRequest) Before() error {
	if p.Username == "" || p.Password == "" || p.BaseCurrency.String == "" {
		return errors.New("`username`, `password` and `base_currency` are required fields")
	}

	// validate username
	usernameMatchesPattern, err := regexp.MatchString(usernameRegex, p.Username) // todo
	if err != nil {
		panic(err) // todo: verify that panic should be called, not error logging
	}
	if !usernameMatchesPattern {
		return errors.New("invalid username format")
	}

	// validate password
	if len(p.Password) < minPasswordLen || len(p.Password) > maxPasswordLen {
		return errors.New(fmt.Sprintf(
			"password must contain from %v to %v symbols", minPasswordLen, maxPasswordLen,
		))
	}

	// validate base currency
	if !p.BaseCurrency.IsKnown {
		return errors.New(schema.UnknownCurrencyErrorDetail)
	}

	return nil
}

type userCreateResponse struct {
	Username string `json:"username"`
}

// UserCreate creates new user.
func UserCreate(request *ghttp.Request, _ *database.User) {
	params := userCreateRequest{}
	if ok := request.Decode(&params); !ok {
		return
	}

	if err := params.Before(); err != nil {
		request.SendClientSideErrorResponse(
			schema.InvalidRequestErrorTag, err.Error(),
		)
		return
	}

	userExists, err := database.UserExists(params.Username)
	if err != nil {
		request.SendServerSideErrorResponse(
			"could not check if user already exists", err,
		)
		return
	}
	if userExists {
		request.SendClientSideErrorResponse(
			schema.ConflictErrorTag, "user with this username already exists",
		)
		return
	}

	passwordHash, err := passhash.HashPassword(params.Password)
	if err != nil {
		request.SendServerSideErrorResponse(
			"could not generate password hash", err,
		)
		return
	}

	user := database.User{
		Username:     params.Username,
		Password:     passwordHash,
		BaseCurrency: params.BaseCurrency.String,
	}

	_, err = database.Db.NewInsert().Model(&user).Exec(database.Ctx)
	if err != nil {
		request.SendServerSideErrorResponse(
			"could not insert new user", err,
		)
		return
	}
	response := userCreateResponse{Username: user.Username}
	request.SendSuccessfulResponse(&response)
}

type userAuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (p *userAuthRequest) Before() error {
	if p.Username == "" || p.Password == "" {
		return errors.New("`username` and `password` are required fields")
	}
	return nil
}

type userAuthResponse struct {
	Token string `json:"token"`
}

// UserAuth authorizes user (generates and returns JWT).
func UserAuth(request *ghttp.Request, _ *database.User) {
	params := userAuthRequest{}
	if ok := request.Decode(&params); !ok {
		return
	}

	if err := params.Before(); err != nil {
		request.SendClientSideErrorResponse(
			schema.InvalidRequestErrorTag, err.Error(),
		)
		return
	}

	user, err := database.FetchUserByUsername(params.Username)
	if err != nil {
		request.SendClientSideErrorResponse(
			schema.ObjectNotFoundErrorTag, "user with this username does not exist",
		)
		return
	}

	if !passhash.ValidatePassword(params.Password, user.Password) {
		request.SendClientSideErrorResponse(
			schema.AccessDeniedErrorTag, "invalid password",
		)
		return
	}

	token, err := jwt.GenerateJWT(params.Username)
	if err != nil {
		request.SendServerSideErrorResponse("could not generate JWT", err)
		return
	}

	response := userAuthResponse{Token: token}
	request.SendSuccessfulResponse(&response)
}

type userReadResponse struct {
	Username     string `json:"username"`
	BaseCurrency string `json:"base_currency"`
}

// UserRead returns information about current user.
func UserRead(request *ghttp.Request, currentUser *database.User) {
	response := userReadResponse{
		Username:     currentUser.Username,
		BaseCurrency: currentUser.BaseCurrency,
	}
	request.SendSuccessfulResponse(&response)
}

type userUpdateRequest struct {
	NewUsername     string             `json:"new_username"`
	NewPassword     string             `json:"new_password"`
	NewBaseCurrency datatypes.Currency `json:"new_base_currency"`
}

func (p *userUpdateRequest) Before() error {
	if p.NewUsername == "" && p.NewPassword == "" && p.NewBaseCurrency.String == "" {
		return errors.New(
			"at least one of these fields is required: " +
				"`new_username`, `new_password` and `new_base_currency`",
		)
	}
	if !p.NewBaseCurrency.IsKnown {
		return errors.New(schema.UnknownCurrencyErrorDetail)
	}
	return nil
}

type userUpdateResponse struct {
	Username           string `json:"username"`
	PasswordWasUpdated bool   `json:"password_was_updated"`
	BaseCurrency       string `json:"base_currency"`
}

// UserUpdate updates current user.
func UserUpdate(request *ghttp.Request, currentUser *database.User) {
	params := userUpdateRequest{}
	if ok := request.Decode(&params); !ok {
		return
	}

	if err := params.Before(); err != nil {
		request.SendClientSideErrorResponse(
			schema.InvalidRequestErrorTag, err.Error(),
		)
		return
	}

	if params.NewUsername != "" {
		newUsernameTaken, err := database.UserExists(params.NewUsername)
		if err != nil {
			request.SendServerSideErrorResponse("could not check if user exists", err)
			return
		}

		if newUsernameTaken {
			request.SendClientSideErrorResponse(
				schema.ConflictErrorTag, "new username chosen by you is already taken",
			)
			return
		}
		currentUser.Username = params.NewUsername
	}

	passwordWasUpdated := false
	if params.NewPassword != "" {
		passwordHash, err := passhash.HashPassword(params.NewPassword)
		if err != nil {
			request.SendServerSideErrorResponse(
				"could not generate password hash", err,
			)
			return
		}
		currentUser.Password = passwordHash
		passwordWasUpdated = true
	}

	if params.NewBaseCurrency.String != "" {
		currentUser.BaseCurrency = params.NewBaseCurrency.String
	}

	if _, err := database.Db.NewUpdate().Model(currentUser).WherePK().Exec(database.Ctx); err != nil {
		request.SendServerSideErrorResponse("could not update user", err)
		return
	}
	response := userUpdateResponse{
		Username:           currentUser.Username,
		PasswordWasUpdated: passwordWasUpdated,
		BaseCurrency:       params.NewBaseCurrency.String,
	}
	request.SendSuccessfulResponse(response)
}

type userDeleteResponse struct {
	Username string `json:"username"`
}

// UserDelete deletes current user.
func UserDelete(request *ghttp.Request, currentUser *database.User) {
	_, err := database.Db.NewDelete().Model(currentUser).WherePK().Exec(database.Ctx)
	if err != nil {
		request.SendServerSideErrorResponse("could not delete user", err)
		return
	}

	response := userDeleteResponse{currentUser.Username}
	request.SendSuccessfulResponse(&response)
}
