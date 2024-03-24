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

func TestHandler_UserCreate(t *testing.T) {
	const (
		testUsername     = "test-user"
		testPassword     = "test-password"
		testPasswordHash = "hash(test-password)"
	)

	var (
		ctx = context.Background() // context which will be used for all test requests

		handler *Handler
		params  *userCreateParams
		rec     *httptest.ResponseRecorder
	)

	// case: create a new user with a unique non-existing username:
	handler = newTestHandler()

	params = &userCreateParams{
		Username: testUsername,
		Password: testPassword,
	}
	rec = testRequest(ctx, params, handler.UserCreate)
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

	// case: create a user with username which is already taken:
	handler = newTestHandler()

	if err := handler.database.CreateUser(ctx, &database.User{
		Username: testUsername,
	}); err != nil {
		panic(err)
	}

	params = &userCreateParams{
		Username: testUsername,
		Password: testPassword,
	}
	rec = testRequest(ctx, params, handler.UserCreate)
	assert.Equal(t, http.StatusConflict, rec.Code)

	// case: empty params:
	handler = newTestHandler()

	params = nil
	rec = testRequest(ctx, params, handler.UserCreate)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestHandler_UserGet(t *testing.T) {
	const (
		testUsername = "test-user"
	)

	var (
		handler = newTestHandler() // handler which will be used to test all requests

		ctx context.Context
		rec *httptest.ResponseRecorder
	)

	if err := handler.database.CreateUser(context.Background(), &database.User{
		Username: testUsername,
	}); err != nil {
		panic(err)
	}

	// case: get existing user:
	ctx = context.WithValue(context.Background(), "username", testUsername)
	rec = testRequest(ctx, nil, handler.UserGet)
	if assert.Equal(t, http.StatusOK, rec.Code) {
		resp := &userGetResponse{}
		err := json.NewDecoder(rec.Body).Decode(&resp)
		assert.NoError(t, err)

		if assert.NotEmpty(t, resp) {
			assert.Equal(t, testUsername, resp.Username)
		}
	}

	// case: get non-existing user:
	ctx = context.WithValue(context.Background(), "username", "i-do-not-exist")
	rec = testRequest(ctx, nil, handler.UserGet)
	assert.Equal(t, http.StatusNotFound, rec.Code)

	// case: call handler without "username" context var:
	ctx = context.Background()
	rec = testRequest(ctx, nil, handler.UserGet)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

// todo: test user update?

func TestHandler_UserDelete(t *testing.T) {
	const (
		testUsername = "test-user"
	)

	var (
		handler = newTestHandler()

		ctx context.Context
		rec *httptest.ResponseRecorder
	)

	// test deleting existing user:
	if err := handler.database.CreateUser(context.Background(), &database.User{
		Username: testUsername,
	}); err != nil {
		panic(err)
	}

	ctx = context.WithValue(context.Background(), "username", testUsername)
	rec = testRequest(ctx, nil, handler.UserDelete)
	if assert.Equal(t, http.StatusOK, rec.Code) {
		resp := &userDeleteResponse{}
		err := json.NewDecoder(rec.Body).Decode(&resp)
		if assert.NoError(t, err) {
			if assert.NotEmpty(t, resp) {
				assert.Equal(t, testUsername, resp.Username)
			}
		}

		exists, err := handler.database.UserExistsByUsername(ctx, testUsername)
		if err != nil {
			panic(err)
		}
		assert.Equal(t, false, exists)
	}

	// test deleting non-existing user:
	ctx = context.WithValue(context.Background(), "username", "i-do-not-exist")
	rec = testRequest(ctx, nil, handler.UserDelete)
	assert.Equal(t, http.StatusNotFound, rec.Code)

	// test without "username" context var:
	ctx = context.Background()
	rec = testRequest(ctx, nil, handler.UserDelete)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
