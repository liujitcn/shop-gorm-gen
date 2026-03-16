package data

import (
	"context"

	baseRepo "github.com/liujitcn/gorm-kit/repo"
	"github.com/liujitcn/shop-gorm-gen/models"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

// BaseDictItemRepo 定义 BaseDictItem 的基础仓储能力。
type BaseDictItemRepo struct {
	baseRepo.BaseRepo[models.BaseDictItem]
	*Data
}

// NewBaseDictItemRepo 创建 BaseDictItem 基础仓储实例。
func NewBaseDictItemRepo(data *Data) *BaseDictItemRepo {
	base := baseRepo.NewBaseRepo[models.BaseDictItem](
		func(ctx context.Context) gen.Dao {
			return new(data.Query(ctx).BaseDictItem.WithContext(ctx).DO)
		},
		func(ctx context.Context) field.Int64 {
			return data.Query(ctx).BaseDictItem.ID
		},
		func(entity *models.BaseDictItem) int64 {
			return entity.ID
		},
	)
	return &BaseDictItemRepo{
		BaseRepo: base,
		Data:     data,
	}
}
