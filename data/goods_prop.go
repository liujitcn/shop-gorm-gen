package data

import (
	"context"

	baseRepo "github.com/liujitcn/gorm-kit/repo"
	"github.com/liujitcn/shop-gorm-gen/models"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

// GoodsPropRepo 定义 GoodsProp 的基础仓储能力。
type GoodsPropRepo interface {
	baseRepo.BaseRepo[models.GoodsProp]
}

type goodsPropRepo struct {
	baseRepo.BaseRepo[models.GoodsProp]
}

// NewGoodsPropRepo 创建 GoodsProp 基础仓储实例。
func NewGoodsPropRepo(data *Data) GoodsPropRepo {
	base := baseRepo.NewBaseRepo[models.GoodsProp](
		func(ctx context.Context) gen.Dao {
			return new(data.Query(ctx).GoodsProp.WithContext(ctx).DO)
		},
		func(ctx context.Context) field.Int64 {
			return data.Query(ctx).GoodsProp.ID
		},
		func(entity *models.GoodsProp) int64 {
			return entity.ID
		},
	)
	return &goodsPropRepo{
		BaseRepo: base,
	}
}
