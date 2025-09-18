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

// 在数组中查找指定元素，返回元素索引和指针
func ArrayFind[T any](array []T, find func(item T) bool) (item T, index int, err error) {
	if len(array) == 0 || array == nil {
		return item, -1, fmt.Errorf("empty array")
	}
	for i, v := range array {
		if find(v) {
			return v, i, nil
		}
	}
	return item, -1, fmt.Errorf("not found")
}

// 根据下标删除数组元素
func ArrayDeleteOfIndex[T any](array []T, index int) []T {
	if index < 0 || index >= len(array) || array == nil {
		return array
	}
	return append(array[:index], array[index+1:]...)
}

// 合并多个数组
func ArrayMerge[T any](slices ...[]T) []T {
	var result []T
	for _, slice := range slices {
		result = append(result, slice...)
	}
	return result
}
