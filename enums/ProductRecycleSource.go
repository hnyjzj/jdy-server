package enums

import (
	"errors"
)

/* 回收来源 */
// 报损转旧料、回收
type ProductRecycleSource int

const (
	ProductRecycleSource_BAOSONGZHENGOLD ProductRecycleSource = iota + 1 // 报损转旧料
	ProductRecycleSource_HUISHOU                                         // 回收
)

var ProductRecycleSourceMap = map[ProductRecycleSource]string{
	ProductRecycleSource_BAOSONGZHENGOLD: "报损转旧料",
	ProductRecycleSource_HUISHOU:         "回收",
}

func (p ProductRecycleSource) ToMap() any {
	return ProductRecycleSourceMap
}

func (p ProductRecycleSource) InMap() error {
	if _, ok := ProductRecycleSourceMap[p]; !ok {
		return errors.New("not in enum")
	}
	return nil
}
