package utils

import (
	"fmt"
	"reflect"

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
		WeaklyTypedInput: true,  // 支持弱类型转换（如字符串转数字）
		ErrorUnused:      false, // 忽略未使用的字段
		ZeroFields:       false, // 不将目标字段默认置零（保留目标类型的零值）
		Squash:           true,  // 合并嵌套结构体
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

// StructMerge 将 src 的字段值合并到 dst 中，如果 src 中存在 dst 中不存在的字段，则忽略
func StructMerge[T any](dst *T, src T) error {
	dstValue := reflect.ValueOf(dst)
	if dstValue.Kind() != reflect.Ptr || dstValue.IsNil() {
		return fmt.Errorf("dst must be a non - nil pointer")
	}
	dstValue = dstValue.Elem()

	srcValue := reflect.ValueOf(src)
	if srcValue.Kind() != reflect.Struct {
		return fmt.Errorf("src must be a struct")
	}

	if dstValue.Type() != srcValue.Type() {
		return fmt.Errorf("dst and src must be of the same type")
	}

	for i := range dstValue.NumField() {
		srcField := srcValue.Field(i)
		if srcField.CanInterface() && !reflect.DeepEqual(srcField.Interface(), reflect.Zero(srcField.Type()).Interface()) {
			dstValue.Field(i).Set(srcField)
		}
	}

	return nil
}
