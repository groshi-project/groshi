package handler

import (
	"context"
	"encoding/json"
	"github.com/groshi-project/groshi/internal/database"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestHandler_AuthLogin(t *testing.T) {
	const (
		testUsername     = "test-username"
		testPassword     = "test-password"
		testPasswordHash = "hash(test-password)"
	)

	t.Run("log in with correct password as an existing user", func(t *testing.T) {
		var (
			handler = newTestHandler()
			ctx     = context.Background()
		)

		// create a test user:
		if err := handler.database.CreateUser(ctx, &database.User{
			Username: testUsername,
			Password: testPasswordHash,
		}); err != nil {
			panic(err)
		}

		params := &authLoginParams{
			Username: testUsername,
			Password: testPassword,
		}
		rec := testRequest(ctx, params, handler.AuthLogin)
		if assert.Equal(t, http.StatusOK, rec.Code) {
			resp := &authLoginResponse{}
			err := json.NewDecoder(rec.Body).Decode(resp)
			if assert.NoError(t, err) {
				if assert.NotEmpty(t, resp) {
					assert.NotEmpty(t, resp.Token)
					assert.NotEmpty(t, resp.Expires)
				}
			}
		}
	})

	t.Run("log in as an existing user with wrong password", func(t *testing.T) {
		var (
			handler = newTestHandler()
			ctx     = context.Background()
		)

		params := &authLoginParams{
			Username: testUsername,
			Password: "wrong-password",
		}
		rec := testRequest(ctx, params, handler.AuthLogin)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("log in as non-existent user", func(t *testing.T) {
		var (
			handler = newTestHandler()
			ctx     = context.Background()
		)

		params := &authLoginParams{
			Username: "i-dont-exist",
			Password: "password",
		}
		rec := testRequest(ctx, params, handler.AuthLogin)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("call the handler with no params", func(t *testing.T) {
		var (
			handler = newTestHandler()
			ctx     = context.Background()
		)

		rec := testRequest(ctx, nil, handler.AuthLogin)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}
