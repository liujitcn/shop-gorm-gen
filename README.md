# shop-gorm-gen

用于根据数据库表结构自动生成 GORM 的 `models`、`query` 和 `data` 代码。

## 生成方式

1. 安装生成工具：

```bash
make cli
```

说明：`make cli` 会执行 `go install github.com/liujitcn/gorm-kit/gen@latest`，将 `gen` 安装到你的 `GOBIN` 或 `$(go env GOPATH)/bin`。

2. 在项目根目录准备 `config.yaml`：

```yaml
driver: mysql
source: root:112233@tcp(127.0.0.1:3306)/shop?charset=utf8&parseTime=True&loc=Local&timeout=1000ms
out_path: query
model_pkg_path: models
data_path: data
acronyms:
  api: API
  sku: SKU
```

3. 在项目目录执行生成：

```bash
make gen
```

`make gen` 的实际执行命令为：

```bash
gen -config config.yaml
```

也可以直接手工执行：

```bash
gen -config ./config.yaml
```

如需覆盖配置文件中的单项字段，可追加 `-set key=value`，例如：

```bash
gen -config ./config.yaml -set out_path=query_tmp -set data_path=data_tmp
```

## 打 Tag

```bash
make tag                    # 默认从仓库根目录递归检查 go.mod 并自动打/推送 tag（含根模块）
make tag MODULE=repo        # 从 repo 目录开始递归检查 go.mod 并打 tag
```

说明：上述命令通过 `python3 scripts/tag_release.py` 执行统一的版本计算与远程更新检测逻辑。

## 配置说明

- 默认配置文件：`config.yaml`
- 输出目录：`out_path`
- 模型包目录：`model_pkg_path`
- `data` 输出目录：`data_path`
- 默认会同时生成 `models`、`query` 与 `data`
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
- `data/*.go`：基础仓储代码

注意：请优先在 `config.yaml` 中填写你自己的数据库连接信息后再执行生成。

> 说明：仓库当前仅承担代码生成职责，不包含 `repo/base_repo.go` 相关实现。
