package data

import (
	"context"

	baseRepo "github.com/liujitcn/gorm-kit/repo"
	"github.com/liujitcn/shop-gorm-gen/models"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

// ShopHotGoodsRepo 定义 ShopHotGoods 的基础仓储能力。
type ShopHotGoodsRepo struct {
	baseRepo.BaseRepo[models.ShopHotGoods]
	*Data
}

// NewShopHotGoodsRepo 创建 ShopHotGoods 基础仓储实例。
func NewShopHotGoodsRepo(data *Data) *ShopHotGoodsRepo {
	base := baseRepo.NewBaseRepo[models.ShopHotGoods](
		func(ctx context.Context) gen.Dao {
			return new(data.Query(ctx).ShopHotGoods.WithContext(ctx).DO)
		},
		func(ctx context.Context) field.Int64 {
			// 当前通用 BaseRepo 仅支持单个 int64 主键，这里先复用 HotItemID 以保持生成代码可编译。
			return data.Query(ctx).ShopHotGoods.HotItemID
		},
		func(entity *models.ShopHotGoods) int64 {
			// 联合主键场景下无法完整表达，这里返回 HotItemID 作为兼容值。
			return entity.HotItemID
		},
	)
	return &ShopHotGoodsRepo{
		BaseRepo: base,
		Data:     data,
	}
}
