package xstring

import (
	"regexp"
	"strings"
)

var snakeReg = regexp.MustCompile("[A-Z][a-z]")
var ColumnReg = regexp.MustCompile(`^([a-zA-Z_][a-zA-Z0-9_]*\.[a-zA-Z_][a-zA-Z0-9_]*)$|^([a-zA-Z_][a-zA-Z0-9_]*)$`) //字母数字下划线
var DirectReg = regexp.MustCompile("^asc|desc|ASC|DESC$")

const underline = "_"

// SnakeCase 驼峰转下划线
func SnakeCase(src string) string {
	str := snakeReg.ReplaceAllStringFunc(src, func(s string) string {
		return underline + s
	})
	return strings.ToLower(strings.TrimLeft(str, underline))
}
