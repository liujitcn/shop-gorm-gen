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

当前生成入口在 `internal/cmd/gen/main.go`，默认使用：

- `WithDriver("mysql")`
- `WithSource(defaultDSN)`

## 打 Tag

```bash
make tag                    # 默认从仓库根目录递归检查 go.mod 并自动打/推送 tag（含根模块）
make tag MODULE=repo        # 从 repo 目录开始递归检查 go.mod 并打 tag
```

说明：上述命令通过 `python3 scripts/tag_release.py` 执行统一的版本计算与远程更新检测逻辑。

## 当前配置说明

- 生成入口：`internal/cmd/gen/main.go`
- 输出目录：`query`
- 模型包名：`models`
- 默认会同时生成 `models` 与 `query`
- 模型命名策略：按表名转 CamelCase，不做单数化（例如 `goods -> Goods`、`order_goods -> OrderGoods`）
- 生成器依赖：`github.com/liujitcn/gorm-kit/gen`

## Codex 文档规则

- 每次修改代码后，必须同步检查并优化 `AGENTS.md` 与 `README.md`。
- 若命令、流程、生成逻辑或行为发生变化，必须在同次变更中更新文档。
- 提交推送前，必须先完成 `README.md` 更新。
- 提交代码时，必须将 `README.md` 改动与本次代码改动一起提交。
- 后续新增或修改代码时，代码注释统一使用中文。

## 生成结果

- `models/*.gen.go`：数据模型代码
- `query/*.gen.go`：查询构造与 DAO 代码
- `query/gen.go`：查询入口

注意：默认 DSN 仅用于本地开发，建议使用 `GORM_GEN_DSN` 指向你自己的数据库后再执行生成。

> 说明：仓库当前仅承担代码生成职责，不包含 `repo/base_repo.go` 相关实现。
