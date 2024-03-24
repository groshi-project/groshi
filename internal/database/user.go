package database

import (
	"context"
	"github.com/uptrace/bun"
)

// User database model.
type User struct {
	bun.BaseModel `bun:"table:users"`

	ID int64 `bun:"id,pk,autoincrement"`

	Username string `bun:"username,notnull"`
	Password string `bun:"password,notnull"`
}

// UserQuerier interface describes a type which executes database queries related to the [User] model.
type UserQuerier interface {
	CreateUser(ctx context.Context, u *User) error
	UserExistsByUsername(ctx context.Context, username string) (bool, error)
	SelectUserByUsername(ctx context.Context, username string, u *User) error
	DeleteUserByUsername(ctx context.Context, username string) error
}

func (d *DefaultDatabase) selectUserByUsernameQuery(username string) *bun.SelectQuery {
	return d.client.NewSelect().Model(sampleUser).Where("username = ?", username)
}

func (d *DefaultDatabase) CreateUser(ctx context.Context, u *User) error {
	if _, err := d.client.NewInsert().Model(u).Exec(ctx); err != nil {
		return err
	}
	return nil
}

func (d *DefaultDatabase) SelectUserByUsername(ctx context.Context, username string, u *User) error {
	if err := d.selectUserByUsernameQuery(username).Scan(ctx, u); err != nil {
		return err
	}
	return nil
}

func (d *DefaultDatabase) UserExistsByUsername(ctx context.Context, username string) (bool, error) {
	exists, err := d.selectUserByUsernameQuery(username).Exists(ctx)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (d *DefaultDatabase) DeleteUserByUsername(ctx context.Context, username string) error {
	if _, err := d.client.NewDelete().Model(sampleUser).Where("username = ?", username).Exec(ctx); err != nil {
		return err
	}
	return nil
}
