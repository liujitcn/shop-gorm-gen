package data

import (
	"context"

	baseRepo "github.com/liujitcn/gorm-kit/repo"
	"github.com/liujitcn/shop-gorm-gen/models"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

// UserCartRepo 定义 UserCart 的基础仓储能力。
type UserCartRepo struct {
	baseRepo.BaseRepo[models.UserCart]
	*Data
}

// NewUserCartRepo 创建 UserCart 基础仓储实例。
func NewUserCartRepo(data *Data) *UserCartRepo {
	base := baseRepo.NewBaseRepo[models.UserCart](
		func(ctx context.Context) gen.Dao {
			return new(data.Query(ctx).UserCart.WithContext(ctx).DO)
		},
		func(ctx context.Context) field.Int64 {
			return data.Query(ctx).UserCart.ID
		},
		func(entity *models.UserCart) int64 {
			return entity.ID
		},
	)
	return &UserCartRepo{
		BaseRepo: base,
		Data:     data,
	}
}
