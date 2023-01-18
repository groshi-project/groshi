package database

import (
	"context"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"time"
)

type Currency string // todo: think about currencies implementation
// e.g: how to check if string is currency???

const (
	CurrencyUSD Currency = "USD"
	CurrencyEUR Currency = "EUR"
	CurrencyRUB Currency = "RUB"
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID int64 `bun:",pk,autoincrement"`

	Username    string `bun:",unique"`
	Password    string
	IsSuperuser bool
}

func FetchUserByUsername(username string) (*User, error) {
	user := User{}
	err := Db.NewSelect().Model(&user).Where("username = ?", username).Scan(Ctx)
	return &user, err
}

func UserExists(username string) (bool, error) {
	return Db.NewSelect().Model((*User)(nil)).Where("username = ?", username).Exists(Ctx)
}

type Transaction struct {
	bun.BaseModel `bun:"table:transactions,alias:t"`

	ID   int64  `bun:",pk,autoincrement"`
	UUID string `bun:",unique,notnull"`

	Amount      float64  `bun:",notnull"`
	Currency    Currency `bun:",notnull"`
	Description string   `bun:",notnull"`

	OwnerId int64 `bun:",notnull"`
	Owner   *User `bun:"rel:belongs-to,join:owner_id=id"`

	Date time.Time `bun:",nullzero,notnull,default:current_timestamp"`
}

func FetchTransactionByUUID(uuid string) (*Transaction, error) {
	transaction := Transaction{}
	err := Db.NewSelect().Model(&transaction).Where("uuid = ?", uuid).Scan(Ctx)
	return &transaction, err
}

var _ bun.BeforeAppendModelHook = (*Transaction)(nil)

func (t *Transaction) BeforeAppendModel(_ context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		t.UUID = uuid.NewString()
	}
	return nil
}
