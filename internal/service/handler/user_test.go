package handler

import (
	"context"
	"encoding/json"
	"github.com/groshi-project/groshi/internal/database"
	"github.com/groshi-project/groshi/internal/middleware"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestHandler_UserCreate(t *testing.T) {
	const (
		testUsername     = "test-user"
		testPassword     = "test-password"
		testPasswordHash = "hash(test-password)"
	)

	t.Run("create a new user with non-taken username", func(t *testing.T) {
		var (
			handler = newTestHandler()
			ctx     = context.Background()
		)

		params := &userCreateParams{
			Username: testUsername,
			Password: testPassword,
		}
		rec := testRequest(ctx, params, handler.UserCreate)
		if assert.Equal(t, http.StatusOK, rec.Code) {
			u := &database.User{}
			err := handler.database.SelectUserByUsername(ctx, testUsername, u)
			if assert.NoError(t, err) {
				if assert.NotEmpty(t, u) {
					assert.Equal(t, testUsername, u.Username)
					assert.Equal(t, testPasswordHash, u.Password)
				}
			}
		}
	})

	t.Run("create a user with username which is already taken", func(t *testing.T) {
		var (
			handler = newTestHandler()
			ctx     = context.Background()
		)

		// create a test user:
		if err := handler.database.CreateUser(ctx, &database.User{
			Username: testUsername,
		}); err != nil {
			panic(err)
		}

		params := &userCreateParams{
			Username: testUsername,
			Password: testPassword,
		}
		rec := testRequest(ctx, params, handler.UserCreate)
		assert.Equal(t, http.StatusConflict, rec.Code)
	})

	t.Run("call the handler without params", func(t *testing.T) {
		var (
			handler = newTestHandler()
			ctx     = context.Background()
		)
		rec := testRequest(ctx, nil, handler.UserCreate)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestHandler_UserGet(t *testing.T) {
	const (
		//testUserID   int64 = 99
		testUsername = "test-user"
	)

	t.Run("get existing user", func(t *testing.T) {
		var (
			handler = newTestHandler()
			ctx     = context.Background()
		)

		// create a test user:
		if err := handler.database.CreateUser(ctx, &database.User{
			Username: testUsername,
		}); err != nil {
			panic(err)
		}

		ctx = context.WithValue(ctx, middleware.UsernameContextKey, testUsername)
		rec := testRequest(ctx, nil, handler.UserGet)
		if assert.Equal(t, http.StatusOK, rec.Code) {
			resp := &userGetResponse{}
			err := json.NewDecoder(rec.Body).Decode(&resp)
			if assert.NoError(t, err) {
				if assert.NotEmpty(t, resp) {
					assert.Equal(t, testUsername, resp.Username)
				}
			}
		}

	})

	t.Run("get a non-existent user", func(t *testing.T) {
		var (
			handler = newTestHandler()
			ctx     = context.WithValue(context.Background(), middleware.UsernameContextKey, "i-dont-exist")
		)
		rec := testRequest(ctx, nil, handler.UserGet)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("call the handler without username context value", func(t *testing.T) {
		var (
			handler = newTestHandler()
			ctx     = context.Background()
		)
		rec := testRequest(ctx, nil, handler.UserGet)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestHandler_UserDelete(t *testing.T) {
	const (
		testUsername = "test-user"
	)

	t.Run("delete existing user", func(t *testing.T) {
		var (
			handler = newTestHandler()
			ctx     = context.Background()
		)

		// create a test user:
		if err := handler.database.CreateUser(ctx, &database.User{
			Username: testUsername,
		}); err != nil {
			panic(err)
		}

		ctx = context.WithValue(ctx, middleware.UsernameContextKey, testUsername)
		rec := testRequest(ctx, nil, handler.UserDelete)
		if assert.Equal(t, http.StatusOK, rec.Code) {
			resp := &userDeleteResponse{}
			err := json.NewDecoder(rec.Body).Decode(&resp)
			if assert.NoError(t, err) {
				if assert.NotEmpty(t, resp) {
					assert.Equal(t, testUsername, resp.Username)
				}
			}

			exists, err := handler.database.UserExistsByUsername(ctx, testUsername)
			if assert.NoError(t, err) {
				assert.False(t, exists)
			}
		}
	})

	t.Run("delete non-existent user", func(t *testing.T) {
		var (
			handler = newTestHandler()
			ctx     = context.WithValue(context.Background(), middleware.UsernameContextKey, "i-dont-exist")
		)
		rec := testRequest(ctx, nil, handler.UserDelete)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("call the handler without username context value", func(t *testing.T) {
		var (
			handler = newTestHandler()
			ctx     = context.Background()
		)
		rec := testRequest(ctx, nil, handler.UserDelete)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}
