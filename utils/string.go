package utils

import (
	"fmt"
	"strings"
	"time"

	"golang.org/x/exp/rand"
)

const (
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

// 随机字符串（字母数字）
func RandomAlphanumeric(length int) string {
	rand.Seed(uint64(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// 获取当前毫秒级时间戳
func GetCurrentMilliseconds() string {
	now := time.Now()
	return fmt.Sprintf("%d%02d%02d%02d%02d%02d%03d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), now.Nanosecond()/1000)
}

// 将数组转换为字符串，使用指定分隔符(sep)
func ArrayToString[T any](array []T, sep string) string {
	var strs []string
	for _, item := range array {
		strs = append(strs, fmt.Sprint(item))
	}
	return strings.Join(strs, sep)
}
