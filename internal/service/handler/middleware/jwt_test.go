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
		assert.ErrorIs(t, err, errEmptyAuthHeader)
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

	// test with valid token:
	recorder1 := httptest.NewRecorder()
	req1 := httptest.NewRequest(http.MethodPost, "/", nil)
	req1.Header.Set("Authorization", "Bearer example-token-123")

	handlerWithMiddleware.ServeHTTP(recorder1, req1)
	assert.Equal(t, http.StatusOK, recorder1.Code)

	// test with wrong token:
	recorder2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodPost, "/", nil)
	req2.Header.Set("Authorization", "Bearer wrong-token")

	handlerWithMiddleware.ServeHTTP(recorder2, req2)
	assert.Equal(t, http.StatusUnauthorized, recorder2.Code)

	// test without "Authorization" header
	recorder3 := httptest.NewRecorder()
	req3 := httptest.NewRequest(http.MethodPost, "/", nil)
	handlerWithMiddleware.ServeHTTP(recorder3, req3)
	assert.Equal(t, http.StatusBadRequest, recorder3.Code)

	// test with invalid "Authorization" header:
	recorder4 := httptest.NewRecorder()
	req4 := httptest.NewRequest(http.MethodPost, "/", nil)
	req2.Header.Set("Authorization", "hello! I am temmie www!")
	handlerWithMiddleware.ServeHTTP(recorder4, req4)
	assert.Equal(t, http.StatusBadRequest, recorder4.Code)
}
