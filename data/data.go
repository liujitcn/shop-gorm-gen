package data

import (
	"context"

	databaseGorm "github.com/liujitcn/kratos-kit/database/gorm"
	"github.com/liujitcn/shop-gorm-gen/query"
	"gorm.io/gorm"
)

type contextTxKey struct{}

type Data struct {
	query *query.Query
	db    *gorm.DB
}

// NewData .
func NewData(c *databaseGorm.Client) *Data {
	db := c.DB
	d := &Data{
		query: query.Use(db),
		db:    db,
	}
	return d
}

type Transaction interface {
	Transaction(context.Context, func(ctx context.Context) error) error
}

func NewTransaction(d *Data) Transaction {
	return d
}

func (d *Data) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return d.query.Transaction(func(tx *query.Query) error {
		ctx = context.WithValue(ctx, contextTxKey{}, tx)
		return fn(ctx)
	})
}

func (d *Data) Query(ctx context.Context) *query.Query {
	tx, ok := ctx.Value(contextTxKey{}).(*query.Query)
	if ok {
		return tx
	}
	return d.query
}
