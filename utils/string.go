package utils

import (
	"fmt"
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
