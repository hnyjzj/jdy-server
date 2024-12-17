package utils

import (
	"errors"
	"fmt"
	"jdy/types"
	"reflect"
	"strconv"
)

// 将结构体根据 tags 转换为查询参数
func StructToWhere[S any](s S) map[string]types.WhereForm {
	params := make(map[string]types.WhereForm)
	t := reflect.TypeOf(s)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		class := field.Type
		tag := field.Tag

		var json string
		if tag.Get("json") != "" {
			json = tag.Get("json")
		}

		whereForm, err := parseTag(class, tag)
		if err != nil {
			fmt.Printf("Error parsing tag for field %s: %v\n", json, err)
			continue
		}
		params[json] = whereForm
	}

	return params
}

func parseTag(class reflect.Type, tga reflect.StructTag) (types.WhereForm, error) {
	var whereForm types.WhereForm

	if tga.Get("json") != "" {
		whereForm.Name = tga.Get("json")
	}
	if tga.Get("label") != "" {
		whereForm.Label = tga.Get("label")
	}
	if tga.Get("sort") != "" {
		v, err := strconv.ParseInt(tga.Get("sort"), 10, 64)
		if err != nil {
			return whereForm, err
		}
		whereForm.Sort = int(v)
	}
	if tga.Get("type") != "" {
		whereForm.Type = tga.Get("type")
	}
	if tga.Get("input") != "" {
		whereForm.Input = tga.Get("input")
	}
	if tga.Get("required") != "" {
		whereForm.Required = tga.Get("required") == "true" || tga.Get("required") == ""
	}
	if tga.Get("show") != "" {
		whereForm.Show = tga.Get("show") == "true" || tga.Get("show") == ""
	}
	if tga.Get("preset") != "" {
		switch tga.Get("preset") {
		case "typeMap":
			enum, ok := reflect.New(class).Interface().(types.EnumMapper)
			if !ok {
				return whereForm, errors.New("class does not implement types.WhereValidate")
			}
			whereForm.Preset = enum.ToMap()
		default:
			whereForm.Preset = tga.Get("preset")
		}
	}

	return whereForm, nil
}
