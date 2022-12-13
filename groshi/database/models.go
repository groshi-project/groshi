package database

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"time"
)

type Currency string

type User struct {
	bun.BaseModel `bun:"table:users"` // todo: alias
	ID            int64               `bun:",pk,autoincrement"`

	Username    string // todo: unique
	Password    string
	IsSuperuser bool

	BaseCurrency Currency
}

type Transaction struct {
	bun.BaseModel `bun:"table:transactions"`
	ID            int64 `bun:",pk,autoincrement"`

	UUID      uuid.UUID `bun:"type:uuid,default:uuid_generate_v4()"` // CREATE EXTENSION IF NOT EXISTS "uuid-ossp"; (https://bun.uptrace.dev/postgres/postgres-uuid-generate.html#uuid-in-postgresql)
	Amount    int
	Currency  Currency
	Timestamp time.Time `bun:",nullzero,notnull,default:current_timestamp"`
}
