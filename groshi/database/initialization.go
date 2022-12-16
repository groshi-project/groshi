package database

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jieggii/groshi/groshi/logger"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

var ctx = context.Background()
var DB *bun.DB

func createSuperuserIfNotExist(username string, password string) error {
	superUserExists, err := DB.NewSelect().Model((*User)(nil)).Where("username = ?", username).Exists(ctx)
	if err != nil {
		return fmt.Errorf("could not check if superuser @%v exists: %v", username, err)
	}
	if !superUserExists {
		user := &User{
			Username:    username,
			Password:    password,
			IsSuperuser: true,
		}
		_, err := DB.NewInsert().Model(user).Exec(ctx)
		if err != nil {
			return fmt.Errorf("could not create new superuser @%v: %v", username, err)
		}
		logger.Info.Printf("Created superuser @%v.", username)
	}
	return nil
}

func createTablesIfNotExist() error {
	_, errUsers := DB.NewCreateTable().IfNotExists().Model((*User)(nil)).Exec(ctx)
	_, errTransactions := DB.NewCreateTable().IfNotExists().Model((*Transaction)(nil)).Exec(ctx)
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
		//logger.Fatal.Fatalf("Could not ping database.")
	}
	DB = bun.NewDB(pgdb, pgdialect.New())
	//db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	return nil
}

func Initialize(superuserUsername string, superuserPassword string) error {
	if err := createTablesIfNotExist(); err != nil {
		return err
	}
	if err := createSuperuserIfNotExist(superuserUsername, superuserPassword); err != nil {
		return err
	}
	return nil
}
