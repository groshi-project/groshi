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

func Initialize() error {
	// todo: check if table exists and then create table
	// don't use IfNotExists() construction
	// notify user about newly created tables
	_, errUsers := Db.NewCreateTable().IfNotExists().Model((*User)(nil)).Exec(Ctx)
	_, errTransactions := Db.NewCreateTable().IfNotExists().Model((*Transaction)(nil)).Exec(Ctx)
	if errUsers != nil || errTransactions != nil {
		return fmt.Errorf("could not create necessary tables (%v; %v)", errUsers, errTransactions)
	} // todo:
	return nil
}

func Connect(host string, port int, username string, password string, dbName string) error {
	dsn := fmt.Sprintf(
		"postgres://%v:%v@%v:%v/%v?sslmode=disable",
		username, password, host, port, dbName,
	)
	pgdb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	if err := pgdb.Ping(); err != nil {
		return err
		//logger.LoggerFatal.Fatalf("Could not ping database.")
	}
	Db = bun.NewDB(pgdb, pgdialect.New())
	return nil
}
