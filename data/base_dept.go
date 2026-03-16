package data

import (
	"context"

	baseRepo "github.com/liujitcn/gorm-kit/repo"
	"github.com/liujitcn/shop-gorm-gen/models"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

// BaseDeptRepo 定义 BaseDept 的基础仓储能力。
type BaseDeptRepo interface {
	baseRepo.BaseRepo[models.BaseDept]
}

type baseDeptRepo struct {
	baseRepo.BaseRepo[models.BaseDept]
}

// NewBaseDeptRepo 创建 BaseDept 基础仓储实例。
func NewBaseDeptRepo(data *Data) BaseDeptRepo {
	base := baseRepo.NewBaseRepo[models.BaseDept](
		func(ctx context.Context) gen.Dao {
			return new(data.Query(ctx).BaseDept.WithContext(ctx).DO)
		},
		func(ctx context.Context) field.Int64 {
			return data.Query(ctx).BaseDept.ID
		},
		func(entity *models.BaseDept) int64 {
			return entity.ID
		},
	)
	return &baseDeptRepo{
		BaseRepo: base,
	}
}
