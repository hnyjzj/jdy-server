package utils

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
)

// StructToStruct 使用 mapstructure 将输入结构体转换为指定的输出类型
// 支持深度转换和自定义类型转换
// 示例:
//
//	type Input struct { Name string }
//	type Output struct { Name string }
//	result, err := StructToStruct[Output](input)
func StructToStruct[Output any](input any) (Output, error) {
	var dstStruct Output
	config := &mapstructure.DecoderConfig{
		Result:           &dstStruct,
		WeaklyTypedInput: true, // 允许类型转换
		ErrorUnused:      true, // 忽略未使用的字段
		Squash:           true, // 合并嵌套结构体
	}
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return dstStruct, fmt.Errorf("创建解码器失败: %w", err)
	}
	if err := decoder.Decode(input); err != nil {
		return dstStruct, fmt.Errorf("结构体转换失败: %w", err)
	}
	return dstStruct, err
}
