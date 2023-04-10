package database

import (
	"context"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"time"
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID int64 `bun:",pk,autoincrement"`

	Username string `bun:",unique"`
	Password string

	BaseCurrency string `bun:",notnull"`
}

func SelectUser(username string) *bun.SelectQuery {
	return DB.NewSelect().Model((*User)(nil)).Where("username = ?", username)
}

type Transaction struct {
	bun.BaseModel `bun:"table:transactions,alias:t"`

	ID   int64  `bun:",pk,autoincrement"`
	UUID string `bun:",unique,notnull"`

	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt *time.Time

	Date time.Time `bun:",nullzero,notnull,default:current_timestamp"`

	OwnerId int64 `bun:",notnull"`
	Owner   *User `bun:"rel:belongs-to,join:owner_id=id"`

	BaseAmount float64 `bun:",notnull"` // amount in base currency

	Amount   float64 // amount in original currency (optional)
	Currency string  // todo: Currency // original currency

	Description *string
}

var _ bun.BeforeAppendModelHook = (*Transaction)(nil) // compile-time check for BeforeAppendModel hook

func (t *Transaction) BeforeAppendModel(_ context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		t.UUID = uuid.NewString() // generate transaction UUID on INSERT query
	case *bun.UpdateQuery:
		currentTime := time.Now()
		t.UpdatedAt = &currentTime // set transaction update time on UPDATE query
	}

	return nil
}

func SelectTransaction(uuid string) *bun.SelectQuery {
	return DB.NewSelect().Model((*Transaction)(nil)).Where("uuid = ?", uuid)
}
