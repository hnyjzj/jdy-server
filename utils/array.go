package utils

import (
	"fmt"
	"strings"
)

// 将数组转换为字符串，使用指定分隔符(sep)
func ArrayToString[T comparable](array []T, sep string) string {
	var strs []string
	for _, item := range array {
		strs = append(strs, fmt.Sprint(item))
	}
	return strings.Join(strs, sep)
}
