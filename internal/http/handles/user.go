package handles

import (
	"github.com/jieggii/groshi/internal/database"
	"github.com/jieggii/groshi/internal/http/ghttp"
	"github.com/jieggii/groshi/internal/http/jwt"
	"github.com/jieggii/groshi/internal/http/schema"
	"github.com/jieggii/groshi/internal/passhash"
)

type userAuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (p *userAuthRequest) validate() bool {
	return p.Username != "" && p.Password != ""
}

type userAuthResponse struct {
	Token string `json:"token"`
}

func UserAuth(request *ghttp.Request, _ *database.User) {
	params := userAuthRequest{}
	if ok := request.DecodeSafe(&params); !ok {
		return
	}

	if !params.validate() {
		request.SendClientSideErrorResponse(
			schema.InvalidRequestErrorTag, schema.RequestBodyDidNotPassValidation,
		)
		return
	}

	user, err := database.FetchUserByUsername(params.Username)
	if err != nil {
		request.SendClientSideErrorResponse(
			schema.ObjectNotFoundErrorTag, schema.UserNotFoundErrorDetail,
		)
		return
	}

	if !passhash.CheckPasswordHash(params.Password, user.Password) {
		request.SendClientSideErrorResponse(
			schema.AccessDeniedErrorTag, "Invalid password.",
		)
		return
	}

	token, err := jwt.GenerateJWT(params.Username)
	if err != nil {
		request.SendServerSideErrorResponse("could not generate JWT", err)
		return
	}
	response := userAuthResponse{Token: token}
	request.SendSuccessResponse(&response)
}

type userCreateRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`

	IsSuperUser bool `json:"is_superuser"`
}

func (p *userCreateRequest) validate() bool {
	return p.Username != "" && p.Password != ""
}

type userCreateResponse struct {
	Username string `json:"username"`
}

func UserCreate(request *ghttp.Request, currentUser *database.User) {
	if !currentUser.IsSuperuser {
		request.SendClientSideErrorResponse(
			schema.AccessDeniedErrorTag, schema.NoRightToPerformOperationErrorDetail,
		)
		return
	}

	params := userCreateRequest{}
	if ok := request.DecodeSafe(&params); !ok {
		return
	}

	if !params.validate() {
		request.SendClientSideErrorResponse(
			schema.InvalidRequestErrorTag, schema.RequestBodyDidNotPassValidation,
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
			schema.ConflictErrorTag, "User with this username already exists.",
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
		Username:    params.Username,
		Password:    passwordHash,
		IsSuperuser: params.IsSuperUser,
	}

	_, err = database.Db.NewInsert().Model(&user).Exec(database.Ctx)
	if err != nil {
		request.SendServerSideErrorResponse(
			"could not insert new user", err,
		)
		return
	}
	response := userCreateResponse{Username: params.Username}
	request.SendSuccessResponse(&response)
}

type userReadRequest struct {
	Username string `json:"username"`
}

func (p *userReadRequest) validate() bool {
	return p.Username != ""
}

type userReadResponse struct {
	Username    string `json:"username"`
	IsSuperuser bool   `json:"is_superuser"`
}

func UserRead(request *ghttp.Request, currentUser *database.User) {
	params := userReadRequest{}
	if ok := request.DecodeSafe(&params); !ok {
		return
	}

	if !params.validate() {
		request.SendClientSideErrorResponse(
			schema.InvalidRequestErrorTag, schema.RequestBodyDidNotPassValidation,
		)
		return
	}

	if currentUser.Username != params.Username && !currentUser.IsSuperuser {
		request.SendClientSideErrorResponse(
			schema.AccessDeniedErrorTag, schema.NoRightToPerformOperationErrorDetail,
		)
		return
	}

	user, err := database.FetchUserByUsername(params.Username)
	if err != nil {
		request.SendClientSideErrorResponse(
			schema.ObjectNotFoundErrorTag, schema.UserNotFoundErrorDetail,
		)
		return
	}
	response := userReadResponse{
		Username:    user.Username,
		IsSuperuser: user.IsSuperuser,
	}
	request.SendSuccessResponse(&response)
}

type userUpdateRequest struct {
	Username string `json:"username"`

	NewUsername string `json:"new_username"`
	NewPassword string `json:"new_password"`

	Promote bool `json:"promote"`
	Demote  bool `json:"demote"`
}

func (p *userUpdateRequest) validate() bool {
	return p.Username != ""
}

type userUpdateResponse struct {
	Username        string `json:"username"`
	PasswordUpdated bool   `json:"password_updated"`
	IsSuperUser     bool   `json:"is_superuser"`
}

func UserUpdate(request *ghttp.Request, currentUser *database.User) {
	params := userUpdateRequest{}
	if ok := request.DecodeSafe(&params); !ok {
		return
	}

	if !params.validate() {
		request.SendClientSideErrorResponse(
			schema.InvalidRequestErrorTag, schema.RequestBodyDidNotPassValidation,
		)
		return
	}

	if params.Username != currentUser.Username && !currentUser.IsSuperuser {
		request.SendClientSideErrorResponse(
			schema.AccessDeniedErrorTag, schema.NoRightToPerformOperationErrorDetail,
		)
		return
	}

	var user *database.User

	if params.Username == currentUser.Username {
		user = currentUser
	} else {
		var err error
		user, err = database.FetchUserByUsername(params.Username)
		if err != nil {
			request.SendClientSideErrorResponse(
				schema.ObjectNotFoundErrorTag, schema.UserNotFoundErrorDetail,
			)
			return
		}
	}

	if params.NewUsername != "" {
		newUsernameTaken, err := database.UserExists(params.NewUsername)
		if err != nil {
			request.SendServerSideErrorResponse("could not check if user exists", err)
			return
		}

		if newUsernameTaken {
			request.SendClientSideErrorResponse(
				schema.ConflictErrorTag, "New username chosen by you is already taken.",
			)
			return
		}
		user.Username = params.NewUsername
	}

	passwordUpdated := false
	if params.NewPassword != "" {
		passwordHash, err := passhash.HashPassword(params.NewPassword)
		if err != nil {
			request.SendServerSideErrorResponse(
				"could not generate password hash", err,
			)
			return
		}
		user.Password = passwordHash
		passwordUpdated = true
	}

	var newIsSuperuser bool
	if params.Promote || params.Demote {
		if !currentUser.IsSuperuser {
			request.SendClientSideErrorResponse(
				schema.AccessDeniedErrorTag,
				"You have no right to use `promote` and `demote` fields.",
			)
			return
		}
		if params.Promote && params.Demote {
			request.SendClientSideErrorResponse(
				schema.InvalidRequestErrorTag,
				"`promote` and `demote` fields cannot both be used at once.",
			)
			return
		}
		if params.Promote {
			newIsSuperuser = true
		} else if params.Demote {
			newIsSuperuser = false
		}
		user.IsSuperuser = newIsSuperuser
	}

	_, err := database.Db.NewUpdate().Model(user).WherePK().Exec(database.Ctx)
	if err != nil {
		request.SendServerSideErrorResponse("could not update user", err)
		return
	}
	response := userUpdateResponse{
		Username:        user.Username,
		PasswordUpdated: passwordUpdated,
		IsSuperUser:     user.IsSuperuser,
	}
	request.SendSuccessResponse(&response)
}

type userDeleteRequest struct {
	Username string `json:"username"`
}

func (p *userDeleteRequest) validate() bool {
	return p.Username != ""
}

type userDeleteResponse struct {
	Username string `json:"username"`
}

func UserDelete(request *ghttp.Request, currentUser *database.User) {
	params := userDeleteRequest{}
	if ok := request.DecodeSafe(&params); !ok {
		return
	}

	if !params.validate() {
		request.SendClientSideErrorResponse(
			schema.InvalidRequestErrorTag, schema.RequestBodyDidNotPassValidation,
		)
		return
	}

	if currentUser.Username != params.Username && !currentUser.IsSuperuser {
		request.SendClientSideErrorResponse(
			schema.AccessDeniedErrorTag, schema.NoRightToPerformOperationErrorDetail,
		)
		return
	}

	user, err := database.FetchUserByUsername(params.Username)
	if err != nil {
		request.SendClientSideErrorResponse(
			schema.ObjectNotFoundErrorTag, schema.UserNotFoundErrorDetail,
		)
		return
	}

	_, err = database.Db.NewDelete().Model(user).WherePK().Exec(database.Ctx)
	if err != nil {
		request.SendServerSideErrorResponse("could not delete user", err)
		return
	}

	response := userDeleteResponse{Username: user.Username}
	request.SendSuccessResponse(&response)
}
