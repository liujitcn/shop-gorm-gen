# Codex 规则

## 适用范围
- 本规则适用于 `shop-gorm-gen` 仓库全量目录。
- 变更应尽量小且聚焦于代码生成、依赖兼容和发布可用性。

## 语言规范
- 与用户沟通、总结说明统一使用中文。
- 后续新增或修改代码时，代码注释统一使用中文。

## 项目定位
- 本仓库用于根据数据库表结构生成 GORM 的 `models` 与 `query` 代码。
- 生成入口：`internal/cmd/gen/main.go`。
- 生成产物：`models/*.gen.go`、`query/*.gen.go`、`query/gen.go`。

## 标准命令
- 初始化依赖：`make init`
- 执行生成：`make gen`
- 校验编译：`go test ./...`
- 查看目标：`make help`

## 生成流程
- 涉及表结构或生成逻辑的任务，按以下顺序执行：
  1. `make init`
  2. `make gen`
  3. `go test ./...`
- 第 3 步失败时，任务不算完成。

## 发布与打标门禁
- 本模块会打 tag 供其他项目引用。
- 在 commit/tag/push 前，必须保证生成成功且 `go test ./...` 通过。
- 若出现依赖冲突（例如 GORM 插件版本不兼容），先修复 `go.mod/go.sum`，再重新执行生成与测试。

## 提交流程约定
- 用户要求“提交”时，默认执行完整发布动作：`git commit` + `git push`。
- 未明确指定分支时，推送当前分支到同名远程分支。
- `git commit -m` 信息默认使用中文，简洁描述本次变更。
- 若用户未指定提交信息，按变更内容自动生成中文提交信息。

## Tag 规则
- Tag 使用语义化前缀格式：`v主版本.次版本.补丁版本`（例如 `v0.0.1`）。
- 当用户要求“打 tag”且未指定具体版本号时：
  1. 读取当前仓库最新 tag（按版本排序）。
  2. 在最新 tag 基础上仅递增补丁版本。
  3. 例如：`v0.0.1` 的下一次自动 tag 为 `v0.0.2`。
- 若仓库无历史 tag，则从 `v0.0.1` 开始。

## 编辑约束
- 非明确要求下，不手工修改生成文件 `*.gen.go`。
- 优先通过修改 `internal/cmd/gen/main.go`、依赖版本或 Makefile 目标来解决问题。
- Go 版本基线保持为：`go 1.26.0`。

## 文档同步
- 若命令名、流程、生成逻辑或代码行为有变更，需在同次变更中同步更新 `AGENTS.md` 与 `README.md`。
- Codex 每次修改代码后，都必须检查并优化 `AGENTS.md` 与 `README.md`，确保规则与文档始终最新。
- README 中命令示例必须与 Makefile 实际目标保持一致。
- 提交推送前，必须先完成 `README.md` 更新。
- 提交代码时，必须将 `README.md` 改动与本次代码改动一起提交。

## 仓储抽象约定
- `repo/base_repo.go` 是通用仓储实现入口，优先复用 `BaseRepo[T, C]` 与 `NewBaseRepo[T, C]`。
- 条件标签统一使用：`query:"type:<op>;column:<db_column>"`。
- 当前允许的 `type`：`eq`、`in`、`contains`、`order`。
- 新增或调整标签语义时，必须同步更新 `README.md` 中的标签说明与示例。
