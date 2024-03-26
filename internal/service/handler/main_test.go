package handler

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/groshi-project/groshi/internal/database"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"time"
)

type mockPasswordAuthenticator struct{}

func newMockPasswordAuthenticator() *mockPasswordAuthenticator {
	return &mockPasswordAuthenticator{}
}

func (m *mockPasswordAuthenticator) HashPassword(password string) (string, error) {
	return fmt.Sprintf("hash(%s)", password), nil
}

func (m *mockPasswordAuthenticator) VerifyPassword(password string, hash string) (bool, error) {
	givenPasswordHash, err := m.HashPassword(password)
	if err != nil {
		return false, err
	}

	if givenPasswordHash == hash {
		return true, nil
	}

	return false, nil
}

type mockJWTAuthenticator struct{}

func newMockJWTAuthenticator() *mockJWTAuthenticator {
	return &mockJWTAuthenticator{}
}

func (m *mockJWTAuthenticator) CreateToken(username string) (token string, expires time.Time, err error) {
	return "", time.Time{}, nil
}

func (m *mockJWTAuthenticator) VerifyToken(token string) (jwt.MapClaims, error) {
	return nil, nil
}

type mockDatabase struct {
	users []*database.User

	categories []*database.Category
}

func newMockDatabase() *mockDatabase {
	return &mockDatabase{
		users:      make([]*database.User, 0),
		categories: make([]*database.Category, 0),
	}
}

func (m *mockDatabase) TestConnection() error {
	panic("implement me")
}

func (m *mockDatabase) Init(ctx context.Context) error {
	panic("implement me")
}

func (m *mockDatabase) CreateUser(ctx context.Context, u *database.User) error {
	if u.ID == 0 {
		u.ID = int64(rand.Intn(9999) + 1)
	}
	m.users = append(m.users, u)
	return nil
}

func (m *mockDatabase) UserExistsByUsername(ctx context.Context, username string) (bool, error) {
	for _, user := range m.users {
		if username == user.Username {
			return true, nil
		}
	}
	return false, nil
}

func (m *mockDatabase) SelectUserByUsername(ctx context.Context, username string, u *database.User) error {
	for _, user := range m.users {
		if username == user.Username {
			*u = *user
			return nil
		}
	}
	return sql.ErrNoRows
}

func (m *mockDatabase) DeleteUserByUsername(ctx context.Context, username string) error {
	userIndex := -1
	for i, user := range m.users {
		if username == user.Username {
			userIndex = i
			break
		}
	}
	if userIndex == -1 {
		return nil
	}

	m.users[userIndex] = m.users[len(m.users)-1]
	m.users = m.users[:len(m.users)-1]

	return nil
}

func (m *mockDatabase) CreateCategory(ctx context.Context, c *database.Category) error {
	if c.UUID.String() == "" {
		c.UUID = uuid.New()
	} else {
		exists, err := m.CategoryExistsByUUID(ctx, c.UUID.String())
		if err != nil {
			panic(err)
		}
		if exists {
			panic(fmt.Errorf("category with UUID %s already exists", c.UUID.String()))
		}
	}

	m.categories = append(m.categories, c)
	return nil
}

func (m *mockDatabase) CategoryExistsByUUID(ctx context.Context, uuid string) (bool, error) {
	for _, category := range m.categories {
		if category.UUID.String() == uuid {
			return true, nil
		}
	}
	return false, nil
}

func (m *mockDatabase) SelectCategoryByUUID(ctx context.Context, uuid string, c *database.Category) error {
	for _, category := range m.categories {
		if category.UUID.String() == uuid {
			*c = *category
			return nil
		}
	}
	return sql.ErrNoRows
}

func (m *mockDatabase) SelectCategoriesByOwnerID(ctx context.Context, ownerID int64, c *[]database.Category) error {
	panic("implement me")
}

func (m *mockDatabase) UpdateCategory(ctx context.Context, c *database.Category) error {
	panic("implement me")
}

func (m *mockDatabase) DeleteCategoryByID(ctx context.Context, id int64) error {
	panic("implement me")
}

func (m *mockDatabase) SelectCurrencyByCode(ctx context.Context, code string, c *database.Currency) error {
	panic("implement me")
}

func (m *mockDatabase) CreateTransaction(ctx context.Context, transaction *database.Transaction) error {
	panic("implement me")
}

func newTestHandler() *Handler {
	return New(
		newMockDatabase(),
		newMockJWTAuthenticator(),
		newMockPasswordAuthenticator(),
		log.New(io.Discard, "", 0),
	)
}

func testRequest(ctx context.Context, params any, handlerFunc http.HandlerFunc) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	body, err := json.Marshal(params)
	if err != nil {
		panic(err)
	}

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	handlerFunc(rec, req.WithContext(ctx))

	return rec
}
