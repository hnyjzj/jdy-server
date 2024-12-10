package utils

import (
	"errors"
	"fmt"
	"jdy/types"
	"reflect"
	"strings"
)

// ModelToWhere 将模型根据 tags 转换为查询参数
func ModelToWhere[M any](model M, values map[string]any) map[string]types.WhereForm {
	params := make(map[string]types.WhereForm)
	t := reflect.TypeOf(model)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		json := field.Tag.Get("json")
		where := field.Tag.Get("where")
		if where == "" {
			continue
		}

		whereForm, err := parseTag(where, values)
		if err != nil {
			fmt.Printf("Error parsing tag for field %s: %v\n", json, err)
			continue
		}
		params[json] = whereForm
	}

	return params
}

func parseTag(where string, values map[string]any) (types.WhereForm, error) {
	var whereForm types.WhereForm
	parts := strings.Split(where, ";")
	for _, part := range parts {
		kv := strings.Split(part, ":")
		var (
			key   string
			value string
		)

		switch len(kv) {
		case 1:
			key = kv[0]
			value = ""
		case 2:
			key = kv[0]
			value = kv[1]
		default:
			return whereForm, errors.New("invalid tag format")
		}

		switch key {
		case "label":
			whereForm.Label = value
		case "type":
			whereForm.Type = value
		case "required":
			whereForm.Required = value == "true" || value == ""
		case "preset":
			if strings.HasPrefix(value, "{{.") && strings.HasSuffix(value, "}}") {
				variableName := strings.TrimPrefix(strings.TrimSuffix(value, "}}"), "{{.")
				if val, ok := values[variableName]; ok {
					whereForm.Preset = val
				} else {
					return whereForm, errors.New("variable not found in provided values for preset")
				}
			} else {
				whereForm.Preset = value
			}
		}
	}
	return whereForm, nil
}
