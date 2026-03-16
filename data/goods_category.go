package data

import (
	"context"

	baseRepo "github.com/liujitcn/gorm-kit/repo"
	"github.com/liujitcn/shop-gorm-gen/models"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

// GoodsCategoryRepo 定义 GoodsCategory 的基础仓储能力。
type GoodsCategoryRepo interface {
	baseRepo.BaseRepo[models.GoodsCategory]
}

type goodsCategoryRepo struct {
	baseRepo.BaseRepo[models.GoodsCategory]
}

// NewGoodsCategoryRepo 创建 GoodsCategory 基础仓储实例。
func NewGoodsCategoryRepo(data *Data) GoodsCategoryRepo {
	base := baseRepo.NewBaseRepo[models.GoodsCategory](
		func(ctx context.Context) gen.Dao {
			return new(data.Query(ctx).GoodsCategory.WithContext(ctx).DO)
		},
		func(ctx context.Context) field.Int64 {
			return data.Query(ctx).GoodsCategory.ID
		},
		func(entity *models.GoodsCategory) int64 {
			return entity.ID
		},
	)
	return &goodsCategoryRepo{
		BaseRepo: base,
	}
}
