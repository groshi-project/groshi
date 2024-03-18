package database

import "github.com/uptrace/bun"

// User represents a user of the service.
type User struct {
	bun.BaseModel `bun:"table:users"`

	ID int64 `bun:"id,pk,autoincrement"`

	Username string `bun:"username,notnull"`
	Password string `bun:"password,notnull"`
}

func (d *Database) selectUserByUsernameQuery(username string) *bun.SelectQuery {
	return d.Client.NewSelect().Model(sampleUser).Where("username = ?", username)
}

func (d *Database) CreateUser(u *User) error {
	if _, err := d.Client.NewInsert().Model(u).Exec(d.Ctx); err != nil {
		return err
	}
	return nil
}

func (d *Database) SelectUserByUsername(username string, u *User) error {
	if err := d.selectUserByUsernameQuery(username).Scan(d.Ctx, u); err != nil {
		return err
	}
	return nil
}

func (d *Database) UserExistsByUsername(username string) (bool, error) {
	exists, err := d.selectUserByUsernameQuery(username).Exists(d.Ctx)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (d *Database) DeleteUserByUsername(username string) error {
	if _, err := d.Client.NewDelete().Model(sampleUser).Where("username = ?", username).Exec(d.Ctx); err != nil {
		return err
	}
	return nil
}
