package main

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"text/template"
)

type tableMeta struct {
	TableName string
	ModelName string
	RepoName  string
}

// generateDataFiles 根据 gorm/gen 导出的全部表结果生成 data 包代码。
func generateDataFiles(tableModels []interface{}) error {
	rootDir, err := projectRoot()
	if err != nil {
		return err
	}
	tables, err := loadTables(tableModels)
	if err != nil {
		return err
	}
	if err := generateDataLayer(rootDir, tables); err != nil {
		return err
	}
	return nil
}

// loadTables 从 gorm/gen 导出的全部表结果中提取表元信息，并按字典序返回。
func loadTables(tableModels []interface{}) ([]tableMeta, error) {
	metas := make([]tableMeta, 0, len(tableModels))
	for _, tableModel := range tableModels {
		meta, ok := extractTableMeta(tableModel)
		if !ok {
			return nil, fmt.Errorf("解析表元信息失败，类型=%T", tableModel)
		}
		metas = append(metas, meta)
	}
	sort.Slice(metas, func(i, j int) bool {
		return metas[i].TableName < metas[j].TableName
	})
	return metas, nil
}

// extractTableMeta 从 gorm/gen 返回对象中提取 data 层所需的表信息。
func extractTableMeta(tableModel any) (tableMeta, bool) {
	if tableModel == nil {
		return tableMeta{}, false
	}
	value := reflect.ValueOf(tableModel)
	if value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return tableMeta{}, false
		}
		value = value.Elem()
	}
	if value.Kind() != reflect.Struct {
		return tableMeta{}, false
	}

	tableName, ok := readStringField(value, "TableName")
	if !ok || tableName == "" {
		return tableMeta{}, false
	}
	modelName, ok := readStringField(value, "ModelStructName")
	if !ok || modelName == "" {
		// 回退到表名推导，兼容返回对象字段变化场景。
		modelName = buildModelName(tableName)
	}
	return tableMeta{
		TableName: tableName,
		ModelName: modelName,
		RepoName:  buildRepoName(tableName),
	}, true
}

// readStringField 读取结构体中的字符串字段。
func readStringField(value reflect.Value, fieldName string) (string, bool) {
	fieldValue := value.FieldByName(fieldName)
	if !fieldValue.IsValid() || fieldValue.Kind() != reflect.String {
		return "", false
	}
	return fieldValue.String(), true
}

// generateDataLayer 生成 data 包中的基础仓储、迁移注册与 ProviderSet。
func generateDataLayer(rootDir string, tables []tableMeta) error {
	dataDir := filepath.Join(rootDir, "data")
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		return err
	}
	if err := writeTemplateFile(filepath.Join(dataDir, "data.go"), dataFileTemplate, tables); err != nil {
		return err
	}
	if err := writeTemplateFile(filepath.Join(dataDir, "init.go"), initFileTemplate, tables); err != nil {
		return err
	}
	for _, table := range tables {
		if err := writeTemplateFile(filepath.Join(dataDir, table.TableName+".go"), repoFileTemplate, table); err != nil {
			return err
		}
	}
	return nil
}

// writeTemplateFile 根据模板渲染 Go 文件，并自动格式化后写入磁盘。
func writeTemplateFile(filename, tpl string, data any) error {
	t, err := template.New(filepath.Base(filename)).Funcs(template.FuncMap{
		"lowerFirst": lowerFirst,
	}).Parse(tpl)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	if err = t.Execute(&buf, data); err != nil {
		return err
	}
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return fmt.Errorf("格式化文件%s失败: %w", filename, err)
	}
	return os.WriteFile(filename, formatted, 0o644)
}

// projectRoot 返回 shop-gorm-gen 仓库根目录。
func projectRoot() (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("获取当前文件路径失败")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(filename), "..", "..", "..")), nil
}
