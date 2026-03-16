package data

import (
	"context"

	baseRepo "github.com/liujitcn/gorm-kit/repo"
	"github.com/liujitcn/shop-gorm-gen/models"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

// OrderGoodsRepo 定义 OrderGoods 的基础仓储能力。
type OrderGoodsRepo interface {
	baseRepo.BaseRepo[models.OrderGoods]
}

type orderGoodsRepo struct {
	baseRepo.BaseRepo[models.OrderGoods]
}

// NewOrderGoodsRepo 创建 OrderGoods 基础仓储实例。
func NewOrderGoodsRepo(data *Data) OrderGoodsRepo {
	base := baseRepo.NewBaseRepo[models.OrderGoods](
		func(ctx context.Context) gen.Dao {
			return new(data.Query(ctx).OrderGoods.WithContext(ctx).DO)
		},
		func(ctx context.Context) field.Int64 {
			return data.Query(ctx).OrderGoods.ID
		},
		func(entity *models.OrderGoods) int64 {
			return entity.ID
		},
	)
	return &orderGoodsRepo{
		BaseRepo: base,
	}
}
