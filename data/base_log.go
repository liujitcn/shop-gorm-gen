package data

import (
	"context"

	baseRepo "github.com/liujitcn/gorm-kit/repo"
	"github.com/liujitcn/shop-gorm-gen/models"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

// BaseLogRepo 定义 BaseLog 的基础仓储能力。
type BaseLogRepo interface {
	baseRepo.BaseRepo[models.BaseLog]
}

type baseLogRepo struct {
	baseRepo.BaseRepo[models.BaseLog]
}

// NewBaseLogRepo 创建 BaseLog 基础仓储实例。
func NewBaseLogRepo(data *Data) BaseLogRepo {
	base := baseRepo.NewBaseRepo[models.BaseLog](
		func(ctx context.Context) gen.Dao {
			return new(data.Query(ctx).BaseLog.WithContext(ctx).DO)
		},
		func(ctx context.Context) field.Int64 {
			return data.Query(ctx).BaseLog.ID
		},
		func(entity *models.BaseLog) int64 {
			return entity.ID
		},
	)
	return &baseLogRepo{
		BaseRepo: base,
	}
}
