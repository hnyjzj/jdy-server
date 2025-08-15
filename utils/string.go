package utils

import (
	"fmt"
	"time"

	"golang.org/x/exp/rand"
)

const (
	digits       = "0123456789"                 // 数字
	letter       = "abcdefghijklmnopqrstuvwxyz" // 小写字母
	letter_upper = "ABCDEFGHIJKLMNOPQRSTUVWXYZ" // 大写字母
)

// 随机字符串（字母数字）
func RandomAlphanumeric(length int) string {
	if length <= 0 {
		return ""
	}

	rand.Seed(uint64(time.Now().UnixNano()))
	// 小写字母+数字
	bytes := digits + letter
	b := make([]byte, length)
	for i := range b {
		b[i] = bytes[rand.Intn(len(bytes))]
	}
	return string(b)
}

// 随机条码【数字+字母（字母最多 2 位且不能字母开头）】
func RandomCode(length int) string {
	if length <= 0 {
		return ""
	}

	rand.Seed(uint64(time.Now().UnixNano()))

	b := make([]byte, length)

	// 确保第一个字符是数字
	b[0] = digits[rand.Intn(len(digits))]

	// 计算允许的最大字母数量（不超过2且不超过剩余位置数）
	maxLetters := 2
	remaining := length - 1
	if remaining < maxLetters {
		maxLetters = remaining
	}
	letterCount := rand.Intn(maxLetters + 1) // 随机生成0~maxLetters的字母数量

	// 随机选择字母出现的位置（在非首字符的位置）
	if letterCount > 0 {
		positions := rand.Perm(remaining)[:letterCount] // 生成不重复的随机位置
		for _, pos := range positions {
			// 注意：positions是相对于b[1:]的索引，实际位置需+1
			b[pos+1] = letter_upper[rand.Intn(len(letter_upper))]
		}
	}

	// 填充剩余位置为数字
	for i := 1; i < length; i++ {
		if b[i] == 0 { // 未设置的位置填充数字
			b[i] = digits[rand.Intn(len(digits))]
		}
	}

	return string(b)
}

// 获取当前毫秒级时间戳
func GetCurrentMilliseconds() string {
	now := time.Now()
	return fmt.Sprintf("%d%02d%02d%02d%02d%02d%03d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), now.Nanosecond()/1000)
}
