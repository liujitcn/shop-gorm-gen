package data

import (
	"context"

	baseRepo "github.com/liujitcn/gorm-kit/repo"
	"github.com/liujitcn/shop-gorm-gen/models"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

// UserAddressRepo 定义 UserAddress 的基础仓储能力。
type UserAddressRepo struct {
	baseRepo.BaseRepo[models.UserAddress]
	*Data
}

// NewUserAddressRepo 创建 UserAddress 基础仓储实例。
func NewUserAddressRepo(data *Data) *UserAddressRepo {
	base := baseRepo.NewBaseRepo[models.UserAddress](
		func(ctx context.Context) gen.Dao {
			return new(data.Query(ctx).UserAddress.WithContext(ctx).DO)
		},
		func(ctx context.Context) field.Int64 {
			return data.Query(ctx).UserAddress.ID
		},
		func(entity *models.UserAddress) int64 {
			return entity.ID
		},
	)
	return &UserAddressRepo{
		BaseRepo: base,
		Data:     data,
	}
}
