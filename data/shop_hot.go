package data

import (
	"context"

	baseRepo "github.com/liujitcn/gorm-kit/repo"
	"github.com/liujitcn/shop-gorm-gen/models"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

// ShopHotRepo 定义 ShopHot 的基础仓储能力。
type ShopHotRepo struct {
	baseRepo.BaseRepo[models.ShopHot]
	*Data
}

// NewShopHotRepo 创建 ShopHot 基础仓储实例。
func NewShopHotRepo(data *Data) *ShopHotRepo {
	base := baseRepo.NewBaseRepo[models.ShopHot](
		func(ctx context.Context) gen.Dao {
			return new(data.Query(ctx).ShopHot.WithContext(ctx).DO)
		},
		func(ctx context.Context) field.Int64 {
			return data.Query(ctx).ShopHot.ID
		},
		func(entity *models.ShopHot) int64 {
			return entity.ID
		},
	)
	return &ShopHotRepo{
		BaseRepo: base,
		Data:     data,
	}
}
