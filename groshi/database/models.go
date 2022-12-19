package database

import (
	"fmt"
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

	Currency Currency
}

func (u User) String() string {
	return fmt.Sprintf(
		"User<id=%v, @%v, isSuperuser=%v>",
		u.ID, u.Username, u.IsSuperuser,
	)
}

type Transaction struct {
	bun.BaseModel `bun:"table:transactions,alias:t"`

	ID   int64 `bun:",pk,autoincrement"`
	UUID string

	Amount   int64
	Currency Currency

	OwnerId int64
	Owner   *User `bun:"rel:belongs-to,join:owner_id=id"`

	Timestamp time.Time `bun:",nullzero,notnull,default:current_timestamp"`
}

func (t Transaction) String() string {
	return fmt.Sprintf(
		"Transaction<id=%v, by=%v, amount=%v (%v)>",
		t.ID, t.Owner, t.Amount, t.Currency,
	)
}
