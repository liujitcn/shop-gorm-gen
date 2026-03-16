package data

import (
	"context"

	baseRepo "github.com/liujitcn/gorm-kit/repo"
	"github.com/liujitcn/shop-gorm-gen/models"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

// ShopBannerRepo 定义 ShopBanner 的基础仓储能力。
type ShopBannerRepo struct {
	baseRepo.BaseRepo[models.ShopBanner]
	*Data
}

// NewShopBannerRepo 创建 ShopBanner 基础仓储实例。
func NewShopBannerRepo(data *Data) *ShopBannerRepo {
	base := baseRepo.NewBaseRepo[models.ShopBanner](
		func(ctx context.Context) gen.Dao {
			return new(data.Query(ctx).ShopBanner.WithContext(ctx).DO)
		},
		func(ctx context.Context) field.Int64 {
			return data.Query(ctx).ShopBanner.ID
		},
		func(entity *models.ShopBanner) int64 {
			return entity.ID
		},
	)
	return &ShopBannerRepo{
		BaseRepo: base,
		Data:     data,
	}
}
