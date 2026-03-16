package data

import (
	"context"

	baseRepo "github.com/liujitcn/gorm-kit/repo"
	"github.com/liujitcn/shop-gorm-gen/models"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

// GoodsSpecRepo 定义 GoodsSpec 的基础仓储能力。
type GoodsSpecRepo struct {
	baseRepo.BaseRepo[models.GoodsSpec]
	*Data
}

// NewGoodsSpecRepo 创建 GoodsSpec 基础仓储实例。
func NewGoodsSpecRepo(data *Data) *GoodsSpecRepo {
	base := baseRepo.NewBaseRepo[models.GoodsSpec](
		func(ctx context.Context) gen.Dao {
			return new(data.Query(ctx).GoodsSpec.WithContext(ctx).DO)
		},
		func(ctx context.Context) field.Int64 {
			return data.Query(ctx).GoodsSpec.ID
		},
		func(entity *models.GoodsSpec) int64 {
			return entity.ID
		},
	)
	return &GoodsSpecRepo{
		BaseRepo: base,
		Data:     data,
	}
}
