package data

import (
	"context"

	baseRepo "github.com/liujitcn/gorm-kit/repo"
	"github.com/liujitcn/shop-gorm-gen/models"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

// BaseJobRepo 定义 BaseJob 的基础仓储能力。
type BaseJobRepo struct {
	baseRepo.BaseRepo[models.BaseJob]
	*Data
}

// NewBaseJobRepo 创建 BaseJob 基础仓储实例。
func NewBaseJobRepo(data *Data) *BaseJobRepo {
	base := baseRepo.NewBaseRepo[models.BaseJob](
		func(ctx context.Context) gen.Dao {
			return new(data.Query(ctx).BaseJob.WithContext(ctx).DO)
		},
		func(ctx context.Context) field.Int64 {
			return data.Query(ctx).BaseJob.ID
		},
		func(entity *models.BaseJob) int64 {
			return entity.ID
		},
	)
	return &BaseJobRepo{
		BaseRepo: base,
		Data:     data,
	}
}
