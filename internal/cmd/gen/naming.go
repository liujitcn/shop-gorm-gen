package main

import (
	"strings"
	"unicode"
)

var defaultAcronyms = map[string]string{
	"api":   "API",
	"id":    "ID",
	"ip":    "IP",
	"url":   "URL",
	"uri":   "URI",
	"http":  "HTTP",
	"https": "HTTPS",
	"tcp":   "TCP",
	"udp":   "UDP",
	"rpc":   "RPC",
	"sql":   "SQL",
	"db":    "DB",
	"uid":   "UID",
	"uuid":  "UUID",
	"sku":   "SKU",
	"sn":    "SN",
}

// buildModelName 将表名转换为模型名，保留常见缩写的大写形式。
func buildModelName(tableName string) string {
	parts := strings.Split(tableName, "_")
	for i, part := range parts {
		lowerPart := strings.ToLower(part)
		if acronym, ok := defaultAcronyms[lowerPart]; ok {
			parts[i] = acronym
			continue
		}
		parts[i] = upperFirst(lowerPart)
	}
	return strings.Join(parts, "")
}

// buildRepoName 将表名转换为仓储名称，统一使用普通驼峰命名。
func buildRepoName(tableName string) string {
	parts := strings.Split(tableName, "_")
	for i, part := range parts {
		parts[i] = upperFirst(strings.ToLower(part))
	}
	return strings.Join(parts, "")
}

// upperFirst 将字符串首字母转换为大写。
func upperFirst(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

// lowerFirst 将字符串首字母转换为小写。
func lowerFirst(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}
