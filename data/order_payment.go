package data

import (
	"context"

	baseRepo "github.com/liujitcn/gorm-kit/repo"
	"github.com/liujitcn/shop-gorm-gen/models"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

// OrderPaymentRepo 定义 OrderPayment 的基础仓储能力。
type OrderPaymentRepo struct {
	baseRepo.BaseRepo[models.OrderPayment]
	*Data
}

// NewOrderPaymentRepo 创建 OrderPayment 基础仓储实例。
func NewOrderPaymentRepo(data *Data) *OrderPaymentRepo {
	base := baseRepo.NewBaseRepo[models.OrderPayment](
		func(ctx context.Context) gen.Dao {
			return new(data.Query(ctx).OrderPayment.WithContext(ctx).DO)
		},
		func(ctx context.Context) field.Int64 {
			return data.Query(ctx).OrderPayment.ID
		},
		func(entity *models.OrderPayment) int64 {
			return entity.ID
		},
	)
	return &OrderPaymentRepo{
		BaseRepo: base,
		Data:     data,
	}
}
