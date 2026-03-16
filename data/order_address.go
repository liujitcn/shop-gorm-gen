package data

import (
	"context"

	baseRepo "github.com/liujitcn/gorm-kit/repo"
	"github.com/liujitcn/shop-gorm-gen/models"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

// OrderAddressRepo 定义 OrderAddress 的基础仓储能力。
type OrderAddressRepo struct {
	baseRepo.BaseRepo[models.OrderAddress]
	*Data
}

// NewOrderAddressRepo 创建 OrderAddress 基础仓储实例。
func NewOrderAddressRepo(data *Data) *OrderAddressRepo {
	base := baseRepo.NewBaseRepo[models.OrderAddress](
		func(ctx context.Context) gen.Dao {
			return new(data.Query(ctx).OrderAddress.WithContext(ctx).DO)
		},
		func(ctx context.Context) field.Int64 {
			return data.Query(ctx).OrderAddress.ID
		},
		func(entity *models.OrderAddress) int64 {
			return entity.ID
		},
	)
	return &OrderAddressRepo{
		BaseRepo: base,
		Data:     data,
	}
}
