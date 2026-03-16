package data

import (
	"context"

	baseRepo "github.com/liujitcn/gorm-kit/repo"
	"github.com/liujitcn/shop-gorm-gen/models"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

// ShopHotGoodsRepo 定义 ShopHotGoods 的基础仓储能力。
type ShopHotGoodsRepo interface {
	baseRepo.BaseRepo[models.ShopHotGoods]
}

type shopHotGoodsRepo struct {
	baseRepo.BaseRepo[models.ShopHotGoods]
}

// NewShopHotGoodsRepo 创建 ShopHotGoods 基础仓储实例。
func NewShopHotGoodsRepo(data *Data) ShopHotGoodsRepo {
	base := baseRepo.NewBaseRepo[models.ShopHotGoods](
		func(ctx context.Context) gen.Dao {
			return new(data.Query(ctx).ShopHotGoods.WithContext(ctx).DO)
		},
		func(ctx context.Context) field.Int64 {
			// 当前通用仓储仅支持单一 int64 主键，这里先兼容使用复合主键中的 HotItemID。
			return data.Query(ctx).ShopHotGoods.HotItemID
		},
		func(entity *models.ShopHotGoods) int64 {
			// 当前通用仓储仅支持单一 int64 主键，这里先兼容使用复合主键中的 HotItemID。
			return entity.HotItemID
		},
	)
	return &shopHotGoodsRepo{
		BaseRepo: base,
	}
}
