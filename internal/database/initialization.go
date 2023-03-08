package database

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jieggii/groshi/internal/loggers"
	"github.com/jieggii/groshi/internal/passhash"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

var Ctx = context.Background()
var Db *bun.DB

func createSuperuserIfNotExists(username string, password string) error {
	superUserExists, err := Db.NewSelect().Model((*User)(nil)).Where("username = ?", username).Exists(Ctx)
	if err != nil {
		return fmt.Errorf("could not check if superuser @%v exists: %v", username, err)
	}
	if !superUserExists {
		passwordHash, err := passhash.HashPassword(password)
		if err != nil {
			return fmt.Errorf("could not generate password hash for superuser @%v: %v\n", username, err)
		}
		user := &User{
			Username:    username,
			Password:    passwordHash,
			IsSuperuser: true,
		}
		_, err = Db.NewInsert().Model(user).Exec(Ctx)
		if err != nil {
			return fmt.Errorf("could not create new superuser @%v: %v", username, err)
		}
		loggers.Info.Printf("Created superuser @%v.", username)
	}
	return nil
}

func createTablesIfNotExist() error {
	// todo: check if table exists and then create table
	// don't use IfNotExists() construction
	// notify user about newly created tables
	_, errUsers := Db.NewCreateTable().IfNotExists().Model((*User)(nil)).Exec(Ctx)
	_, errTransactions := Db.NewCreateTable().IfNotExists().Model((*Transaction)(nil)).Exec(Ctx)
	if errUsers != nil || errTransactions != nil {
		return fmt.Errorf("could not create necessary tables (%v; %v)", errUsers, errTransactions)
	}
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
	//db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	return nil
}

func Initialize(superuserUsername string, superuserPassword string) error {
	if err := createTablesIfNotExist(); err != nil {
		return err
	}
	if err := createSuperuserIfNotExists(superuserUsername, superuserPassword); err != nil {
		return err
	}
	return nil
}
