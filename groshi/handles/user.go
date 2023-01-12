package handles

import (
	"github.com/jieggii/groshi/groshi/auth"
	"github.com/jieggii/groshi/groshi/database"
	"github.com/jieggii/groshi/groshi/ghttp"
	"github.com/jieggii/groshi/groshi/handles/jwt"
	"github.com/jieggii/groshi/groshi/handles/schema"
)

type userAuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type userAuthResponse struct {
	Token string `json:"token"`
}

func UserAuth(req *ghttp.Request, claims *jwt.Claims) {
	var credentials userAuthRequest
	if ok := req.DecodeSafe(&credentials); !ok {
		return
	}
	if ok := req.WrapCondition(len(credentials.Username) != 0 && len(credentials.Password) != 0, schema.InvalidRequestBody); ok {
		return
	}

	user := new(database.User)
	err := database.Db.NewSelect().Model(user).Where("username = ?", credentials.Username).Scan(database.Ctx)
	if err != nil {
		req.SendErrorResponse(schema.ClientSideError, schema.UserNotFound, nil)
		return
	}

	if !auth.CheckPasswordHash(credentials.Password, user.Password) {
		req.SendErrorResponse(schema.ClientSideError, "Invalid password.", nil)
		return
	}

	token, err := jwt.GenerateJWT(credentials.Username)
	if err != nil {
		req.SendErrorResponse(schema.ServerSideError, "Could not generate JWT.", err)
		return
	}
	response := userAuthResponse{Token: token}
	req.SendSuccessResponse(&response)
}

type userCreateRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type userCreateResponse struct{}

func UserCreate(req *ghttp.Request, claims *jwt.Claims) {
	currentUser, err := database.FetchUserByUsername(claims.Username)
	if err != nil {
		req.SendErrorResponse(schema.ClientSideError, schema.UserNotFound, nil)
		return
	}

	if !currentUser.IsSuperuser {
		req.SendErrorResponse(
			schema.ClientSideError,
			"You are not superuser (only superusers are allowed to create new users).",
			nil,
		)
		return
	}

	newUserData := userCreateRequest{}
	if ok := req.DecodeSafe(&newUserData); !ok {
		return
	}
	passwordHash, err := auth.HashPassword(newUserData.Password)
	if err != nil {
		req.SendErrorResponse(schema.ServerSideError, "Could not generate password hash.", err)
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
		req.SendErrorResponse(schema.ServerSideError, "Could not check if user already exists.", err)
		return
	}
	if newUserExists {
		req.SendErrorResponse(schema.ClientSideError, "User already exists.", nil)
		return
	}
	_, err = database.Db.NewInsert().Model(&newUser).Exec(database.Ctx)
	if err != nil {
		req.SendErrorResponse(schema.ServerSideError, "Could not create new user.", err)
		return
	}
	response := userCreateResponse{}
	req.SendSuccessResponse(&response)
}

type userReadRequest struct {
	Username string `json:"username"`
}

type userReadResponse struct {
	Username    string `json:"username"`
	IsSuperuser bool   `json:"is_superuser"`
}

func UserRead(req *ghttp.Request, claims *jwt.Claims) {
	//currentUser, err := database.FetchUserByUsername(claims.Username)
	//if err != nil {
	//	req.SendErrorResponse(schema.ClientSideError, schema.UserNotFound, nil)
	//	return
	//}
	params := userReadRequest{}
	if ok := req.DecodeSafe(&params); !ok {
		return
	}
	//if ok := req.WrapCondition(currentUser.Username == params.Username || currentUser.IsSuperuser, ""); !ok {
	//	return
	//}
	targetUser, err := database.FetchUserByUsername(params.Username)
	if err != nil {
		req.SendErrorResponse(schema.ClientSideError, schema.UserNotFound, nil)
	}
	response := userReadResponse{
		Username:    targetUser.Username,
		IsSuperuser: targetUser.IsSuperuser,
	}
	req.SendSuccessResponse(&response)
}

type userUpdateRequest struct {
}

type userUpdateResponse struct {
}

func UserUpdate(req *ghttp.Request, claims *jwt.Claims) {

}

type userDeleteRequest struct {
	Username string `json:"username"`
}
type userDeleteResponse struct{}

func UserDelete(req *ghttp.Request, claims *jwt.Claims) {
	currentUser, err := database.FetchUserByUsername(claims.Username)
	if err != nil {
		req.SendErrorResponse(schema.ClientSideError, schema.UserNotFound, nil)
		return
	}
	params := userDeleteRequest{}
	if ok := req.DecodeSafe(&params); !ok {
		return
	}
	if currentUser.Username != params.Username && !currentUser.IsSuperuser {
		req.SendErrorResponse(schema.ClientSideError, schema.AccessDenied, nil)
		return
	}
	targetUser, err := database.FetchUserByUsername(params.Username)
	if err != nil {
		req.SendErrorResponse(schema.ClientSideError, schema.UserNotFound, nil)
		return
	}
	_, err = database.Db.NewDelete().Model(targetUser).Exec(database.Ctx)
	if err != nil {
		req.SendErrorResponse(schema.ServerSideError, "Could not delete user.", err)
		return
	}
	response := userDeleteResponse{}
	req.SendSuccessResponse(&response)
}
