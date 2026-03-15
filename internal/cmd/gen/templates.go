package main

const dataFileTemplate = `package data

import (
	"context"

	databaseGorm "github.com/liujitcn/kratos-kit/database/gorm"
	"github.com/liujitcn/shop-gorm-gen/models"
	"github.com/liujitcn/shop-gorm-gen/query"
)

func init() {
	databaseGorm.RegisterMigrateModels(
{{- range . }}
		new(models.{{ .ModelName }}),
{{- end }}
	)
}

type contextTxKey struct{}

var txQueryKey = contextTxKey{}

type Data struct {
	query *query.Query
}

// NewData 初始化数据访问对象，并构建默认查询入口。
func NewData(c *databaseGorm.Client) *Data {
	d := &Data{
		query: query.Use(c.DB),
	}
	return d
}

// Transaction 定义事务执行能力，便于业务层按接口依赖。
type Transaction interface {
	Transaction(context.Context, func(ctx context.Context) error) error
}

// NewTransaction 创建事务执行器。
func NewTransaction(d *Data) Transaction {
	return d
}

// Transaction 在事务中执行传入函数，并将事务查询对象写入上下文。
func (d *Data) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return d.query.Transaction(func(tx *query.Query) error {
		// 将事务态查询对象注入上下文，仓储层可透明复用当前事务。
		ctx = context.WithValue(ctx, txQueryKey, tx)
		return fn(ctx)
	})
}

// Query 返回当前上下文对应的查询入口；若存在事务则优先返回事务查询对象。
func (d *Data) Query(ctx context.Context) *query.Query {
	if ctx == nil {
		return d.query
	}
	tx, ok := ctx.Value(txQueryKey).(*query.Query)
	if ok {
		return tx
	}
	return d.query
}
`

const initFileTemplate = `package data

import "github.com/google/wire"

// ProviderSet 定义 data 包依赖注入提供者集合。
var ProviderSet = wire.NewSet(
	NewData,
	NewTransaction,
{{- range . }}
	New{{ .RepoName }}Repo,
{{- end }}
)
`

const repoFileTemplate = `package data

import (
	"context"

	baseRepo "github.com/liujitcn/gorm-kit/repo"
	"github.com/liujitcn/shop-gorm-gen/models"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

// {{ .RepoName }}Repo 定义 {{ .ModelName }} 的基础仓储能力。
type {{ .RepoName }}Repo interface {
	baseRepo.BaseRepo[models.{{ .ModelName }}]
}

type {{ .RepoName | lowerFirst }}Repo struct {
	baseRepo.BaseRepo[models.{{ .ModelName }}]
}

// New{{ .RepoName }}Repo 创建 {{ .ModelName }} 基础仓储实例。
func New{{ .RepoName }}Repo(data *Data) {{ .RepoName }}Repo {
	base := baseRepo.NewBaseRepo[models.{{ .ModelName }}](
		func(ctx context.Context) gen.Dao {
			return new(data.Query(ctx).{{ .ModelName }}.WithContext(ctx).DO)
		},
		func(ctx context.Context) field.Int64 {
			return data.Query(ctx).{{ .ModelName }}.ID
		},
		func(entity *models.{{ .ModelName }}) int64 {
			return entity.ID
		},
	)
	return &{{ .RepoName | lowerFirst }}Repo{
		BaseRepo: base,
	}
}
`
