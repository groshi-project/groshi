package integration_tests

import (
	"fmt"
	groshi "github.com/groshi-project/go-groshi"
	"math/rand"
)

func GenerateRandomString(length int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func GenerateCredentials() (username string, password string) {
	username = GenerateRandomString(5)
	password = "test-password-1234"
	return username, password
}

func NewPureGroshiClient(groshiSocket string) *groshi.GroshiAPIClient {
	return groshi.NewGroshiAPIClient(groshiSocket, "")
}

func NewGroshiClientWithUser(groshiSocket string) (username string, password string, client *groshi.GroshiAPIClient) {
	client = NewPureGroshiClient(groshiSocket)
	username, password = GenerateCredentials()
	if _, err := client.UserCreate(username, password); err != nil {
		panic(
			fmt.Sprintf("helper was unable to create user: %v", err),
		)
	}
	return username, password, client
}

func NewAuthorizedGroshiClientWithUser(groshiSocket string) (username string, password string, client *groshi.GroshiAPIClient) {
	username, password, client = NewGroshiClientWithUser(groshiSocket)
	if err := client.Auth(username, password); err != nil {
		panic(
			fmt.Sprintf("helper was unable to authorize user: %v", err),
		)
	}
	return username, password, client
}
