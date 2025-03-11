package enums

import (
	"errors"
)

/* 回收来源 */
// 报损转旧料、回收
type ProductRecycleSource int

const (
	ProductRecycleSourceBaoSongZhenGold ProductRecycleSource = iota + 1 // 报损转旧料
	ProductRecycleSourceHuiShou                                         // 回收
)

var ProductRecycleSourceMap = map[ProductRecycleSource]string{
	ProductRecycleSourceBaoSongZhenGold: "报损转旧料",
	ProductRecycleSourceHuiShou:         "回收",
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
