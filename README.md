# shop-gorm-gen

用于根据数据库表结构自动生成 GORM 的 `models` 和 `query` 代码。

## 生成方式

1. 准备依赖：

```bash
make init
```

2. 在项目目录执行生成：

```bash
make gen
```

如需指定数据库连接，可通过环境变量覆盖默认 DSN：

```bash
GORM_GEN_DSN="<你的dsn>" make gen
```

## 当前配置说明

- 生成入口：`internal/cmd/gen/main.go`
- 输出目录：`query`
- 模型包名：`models`
- 默认会同时生成 `models` 与 `query`
- 模型命名策略：按表名转 CamelCase，不做单数化（例如 `goods -> Goods`、`order_goods -> OrderGoods`）

## Codex 文档规则

- 每次修改代码后，必须同步检查并优化 `AGENTS.md` 与 `README.md`。
- 若命令、流程、生成逻辑或行为发生变化，必须在同次变更中更新文档。
- 后续新增或修改代码时，代码注释统一使用中文。

## 生成结果

- `models/*.gen.go`：数据模型代码
- `query/*.gen.go`：查询构造与 DAO 代码
- `query/gen.go`：查询入口

注意：默认 DSN 仅用于本地开发，建议使用 `GORM_GEN_DSN` 指向你自己的数据库后再执行生成。

## 泛型仓储（repo/base_repo.go）

当前仓库提供了可复用的泛型仓储接口与构造方法：

- `BaseRepo[T, C]`
- `NewBaseRepo[T, C](...)`

说明：
- `baseRepo` 为包内实现（小写），业务侧通过 `NewBaseRepo` 获取 `BaseRepo` 即可。
- `C`（Condition）通过结构体标签声明查询行为。

### Condition 标签规范

格式：

```go
query:"type:eq;column:id"
```

支持的 `type`：
- `eq`：等值查询
- `in`：IN 查询
- `contains`：LIKE `%v%`
- `order`：排序（按字段值控制升降序）

`order` 示例：

```go
type BaseRoleCondition struct {
    CodeOrder string `query:"type:order;column:code"` // "desc" 为降序，其它非空值为升序
}
```

### 接入示例

```go
repo := repo.NewBaseRepo[models.BaseAPI, BaseApiCondition](
    func(ctx context.Context) gen.Dao {
        dao := q.BaseAPI.WithContext(ctx).DO
        return &dao
    },
    func(ctx context.Context) field.Int64 { return q.BaseAPI.ID },
    func(entity *models.BaseAPI) int64 { return entity.ID },
    func(ctx context.Context) any { return q.BaseAPI },
    100,
)
```
