package handler

import (
	"context"
	"encoding/json"
	"github.com/groshi-project/groshi/internal/database"
	"github.com/groshi-project/groshi/internal/middleware"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_CategoriesCreate(t *testing.T) {
	const (
		testUserID       int64 = 5
		testUsername           = "test-username"
		testCategoryName       = "test-category-name"
	)

	var (
		handler = newTestHandler()

		ctx    context.Context
		rec    *httptest.ResponseRecorder
		params *categoriesCreateParams
	)

	// create a test user:
	ctx = context.Background()
	if err := handler.database.CreateUser(ctx, &database.User{
		ID:       testUserID,
		Username: testUsername,
	}); err != nil {
		panic(err)
	}

	// case: create a new category with existing user:
	ctx = context.WithValue(context.Background(), middleware.UsernameContextVar, testUsername)

	// call handler to create a new category:
	params = &categoriesCreateParams{
		Name: testCategoryName,
	}
	rec = testRequest(ctx, params, handler.CategoriesCreate)
	var createdCategoryUUID string
	if assert.Equal(t, http.StatusOK, rec.Code) {
		resp := &categoriesCreateResponse{}
		err := json.NewDecoder(rec.Body).Decode(resp)
		if assert.NoError(t, err) {
			if assert.NotEmpty(t, resp.UUID) {
				// check if returned UUID matches with the UUID of the created category:
				createdCategoryUUID = resp.UUID
			}
		}
	}

	// check if the category was added to the database:
	c := &database.Category{}
	err := handler.database.SelectCategoryByUUID(ctx, createdCategoryUUID, c)
	if assert.NoError(t, err) {
		if assert.NotEmpty(t, c) {
			assert.Equal(t, testCategoryName, c.Name)
			assert.Equal(t, testUserID, c.OwnerID)
		}
	}

	// case: create a new category with non-existing user:
	ctx = context.WithValue(context.Background(), middleware.UsernameContextVar, "i-do-not-exist")

	params = &categoriesCreateParams{
		Name: testCategoryName,
	}
	rec = testRequest(ctx, params, handler.CategoriesCreate)
	assert.Equal(t, http.StatusNotFound, rec.Code)

	// case call handler without params:
	ctx = context.WithValue(context.Background(), middleware.UsernameContextVar, testUsername)
	params = nil
	rec = testRequest(ctx, params, handler.CategoriesCreate)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	// case: call handler without username context key
	ctx = context.Background()
	params = &categoriesCreateParams{
		Name: testCategoryName,
	}
	rec = testRequest(ctx, params, handler.CategoriesCreate)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
