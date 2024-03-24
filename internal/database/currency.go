package database

import (
	"context"
	"github.com/uptrace/bun"
	"time"
)

var _ bun.BeforeAppendModelHook = (*Currency)(nil)

// Currency database model.
type Currency struct {
	bun.BaseModel `bun:"table:currencies"`

	ID int64 `bun:"id,pk,autoincrement"`

	Code   string  `bun:"code,notnull"`
	Symbol string  `bun:"symbol,notnull"`
	Rate   float64 `bun:"rate,notnull"`

	UpdatedAt time.Time `bun:",notnull,default:current_timestamp"`
}

func (c *Currency) BeforeAppendModel(_ context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		c.UpdatedAt = time.Now()
	case *bun.UpdateQuery:
		c.UpdatedAt = time.Now()
	}

	return nil
}

type CurrencyQuerier interface {
	SelectCurrencyByCode(ctx context.Context, code string, c *Currency) error
}

func (d *DefaultDatabase) selectCurrencyByCodeQuery(code string) *bun.SelectQuery {
	return d.client.NewSelect().Model(sampleCurrency).Where("code = ?", code)
}

func (d *DefaultDatabase) SelectCurrencyByCode(ctx context.Context, code string, c *Currency) error {
	if err := d.selectCurrencyByCodeQuery(code).Scan(ctx, c); err != nil {
		return err
	}
	return nil
}
