package data

import (
	"context"

	baseRepo "github.com/liujitcn/gorm-kit/repo"
	"github.com/liujitcn/shop-gorm-gen/models"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

// ShopHotItemRepo 定义 ShopHotItem 的基础仓储能力。
type ShopHotItemRepo struct {
	baseRepo.BaseRepo[models.ShopHotItem]
	*Data
}

// NewShopHotItemRepo 创建 ShopHotItem 基础仓储实例。
func NewShopHotItemRepo(data *Data) *ShopHotItemRepo {
	base := baseRepo.NewBaseRepo[models.ShopHotItem](
		func(ctx context.Context) gen.Dao {
			return new(data.Query(ctx).ShopHotItem.WithContext(ctx).DO)
		},
		func(ctx context.Context) field.Int64 {
			return data.Query(ctx).ShopHotItem.ID
		},
		func(entity *models.ShopHotItem) int64 {
			return entity.ID
		},
	)
	return &ShopHotItemRepo{
		BaseRepo: base,
		Data:     data,
	}
}
