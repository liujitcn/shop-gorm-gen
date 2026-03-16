package data

import (
	"context"

	baseRepo "github.com/liujitcn/gorm-kit/repo"
	"github.com/liujitcn/shop-gorm-gen/models"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

// GoodsSkuRepo 定义 GoodsSKU 的基础仓储能力。
type GoodsSkuRepo struct {
	baseRepo.BaseRepo[models.GoodsSKU]
	*Data
}

// NewGoodsSkuRepo 创建 GoodsSKU 基础仓储实例。
func NewGoodsSkuRepo(data *Data) *GoodsSkuRepo {
	base := baseRepo.NewBaseRepo[models.GoodsSKU](
		func(ctx context.Context) gen.Dao {
			return new(data.Query(ctx).GoodsSKU.WithContext(ctx).DO)
		},
		func(ctx context.Context) field.Int64 {
			return data.Query(ctx).GoodsSKU.ID
		},
		func(entity *models.GoodsSKU) int64 {
			return entity.ID
		},
	)
	return &GoodsSkuRepo{
		BaseRepo: base,
		Data:     data,
	}
}
