package handler

import (
	"encoding/json"
	"github.com/groshi-project/groshi/internal/service/handler/model"
	"github.com/groshi-project/groshi/internal/service/handler/response"
	"net/http"

	"github.com/groshi-project/groshi/pkg/httpresp"

	"github.com/groshi-project/groshi/internal/database"
)

type userCreateParams struct {
	Username string `json:"username" example:"jieggii"`
	Password string `json:"password" example:"my-secret-password"`
}

type userCreateResponse struct {
	Username string `json:"username" example:"jieggii"`
}

// UserCreate creates a new user and returns its username.
//
//	@Summary		Create a new user
//	@Summary		Create a new user
//	@Description	Creates a new user and returns its username
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			user	body		userCreateParams	true	"Username and password"
//	@Success		200		{object}	userCreateResponse	"Successful operation"
//	@Failure		409		{object}	model.Error			"User with such username already exists"
//	@Failure		400		{object}	model.Error			"Invalid request body format or invalid request params"
//	@Failure		500		{object}	model.Error			"Internal server error"
//	@Router			/user [post]
func (h *Handler) UserCreate(w http.ResponseWriter, r *http.Request) {
	// parse request params:
	params := &userCreateParams{}
	if err := json.NewDecoder(r.Body).Decode(params); err != nil {
		httpresp.Render(w, response.InvalidRequestBodyFormat)
		return
	}

	// validate request params:
	if err := h.paramsValidate.Struct(params); err != nil {
		httpresp.Render(w, response.InvalidRequestParams)
		return
	}

	// check if user with such username already exist:
	exists, err := h.database.UserExistsByUsername(params.Username)
	if err != nil {
		h.internalServerErrorLogger.Println(err)
		httpresp.Render(w, response.InternalServerError)
		return
	}
	if exists {
		httpresp.Render(w, httpresp.New(http.StatusConflict, model.NewError("user already exists")))
		return
	}

	// create a new user:
	passwordHash, err := h.passwordAuthority.HashPassword(params.Password)
	if err != nil {
		h.internalServerErrorLogger.Println(err)
		httpresp.Render(w, response.InternalServerError)
		return
	}
	user := &database.User{
		Username: params.Username,
		Password: passwordHash,
	}
	if err := h.database.CreateUser(user); err != nil {
		h.internalServerErrorLogger.Println(err)
		httpresp.Render(w, response.InternalServerError)
		return
	}

	// respond:
	resp := &userCreateResponse{
		Username: user.Username,
	}
	httpresp.Render(w, httpresp.NewOK(resp))
}

type userGetResponse struct {
	Username string `json:"username" example:"jieggii"`
}

// UserGet returns information about the current user.
//
//	@Summary		Retrieve information about the current user
//	@Description	Returns information about the current user
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	userGetResponse	"Successful operation"
//	@Failure		404	{object}	model.Error		"User not found"
//	@Failure		400	{object}	model.Error		"Invalid request body format or invalid request params"
//	@Failure		500	{object}	model.Error		"Internal server error"
//	@Security		Bearer
//	@Router			/user [get]
func (h *Handler) UserGet(w http.ResponseWriter, r *http.Request) {
	// extract current user's username from context:
	username := r.Context().Value("username").(string)

	// todo?: fetch user from the database to check if it exists.

	// respond:
	resp := &userGetResponse{Username: username}
	httpresp.Render(w, httpresp.NewOK(resp))
}

func (h *Handler) UserUpdate(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}

type userDeleteResponse struct {
	Username string `json:"username" example:"jieggii"`
}

// UserDelete deletes the current user.
//
//	@Summary		Delete the current user
//	@Description	Deletes the current user and returns its username
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	userDeleteResponse	"Successful operation"
//	@Failure		404	{object}	model.Error			"User not found"
//	@Failure		400	{object}	model.Error			"Invalid request body format or invalid request params"
//	@Failure		500	{object}	model.Error			"Internal server error"
//	@Router			/user [delete]
func (h *Handler) UserDelete(w http.ResponseWriter, r *http.Request) {
	// extract current user's username from context:
	username := r.Context().Value("username").(string)

	// check if the user exists:
	exists, err := h.database.UserExistsByUsername(username)
	if err != nil {
		h.internalServerErrorLogger.Println(err)
		httpresp.Render(w, response.InternalServerError)
		return
	}
	if !exists {
		httpresp.Render(w, response.UserNotFound)
		return
	}

	// delete the user from the database:
	if err := h.database.DeleteUserByUsername(username); err != nil {
		h.internalServerErrorLogger.Println(err)
		httpresp.Render(w, response.InternalServerError)
		return
	}

	// respond:
	resp := &userDeleteResponse{
		Username: username,
	}
	httpresp.Render(w, httpresp.NewOK(resp))
}
