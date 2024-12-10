package utils

import "github.com/mitchellh/mapstructure"

// 将一个结构体转换为另一个结构体
func StructToStruct[Output any](input any) (Output, error) {
	var dstStruct Output
	err := mapstructure.Decode(input, &dstStruct)
	return dstStruct, err
}
