package service

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/groshi-project/groshi/internal/database"
	"net/http"
	"time"
)

type AuthLoginParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthLoginResponse struct {
	Token   string    `json:"token"`
	Expires time.Time `json:"expires"`
}

func (s *Service) AuthLogin(w http.ResponseWriter, r *http.Request) {
	credentials := AuthLoginParams{}
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// fetch the user from the database:
	user := database.User{}
	if err := s.Database.Client.NewSelect().Model(&user).Where("username = ?").Scan(s.Database.Ctx, &user); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// validate provided password:
	if !s.PasswordAuthority.ValidatePassword(credentials.Password, user.Password) {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	// generate a new jwt:
	token, expires, err := s.JWTAuthority.GenerateToken(user.Username)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// respond:
	response := AuthLoginResponse{
		Token:   token,
		Expires: expires,
	}
	if err := json.NewEncoder(w).Encode(&response); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
