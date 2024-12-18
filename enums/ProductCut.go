package enums

import "errors"

/* 产品切工 */
// 全部、EX、VG、G、P
type ProductCut int

const (
	ProductCutAll ProductCut = iota // 全部
	ProductCutEX                    // EX
	ProductCutVG                    // VG
	ProductCutG                     // G
	ProductCutP                     // P
)

var ProductCutMap = map[ProductCut]string{
	ProductCutAll: "全部",
	ProductCutEX:  "EX",
	ProductCutVG:  "VG",
	ProductCutG:   "G",
	ProductCutP:   "P",
}

func (p ProductCut) ToMap() any {
	return ProductCutMap
}

func (p ProductCut) InMap() error {
	if _, ok := ProductCutMap[p]; !ok {
		return errors.New("not in enum")
	}
	return nil
}
