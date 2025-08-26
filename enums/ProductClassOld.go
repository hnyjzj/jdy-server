package enums

import "errors"

/* 旧料大类 */
// 黄金旧料、K金旧料、铂金旧料、银旧料、足金镶嵌旧料、镶嵌旧料、其他
type ProductClassOld int

const (
	ProductClassOldGold      ProductClassOld = iota + 1 // 黄金旧料
	ProductClassOldKGold                                // K金旧料
	ProductClassOldPlatinum                             // 铂金旧料
	ProductClassOldSilver                               // 银旧料
	ProductClassOldInlayGold                            // 足金镶嵌旧料
	ProductClassOldInlay                                // 镶嵌旧料
	ProductClassOldOther                                // 其他
)

var ProductClassOldMap = map[ProductClassOld]string{
	ProductClassOldGold:      "黄金旧料",
	ProductClassOldKGold:     "K金旧料",
	ProductClassOldPlatinum:  "铂金旧料",
	ProductClassOldSilver:    "银旧料",
	ProductClassOldInlayGold: "足金镶嵌旧料",
	ProductClassOldInlay:     "镶嵌旧料",
	ProductClassOldOther:     "其他",
}

func (p ProductClassOld) ToMap() any {
	return ProductClassOldMap
}

func (p ProductClassOld) InMap() error {
	if _, ok := ProductClassOldMap[p]; !ok {
		return errors.New("not in enum")
	}
	return nil
}

func (p ProductClassOld) String() string {
	return ProductClassOldMap[p]
}
