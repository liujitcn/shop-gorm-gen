package data

import (
	"context"

	baseRepo "github.com/liujitcn/gorm-kit/repo"
	"github.com/liujitcn/shop-gorm-gen/models"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

// OrderLogisticsRepo 定义 OrderLogistics 的基础仓储能力。
type OrderLogisticsRepo struct {
	baseRepo.BaseRepo[models.OrderLogistics]
	*Data
}

// NewOrderLogisticsRepo 创建 OrderLogistics 基础仓储实例。
func NewOrderLogisticsRepo(data *Data) *OrderLogisticsRepo {
	base := baseRepo.NewBaseRepo[models.OrderLogistics](
		func(ctx context.Context) gen.Dao {
			return new(data.Query(ctx).OrderLogistics.WithContext(ctx).DO)
		},
		func(ctx context.Context) field.Int64 {
			return data.Query(ctx).OrderLogistics.ID
		},
		func(entity *models.OrderLogistics) int64 {
			return entity.ID
		},
	)
	return &OrderLogisticsRepo{
		BaseRepo: base,
		Data:     data,
	}
}
