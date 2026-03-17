package data

import (
	"context"

	baseRepo "github.com/liujitcn/gorm-kit/repo"
	"github.com/liujitcn/shop-gorm-gen/models"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

// BaseApiRepo 定义 BaseApi 的基础仓储能力。
type BaseApiRepo struct {
	baseRepo.BaseRepo[models.BaseApi]
	*Data
}

// NewBaseApiRepo 创建 BaseApi 基础仓储实例。
func NewBaseApiRepo(data *Data) *BaseApiRepo {
	base := baseRepo.NewBaseRepo[models.BaseApi](
		func(ctx context.Context) gen.Dao {
			return new(data.Query(ctx).BaseApi.WithContext(ctx).DO)
		},
		func(ctx context.Context) field.Int64 {
			return data.Query(ctx).BaseApi.ID
		},
		func(entity *models.BaseApi) int64 {
			return entity.ID
		},
	)
	return &BaseApiRepo{
		BaseRepo: base,
		Data:     data,
	}
}
