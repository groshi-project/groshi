package database

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

// Database represents postgres database for the service.
type Database struct {
	Client *bun.DB
	Ctx    context.Context
}

// New creates a new instance of [Database] and returns a pointer to it.
func New() *Database {
	return &Database{
		Client: nil,
		Ctx:    context.TODO(),
	}
}

// Connect initializes connection to a database with provided credentials and verifies the connection.
func (d *Database) Connect(host string, port int, username string, password string, dbName string) error {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", username, password, host, port, dbName)
	db := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	// create a new database client:
	d.Client = bun.NewDB(db, pgdialect.New())

	// test the connection:
	if err := d.Client.Ping(); err != nil {
		return err
	}

	return nil
}

// Init creates all necessary tables and extensions if they do not exist.
func (d *Database) Init() error {
	var models = []any{ZeroUser, ZeroCategory, ZeroCurrency, ZeroTransaction}
	var extensions = []string{"uuid-ossp"}

	// create necessary extensions if they do not exist:
	for _, extension := range extensions {
		if _, err := d.Client.NewRaw("CREATE EXTENSION IF NOT EXISTS ?;", bun.Ident(extension)).Exec(d.Ctx); err != nil {
			return err
		}
	}

	// create necessary tables if they do not exist:
	for _, model := range models {
		_, err := d.Client.NewCreateTable().Model(model).IfNotExists().Exec(d.Ctx)
		if err != nil {
			return err
		}
	}

	return nil
}
