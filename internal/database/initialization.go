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
var Db *bun.DB

func Initialize() []error {
	var errors []error
	if _, err := Db.NewCreateTable().IfNotExists().Model((*User)(nil)).Exec(Ctx); err != nil {
		errors = append(errors, err)
	}
	if _, err := Db.NewCreateTable().IfNotExists().Model((*Transaction)(nil)).Exec(Ctx); err != nil {
		errors = append(errors, err)
	}
	return errors
}

func Connect(host string, port int, username string, password string, dbName string) error {
	dsn := fmt.Sprintf(
		"postgres://%v:%v@%v:%v/%v?sslmode=disable",
		username, password, host, port, dbName,
	)
	postgres := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	if err := postgres.Ping(); err != nil {
		return err
	}
	Db = bun.NewDB(postgres, pgdialect.New())
	return nil
}
