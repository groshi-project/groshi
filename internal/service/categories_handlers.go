package service

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/groshi-project/groshi/internal/database"
	"github.com/groshi-project/groshi/internal/service/httputil"
	"net/http"
)

func (s *Service) CategoriesCreate(w http.ResponseWriter, r *http.Request) {

}

type categoryResponse struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

type categoriesGetResponse []categoryResponse

func (s *Service) CategoriesGet(w http.ResponseWriter, r *http.Request) {
	// extract current user's username from claims
	claims, err := s.JWTAuthority.ExtractClaims(r.Context())
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	username := claims["username"].(string)

	// fetch current user from the database:
	user := database.User{}
	if err := s.Database.Client.NewSelect().Model(&user).Where("username = ?", username).Scan(s.Database.Ctx, &user); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// fetch categories that belong to this user from the database:
	categories := make([]database.Category, 0)
	if err := s.Database.Client.NewSelect().Model(&categories).Where("owner_id = ?", user.ID).Scan(s.Database.Ctx, &categories); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}

	// respond:
	response := make([]categoryResponse, 0)
	for _, category := range categories {
		response = append(response, categoryResponse{
			UUID: category.UUID.String(),
			Name: category.Name,
		})
	}
	if err := json.NewEncoder(w).Encode(&response); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

type categoriesDeleteParams struct {
	UUID string `json:"uuid"`
}

type categoriesDeleteResponse struct {
	UUID string `json:"uuid"`
}

func (s *Service) CategoriesDelete(w http.ResponseWriter, r *http.Request) {
	// decode the request body:
	params := categoriesDeleteParams{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// fetch the given category from the database:
	category := database.Category{}
	if err := s.Database.Client.NewSelect().Model(&category).Where("uuid = ?", params.UUID).Scan(s.Database.Ctx, &category); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// extract current user's username from claims
	claims, err := s.JWTAuthority.ExtractClaims(r.Context())
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	username := claims["username"].(string)

	// fetch current user from the database:
	user := database.User{}
	if err := s.Database.Client.NewSelect().Model(&user).Where("username = ?", username).Scan(s.Database.Ctx, &user); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, httputil.UserNotFound, http.StatusNotFound) // todo
			return
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// check if the category belongs to the user:
	if category.OwnerID != user.ID {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	// delete the given category from the database:
	if _, err := s.Database.Client.NewDelete().Model(database.ZeroCategory).Exec(s.Database.Ctx); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// respond:
	response := categoriesDeleteResponse{UUID: params.UUID}
	if err := json.NewEncoder(w).Encode(&response); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
