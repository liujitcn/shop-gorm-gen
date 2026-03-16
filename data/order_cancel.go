package data

import (
	"context"

	baseRepo "github.com/liujitcn/gorm-kit/repo"
	"github.com/liujitcn/shop-gorm-gen/models"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

// OrderCancelRepo 定义 OrderCancel 的基础仓储能力。
type OrderCancelRepo struct {
	baseRepo.BaseRepo[models.OrderCancel]
	*Data
}

// NewOrderCancelRepo 创建 OrderCancel 基础仓储实例。
func NewOrderCancelRepo(data *Data) *OrderCancelRepo {
	base := baseRepo.NewBaseRepo[models.OrderCancel](
		func(ctx context.Context) gen.Dao {
			return new(data.Query(ctx).OrderCancel.WithContext(ctx).DO)
		},
		func(ctx context.Context) field.Int64 {
			return data.Query(ctx).OrderCancel.ID
		},
		func(entity *models.OrderCancel) int64 {
			return entity.ID
		},
	)
	return &OrderCancelRepo{
		BaseRepo: base,
		Data:     data,
	}
}
