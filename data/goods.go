package data

import (
	"context"

	baseRepo "github.com/liujitcn/gorm-kit/repo"
	"github.com/liujitcn/shop-gorm-gen/models"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

// GoodsRepo 定义 Goods 的基础仓储能力。
type GoodsRepo struct {
	baseRepo.BaseRepo[models.Goods]
	*Data
}

// NewGoodsRepo 创建 Goods 基础仓储实例。
func NewGoodsRepo(data *Data) *GoodsRepo {
	base := baseRepo.NewBaseRepo[models.Goods](
		func(ctx context.Context) gen.Dao {
			return new(data.Query(ctx).Goods.WithContext(ctx).DO)
		},
		func(ctx context.Context) field.Int64 {
			return data.Query(ctx).Goods.ID
		},
		func(entity *models.Goods) int64 {
			return entity.ID
		},
	)
	return &GoodsRepo{
		BaseRepo: base,
		Data:     data,
	}
}
