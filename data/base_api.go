package data

import (
	"context"

	baseRepo "github.com/liujitcn/gorm-kit/repo"
	"github.com/liujitcn/shop-gorm-gen/models"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

// BaseApiRepo 定义 BaseAPI 的基础仓储能力。
type BaseApiRepo struct {
	baseRepo.BaseRepo[models.BaseAPI]
	*Data
}

// NewBaseApiRepo 创建 BaseAPI 基础仓储实例。
func NewBaseApiRepo(data *Data) *BaseApiRepo {
	base := baseRepo.NewBaseRepo[models.BaseAPI](
		func(ctx context.Context) gen.Dao {
			return new(data.Query(ctx).BaseAPI.WithContext(ctx).DO)
		},
		func(ctx context.Context) field.Int64 {
			return data.Query(ctx).BaseAPI.ID
		},
		func(entity *models.BaseAPI) int64 {
			return entity.ID
		},
	)
	return &BaseApiRepo{
		BaseRepo: base,
		Data:     data,
	}
}
