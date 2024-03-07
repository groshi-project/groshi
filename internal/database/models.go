package database

import (
	"context"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"time"
)

var EmptyUser = (*User)(nil)
var EmptyTransaction = (*Transaction)(nil)
var EmptyCurrency = (*Currency)(nil)

// User represents a user of the service.
type User struct {
	bun.BaseModel `bun:"table:users"`

	ID int64 `bun:"id,pk,autoincrement"`

	Username string `bun:"username,notnull"`
	Password string `bun:"password,notnull"`
}

// Transaction represents financial transaction created by the service user.
type Transaction struct {
	bun.BaseModel `bun:"table:transactions"`

	ID   int64     `bun:"id,pk,autoincrement"`
	UUID uuid.UUID `bun:"uuid,type:uuid,notnull,default:uuid_generate_v4()"`

	Amount int32 `bun:"amount,notnull"`

	CurrencyID int64    `bun:"currency_id,notnull"`
	Currency   Currency `bun:"rel:belongs-to,join:currency_id=id"`

	Description string `bun:"description,nullzero"`

	OwnerID int64 `bun:"owner_id,notnull"`
	Owner   User  `bun:"rel:belongs-to,join:owner_id=id"`

	Timestamp time.Time `bun:",notnull"`
	Timezone  string    `bun:",notnull"`

	CreatedAt time.Time `bun:",notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:",notnull,default:current_timestamp"`
}

var _ bun.BeforeAppendModelHook = (*Transaction)(nil)

func (t *Transaction) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		t.CreatedAt = time.Now()
	case *bun.UpdateQuery:
		t.UpdatedAt = time.Now()
	}
	return nil
}

// Currency represents currency supported by the service.
type Currency struct {
	bun.BaseModel `bun:"table:currencies"`

	ID int64 `bun:"id,pk,autoincrement"`

	Code   string  `bun:"code,notnull"`
	Symbol string  `bun:"symbol,notnull"`
	Rate   float64 `bun:"rate,notnull"`

	UpdatedAt time.Time `bun:",notnull,default:current_timestamp"`
}

var _ bun.BeforeAppendModelHook = (*Currency)(nil)

func (c *Currency) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) { // todo: is this switch needed? are other queries provided to before append model?
	case *bun.InsertQuery:
		c.UpdatedAt = time.Now()
	case *bun.UpdateQuery:
		c.UpdatedAt = time.Now()
	}

	return nil
}
