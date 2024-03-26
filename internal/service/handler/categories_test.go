package handler

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/groshi-project/groshi/internal/database"
	"github.com/groshi-project/groshi/internal/middleware"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestHandler_CategoriesCreate(t *testing.T) {
	const (
		testUserID       int64 = 5
		testUsername           = "test-username"
		testCategoryName       = "test-category-name"
	)

	t.Run("create a new category owned by an existing user", func(t *testing.T) {
		var (
			handler = newTestHandler()
			ctx     = context.Background()
		)

		// create a test user:
		if err := handler.database.CreateUser(ctx, &database.User{
			ID:       testUserID,
			Username: testUsername,
		}); err != nil {
			panic(err)
		}

		ctx = context.WithValue(context.Background(), middleware.UsernameContextKey, testUsername)
		params := &categoriesCreateParams{
			Name: testCategoryName,
		}
		rec := testRequest(ctx, params, handler.CategoriesCreate)
		var createdCategoryUUID string
		if assert.Equal(t, http.StatusOK, rec.Code) {
			resp := &categoriesCreateResponse{}
			err := json.NewDecoder(rec.Body).Decode(resp)
			if assert.NoError(t, err) {
				if assert.NotEmpty(t, resp.UUID) {
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
	})

	t.Run("create a new category owned by a non-existed user", func(t *testing.T) {
		var (
			handler = newTestHandler()
			ctx     = context.WithValue(context.Background(), middleware.UsernameContextKey, "i-dont-exist")
		)

		params := &categoriesCreateParams{
			Name: testCategoryName,
		}
		rec := testRequest(ctx, params, handler.CategoriesCreate)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("call the handler with no params", func(t *testing.T) {
		var (
			handler = newTestHandler()
			ctx     = context.WithValue(context.Background(), middleware.UsernameContextKey, testUsername)
		)

		rec := testRequest(ctx, nil, handler.CategoriesCreate)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("call the handler without username context value", func(t *testing.T) {
		var (
			handler = newTestHandler()
			ctx     = context.Background()
		)

		params := &categoriesCreateParams{
			Name: testCategoryName,
		}
		rec := testRequest(ctx, params, handler.CategoriesCreate)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestHandler_CategoriesGet(t *testing.T) {
	const (
		testUserID   int64 = 5
		testUsername       = "test-username"
	)

	t.Run("get categories owned by existent user, when there are few categories", func(t *testing.T) {
		var (
			handler    = newTestHandler()
			ctx        = context.WithValue(context.Background(), middleware.UsernameContextKey, testUsername)
			categories = []*database.Category{
				{UUID: uuid.New(), Name: "Food", OwnerID: testUserID},
				{UUID: uuid.New(), Name: "Transport", OwnerID: testUserID},
				{UUID: uuid.New(), Name: "Sports", OwnerID: testUserID},
			}
		)

		// create a test user:
		if err := handler.database.CreateUser(ctx, &database.User{
			ID:       testUserID,
			Username: testUsername,
		}); err != nil {
			panic(err)
		}

		// create categories owned by the test user:
		for _, category := range categories {
			if err := handler.database.CreateCategory(ctx, category); err != nil {
				panic(err)
			}
		}

		rec := testRequest(ctx, nil, handler.CategoriesGet)
		if assert.Equal(t, http.StatusOK, rec.Code) {
			body := &categoriesGetResponse{}
			err := json.NewDecoder(rec.Body).Decode(body)
			if assert.NoError(t, err) {
				assert.Len(t, *body, len(categories))
			}
		}

	})

	t.Run("get categories owned by existent user, when there are no categories", func(t *testing.T) {
		var (
			handler = newTestHandler()
			ctx     = context.WithValue(context.Background(), middleware.UsernameContextKey, testUsername)
		)

		// create a test user:
		if err := handler.database.CreateUser(ctx, &database.User{
			ID:       testUserID,
			Username: testUsername,
		}); err != nil {
			panic(err)
		}

		rec := testRequest(ctx, nil, handler.CategoriesGet)
		if assert.Equal(t, http.StatusOK, rec.Code) {
			body := &categoriesGetResponse{}
			err := json.NewDecoder(rec.Body).Decode(body)
			if assert.NoError(t, err) {
				assert.Len(t, *body, 0)
			}
		}
	})

	t.Run("get categories owned by non-existent user", func(t *testing.T) {
		var (
			handler = newTestHandler()
			ctx     = context.WithValue(context.Background(), middleware.UsernameContextKey, "i-dont-exist")
		)
		rec := testRequest(ctx, nil, handler.CategoriesGet)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}
