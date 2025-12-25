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
type findRes[T any] struct {
	Index int
	Item  T
	Has   bool
	Err   error
}

func ArrayFind[T comparable](array []T, find func(item T) bool) (res findRes[T]) {
	if len(array) == 0 || array == nil {
		return findRes[T]{
			Has:   false,
			Index: -1,
			Err:   fmt.Errorf("array is empty"),
		}
	}
	for i, v := range array {
		if find(v) {
			return findRes[T]{
				Has:   true,
				Item:  v,
				Index: i,
				Err:   nil,
			}
		}
	}

	return findRes[T]{
		Has:   false,
		Index: -1,
		Err:   fmt.Errorf("array is empty"),
	}
}

func ArrayFindIn[T comparable](array []T, item T) bool {
	res := ArrayFind(array, func(i T) bool {
		return i == item
	})

	return res.Has
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

// ArrayUnique 对切片去重（保留原顺序，根据 fun 返回的 key 判断唯一性）
// T：切片元素类型（需可比较）
// K：用于判断唯一性的 key 类型（需可比较）
// fun：从元素中提取用于去重的 key（例如：func(item User) int { return item.ID }）
func ArrayUnique[T any, K comparable](array []T, fun func(item T) K) []T {
	if len(array) == 0 {
		return array
	}

	seen := make(map[K]struct{}, len(array)) // 记录已出现的 key
	result := make([]T, 0, len(array))       // 结果切片（预分配容量）

	for _, item := range array {
		key := fun(item)             // 提取当前元素的 key
		if _, ok := seen[key]; !ok { // 如果 key 未出现过
			seen[key] = struct{}{}        // 标记为已出现
			result = append(result, item) // 加入结果切片（保留原顺序）
		}
	}

	return result
}
