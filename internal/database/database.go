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

// Connect initializes connection to the given database with provided credentials and checks the connection.
func (d *Database) Connect(host string, port int, username string, password string, dbName string) error {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", username, password, host, port, dbName)
	db := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	d.Client = bun.NewDB(db, pgdialect.New())
	if err := d.Client.Ping(); err != nil {
		return err
	}
	return nil
}

// InitSchema creates all necessary tables if they do not exist,
func (d *Database) InitSchema() error {
	if err := d.createTableIfNotExists((*User)(nil)); err != nil {
		return err
	}
	if err := d.createTableIfNotExists((*Currency)(nil)); err != nil {
		return err
	}
	if err := d.createTableIfNotExists((*Transaction)(nil)); err != nil {
		return err
	}
	return nil
}

func (d *Database) createTableIfNotExists(model any) error {
	_, err := d.Client.NewCreateTable().Model((*User)(nil)).IfNotExists().Exec(d.Ctx)
	if err != nil {
		return err
	}
	return nil
}
