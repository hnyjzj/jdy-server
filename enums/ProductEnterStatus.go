package enums

import "errors"

/* 入库单状态 */
// 草稿、已完成、已撤销
type ProductEnterStatus int

const (
	ProductEnterStatusDraft     ProductEnterStatus = iota + 1 // 草稿
	ProductEnterStatusCompleted                               // 已完成
	ProductEnterStatusCanceled                                // 已撤销
)

var ProductEnterStatusMap = map[ProductEnterStatus]string{
	ProductEnterStatusDraft:     "草稿",
	ProductEnterStatusCompleted: "已完成",
	ProductEnterStatusCanceled:  "已撤销",
}

func (p ProductEnterStatus) ToMap() any {
	return ProductEnterStatusMap
}

func (p ProductEnterStatus) InMap() error {
	if _, ok := ProductEnterStatusMap[p]; !ok {
		return errors.New("not in enum")
	}
	return nil
}
