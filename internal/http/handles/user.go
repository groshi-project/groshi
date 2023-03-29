package handles

import (
	"errors"
	"github.com/jieggii/groshi/internal/database"
	"github.com/jieggii/groshi/internal/http/ghttp"
	"github.com/jieggii/groshi/internal/http/ghttp/schema"
	"github.com/jieggii/groshi/internal/http/handles/datatypes"
	"github.com/jieggii/groshi/internal/http/handles/validators"
	"github.com/jieggii/groshi/internal/http/jwt"
	"github.com/jieggii/groshi/internal/passhash"
)

type userCreateRequest struct {
	Username     *string             `json:"username"`
	Password     *string             `json:"password"`
	BaseCurrency *datatypes.Currency `json:"base_currency"`
}

func (p *userCreateRequest) Before() error {
	if p.Username == nil || p.Password == nil || p.BaseCurrency == nil {
		return errors.New(
			schema.MissingRequiredFieldsErrorDetail("username", "password", "base_currency"),
		)
	}

	// validate username
	if err := validators.ValidateUserUsername(*p.Username); err != nil {
		return err
	}

	// validate password
	if err := validators.ValidateUserPassword(*p.Password); err != nil {
		return err
	}

	// validate base currency
	if err := validators.ValidateCurrency(*p.BaseCurrency); err != nil {
		return err
	}

	return nil
}

type userCreateResponse struct{}

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

	userExists, err := database.SelectUser(*params.Username).Exists(database.Ctx)
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

	passwordHash, err := passhash.HashPassword(*params.Password)
	if err != nil {
		request.SendServerSideErrorResponse(
			"could not generate password hash", err,
		)
		return
	}

	user := database.User{
		Username:     *params.Username,
		Password:     passwordHash,
		BaseCurrency: (*params.BaseCurrency).Name,
	}

	_, err = database.Db.NewInsert().Model(&user).Exec(database.Ctx)
	if err != nil {
		request.SendServerSideErrorResponse(
			"could not insert new user", err,
		)
		return
	}
	response := userCreateResponse{}
	request.SendSuccessfulResponse(&response)
}

type userAuthRequest struct {
	Username *string `json:"username"`
	Password *string `json:"password"`
}

func (p *userAuthRequest) Before() error {
	if p.Username == nil || p.Password == nil {
		return errors.New(
			schema.MissingRequiredFieldsErrorDetail("username", "password"),
		)
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

	user, err := database.GetUser(*params.Username)
	if err != nil {
		request.SendClientSideErrorResponse(
			schema.ObjectNotFoundErrorTag, "user with this username does not exist",
		)
		return
	}

	if !passhash.ValidatePassword(*params.Password, user.Password) {
		request.SendClientSideErrorResponse(
			schema.AccessDeniedErrorTag, "invalid password",
		)
		return
	}

	token, err := jwt.GenerateJWT(*params.Username)
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
	NewUsername *string `json:"new_username"`
	NewPassword *string `json:"new_password"`
}

func (p *userUpdateRequest) Before() error {
	if p.NewUsername == nil || p.NewPassword == nil {
		return errors.New(
			schema.AtLeastOneOfFieldsIsRequiredErrorDetail("new_username", "new_password"),
		)
	}

	// validate new username
	if err := validators.ValidateUserUsername(*p.NewUsername); err != nil {
		return err
	}

	// validate new password
	if err := validators.ValidateUserPassword(*p.NewPassword); err != nil {
		return err
	}
	return nil
}

type userUpdateResponse struct{}

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

	if params.NewUsername != nil {
		newUsernameTaken, err := database.SelectUser(*params.NewUsername).Exists(database.Ctx)
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
		currentUser.Username = *params.NewUsername
	}

	if params.NewPassword != nil {
		passwordHash, err := passhash.HashPassword(*params.NewPassword)
		if err != nil {
			request.SendServerSideErrorResponse(
				"could not generate password hash", err,
			)
			return
		}
		currentUser.Password = passwordHash
	}

	if _, err := database.Db.NewUpdate().Model(currentUser).WherePK().Exec(database.Ctx); err != nil {
		request.SendServerSideErrorResponse("could not update user", err)
		return
	}
	response := userUpdateResponse{}
	request.SendSuccessfulResponse(&response)
}

type userDeleteResponse struct{}

// UserDelete deletes current user.
func UserDelete(request *ghttp.Request, currentUser *database.User) {
	_, err := database.Db.NewDelete().Model(currentUser).WherePK().Exec(database.Ctx)
	if err != nil {
		request.SendServerSideErrorResponse("could not delete user", err)
		return
	}

	response := userDeleteResponse{}
	request.SendSuccessfulResponse(&response)
}
