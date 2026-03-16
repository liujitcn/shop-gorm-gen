package data

import (
	"context"

	baseRepo "github.com/liujitcn/gorm-kit/repo"
	"github.com/liujitcn/shop-gorm-gen/models"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

// BaseMenuRepo 定义 BaseMenu 的基础仓储能力。
type BaseMenuRepo interface {
	baseRepo.BaseRepo[models.BaseMenu]
}

type baseMenuRepo struct {
	baseRepo.BaseRepo[models.BaseMenu]
}

// NewBaseMenuRepo 创建 BaseMenu 基础仓储实例。
func NewBaseMenuRepo(data *Data) BaseMenuRepo {
	base := baseRepo.NewBaseRepo[models.BaseMenu](
		func(ctx context.Context) gen.Dao {
			return new(data.Query(ctx).BaseMenu.WithContext(ctx).DO)
		},
		func(ctx context.Context) field.Int64 {
			return data.Query(ctx).BaseMenu.ID
		},
		func(entity *models.BaseMenu) int64 {
			return entity.ID
		},
	)
	return &baseMenuRepo{
		BaseRepo: base,
	}
}
