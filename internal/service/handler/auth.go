package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/groshi-project/groshi/internal/database"
	"github.com/groshi-project/groshi/internal/service/handler/httpresp"
	"github.com/groshi-project/groshi/internal/service/handler/response"
	"net/http"
	"time"
)

type authLoginParams struct {
	Username string `json:"username" username:"username" validate:"required"`
	Password string `json:"password" password:"my-secret-password" validate:"required"`
}

type authLoginResponse struct {
	Token   string    `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6IldlbGwsIGlmIHlvdSBjYW4gcmVhZCB0aGlzLCB0aGVuIHlvdSBkZWZpbmV0ZWx5IHdvdWxkIGxpa2UgdGhpcyBvbmU6IGh0dHBzOi8veW91dHUuYmUvZFF3NHc5V2dYY1EifQ.1ervhGZz1m6xiHR447rbwh8W1sfATF2qYudOtNWhkkw"`
	Expires time.Time `json:"expires" example:"2034-03-20T12:57:38+02:00"`
}

// AuthLogin authenticates user, generates and returns JWT.
//
//	@Summary		Authenticate user
//	@Description	Authenticates user, generates and returns valid JSON Web Token
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			credentials	body		authLoginParams		true	"Username and password"
//	@Success		200			{object}	authLoginResponse	"Successful operation"
//	@Failure		403			{object}	model.Error			"Invalid credentials"
//	@Failure		400			{object}	model.Error			"Invalid request body format or invalid request params"
//	@Failure		500			{object}	model.Error			"Internal server error"
//	@Router			/auth/login [post]
func (h *Handler) AuthLogin(w http.ResponseWriter, r *http.Request) {
	// decode request params:
	params := &authLoginParams{}
	if err := json.NewDecoder(r.Body).Decode(params); err != nil {
		httpresp.Render(w, response.InvalidRequestBodyFormat)
		return
	}

	// validate request params:
	if err := h.paramsValidate.Struct(params); err != nil {
		httpresp.Render(w, response.InvalidRequestParams)
		return
	}

	// fetch the user from the database:
	user := &database.User{}
	if err := h.database.SelectUserByUsername(r.Context(), params.Username, user); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			httpresp.Render(w, response.InvalidCredentials)
			return
		}
		httpresp.Render(w, response.InternalServerError)
		return
	}

	// validate provided password:
	ok, err := h.passwordAuth.VerifyPassword(params.Password, user.Password)
	if err != nil {
		httpresp.Render(w, response.InternalServerError)
		return
	}
	if !ok {
		httpresp.Render(w, response.InvalidCredentials)
		return
	}

	// generate a new jwt:
	token, expires, err := h.JWTAuth.CreateToken(user.Username)
	if err != nil {
		httpresp.Render(w, response.InternalServerError)
		return
	}

	// respond:
	resp := &authLoginResponse{
		Token:   token,
		Expires: expires,
	}
	httpresp.Render(w, httpresp.NewOK(resp))
}
