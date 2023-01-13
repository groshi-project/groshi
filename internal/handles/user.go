package handles

import (
	"github.com/jieggii/groshi/internal/auth/jwt"
	"github.com/jieggii/groshi/internal/auth/passwords"
	"github.com/jieggii/groshi/internal/database"
	"github.com/jieggii/groshi/internal/ghttp"
	"github.com/jieggii/groshi/internal/handles/schema"
)

type userAuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type userAuthResponse struct {
	Token string `json:"token"`
}

func UserAuth(request *ghttp.Request, _ *database.User) {
	var credentials userAuthRequest
	if ok := request.DecodeSafe(&credentials); !ok {
		return
	}
	if len(credentials.Username) == 0 || len(credentials.Password) == 0 {
		request.SendErrorResponse(schema.ClientSideError, schema.InvalidRequestBody, nil)
		return
	}
	user := new(database.User)
	err := database.Db.NewSelect().Model(user).Where("username = ?", credentials.Username).Scan(database.Ctx)
	if err != nil {
		request.SendErrorResponse(schema.ClientSideError, schema.UserNotFound, nil)
		return
	}

	if !passwords.CheckPasswordHash(credentials.Password, user.Password) {
		request.SendErrorResponse(schema.ClientSideError, "Invalid password.", nil)
		return
	}

	token, err := jwt.GenerateJWT(credentials.Username)
	if err != nil {
		request.SendErrorResponse(schema.ServerSideError, "Could not generate JWT.", err)
		return
	}
	response := userAuthResponse{Token: token}
	request.SendSuccessResponse(&response)
}

type userCreateRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type userCreateResponse struct{}

func UserCreate(request *ghttp.Request, currentUser *database.User) {

	if !currentUser.IsSuperuser {
		request.SendErrorResponse(
			schema.ClientSideError,
			"You are not superuser (only superusers are allowed to create new users).",
			nil,
		)
		return
	}

	newUserData := userCreateRequest{}
	if ok := request.DecodeSafe(&newUserData); !ok {
		return
	}

	passwordHash, err := passwords.HashPassword(newUserData.Password)
	if err != nil {
		request.SendErrorResponse(schema.ServerSideError, "Could not generate password hash.", err)
		return
	}

	newUser := database.User{
		Username: newUserData.Username,
		Password: passwordHash,
	}

	newUserExists, err := database.Db.NewSelect().
		Model(newUser).
		Where("username = ?", newUserData.Username).
		Exists(database.Ctx)
	if err != nil {
		request.SendErrorResponse(schema.ServerSideError, "Could not check if user already exists.", err)
		return
	}

	if newUserExists {
		request.SendErrorResponse(schema.ClientSideError, "User already exists.", nil)
		return
	}

	_, err = database.Db.NewInsert().Model(&newUser).Exec(database.Ctx)
	if err != nil {
		request.SendErrorResponse(schema.ServerSideError, "Could not create new user.", err)
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
	//if ok := request.WrapCondition(currentUser.Username == params.Username || currentUser.IsSuperuser, ""); !ok {
	//	return
	//}
	targetUser, err := database.FetchUserByUsername(params.Username)
	if err != nil {
		request.SendErrorResponse(schema.ClientSideError, schema.UserNotFound, nil)
	}
	response := userReadResponse{
		Username:    targetUser.Username,
		IsSuperuser: targetUser.IsSuperuser,
	}
	request.SendSuccessResponse(&response)
}

type userUpdateRequest struct {
}

type userUpdateResponse struct {
}

func UserUpdate(request *ghttp.Request, currentUser *database.User) {

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
		request.SendErrorResponse(schema.ClientSideError, schema.AccessDenied, nil)
		return
	}
	targetUser, err := database.FetchUserByUsername(params.Username)
	if err != nil {
		request.SendErrorResponse(schema.ClientSideError, schema.UserNotFound, nil)
		return
	}
	_, err = database.Db.NewDelete().Model(targetUser).Where("username = ?", targetUser.Username).Exec(database.Ctx)
	if err != nil {
		request.SendErrorResponse(schema.ServerSideError, "Could not delete user.", err)
		return
	}
	response := userDeleteResponse{}
	request.SendSuccessResponse(&response)
}
