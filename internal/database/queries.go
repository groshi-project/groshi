package database

import "github.com/uptrace/bun"

var sampleUser = (*User)(nil)
var sampleCategory = (*Category)(nil)
var sampleCurrency = (*Currency)(nil)
var sampleTransaction = (*Transaction)(nil)

// todo: naming (userCreate or createUser)

func (d *Database) createUserQuery(user *User) *bun.InsertQuery {
	return d.Client.NewInsert().Model(user)
}

func (d *Database) CreateUser(user *User) error {
	if _, err := d.createUserQuery(user).Exec(d.Ctx); err != nil {
		return err
	}
	return nil
}

func (d *Database) selectUserByUsernameQuery(username string) *bun.SelectQuery {
	return d.Client.NewSelect().Model(sampleUser).Where("username = ?", username)
}

func (d *Database) SelectUserByUsername(username string, user *User) error {
	if err := d.selectUserByUsernameQuery(username).Scan(d.Ctx, user); err != nil {
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

func (d *Database) deleteUserByUsernameQuery(username string) *bun.DeleteQuery {
	return d.Client.NewDelete().Model(sampleUser).Where("username = ?", username)
}

func (d *Database) DeleteUserByUsername(username string) error {
	if _, err := d.deleteUserByUsernameQuery(username).Exec(d.Ctx); err != nil {
		return err
	}
	return nil
}

func (d *Database) createCategoryQuery(category *Category) *bun.InsertQuery {
	return d.Client.NewInsert().Model(category)
}

func (d *Database) CreateCategory(category *Category) error {
	if _, err := d.createCategoryQuery(category).Exec(d.Ctx); err != nil {
		return err
	}
	return nil
}

func (d *Database) selectCategoryByUUIDQuery(uuid string) *bun.SelectQuery {
	return d.Client.NewSelect().Model(sampleCategory).Where("uuid = ?", uuid)
}

func (d *Database) SelectCategoryByUUID(uuid string, v *Category) error {
	if err := d.selectCategoryByUUIDQuery(uuid).Scan(d.Ctx, v); err != nil {
		return err
	}
	return nil
}

func (d *Database) selectCategoriesByOwnerIDQuery(ownerId int64) *bun.SelectQuery {
	return d.Client.NewSelect().Model(sampleCategory).Where("owner_id = ?", ownerId)
}

func (d *Database) SelectCategoriesByOwnerID(ownerID int64, categories *[]Category) error {
	if err := d.selectCategoriesByOwnerIDQuery(ownerID).Scan(d.Ctx, categories); err != nil {
		return err
	}
	return nil
}

func (d *Database) updateCategoryQuery(category *Category) *bun.UpdateQuery {
	return d.Client.NewUpdate().Model(category).WherePK()
}

func (d *Database) UpdateCategory(category *Category) error {
	if _, err := d.updateCategoryQuery(category).Exec(d.Ctx); err != nil {
		return err
	}
	return nil
}

func (d *Database) deleteCategoryByIDQuery(id int64) *bun.DeleteQuery {
	return d.Client.NewDelete().Model(sampleCategory).Where("id = ?", id)
}

func (d *Database) DeleteCategoryByID(id int64) error {
	if _, err := d.deleteCategoryByIDQuery(id).Exec(d.Ctx); err != nil {
		return err
	}
	return nil
}

func (d *Database) selectCurrencyByCodeQuery(code string) *bun.SelectQuery {
	return d.Client.NewSelect().Model(sampleCurrency).Where("code = ?", code)
}

func (d *Database) SelectCurrencyByCode(code string, currency *Currency) error {
	if err := d.selectCurrencyByCodeQuery(code).Scan(d.Ctx, currency); err != nil {
		return err
	}
	return nil
}

func (d *Database) insertTransactionQuery(t *Transaction) *bun.InsertQuery {
	return d.Client.NewInsert().Model(t)
}

func (d *Database) CreateTransaction(t *Transaction) error {
	if _, err := d.insertTransactionQuery(t).Exec(d.Ctx); err != nil {
		return err
	}
	return nil
}
