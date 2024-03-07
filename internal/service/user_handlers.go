package service

import (
	"encoding/json"
	"github.com/groshi-project/groshi/internal/database"
	"net/http"
)

type userCreateParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type userCreateResponse struct {
	Username string `json:"username"`
}

func (s *Service) UserCreate(w http.ResponseWriter, r *http.Request) {
	params := userCreateParams{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// check if user with such username already exist:
	exists, err := s.Database.Client.NewSelect().Model(database.EmptyUser).Where("username = ?", params.Username).Exists(s.Database.Ctx)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
		return
	}

	// create a new user:
	passwordHash, err := s.PasswordAuthority.HashPassword(params.Password)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	user := database.User{
		Username: params.Username,
		Password: passwordHash,
	}
	if _, err := s.Database.Client.NewInsert().Model(&user).Exec(s.Database.Ctx); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// respond:
	response := userCreateResponse{
		Username: user.Username,
	}
	if err := json.NewEncoder(w).Encode(&response); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

type userGetResponse struct {
	Username string `json:"username"`
}

func (s *Service) UserGet(w http.ResponseWriter, r *http.Request) {
	// extract current user's username from claims:
	claims, err := s.JWTAuthority.ExtractClaims(r.Context())
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	username := claims["username"].(string)

	// respond:
	response := userGetResponse{Username: username}
	if err := json.NewEncoder(w).Encode(&response); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (s *Service) UserUpdate(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}

type userDeleteResponse struct {
	Username string `json:"username"`
}

func (s *Service) UserDelete(w http.ResponseWriter, r *http.Request) {
	// extract current user's username from claims:
	claims, err := s.JWTAuthority.ExtractClaims(r.Context())
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	username := claims["username"].(string)

	// delete the user from the database:
	if _, err := s.Database.Client.NewDelete().Model(&database.EmptyUser).Where("username = ?", username).Exec(s.Database.Ctx); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// respond:
	response := &userDeleteResponse{
		Username: username,
	}
	if err := json.NewEncoder(w).Encode(&response); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
