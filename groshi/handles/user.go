package handles

import (
	"github.com/jieggii/groshi/groshi/auth"
	"github.com/jieggii/groshi/groshi/database"
	"github.com/jieggii/groshi/groshi/handles/jwt"
	"github.com/jieggii/groshi/groshi/handles/schema"
	"github.com/jieggii/groshi/groshi/handles/util"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type _newUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func UserCreate(
	writer http.ResponseWriter,
	request *http.Request,
	_ httprouter.Params,
	claims *jwt.Claims,
) {
	currentUser := new(database.User)
	err := database.Db.NewSelect().
		Model(currentUser).
		Where("username = ?", claims.Username).
		Scan(database.Ctx)
	if err != nil {
		util.ReturnErrorResponse(
			writer,
			schema.ServerSideError,
			"Could not fetch information about you.",
			nil,
		)
		return
	}

	if !currentUser.IsSuperuser {
		util.ReturnErrorResponse(
			writer,
			schema.ClientSideError,
			"You are not superuser (only superusers are allowed to create new users).",
			nil,
		)
		return
	}

	newUserData := _newUser{}
	if ok := util.DecodeBodyJSON(writer, request, newUserData); !ok {
		return
	}
	passwordHash, err := auth.HashPassword(newUserData.Password)
	if err != nil {
		util.ReturnErrorResponse(
			writer,
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
			writer,
			schema.ServerSideError,
			"Could not check if user already exists.",
			err,
		)
		return
	}
	if newUserExists {
		util.ReturnErrorResponse(writer, schema.ClientSideError, "User already exists", nil)
		return
	}
	if _, err = database.Db.NewInsert().Model(&newUser).Exec(database.Ctx); err != nil {
		util.ReturnErrorResponse(
			writer,
			schema.ServerSideError,
			"Could not create new user",
			err,
		)
		return
	}
	util.ReturnSuccessResponse(writer, nil)
	return
}

func UserRead(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func UserUpdate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func UserDelete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}
