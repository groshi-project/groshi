package middleware

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestTokenFromHeader(t *testing.T) {
	// test empty header value:
	token1, err := tokenFromHeader("")
	if assert.Error(t, err) {
		assert.ErrorIs(t, err, errEmptyOrMissingAuthHeader)
	}
	assert.Empty(t, token1)

	// test malformed header value:
	token2, err := tokenFromHeader("something bla bla blabla")
	if assert.Error(t, err) {
		assert.ErrorIs(t, err, errInvalidAuthHeader)
	}
	assert.Empty(t, token2)

	// test "Bearer" header value:
	token3, err := tokenFromHeader("Bearer")
	if assert.Error(t, err) {
		assert.ErrorIs(t, err, errInvalidAuthHeader)
	}
	assert.Empty(t, token3)

	// test "Bearer token extra" header value:
	token4, err := tokenFromHeader("Bearer token extra")
	if assert.Error(t, err) {
		assert.ErrorIs(t, err, errInvalidAuthHeader)
	}
	assert.Empty(t, token4)

	// finally, test valid header value:
	token5, err := tokenFromHeader("Bearer some-token")
	assert.NoError(t, err)
	assert.Equal(t, "some-token", token5)
}

// mockAuthority is a mock implementation of the AuthorityInterface interface.
type mockAuthority struct {
	Token  string
	Claims jwt.MapClaims
}

// CreateToken stub.
func (m *mockAuthority) CreateToken(username string) (string, time.Time, error) {
	return m.Token, time.Time{}, nil
}

// VerifyToken stub. Simply checks if provided tokenString equal to the struct field [mockAuthority.Token].
func (m *mockAuthority) VerifyToken(tokenString string) (jwt.MapClaims, error) {
	if tokenString == m.Token {
		return m.Claims, nil
	} else {
		return nil, jwt.ErrTokenSignatureInvalid
	}
}

// testRequest creates a new [httptest.Recorder] and makes a test request to handler using it,
// then returns pointer to this recorder.
func testRequest(setAuthHeader bool, authHeaderValue string, handler http.Handler) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/", nil)
	if setAuthHeader {
		request.Header.Set("Authorization", authHeaderValue)
	}

	handler.ServeHTTP(recorder, request)
	return recorder
}

func TestNewJWT(t *testing.T) {
	// username of a test user:
	username := "my-username-123"

	authority := &mockAuthority{
		Token:  "example-token-123",
		Claims: map[string]any{"username": username},
	}

	// test handler which checks username from context:
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		actualUsername := r.Context().Value("username").(string)
		assert.Equal(t, username, actualUsername)
	})

	middleware := NewJWT(authority)
	handlerWithMiddleware := middleware(testHandler)

	// test with correct token:
	recorder := testRequest(true, "Bearer example-token-123", handlerWithMiddleware)
	assert.Equal(t, http.StatusOK, recorder.Code)

	// test with wrong token:
	recorder = testRequest(true, "Bearer wrong-token", handlerWithMiddleware)
	assert.Equal(t, http.StatusUnauthorized, recorder.Code)

	// test with empty authorization header value:
	recorder = testRequest(true, "", handlerWithMiddleware)
	assert.Equal(t, http.StatusUnauthorized, recorder.Code)

	// test without authorization header at all:
	recorder = testRequest(false, "", handlerWithMiddleware)
	assert.Equal(t, http.StatusUnauthorized, recorder.Code)

	// test with invalid authorization header:
	recorder = testRequest(true, "Bearer is... who the hell is Bearer???", handlerWithMiddleware)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}
