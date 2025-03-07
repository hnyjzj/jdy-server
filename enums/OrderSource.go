package enums

import (
	"errors"
)

/* 订单来源 */
// 自然客流、回访邀约、营销转化
type OrderSource int

const (
	OrderSourceNatural   OrderSource = iota + 1 // 自然客流
	OrderSourceVisit                            // 回访邀约
	OrderSourceMarketing                        // 营销转化
)

var OrderSourceMap = map[OrderSource]string{
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
