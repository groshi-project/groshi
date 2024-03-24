package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/groshi-project/groshi/internal/database"
	"github.com/groshi-project/groshi/internal/middleware"
	"github.com/groshi-project/groshi/internal/service/handler/httpresp"
	"github.com/groshi-project/groshi/internal/service/handler/response"
	"net/http"
)

type categoriesCreateParams struct {
	Name string `json:"name" example:"Transport" validate:"required"`
}

type categoriesCreateResponse struct {
	UUID string `json:"uuid" example:"c319d169-c7bd-4768-b61c-07f796dce3a2"`
}

// CategoriesCreate creates a new category and returns its UUID.
//
//	@Summary		Create a new category
//	@Description	Creates a new category and returns its UUID
//	@Tags			categories
//	@Accept			json
//	@Produce		json
//	@Param			user	body		categoriesCreateParams		true	"Category name"
//	@Success		200		{object}	categoriesCreateResponse	"Successful operation"
//	@Failure		400		{object}	model.Error					"Invalid request body format or invalid request params"
//	@Failure		404		{object}	model.Error					"User not found"
//	@Failure		500		{object}	model.Error					"Internal server error"
//	@Security		Bearer
//	@Router			/categories [post]
func (h *Handler) CategoriesCreate(w http.ResponseWriter, r *http.Request) {
	// decode request params:
	params := &categoriesCreateParams{}
	if err := json.NewDecoder(r.Body).Decode(params); err != nil {
		httpresp.Render(w, response.InvalidRequestBodyFormat)
		return
	}

	// validate request params:
	if err := h.paramsValidate.Struct(params); err != nil {
		httpresp.Render(w, response.InvalidRequestParams)
		return
	}

	// extract current user's username from context:
	username, ok := r.Context().Value(middleware.UsernameContextVar).(string)
	if !ok {
		h.internalServerErrorLogger.Println(errMissingUsernameContextValue)
		httpresp.Render(w, response.InternalServerError)
	}

	// fetch the current user from the database:
	user := &database.User{}
	if err := h.database.SelectUserByUsername(r.Context(), username, user); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			httpresp.Render(w, response.UserNotFound)
			return
		}
		h.internalServerErrorLogger.Println(err)
		httpresp.Render(w, response.InternalServerError)
		return
	}

	// create a new category owned by the current user:
	category := &database.Category{
		Name:    params.Name,
		OwnerID: user.ID,
	}
	if err := h.database.CreateCategory(r.Context(), category); err != nil {
		h.internalServerErrorLogger.Println(err)
		httpresp.Render(w, response.InternalServerError)
		return
	}

	// respond:
	resp := &categoriesCreateResponse{
		UUID: category.UUID.String(),
	}
	httpresp.Render(w, httpresp.NewOK(resp))
}

type categoriesGetResponseItem struct {
	UUID string `json:"uuid" example:"8b95b038-8a7a-4cdc-96b5-506101ed3a73"`
	Name string `json:"name" example:"Transport"`
}

type categoriesGetResponse []categoriesGetResponseItem

// CategoriesGet returns all categories created by user.
//
//	@Summary		Fetch all categories
//	@Description	Returns all categories created by user.
//	@Tags			categories
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	categoriesGetResponse	"Successful operation"
//	@Failure		404	{object}	model.Error				"User not found"
//	@Failure		500	{object}	model.Error				"Internal server error"
//	@Security		Bearer
//	@Router			/categories [get]
func (h *Handler) CategoriesGet(w http.ResponseWriter, r *http.Request) {
	// extract current user's username from context
	username, ok := r.Context().Value(middleware.UsernameContextVar).(string)
	if !ok {
		h.internalServerErrorLogger.Println(errMissingUsernameContextValue)
		httpresp.Render(w, response.InternalServerError)
		return
	}

	// fetch current user from the database:
	user := &database.User{}
	if err := h.database.SelectUserByUsername(r.Context(), username, user); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			httpresp.Render(w, response.UserNotFound)
			return
		}
		h.internalServerErrorLogger.Println(err)
		httpresp.Render(w, response.InternalServerError)
		return
	}

	// fetch categories that belong to this user from the database:
	// todo: should slice of pointers be used instead of a slice of objects?
	categories := make([]database.Category, 0)
	if err := h.database.SelectCategoriesByOwnerID(r.Context(), user.ID, &categories); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			h.internalServerErrorLogger.Println(err)
			httpresp.Render(w, response.InternalServerError)
			return
		}
	}

	// respond:
	resp := make([]categoriesGetResponseItem, 0)
	for _, category := range categories {
		resp = append(resp, categoriesGetResponseItem{
			UUID: category.UUID.String(),
			Name: category.Name,
		})
	}
	httpresp.NewOK(&resp)
}

type categoriesUpdateParams struct {
	Name string `json:"name" example:"Food" validate:"required"`
}

type categoriesUpdateResponse struct {
	UUID string `json:"uuid" example:"9d1a6ba2-d2e1-4ca4-b8d3-164f2009c823"`
}

func (h *Handler) CategoriesUpdate(w http.ResponseWriter, r *http.Request) {
	// decode request params:
	params := &categoriesUpdateParams{}
	if err := json.NewDecoder(r.Body).Decode(params); err != nil {
		httpresp.Render(w, response.InvalidRequestBodyFormat)
		return
	}

	// validate request params:
	if err := h.paramsValidate.Struct(params); err != nil {
		httpresp.Render(w, response.InvalidRequestParams)
		return
	}

	// parse URL params:
	uuid := chi.URLParam(r, "uuid")

	// fetch category:
	category := &database.Category{}
	if err := h.database.SelectCategoryByUUID(r.Context(), uuid, category); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			httpresp.Render(w, response.CategoryNotFound)
			return
		}
		h.internalServerErrorLogger.Println(err)
		httpresp.Render(w, response.InternalServerError)
		return
	}

	// extract current user's username from context:
	username := r.Context().Value("username").(string)

	// fetch current user:
	user := &database.User{}
	if err := h.database.SelectUserByUsername(r.Context(), username, user); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			httpresp.Render(w, response.UserNotFound)
			return
		}
		h.internalServerErrorLogger.Println(err)
		httpresp.Render(w, response.InternalServerError)
		return
	}

	// check if the category owned by the current user:
	if category.OwnerID != user.ID {
		httpresp.Render(w, response.CategoryForbidden)
		return
	}

	// update category name:
	category.Name = params.Name
	if err := h.database.UpdateCategory(r.Context(), category); err != nil {
		h.internalServerErrorLogger.Println(err)
		httpresp.Render(w, response.InternalServerError)
		return
	}

	// respond:
	resp := &categoriesUpdateResponse{
		UUID: category.UUID.String(),
	}
	httpresp.NewOK(resp)
}

type categoriesDeleteResponse struct {
	UUID string `json:"uuid" example:"9d1a6ba2-d2e1-4ca4-b8d3-164f2009c823"`
}

func (h *Handler) CategoriesDelete(w http.ResponseWriter, r *http.Request) {
	// parse URL params:
	uuid := chi.URLParam(r, "uuid")

	// fetch the given category from the database:
	category := &database.Category{}
	if err := h.database.SelectCategoryByUUID(r.Context(), uuid, category); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			httpresp.Render(w, response.CategoryNotFound)
			return
		}
		h.internalServerErrorLogger.Println(err)
		httpresp.Render(w, response.InternalServerError)
		return
	}

	// extract current user's username from context
	username := r.Context().Value("username").(string)

	// fetch current user from the database:
	user := &database.User{}
	if err := h.database.SelectUserByUsername(r.Context(), username, user); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			httpresp.Render(w, response.UserNotFound)
			return
		}
		h.internalServerErrorLogger.Println(err)
		httpresp.Render(w, response.InternalServerError)
		return
	}

	// check if the category belongs to the user:
	if category.OwnerID != user.ID {
		httpresp.Render(w, response.CategoryForbidden)
		return
	}

	// delete the given category from the database:
	if err := h.database.DeleteCategoryByID(r.Context(), category.ID); err != nil {
		h.internalServerErrorLogger.Println(err)
		httpresp.Render(w, response.InternalServerError)
		return
	}

	// respond:
	resp := &categoriesDeleteResponse{UUID: uuid}
	httpresp.NewOK(resp)
}
