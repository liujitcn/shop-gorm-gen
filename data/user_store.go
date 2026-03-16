package data

import (
	"context"

	baseRepo "github.com/liujitcn/gorm-kit/repo"
	"github.com/liujitcn/shop-gorm-gen/models"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

// UserStoreRepo 定义 UserStore 的基础仓储能力。
type UserStoreRepo interface {
	baseRepo.BaseRepo[models.UserStore]
}

type userStoreRepo struct {
	baseRepo.BaseRepo[models.UserStore]
}

// NewUserStoreRepo 创建 UserStore 基础仓储实例。
func NewUserStoreRepo(data *Data) UserStoreRepo {
	base := baseRepo.NewBaseRepo[models.UserStore](
		func(ctx context.Context) gen.Dao {
			return new(data.Query(ctx).UserStore.WithContext(ctx).DO)
		},
		func(ctx context.Context) field.Int64 {
			return data.Query(ctx).UserStore.ID
		},
		func(entity *models.UserStore) int64 {
			return entity.ID
		},
	)
	return &userStoreRepo{
		BaseRepo: base,
	}
}
