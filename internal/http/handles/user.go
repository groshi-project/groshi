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
	return len(p.Username)*len(p.Password) != 0 // todo: use optimal condition
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
		request.SendClientSideErrorResponse(schema.InvalidRequestBody)
		return
	}

	user, err := database.FetchUserByUsername(params.Username)
	if err != nil {
		request.SendClientSideErrorResponse(schema.UserNotFound)
		return
	}

	if !passhash.CheckPasswordHash(params.Password, user.Password) {
		request.SendClientSideErrorResponse("Invalid password.") // todo?
		return
	}

	token, err := jwt.GenerateJWT(params.Username)
	if err != nil {
		request.SendServerSideErrorResponse("Could not generate JWT.", err)
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

type userCreateResponse struct{}

func UserCreate(request *ghttp.Request, currentUser *database.User) {
	if !currentUser.IsSuperuser {
		request.SendClientSideErrorResponse(schema.AccessDenied)
		return
	}

	params := userCreateRequest{}
	if ok := request.DecodeSafe(&params); !ok {
		return
	}

	userExists, err := database.UserExists(params.Username)
	if err != nil {
		request.SendServerSideErrorResponse(
			"Could not check if user already exists.", err,
		)
		return
	}
	if userExists {
		request.SendClientSideErrorResponse(
			"User with such username already exists.",
		)
		return
	}

	passwordHash, err := passhash.HashPassword(params.Password)
	if err != nil {
		request.SendServerSideErrorResponse(
			"Could not generate password hash.", err,
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
			"Could not insert new user.", err,
		)
		return
	}
	response := userCreateResponse{}
	request.SendSuccessResponse(&response)
}

type userReadRequest struct {
	Username string `json:"username"`
}

type userReadResponse struct {
	Username    string `json:"username"`
	IsSuperuser bool   `json:"is_superuser"`
}

func UserRead(request *ghttp.Request, _ *database.User) {
	params := userReadRequest{}
	if ok := request.DecodeSafe(&params); !ok {
		return
	}

	user, err := database.FetchUserByUsername(params.Username)
	if err != nil {
		request.SendClientSideErrorResponse(schema.UserNotFound)
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

type userUpdateResponse struct {
}

func UserUpdate(request *ghttp.Request, currentUser *database.User) {
	params := userUpdateRequest{}
	if ok := request.DecodeSafe(&params); !ok {
		return
	}

	if params.Username != currentUser.Username && !currentUser.IsSuperuser {
		request.SendClientSideErrorResponse(schema.AccessDenied)
		return
	}

	var user *database.User

	if params.Username == currentUser.Username {
		user = currentUser
	} else {
		var err error
		user, err = database.FetchUserByUsername(params.Username)
		if err != nil {
			request.SendClientSideErrorResponse(schema.UserNotFound)
			return
		}
	}

	if params.NewUsername != "" {
		newUsernameTaken, err := database.UserExists(params.NewUsername)
		if err != nil {
			request.SendServerSideErrorResponse("Could not check if user exists.", err)
			return
		}
		if newUsernameTaken {
			request.SendClientSideErrorResponse("New username is already taken.")
			return
		}
		user.Username = params.NewUsername
	}

	if params.NewPassword != "" {
		passwordHash, err := passhash.HashPassword(params.NewPassword)
		if err != nil {
			request.SendServerSideErrorResponse(
				"Could not generate password hash.", err,
			)
			return
		}
		user.Password = passwordHash
	}

	var newIsSuperuser bool
	if params.Promote || params.Demote {
		if !currentUser.IsSuperuser {
			request.SendClientSideErrorResponse(
				"You are not allowed to affect `promote` and `demote` fields.",
			) // todo
			return
		}
		if params.Promote && params.Demote {
			request.SendClientSideErrorResponse(
				"`promote` and `demote` fields cannot be used at once.",
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

	_, err := database.Db.NewUpdate().Model(&user).WherePK().Exec(database.Ctx)
	if err != nil {
		request.SendServerSideErrorResponse("Could not update user.", err)
		return
	}
	response := userUpdateResponse{}
	request.SendSuccessResponse(&response)
}

type userDeleteRequest struct {
	Username string `json:"username"`
}
type userDeleteResponse struct{}

func UserDelete(request *ghttp.Request, currentUser *database.User) {
	params := userDeleteRequest{}
	if ok := request.DecodeSafe(&params); !ok {
		return
	}

	if currentUser.Username != params.Username && !currentUser.IsSuperuser {
		request.SendClientSideErrorResponse(schema.AccessDenied)
		return
	}
	user, err := database.FetchUserByUsername(params.Username)
	if err != nil {
		request.SendClientSideErrorResponse(schema.UserNotFound)
		return
	}

	_, err = database.Db.NewDelete().Model(user).WherePK().Exec(database.Ctx)
	if err != nil {
		request.SendServerSideErrorResponse("Could not delete user.", err)
		return
	}

	response := userDeleteResponse{}
	request.SendSuccessResponse(&response)
}
