package database

import (
	"context"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"time"
)

// Transaction represents financial transaction created by the service user.
type Transaction struct {
	bun.BaseModel `bun:"table:transactions"`

	ID   int64     `bun:"id,pk,autoincrement"`
	UUID uuid.UUID `bun:"uuid,type:uuid,notnull,default:uuid_generate_v4()"`

	Amount int32 `bun:"amount,notnull"`

	Currency   Currency `bun:"rel:belongs-to,join:currency_id=id"`
	CurrencyID int64    `bun:"currency_id,notnull"`

	Description string `bun:"description,nullzero"`

	Category   Category `bun:"rel:belongs-to,join:category_id=id"`
	CategoryID int64    `bun:"category_id,notnull"`

	Owner   User  `bun:"rel:belongs-to,join:owner_id=id"`
	OwnerID int64 `bun:"owner_id,notnull"`

	// Timestamp of the transaction (when it happened)
	// UTC is always used as the timezone.
	Timestamp time.Time `bun:",notnull"`

	CreatedAt time.Time `bun:",notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:",notnull,default:current_timestamp"`
}

var _ bun.BeforeAppendModelHook = (*Transaction)(nil)

func (t *Transaction) BeforeAppendModel(_ context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		t.CreatedAt = time.Now()
	case *bun.UpdateQuery:
		t.UpdatedAt = time.Now()
	}
	return nil
}

func (d *Database) CreateTransaction(t *Transaction) error {
	if _, err := d.Client.NewInsert().Model(t).Exec(d.Ctx); err != nil {
		return err
	}
	return nil
}
