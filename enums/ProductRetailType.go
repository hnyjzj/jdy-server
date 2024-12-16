package enums

import "errors"

/* 零售方式 */
// 全部、计件、计重工费按克、计重工费按件
type ProductRetailType int

const (
	ProductRetailTypeAll       ProductRetailType = iota // 全部
	ProductRetailTypePiece                              // 计件
	ProductRetailTypeGoldKg                             // 计重工费按克
	ProductRetailTypeGoldPiece                          // 计重工费按件
)

var ProductRetailTypeMap = map[ProductRetailType]string{
	ProductRetailTypeAll:       "全部",
	ProductRetailTypePiece:     "计件",
	ProductRetailTypeGoldKg:    "计重工费按克",
	ProductRetailTypeGoldPiece: "计重工费按件",
}

func (p ProductRetailType) ToMap() any {
	return ProductRetailTypeMap
}

func (p ProductRetailType) InMap() error {
	if _, ok := ProductRetailTypeMap[p]; !ok {
		return errors.New("not in enum")
	}
	return nil
}
