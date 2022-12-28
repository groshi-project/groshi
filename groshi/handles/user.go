package handles

import (
	"github.com/jieggii/groshi/groshi/auth"
	"github.com/jieggii/groshi/groshi/database"
	"github.com/jieggii/groshi/groshi/handles/jwt"
	"github.com/jieggii/groshi/groshi/handles/schema"
	"github.com/jieggii/groshi/groshi/handles/util"
	"net/http"
)

type _userCreateParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type _userCreateResponse struct {
	Created bool `json:"created"`
}

func UserCreate(w http.ResponseWriter, r *http.Request, claims *jwt.Claims) {
	currentUser, err := database.FetchUserByUsername(claims.Username)
	if err != nil {
		util.ReturnErrorResponse(w, schema.ClientSideError, schema.UserNotFound, nil)
		return
	}

	if !currentUser.IsSuperuser {
		util.ReturnErrorResponse(
			w,
			schema.ClientSideError,
			"You are not superuser (only superusers are allowed to create new users).",
			nil,
		)
		return
	}

	newUserData := _userCreateParams{}
	if ok := util.DecodeBodyJSON(w, r, newUserData); !ok {
		return
	}
	passwordHash, err := auth.HashPassword(newUserData.Password)
	if err != nil {
		util.ReturnErrorResponse(w, schema.ServerSideError, "Could not generate password hash.", err)
		return
	}

	newUser := database.User{
		Username: newUserData.Username,
		Password: passwordHash,
	}

	newUserExists, err := database.Db.NewSelect().
		Model((*database.User)(nil)).
		Where("username = ?", newUserData.Username).
		Exists(database.Ctx)
	if err != nil {
		util.ReturnErrorResponse(w, schema.ServerSideError, "Could not check if user already exists.", err)
		return
	}
	if newUserExists {
		util.ReturnErrorResponse(w, schema.ClientSideError, "User already exists", nil)
		return
	}
	_, err = database.Db.NewInsert().Model(&newUser).Exec(database.Ctx)
	if err != nil {
		util.ReturnErrorResponse(w, schema.ServerSideError, "Could not create new user.", err)
		return
	}
	response := _userCreateResponse{Created: true}
	util.ReturnSuccessResponse(w, &response)
	return
}

type _userReadParams struct {
	Username string `json:"username"`
}

type _userReadResponse struct {
	Username    string `json:"username"`
	IsSuperuser bool   `json:"is_superuser"`
}

func UserRead(w http.ResponseWriter, r *http.Request, claims *jwt.Claims) {
	currentUser, err := database.FetchUserByUsername(claims.Username)
	if err != nil {
		util.ReturnErrorResponse(w, schema.ClientSideError, schema.UserNotFound, nil)
		return
	}
	params := _userReadParams{}
	if ok := util.DecodeBodyJSON(w, r, &params); !ok {
		return
	}
	if currentUser.Username != params.Username && !currentUser.IsSuperuser {
		util.ReturnErrorResponse(w, schema.ClientSideError, schema.AccessDenied, nil)
		return
	}
	targetUser, err := database.FetchUserByUsername(params.Username)
	if err != nil {
		util.ReturnErrorResponse(w, schema.ClientSideError, schema.UserNotFound, nil)
	}
	response := _userReadResponse{
		Username:    targetUser.Username,
		IsSuperuser: targetUser.IsSuperuser,
	}
	util.ReturnSuccessResponse(w, &response)
}

func UserUpdate(w http.ResponseWriter, r *http.Request, claims *jwt.Claims) {

}

type _userDeleteParams = _userReadParams // todo
type _userDeleteResponse struct {
	Deleted bool `json:"deleted"`
}

func UserDelete(w http.ResponseWriter, r *http.Request, claims *jwt.Claims) {
	currentUser, err := database.FetchUserByUsername(claims.Username)
	if err != nil {
		util.ReturnErrorResponse(w, schema.ClientSideError, schema.UserNotFound, nil)
		return
	}
	params := _userDeleteParams{}
	if ok := util.DecodeBodyJSON(w, r, &params); !ok {
		return
	}
	if currentUser.Username != params.Username && !currentUser.IsSuperuser {
		util.ReturnErrorResponse(w, schema.ClientSideError, schema.AccessDenied, nil)
		return
	}
	targetUser, err := database.FetchUserByUsername(params.Username)
	if err != nil {
		util.ReturnErrorResponse(w, schema.ClientSideError, schema.UserNotFound, nil)
		return
	}
	_, err = database.Db.NewDelete().Model(targetUser).Exec(database.Ctx)
	if err != nil {
		util.ReturnErrorResponse(w, schema.ServerSideError, "Could not delete user.", err)
		return
	}
	response := _userDeleteResponse{Deleted: true}
	util.ReturnSuccessResponse(w, &response)
}
