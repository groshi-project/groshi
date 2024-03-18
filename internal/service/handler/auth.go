package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/groshi-project/groshi/internal/database"
	"github.com/groshi-project/groshi/internal/service/handler/errresp"
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
		errresp.InvalidRequest.Render(w)
		return
	}

	// validate request params:
	if err := h.paramsValidate.Struct(params); err != nil {
		errresp.InvalidRequestParams.Render(w)
		return
	}

	// fetch the user from the database:
	user := &database.User{}
	if err := h.database.SelectUserByUsername(params.Username, user); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			errresp.InvalidCredentials.Render(w)
			return
		}
		errresp.InternalServerError.Render(w)
		return
	}

	// validate provided password:
	if !h.passwordAuthority.ValidatePassword(params.Password, user.Password) {
		errresp.InvalidCredentials.Render(w)
		return
	}

	// generate a new jwt:
	token, expires, err := h.JWTAuthority.GenerateToken(user.Username)
	if err != nil {
		errresp.InternalServerError.Render(w)
		return
	}

	// respond:
	resp := &authLoginResponse{
		Token:   token,
		Expires: expires,
	}
	httpresp.NewOK(resp).Render(w)
}
