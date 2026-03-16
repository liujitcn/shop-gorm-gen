package data

import (
	"context"

	baseRepo "github.com/liujitcn/gorm-kit/repo"
	"github.com/liujitcn/shop-gorm-gen/models"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

// UserCollectRepo 定义 UserCollect 的基础仓储能力。
type UserCollectRepo interface {
	baseRepo.BaseRepo[models.UserCollect]
}

type userCollectRepo struct {
	baseRepo.BaseRepo[models.UserCollect]
}

// NewUserCollectRepo 创建 UserCollect 基础仓储实例。
func NewUserCollectRepo(data *Data) UserCollectRepo {
	base := baseRepo.NewBaseRepo[models.UserCollect](
		func(ctx context.Context) gen.Dao {
			return new(data.Query(ctx).UserCollect.WithContext(ctx).DO)
		},
		func(ctx context.Context) field.Int64 {
			return data.Query(ctx).UserCollect.ID
		},
		func(entity *models.UserCollect) int64 {
			return entity.ID
		},
	)
	return &userCollectRepo{
		BaseRepo: base,
	}
}
