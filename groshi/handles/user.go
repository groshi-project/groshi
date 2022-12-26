package handles

import (
	"github.com/jieggii/groshi/groshi/auth"
	"github.com/jieggii/groshi/groshi/database"
	"github.com/jieggii/groshi/groshi/handles/jwt"
	"github.com/jieggii/groshi/groshi/handles/schema"
	"github.com/jieggii/groshi/groshi/handles/util"
	"net/http"
)

type _newUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func UserCreate(w http.ResponseWriter, r *http.Request, claims *jwt.Claims) {
	currentUser := new(database.User)
	err := database.Db.NewSelect().
		Model(currentUser).
		Where("username = ?", claims.Username).
		Scan(database.Ctx)
	if err != nil {
		util.ReturnErrorResponse(
			w,
			schema.ServerSideError,
			"Could not fetch information about you.",
			nil,
		)
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

	newUserData := _newUser{}
	if ok := util.DecodeBodyJSON(w, r, newUserData); !ok {
		return
	}
	passwordHash, err := auth.HashPassword(newUserData.Password)
	if err != nil {
		util.ReturnErrorResponse(
			w,
			schema.ServerSideError,
			"Could not generate password hash.",
			err,
		)
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
		util.ReturnErrorResponse(
			w,
			schema.ServerSideError,
			"Could not check if user already exists.",
			err,
		)
		return
	}
	if newUserExists {
		util.ReturnErrorResponse(w, schema.ClientSideError, "User already exists", nil)
		return
	}
	if _, err = database.Db.NewInsert().Model(&newUser).Exec(database.Ctx); err != nil {
		util.ReturnErrorResponse(
			w,
			schema.ServerSideError,
			"Could not create new user",
			err,
		)
		return
	}
	util.ReturnSuccessResponse(w, nil)
	return
}

func UserRead(w http.ResponseWriter, r *http.Request, claims *jwt.Claims) {

}

func UserUpdate(w http.ResponseWriter, r *http.Request, claims *jwt.Claims) {

}

func UserDelete(w http.ResponseWriter, r *http.Request, claims *jwt.Claims) {

}
