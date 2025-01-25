package enums

import (
	"errors"
)

/* 订单来源 */
// 全部、自然客流、回访邀约、营销转化
type OrderSource int

const (
	OrderSourceAll       OrderSource = iota // 全部
	OrderSourceNatural                      // 自然客流
	OrderSourceVisit                        // 回访邀约
	OrderSourceMarketing                    // 营销转化
)

var OrderSourceMap = map[OrderSource]string{
	OrderSourceAll:       "全部",
	OrderSourceNatural:   "自然客流",
	OrderSourceVisit:     "回访邀约",
	OrderSourceMarketing: "营销转化",
}

func (p OrderSource) ToMap() any {
	return OrderSourceMap
}

func (p OrderSource) InMap() error {
	if _, ok := OrderSourceMap[p]; !ok {
		return errors.New("not in enum")
	}
	return nil
}
