package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/groshi-project/groshi/internal/database"
	"github.com/groshi-project/groshi/internal/service/handler/errresp"
	"github.com/groshi-project/groshi/pkg/httpresp"
	"net/http"
)

type categoriesCreateParams struct {
	Name string `json:"name"`
}

type categoriesCreateResponse struct {
	UUID string `json:"uuid"`
}

func (h *Handler) CategoriesCreate(w http.ResponseWriter, r *http.Request) {
	// decode request params:
	params := &categoriesCreateParams{}
	if err := json.NewDecoder(r.Body).Decode(params); err != nil {
		errresp.InvalidRequest.Render(w)
		return
	}

	// validate request params:
	if err := h.paramsValidate.Struct(params); err != nil {
		errresp.InvalidRequestParams.Render(w)
		return
	}

	// extract current user's username from claims:
	username, err := h.JWTAuthority.ExtractUsername(r.Context())
	if err != nil {
		h.internalServerErrorLogger.Println(err)
		errresp.InternalServerError.Render(w)
		return
	}

	// fetch current user from the database:
	user := &database.User{}
	if err := h.database.SelectUserByUsername(username, user); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			errresp.UserNotFound.Render(w)
			return
		}
		h.internalServerErrorLogger.Println(err)
		errresp.InternalServerError.Render(w)
		return
	}

	// create a new category owned by the current user:
	category := &database.Category{
		Name:    params.Name,
		OwnerID: user.ID,
	}
	if err := h.database.CreateCategory(category); err != nil {
		h.internalServerErrorLogger.Println(err)
		errresp.InternalServerError.Render(w)
		return
	}

	// respond:
	resp := &categoriesCreateResponse{
		UUID: category.UUID.String(),
	}
	httpresp.NewOK(resp).Render(w)
}

type categoryResponse struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

//type categoriesGetResponse []categoryResponse

func (h *Handler) CategoriesGet(w http.ResponseWriter, r *http.Request) {
	// extract current user's username from claims
	username, err := h.JWTAuthority.ExtractUsername(r.Context())
	if err != nil {
		h.internalServerErrorLogger.Println(err)
		errresp.InternalServerError.Render(w)
		return
	}

	// fetch current user from the database:
	user := &database.User{}
	if err := h.database.SelectUserByUsername(username, user); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			errresp.UserNotFound.Render(w)
			return
		}
		h.internalServerErrorLogger.Println(err)
		errresp.InternalServerError.Render(w)
		return
	}

	// fetch categories that belong to this user from the database:
	categories := make([]database.Category, 0)
	if err := h.database.SelectCategoriesByOwnerID(user.ID, &categories); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			h.internalServerErrorLogger.Println(err)
			errresp.InternalServerError.Render(w)
			return
		}
	}

	// respond:
	resp := make([]categoryResponse, 0)
	for _, category := range categories {
		resp = append(resp, categoryResponse{
			UUID: category.UUID.String(),
			Name: category.Name,
		})
	}
	httpresp.NewOK(&resp)
}

type categoriesUpdateParams struct {
	Name string `json:"name"`
}

type categoriesUpdateResponse struct {
	UUID string `json:"uuid"`
}

func (h *Handler) CategoriesUpdate(w http.ResponseWriter, r *http.Request) {
	// decode request params:
	params := &categoriesUpdateParams{}
	if err := json.NewDecoder(r.Body).Decode(params); err != nil {
		errresp.InvalidRequest.Render(w)
		return
	}

	// validate request params:
	if err := h.paramsValidate.Struct(params); err != nil {
		errresp.InvalidRequestParams.Render(w)
		return
	}

	// parse URL params:
	uuid := chi.URLParam(r, "uuid")

	// fetch category:
	category := &database.Category{}
	if err := h.database.SelectCategoryByUUID(uuid, category); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			errresp.CategoryNotFound.Render(w)
			return
		}
		h.internalServerErrorLogger.Println(err)
		errresp.InternalServerError.Render(w)
		return
	}

	// extract current user's username from claims:
	username, err := h.JWTAuthority.ExtractUsername(r.Context())
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		h.internalServerErrorLogger.Println(err)
		errresp.InternalServerError.Render(w)
		return
	}

	// fetch current user:
	user := &database.User{}
	if err := h.database.SelectUserByUsername(username, user); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			errresp.UserNotFound.Render(w)
			return
		}
		h.internalServerErrorLogger.Println(err)
		errresp.InternalServerError.Render(w)
		return
	}

	// check if the category owned by the current user:
	if category.OwnerID != user.ID {
		errresp.CategoryForbidden.Render(w)
		return
	}

	// update category name:
	category.Name = params.Name
	if err := h.database.UpdateCategory(category); err != nil {
		h.internalServerErrorLogger.Println(err)
		errresp.InternalServerError.Render(w)
		return
	}

	// respond:
	resp := &categoriesUpdateResponse{
		UUID: category.UUID.String(),
	}
	httpresp.NewOK(resp)
}

type categoriesDeleteResponse struct {
	UUID string `json:"uuid"`
}

func (h *Handler) CategoriesDelete(w http.ResponseWriter, r *http.Request) {
	// parse URL params:
	uuid := chi.URLParam(r, "uuid")

	// fetch the given category from the database:
	category := &database.Category{}
	if err := h.database.SelectCategoryByUUID(uuid, category); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			errresp.CategoryNotFound.Render(w)
			return
		}
		h.internalServerErrorLogger.Println(err)
		errresp.InternalServerError.Render(w)
		return
	}

	// extract current user's username from claims
	username, err := h.JWTAuthority.ExtractUsername(r.Context())
	if err != nil {
		h.internalServerErrorLogger.Println(err)
		errresp.InternalServerError.Render(w)
		return
	}

	// fetch current user from the database:
	user := &database.User{}
	if err := h.database.SelectUserByUsername(username, user); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			errresp.UserNotFound.Render(w)
			return
		}
		h.internalServerErrorLogger.Println(err)
		errresp.InternalServerError.Render(w)
		return
	}

	// check if the category belongs to the user:
	if category.OwnerID != user.ID {
		errresp.CategoryForbidden.Render(w)
		return
	}

	// delete the given category from the database:
	if err := h.database.DeleteCategoryByID(category.ID); err != nil {
		h.internalServerErrorLogger.Println(err)
		errresp.InternalServerError.Render(w)
		return
	}

	// respond:
	resp := &categoriesDeleteResponse{UUID: uuid}
	httpresp.NewOK(resp)
}