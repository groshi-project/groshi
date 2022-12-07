package database

import "github.com/uptrace/bun"

type Currency string

//type User struct {
//	bun.BaseModel `bun:"table:users,alias:u"`
//
//	ID	 int64  `bun:",pk,autoincrement"`
//	Name string
//}

type User struct {
	bun.BaseModel `bun:"table:users"`
	ID            int64 `bun:",pk,autoincrement"`

	Username    string
	Password    string
	IsSuperuser bool

	BaseCurrency Currency
}

type Transaction struct {
	bun.BaseModel `bun:"table:transactions"`
	ID            int64 `bun:",pk,autoincrement"`

	UUID      string
	Amount    int
	Currency  Currency
	Timestamp int
}
