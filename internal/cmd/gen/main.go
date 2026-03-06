package main

import (
	"fmt"
	"os"
	"strings"
	"unicode"

	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

const defaultDSN = "root:112233@tcp(127.0.0.1:3306)/shop?charset=utf8&parseTime=True&loc=Local&timeout=1000ms"

var initialisms = map[string]string{
	"api": "API",
}

func main() {
	dsn := os.Getenv("GORM_GEN_DSN")
	if dsn == "" {
		dsn = defaultDSN
	}

	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic(fmt.Errorf("open db failed: %w", err))
	}

	g := gen.NewGenerator(gen.Config{
		OutPath:           "query",
		ModelPkgPath:      "models",
		FieldNullable:     false,
		FieldCoverable:    false,
		FieldSignable:     false,
		FieldWithIndexTag: true,
		FieldWithTypeTag:  true,
		WithUnitTest:      false,
	})
	g.UseDB(db)
	g.WithModelNameStrategy(tableToModelName)

	g.ApplyBasic(g.GenerateAllTable()...)
	g.Execute()
}

func tableToModelName(tableName string) string {
	parts := strings.Split(tableName, "_")
	for i, part := range parts {
		lowerPart := strings.ToLower(part)
		if acronym, ok := initialisms[lowerPart]; ok {
			parts[i] = acronym
			continue
		}
		parts[i] = upperFirst(lowerPart)
	}
	return strings.Join(parts, "")
}

func upperFirst(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}
