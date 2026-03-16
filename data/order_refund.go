package data

import (
	"context"

	baseRepo "github.com/liujitcn/gorm-kit/repo"
	"github.com/liujitcn/shop-gorm-gen/models"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

// OrderRefundRepo 定义 OrderRefund 的基础仓储能力。
type OrderRefundRepo struct {
	baseRepo.BaseRepo[models.OrderRefund]
	*Data
}

// NewOrderRefundRepo 创建 OrderRefund 基础仓储实例。
func NewOrderRefundRepo(data *Data) *OrderRefundRepo {
	base := baseRepo.NewBaseRepo[models.OrderRefund](
		func(ctx context.Context) gen.Dao {
			return new(data.Query(ctx).OrderRefund.WithContext(ctx).DO)
		},
		func(ctx context.Context) field.Int64 {
			return data.Query(ctx).OrderRefund.ID
		},
		func(entity *models.OrderRefund) int64 {
			return entity.ID
		},
	)
	return &OrderRefundRepo{
		BaseRepo: base,
		Data:     data,
	}
}
