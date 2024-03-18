package database

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

// Category represents transaction category.
type Category struct {
	bun.BaseModel `bun:"table:categories"`

	ID   int64     `bun:"id,pk,autoincrement"`
	UUID uuid.UUID `bun:"uuid,type:uuid,notnull,default:uuid_generate_v4()"`

	Owner   User  `bun:"rel:belongs-to,join:owner_id=id"`
	OwnerID int64 `bun:"owner_id,notnull"`

	Name string `bun:",notnull"`
}

func (d *Database) selectCategoryByUUIDQuery(uuid string) *bun.SelectQuery {
	return d.Client.NewSelect().Model(sampleCategory).Where("uuid = ?", uuid)
}

func (d *Database) selectCategoriesByOwnerIDQuery(ownerId int64) *bun.SelectQuery {
	return d.Client.NewSelect().Model(sampleCategory).Where("owner_id = ?", ownerId)
}

func (d *Database) CreateCategory(c *Category) error {
	if _, err := d.Client.NewInsert().Model(c).Exec(d.Ctx); err != nil {
		return err
	}
	return nil
}

func (d *Database) SelectCategoryByUUID(uuid string, c *Category) error {
	if err := d.selectCategoryByUUIDQuery(uuid).Scan(d.Ctx, c); err != nil {
		return err
	}
	return nil
}

func (d *Database) SelectCategoriesByOwnerID(ownerID int64, categories *[]Category) error {
	if err := d.selectCategoriesByOwnerIDQuery(ownerID).Scan(d.Ctx, categories); err != nil {
		return err
	}
	return nil
}

func (d *Database) UpdateCategory(c *Category) error {
	if _, err := d.Client.NewUpdate().Model(c).WherePK().Exec(d.Ctx); err != nil {
		return err
	}
	return nil
}

func (d *Database) DeleteCategoryByID(id int64) error {
	if _, err := d.Client.NewDelete().Model(sampleCategory).Where("id = ?", id).Exec(d.Ctx); err != nil {
		return err
	}
	return nil
}
