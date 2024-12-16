package enums

import "errors"

/* 产品颜色 */
// 全部、D、E、D-E、F、G、F-G、H、I、J、I-J、K、L、K-L、M、N、M-N
type ProductColor int

const (
	ProductColorAll ProductColor = iota // 全部
	ProductColorD                       // D
	ProductColorE                       // E
	ProductColorDE                      // D-E
	ProductColorF                       // F
	ProductColorG                       // G
	ProductColorFG                      // F-G
	ProductColorH                       // H
	ProductColorI                       // I
	ProductColorJ                       // J
	ProductColorIJ                      // I-J
	ProductColorK                       // K
	ProductColorL                       // L
	ProductColorKL                      // K-L
	ProductColorM                       // M
	ProductColorN                       // N
	ProductColorMN                      // M-N
)

var ProductColorMap = map[ProductColor]string{
	ProductColorAll: "全部",
	ProductColorD:   "D",
	ProductColorE:   "E",
	ProductColorDE:  "D-E",
	ProductColorF:   "F",
	ProductColorG:   "G",
	ProductColorFG:  "F-G",
	ProductColorH:   "H",
	ProductColorI:   "I",
	ProductColorJ:   "J",
	ProductColorIJ:  "I-J",
	ProductColorK:   "K",
	ProductColorL:   "L",
	ProductColorKL:  "K-L",
	ProductColorM:   "M",
	ProductColorN:   "N",
	ProductColorMN:  "M-N",
}

func (p ProductColor) ToMap() any {
	return ProductColorMap
}

func (p ProductColor) InMap() error {
	if _, ok := ProductColorMap[p]; !ok {
		return errors.New("not in enum")
	}
	return nil
}
