package data

import (
	"context"

	databaseGorm "github.com/liujitcn/kratos-kit/database/gorm"
	"github.com/liujitcn/shop-gorm-gen/models"
	"github.com/liujitcn/shop-gorm-gen/query"
	"gorm.io/gorm"
)

func init() {
	databaseGorm.RegisterMigrateModels(
		new(models.BaseAPI),
		new(models.BaseArea),
		new(models.BaseConfig),
		new(models.BaseDept),
		new(models.BaseDict),
		new(models.BaseDictItem),
		new(models.BaseJob),
		new(models.BaseJobLog),
		new(models.BaseLog),
		new(models.BaseMenu),
		new(models.BaseRole),
		new(models.BaseUser),
		new(models.CasbinRule),
		new(models.Goods),
		new(models.GoodsCategory),
		new(models.GoodsProp),
		new(models.GoodsSKU),
		new(models.GoodsSpec),
		new(models.Order),
		new(models.OrderAddress),
		new(models.OrderCancel),
		new(models.OrderGoods),
		new(models.OrderLogistics),
		new(models.OrderPayment),
		new(models.OrderRefund),
		new(models.PayBill),
		new(models.ShopBanner),
		new(models.ShopHot),
		new(models.ShopHotGoods),
		new(models.ShopHotItem),
		new(models.ShopService),
		new(models.UserAddress),
		new(models.UserCart),
		new(models.UserCollect),
		new(models.UserStore),
	)
}

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
