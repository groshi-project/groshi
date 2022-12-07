package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

var ctx = context.Background()
var db *bun.DB

func createTables(db *bun.DB) error {
	_, errUsers := db.NewCreateTable().IfNotExists().Model((*User)(nil)).Exec(ctx)
	_, errTransactions := db.NewCreateTable().IfNotExists().Model((*Transaction)(nil)).Exec(ctx)
	if errUsers != nil || errTransactions != nil {
		return errors.New("could not create necessary tables")
	}
	return nil
}

func Initialize(host string, port int, username string, password string, dbName string) error {
	dsn := fmt.Sprintf(
		"postgres://%v:%v@%v:%v/%v?sslmode=disable",
		username, password, host, port, dbName,
	)
	pgdb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	if err := pgdb.Ping(); err != nil {
		return err
		//logger.Fatal.Fatalf("Could not ping database.")
	}
	db = bun.NewDB(pgdb, pgdialect.New())

	// Print all queries to stdout.
	//db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	if err := createTables(db); err != nil {
		return err
	}
	return nil
}
