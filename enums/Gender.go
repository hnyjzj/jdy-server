package enums

import (
	"errors"
	"fmt"
)

/* 性别 */
// 未知、男、女
type Gender int

const (
	GenderUnknown Gender = iota // 未知
	GenderMale                  // 男
	GenderFemale                // 女
)

var GenderMap = map[Gender]string{
	GenderUnknown: "未知",
	GenderMale:    "男",
	GenderFemale:  "女",
}

func (p Gender) ToMap() any {
	return GenderMap
}

func (p Gender) InMap() error {
	if _, ok := GenderMap[p]; !ok {
		return errors.New("not in enum")
	}
	return nil
}

func (p Gender) String() string {
	return GenderMap[p]
}

func (Gender) Convert(v any) Gender {
	switch fmt.Sprintf("%v", v) {
	case "男":
		return GenderMale
	case "女":
		return GenderFemale
	case "1":
		return GenderMale
	case "2":
		return GenderFemale
	default:
		return GenderUnknown
	}
}
