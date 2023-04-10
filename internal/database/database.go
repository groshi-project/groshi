package database

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

var Ctx = context.Background()
var DB *bun.DB

// Connect initializes connection to database.
func Connect(host string, port int, username string, password string, dbName string) error {
	dsn := fmt.Sprintf(
		"postgres://%v:%v@%v:%v/%v?sslmode=disable",
		username, password, host, port, dbName,
	)
	postgres := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	if err := postgres.Ping(); err != nil {
		return err
	}
	DB = bun.NewDB(postgres, pgdialect.New())

	return nil
}

// Init creates necessary tables in database if they don't exist.
func Init() error {
	if _, err := DB.NewCreateTable().
		IfNotExists().
		Model((*User)(nil)).
		Exec(Ctx); err != nil {
		return err
	}
	if _, err := DB.NewCreateTable().
		IfNotExists().
		Model((*Transaction)(nil)).
		Exec(Ctx); err != nil {
		return err
	}

	return nil
}
