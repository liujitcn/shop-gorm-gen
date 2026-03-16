package data

import (
	"context"

	baseRepo "github.com/liujitcn/gorm-kit/repo"
	"github.com/liujitcn/shop-gorm-gen/models"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

// BaseDictRepo 定义 BaseDict 的基础仓储能力。
type BaseDictRepo struct {
	baseRepo.BaseRepo[models.BaseDict]
	*Data
}

// NewBaseDictRepo 创建 BaseDict 基础仓储实例。
func NewBaseDictRepo(data *Data) *BaseDictRepo {
	base := baseRepo.NewBaseRepo[models.BaseDict](
		func(ctx context.Context) gen.Dao {
			return new(data.Query(ctx).BaseDict.WithContext(ctx).DO)
		},
		func(ctx context.Context) field.Int64 {
			return data.Query(ctx).BaseDict.ID
		},
		func(entity *models.BaseDict) int64 {
			return entity.ID
		},
	)
	return &BaseDictRepo{
		BaseRepo: base,
		Data:     data,
	}
}
