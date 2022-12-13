package database

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

var ctx = context.Background()
var DB *bun.DB

//func createSuperuser(username string, password string) error {
//	user := new(User)
//	err := DB.NewSelect().Model(user).Where("username = ?", username).Scan(ctx)
//
//	user := &User{
//		Username:    username,
//		Password:    password,
//		IsSuperuser: true,
//	}
//	res, err := DB.NewInsert().Model(user).Exec(ctx)
//	return nil
//}

func createTables() error {
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

func Initialize() error {
	if err := createTables(); err != nil {
		return err
	}
	user := new(User)
	err := DB.NewSelect().Model(user).Where("username = ?", "root").Scan(ctx)
	fmt.Printf("|%v|\n", err)
	fmt.Println(user)
	return nil
}
