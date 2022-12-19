package jwt

import (
	"errors"
	"net/http"
	"strings"
)

func parseJWTHeader(header http.Header) (string, error) {
	authHeader, ok := header["Authorization"]
	if !ok {
		return "", errors.New("missing authorization header")
	}
	data := authHeader[0]
	parts := strings.Split(data, " ")
	if len(parts) != 2 || parts[0] != "Bearer" { // this check doesn't panic with parts of length 0 'cause the second condition is not being checked
		return "", errors.New("invalid authorization header")
	}
	return parts[1], nil
}
