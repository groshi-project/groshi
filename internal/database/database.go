package database

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

var (
	// Sample of the [User] database model.
	sampleUser = (*User)(nil)

	// Sample of the [Category] database model.
	sampleCategory = (*Category)(nil)

	// Sample of the [Currency] database model.
	sampleCurrency = (*Currency)(nil)

	// Sample of the [Transaction] database model.
	sampleTransaction = (*Transaction)(nil)
)

var (
	// Model samples which are used to create tables.
	models = []any{sampleUser, sampleUser, sampleCurrency, sampleTransaction}

	// PostgreSQL extensions that should be created.
	extensions = []string{"uuid-ossp"}
)

// Credentials represents PostgreSQL database credentials.
type Credentials struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

type Database interface {
	TestConnection() error
	Init(ctx context.Context) error

	UserQuerier
	CategoryQuerier
	CurrencyQuerier
	TransactionQuerier
}

// DefaultDatabase is the default implementation of the [Database] interface
type DefaultDatabase struct {
	client *bun.DB
}

// New creates a new instance of [DefaultDatabase] and returns a pointer to it.
func New(c Credentials) *DefaultDatabase {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", c.User, c.Password, c.Host, c.Port, c.Database)
	sqlDb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	bunDb := bun.NewDB(sqlDb, pgdialect.New())
	return &DefaultDatabase{
		client: bunDb,
	}
}

// TestConnection tests database connection.
func (d *DefaultDatabase) TestConnection() error {
	if err := d.client.Ping(); err != nil {
		return err
	}
	return nil
}

// Init creates all necessary tables and extensions if they do not exist.
func (d *DefaultDatabase) Init(ctx context.Context) error {
	// create necessary extensions if they do not exist:
	for _, extension := range extensions {
		if _, err := d.client.NewRaw("CREATE EXTENSION IF NOT EXISTS ?;", bun.Ident(extension)).Exec(ctx); err != nil {
			return err
		}
	}

	// create necessary tables if they do not exist:
	for _, model := range models {
		_, err := d.client.NewCreateTable().Model(model).IfNotExists().Exec(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}
