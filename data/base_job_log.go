package data

import (
	"context"

	baseRepo "github.com/liujitcn/gorm-kit/repo"
	"github.com/liujitcn/shop-gorm-gen/models"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

// BaseJobLogRepo 定义 BaseJobLog 的基础仓储能力。
type BaseJobLogRepo struct {
	baseRepo.BaseRepo[models.BaseJobLog]
	*Data
}

// NewBaseJobLogRepo 创建 BaseJobLog 基础仓储实例。
func NewBaseJobLogRepo(data *Data) *BaseJobLogRepo {
	base := baseRepo.NewBaseRepo[models.BaseJobLog](
		func(ctx context.Context) gen.Dao {
			return new(data.Query(ctx).BaseJobLog.WithContext(ctx).DO)
		},
		func(ctx context.Context) field.Int64 {
			return data.Query(ctx).BaseJobLog.ID
		},
		func(entity *models.BaseJobLog) int64 {
			return entity.ID
		},
	)
	return &BaseJobLogRepo{
		BaseRepo: base,
		Data:     data,
	}
}
