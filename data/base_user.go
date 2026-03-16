package data

import (
	"context"

	baseRepo "github.com/liujitcn/gorm-kit/repo"
	"github.com/liujitcn/shop-gorm-gen/models"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

// BaseUserRepo 定义 BaseUser 的基础仓储能力。
type BaseUserRepo interface {
	baseRepo.BaseRepo[models.BaseUser]
}

type baseUserRepo struct {
	baseRepo.BaseRepo[models.BaseUser]
}

// NewBaseUserRepo 创建 BaseUser 基础仓储实例。
func NewBaseUserRepo(data *Data) BaseUserRepo {
	base := baseRepo.NewBaseRepo[models.BaseUser](
		func(ctx context.Context) gen.Dao {
			return new(data.Query(ctx).BaseUser.WithContext(ctx).DO)
		},
		func(ctx context.Context) field.Int64 {
			return data.Query(ctx).BaseUser.ID
		},
		func(entity *models.BaseUser) int64 {
			return entity.ID
		},
	)
	return &baseUserRepo{
		BaseRepo: base,
	}
}
