package data

import (
	"context"

	baseRepo "github.com/liujitcn/gorm-kit/repo"
	"github.com/liujitcn/shop-gorm-gen/models"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

// BaseRoleRepo 定义 BaseRole 的基础仓储能力。
type BaseRoleRepo struct {
	baseRepo.BaseRepo[models.BaseRole]
	*Data
}

// NewBaseRoleRepo 创建 BaseRole 基础仓储实例。
func NewBaseRoleRepo(data *Data) *BaseRoleRepo {
	base := baseRepo.NewBaseRepo[models.BaseRole](
		func(ctx context.Context) gen.Dao {
			return new(data.Query(ctx).BaseRole.WithContext(ctx).DO)
		},
		func(ctx context.Context) field.Int64 {
			return data.Query(ctx).BaseRole.ID
		},
		func(entity *models.BaseRole) int64 {
			return entity.ID
		},
	)
	return &BaseRoleRepo{
		BaseRepo: base,
		Data:     data,
	}
}
