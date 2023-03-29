package database

import (
	"context"
	"github.com/google/uuid"
	"github.com/jieggii/groshi/internal/http/handles/datatypes"
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

func GetUser(username string) (*User, error) { // todo: name as SelectUser and remove .Scan(Ctx)
	user := User{}
	err := Db.NewSelect().Model(&user).Where("username = ?", username).Scan(Ctx)
	return &user, err
}

func UserExists(username string) (bool, error) {
	return Db.NewSelect().Model((*User)(nil)).Where("username = ?", username).Exists(Ctx)
}

func SelectUser(username string) *bun.SelectQuery {
	return Db.NewSelect().Model((*User)(nil)).Where("username = ?", username)
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

	Amount   float64            // amount in original currency (optional)
	Currency datatypes.Currency // original currency

	Description *string
}

func (t *Transaction) BeforeAppendModel(_ context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		t.UUID = uuid.NewString()
	case *bun.UpdateQuery:
		currentTime := time.Now()
		t.UpdatedAt = &currentTime
	}
	return nil
}

func GetTransaction(uuid string) (*Transaction, error) {
	transaction := Transaction{}
	err := Db.NewSelect().Model(&transaction).Where("uuid = ?", uuid).Scan(Ctx)
	return &transaction, err
}

// compile-time checks for hooks
var _ bun.BeforeAppendModelHook = (*Transaction)(nil)
