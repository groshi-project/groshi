package middleware

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/groshi-project/groshi/internal/auth"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestTokenFromHeader(t *testing.T) {
	type out struct {
		token string
		err   error
	}
	testCases := []struct {
		input  string
		output out
	}{
		{"", out{"", errEmptyOrMissingAuthHeader}},
		{"a", out{"", errInvalidAuthHeader}},
		{"a b", out{"", errInvalidAuthHeader}},
		{"Bearer", out{"", errInvalidAuthHeader}},
		{"Bearer ", out{"", errInvalidAuthHeader}},
		{"Bearer  ", out{"", errInvalidAuthHeader}},
		{"Bearer   ", out{"", errInvalidAuthHeader}},
		{"Bearer token extra", out{"", errInvalidAuthHeader}},
		{"Bearer some-token", out{"some-token", nil}},
	}

	for i, testCase := range testCases {
		msg := fmt.Sprintf("input=\"%s\", case_index=%d", testCase.input, i)

		token, err := tokenFromHeader(testCase.input)
		assert.Equal(t, testCase.output.err, err, msg)
		assert.Equal(t, testCase.output.token, token, msg)
	}
}

type mockJWTAuthenticator struct {
	secretKey []byte
	*auth.DefaultJWTAuthenticator
}

func newMockJWTAuthenticator(secretKey string) *mockJWTAuthenticator {
	return &mockJWTAuthenticator{
		secretKey:               []byte(secretKey),
		DefaultJWTAuthenticator: auth.NewJWTAuthenticator(secretKey, time.Hour),
	}
}

func (m *mockJWTAuthenticator) CreateExpiredToken(username string) (string, time.Time, error) {
	issued := time.Now()
	expires := time.Now().Add(-time.Hour) // token expired an hour ago
	token := jwt.NewWithClaims(auth.JWTSigningMethod, jwt.MapClaims{
		"username": username,
		"exp":      expires.Unix(),
		"iat":      issued.Unix(),
	})

	tokenString, err := token.SignedString(m.secretKey)
	if err != nil {
		return "", expires, err
	}

	return tokenString, expires, nil
}

func (m *mockJWTAuthenticator) CreateNotYetValidToken(username string) (string, time.Time, error) {
	issued := time.Now().Add(time.Hour) // issued in one hour from now
	expires := time.Now().Add(2 * time.Hour)
	token := jwt.NewWithClaims(auth.JWTSigningMethod, jwt.MapClaims{
		"username": username,
		"exp":      expires.Unix(),
		"iat":      issued.Unix(),
	})

	tokenString, err := token.SignedString(m.secretKey)
	if err != nil {
		return "", expires, err
	}

	return tokenString, expires, nil
}

// testRequest creates a new [httptest.Recorder] and makes a test request to handler using it,
// then returns pointer to this recorder.
func testRequest(setAuthHeader bool, authHeader string, handler http.Handler) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	if setAuthHeader {
		req.Header.Set("Authorization", authHeader)
	}

	handler.ServeHTTP(rec, req)
	return rec
}

func TestNewJWT(t *testing.T) {
	const (
		testUsername  = "test-username"
		testSecretKey = "test-secret-key"
	)

	var (
		// Test handler which checks if username context value is equal to testUsername.
		handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			username, ok := r.Context().Value(UsernameContextKey).(string)
			if !ok {
				panic("username context key is missing")
			}
			assert.Equal(t, testUsername, username)
		})

		// JWT Authenticator used to issue and validate JWTs.
		jwtAuth = newMockJWTAuthenticator(testSecretKey)

		// Middleware which validates JWT.
		middleware = NewJWT(jwtAuth)
	)

	t.Run("call the handler with valid token", func(t *testing.T) {
		// create a new valid token for the test user:
		token, _, err := jwtAuth.CreateToken(testUsername)
		if err != nil {
			panic(err)
		}
		rec := testRequest(true, fmt.Sprintf("Bearer %s", token), middleware(handler))
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("call the handler with expired token", func(t *testing.T) {
		// create a new expired token for the test user:
		token, _, err := jwtAuth.CreateExpiredToken(testUsername)
		if err != nil {
			panic(err)
		}
		rec := testRequest(true, fmt.Sprintf("Bearer %s", token), middleware(handler))
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("call the handler with token which is not yet valid", func(t *testing.T) {
		token, _, err := jwtAuth.CreateNotYetValidToken(testUsername)
		if err != nil {
			panic(err)
		}
		rec := testRequest(true, fmt.Sprintf("Bearer %s", token), middleware(handler))
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("call the handler with invalid token", func(t *testing.T) {
		rec := testRequest(true, "Bearer invalid-token", middleware(handler))
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("call the handler with invalid authorization header", func(t *testing.T) {
		rec := testRequest(true, "Bearer is... who the hell is Bearer???", middleware(handler))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("call the handler with empty authorization header", func(t *testing.T) {
		rec := testRequest(true, "", middleware(handler))
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("call the handler without authorization header", func(t *testing.T) {
		rec := testRequest(false, "", middleware(handler))
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})
}
