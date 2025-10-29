package enums

import "errors"

/* 销售目标统计范围 */
// 大类、分类、配件、全部
type TargetScope int

const (
	TargetScopeClass      TargetScope = iota + 1 // 大类
	TargetScopeCategory                          // 分类
	TargetScopeAccessorie                        // 配件
	TargetScopeAll                               // 全部
)

var TargetScopeMap = map[TargetScope]string{
	TargetScopeClass:      "大类",
	TargetScopeCategory:   "分类",
	TargetScopeAccessorie: "配件",
	TargetScopeAll:        "全部",
}

func (p TargetScope) ToMap() any {
	return TargetScopeMap
}

func (p TargetScope) InMap() error {
	if _, ok := TargetScopeMap[p]; !ok {
		return errors.New("not in enum")
	}
	return nil
}

func (p TargetScope) String() string {
	return TargetScopeMap[p]
}
