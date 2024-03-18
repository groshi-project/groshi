package handler

import (
	"encoding/json"
	"github.com/groshi-project/groshi/internal/service/handler/errresp"
	"net/http"

	"github.com/groshi-project/groshi/pkg/httpresp"

	"github.com/groshi-project/groshi/internal/database"
)

type userCreateParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type userCreateResponse struct {
	Username string `json:"username"`
}

func (h *Handler) UserCreate(w http.ResponseWriter, r *http.Request) {
	// parse request params:
	params := &userCreateParams{}
	if err := json.NewDecoder(r.Body).Decode(params); err != nil {
		errresp.InvalidRequest.Render(w)
		return
	}

	// validate request params:
	if err := h.paramsValidate.Struct(params); err != nil {
		errresp.InvalidRequestParams.Render(w)
		return
	}

	// check if user with such username already exist:
	exists, err := h.database.UserExistsByUsername(params.Username)
	if err != nil {
		h.internalServerErrorLogger.Println(err)
		errresp.InternalServerError.Render(w)
		return
	}
	if exists {
		httpresp.New(http.StatusConflict, errresp.NewErrorData("user already exists"))
		return
	}

	// create a new user:
	passwordHash, err := h.passwordAuthority.HashPassword(params.Password)
	if err != nil {
		h.internalServerErrorLogger.Println(err)
		errresp.InternalServerError.Render(w)
		return
	}
	user := &database.User{
		Username: params.Username,
		Password: passwordHash,
	}
	if err := h.database.CreateUser(user); err != nil {
		h.internalServerErrorLogger.Println(err)
		errresp.InternalServerError.Render(w)
		return
	}

	// respond:
	resp := &userCreateResponse{
		Username: user.Username,
	}
	httpresp.NewOK(resp)
}

type userGetResponse struct {
	Username string `json:"username"`
}

func (h *Handler) UserGet(w http.ResponseWriter, r *http.Request) {
	// extract current user's username from claims:
	username, err := h.JWTAuthority.ExtractUsername(r.Context())
	if err != nil {
		h.internalServerErrorLogger.Println(err)
		errresp.InternalServerError.Render(w)
		return
	}

	// todo?: fetch user from the database to check if it exists.

	// respond:
	resp := &userGetResponse{Username: username}
	httpresp.NewOK(resp)
}

func (h *Handler) UserUpdate(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}

type userDeleteResponse struct {
	Username string `json:"username"`
}

func (h *Handler) UserDelete(w http.ResponseWriter, r *http.Request) {
	// extract current user's username from claims:
	username, err := h.JWTAuthority.ExtractUsername(r.Context())
	if err != nil {
		h.internalServerErrorLogger.Println(err)
		errresp.InternalServerError.Render(w)
		return
	}

	// check if the user exists:
	exists, err := h.database.UserExistsByUsername(username)
	if err != nil {
		h.internalServerErrorLogger.Println(err)
		errresp.InternalServerError.Render(w)
		return
	}
	if !exists {
		errresp.UserNotFound.Render(w)
		return
	}

	// delete the user from the database:
	if err := h.database.DeleteUserByUsername(username); err != nil {
		h.internalServerErrorLogger.Println(err)
		errresp.InternalServerError.Render(w)
		return
	}

	// respond:
	resp := &userDeleteResponse{
		Username: username,
	}
	httpresp.NewOK(resp)
}
