package data

import (
	"context"

	baseRepo "github.com/liujitcn/gorm-kit/repo"
	"github.com/liujitcn/shop-gorm-gen/models"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

// PayBillRepo 定义 PayBill 的基础仓储能力。
type PayBillRepo interface {
	baseRepo.BaseRepo[models.PayBill]
}

type payBillRepo struct {
	baseRepo.BaseRepo[models.PayBill]
}

// NewPayBillRepo 创建 PayBill 基础仓储实例。
func NewPayBillRepo(data *Data) PayBillRepo {
	base := baseRepo.NewBaseRepo[models.PayBill](
		func(ctx context.Context) gen.Dao {
			return new(data.Query(ctx).PayBill.WithContext(ctx).DO)
		},
		func(ctx context.Context) field.Int64 {
			return data.Query(ctx).PayBill.ID
		},
		func(entity *models.PayBill) int64 {
			return entity.ID
		},
	)
	return &payBillRepo{
		BaseRepo: base,
	}
}
