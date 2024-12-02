package utils

import (
	"reflect"
	"strings"
)

func TagWhere(s interface{}) (map[string]map[string]interface{}, error) {
	v := reflect.ValueOf(s)
	t := v.Type()

	result := make(map[string]map[string]interface{})

	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag

		tname := tag.Get("json")
		category := tag.Get("type")
		enum := tag.Get("enum")
		desc := tag.Get("desc")

		if tname != "" {
			fieldMap := make(map[string]interface{})
			fieldMap["type"] = category

			if enum != "" {
				fieldMap["enum"] = strings.Split(enum, ",")
			}
			if desc != "" {
				fieldMap["desc"] = desc
			}

			result[tname] = fieldMap
		}
	}

	return result, nil
}
