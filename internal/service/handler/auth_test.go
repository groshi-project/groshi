package handler

import (
	"context"
	"encoding/json"
	"github.com/groshi-project/groshi/internal/database"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_AuthLogin(t *testing.T) {
	var (
		handler = newTestHandler()
		ctx     = context.Background()

		rec    *httptest.ResponseRecorder
		params *authLoginParams
	)

	if err := handler.database.CreateUser(ctx, &database.User{
		Username: testUsername,
		Password: testPasswordHash,
	}); err != nil {
		panic(err)
	}

	// log in as existing user with correct password:
	params = &authLoginParams{
		Username: testUsername,
		Password: testPassword,
	}
	rec = testRequest(ctx, params, handler.AuthLogin)
	if assert.Equal(t, http.StatusOK, rec.Code) {
		resp := &authLoginResponse{}
		err := json.NewDecoder(rec.Body).Decode(resp)
		if assert.NoError(t, err) {
			if assert.NotEmpty(t, resp) {
				assert.NotEmpty(t, resp.Token)
				//assert.NotZero(t, resp.Expires)
			}
		}
	}

	// try to log in as existing user but with wrong password
	params = &authLoginParams{
		Username: testUsername,
		Password: "wrong-password-123",
	}
	rec = testRequest(ctx, params, handler.AuthLogin)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)

	// try to log in as non-existing user
	params = &authLoginParams{
		Username: "i-do-not-exist",
		Password: "password",
	}
	rec = testRequest(ctx, params, handler.AuthLogin)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)

	// call handler without params:
	params = nil
	rec = testRequest(ctx, nil, handler.AuthLogin)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
