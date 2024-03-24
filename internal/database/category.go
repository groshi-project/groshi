package database

import (
	"context"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

// Category database model.
type Category struct {
	bun.BaseModel `bun:"table:categories"`

	ID   int64     `bun:"id,pk,autoincrement"`
	UUID uuid.UUID `bun:"uuid,type:uuid,notnull,default:uuid_generate_v4()"`

	Owner   User  `bun:"rel:belongs-to,join:owner_id=id"`
	OwnerID int64 `bun:"owner_id,notnull"`

	Name string `bun:",notnull"`
}

// CategoryQuerier interface describes a type which executes database queries related to the [Category] model.
type CategoryQuerier interface {
	CreateCategory(ctx context.Context, c *Category) error
	CategoryExistsByUUID(ctx context.Context, uuid string) (bool, error)
	SelectCategoryByUUID(ctx context.Context, uuid string, c *Category) error
	SelectCategoriesByOwnerID(ctx context.Context, ownerID int64, c *[]Category) error
	UpdateCategory(ctx context.Context, c *Category) error
	DeleteCategoryByID(ctx context.Context, id int64) error
}

func (d *DefaultDatabase) selectCategoryByUUIDQuery(uuid string) *bun.SelectQuery {
	return d.client.NewSelect().Model(sampleCategory).Where("uuid = ?", uuid)
}

func (d *DefaultDatabase) selectCategoriesByOwnerIDQuery(ownerId int64) *bun.SelectQuery {
	return d.client.NewSelect().Model(sampleCategory).Where("owner_id = ?", ownerId)
}

func (d *DefaultDatabase) CreateCategory(ctx context.Context, c *Category) error {
	if _, err := d.client.NewInsert().Model(c).Exec(ctx); err != nil {
		return err
	}
	return nil
}

func (d *DefaultDatabase) CategoryExistsByUUID(ctx context.Context, uuid string) (bool, error) {
	panic("implement me")
}

func (d *DefaultDatabase) SelectCategoryByUUID(ctx context.Context, uuid string, c *Category) error {
	if err := d.selectCategoryByUUIDQuery(uuid).Scan(ctx, c); err != nil {
		return err
	}
	return nil
}

func (d *DefaultDatabase) SelectCategoriesByOwnerID(ctx context.Context, ownerID int64, c *[]Category) error {
	if err := d.selectCategoriesByOwnerIDQuery(ownerID).Scan(ctx, c); err != nil {
		return err
	}
	return nil
}

func (d *DefaultDatabase) UpdateCategory(ctx context.Context, c *Category) error {
	if _, err := d.client.NewUpdate().Model(c).WherePK().Exec(ctx); err != nil {
		return err
	}
	return nil
}

func (d *DefaultDatabase) DeleteCategoryByID(ctx context.Context, id int64) error {
	if _, err := d.client.NewDelete().Model(sampleCategory).Where("id = ?", id).Exec(ctx); err != nil {
		return err
	}
	return nil
}
