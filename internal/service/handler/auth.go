package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/groshi-project/groshi/internal/database"
	"github.com/groshi-project/groshi/internal/service/handler/response"
	"github.com/groshi-project/groshi/pkg/httpresp"
	"net/http"
	"time"
)

type authLoginParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type authLoginResponse struct {
	Token   string    `json:"token"`
	Expires time.Time `json:"expires"`
}

func (h *Handler) AuthLogin(w http.ResponseWriter, r *http.Request) {
	// decode request params:
	params := &authLoginParams{}
	if err := json.NewDecoder(r.Body).Decode(params); err != nil {
		httpresp.Render(w, response.InvalidRequest)
		return
	}

	// validate request params:
	if err := h.paramsValidate.Struct(params); err != nil {
		httpresp.Render(w, response.InvalidRequestParams)
		return
	}

	// fetch the user from the database:
	user := &database.User{}
	if err := h.database.SelectUserByUsername(params.Username, user); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			httpresp.Render(w, response.InvalidCredentials)
			return
		}
		httpresp.Render(w, response.InternalServerError)
		return
	}

	// validate provided password:
	if !h.passwordAuthority.ValidatePassword(params.Password, user.Password) {
		httpresp.Render(w, response.InvalidCredentials)
		return
	}

	// generate a new jwt:
	token, expires, err := h.JWTAuthority.GenerateToken(user.Username)
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
