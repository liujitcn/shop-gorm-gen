package data

import (
	"context"

	baseRepo "github.com/liujitcn/gorm-kit/repo"
	"github.com/liujitcn/shop-gorm-gen/models"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

// BaseConfigRepo 定义 BaseConfig 的基础仓储能力。
type BaseConfigRepo struct {
	baseRepo.BaseRepo[models.BaseConfig]
	*Data
}

// NewBaseConfigRepo 创建 BaseConfig 基础仓储实例。
func NewBaseConfigRepo(data *Data) *BaseConfigRepo {
	base := baseRepo.NewBaseRepo[models.BaseConfig](
		func(ctx context.Context) gen.Dao {
			return new(data.Query(ctx).BaseConfig.WithContext(ctx).DO)
		},
		func(ctx context.Context) field.Int64 {
			return data.Query(ctx).BaseConfig.ID
		},
		func(entity *models.BaseConfig) int64 {
			return entity.ID
		},
	)
	return &BaseConfigRepo{
		BaseRepo: base,
		Data:     data,
	}
}
