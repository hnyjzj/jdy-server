package utils

import (
	"encoding/json"
	"errors"
	"jdy/types"
	"log"
	"reflect"
	"sort"
	"strconv"
)

// 将结构体根据 tags 转换为查询参数
func StructToWhere[S any](s S) map[string]types.WhereForm {
	params := make(map[string]types.WhereForm)
	t := reflect.TypeOf(s)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		// 如果字段是匿名字段，递归处理
		if field.Anonymous {
			nestedParams := StructToWhere(reflect.ValueOf(s).Field(i).Interface())
			for k, v := range nestedParams {
				params[k] = v
			}
			continue
		}
		class := field.Type
		tag := field.Tag

		var name string
		if tag.Get("json") != "" {
			name = tag.Get("json")
		}

		whereForm, err := parseTag(class, tag)
		if err != nil {
			log.Printf("Error parsing tag for field %s: %v\n", name, err)
			continue
		}
		params[name] = whereForm
	}

	return params
}

// 将结构体转换为查询参数数组
func StructWhereToArray[S any](s S) []types.WhereForm {
	// 将结构体转换为 map
	where := StructToWhere(s)
	// 将 map 转换为 []WhereForm
	var dataSlice []types.WhereForm
	for _, value := range where {
		dataSlice = append(dataSlice, value)
	}
	// 对 []WhereForm 进行排序
	sort.Slice(dataSlice, func(i, j int) bool {
		return dataSlice[i].Sort < dataSlice[j].Sort
	})

	return dataSlice
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
	if tga.Get("find") != "" {
		whereForm.Find = tga.Get("find") == "true" || tga.Get("find") == ""
	}
	if tga.Get("create") != "" {
		whereForm.Create = tga.Get("create") == "true" || tga.Get("create") == ""
	}
	if tga.Get("update") != "" {
		whereForm.Update = tga.Get("update") == "true" || tga.Get("update") == ""
	}
	if tga.Get("info") != "" {
		whereForm.Info = tga.Get("info") == "true" || tga.Get("info") == ""
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
	if tga.Get("condition") != "" {
		// 字符串转json
		var tempConditions []types.WhereCondition
		err := json.Unmarshal([]byte(tga.Get("condition")), &tempConditions)
		if err != nil {
			log.Printf("err.Error(): %v\n", err.Error())
			return whereForm, err
		}
		whereForm.Condition = tempConditions
	}

	return whereForm, nil
}
