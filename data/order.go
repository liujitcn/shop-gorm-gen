package data

import (
	"context"

	baseRepo "github.com/liujitcn/gorm-kit/repo"
	"github.com/liujitcn/shop-gorm-gen/models"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

// OrderRepo 定义 Order 的基础仓储能力。
type OrderRepo interface {
	baseRepo.BaseRepo[models.Order]
}

type orderRepo struct {
	baseRepo.BaseRepo[models.Order]
}

// NewOrderRepo 创建 Order 基础仓储实例。
func NewOrderRepo(data *Data) OrderRepo {
	base := baseRepo.NewBaseRepo[models.Order](
		func(ctx context.Context) gen.Dao {
			return new(data.Query(ctx).Order.WithContext(ctx).DO)
		},
		func(ctx context.Context) field.Int64 {
			return data.Query(ctx).Order.ID
		},
		func(entity *models.Order) int64 {
			return entity.ID
		},
	)
	return &orderRepo{
		BaseRepo: base,
	}
}
