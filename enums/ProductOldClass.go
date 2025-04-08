package enums

import "errors"

/* 旧料大类 */
// 黄金旧料、K金旧料、铂金旧料、银旧料、足金镶嵌旧料、镶嵌旧料、其他
type ProductOldClass int

const (
	ProductOldClassGold      ProductOldClass = iota + 1 // 黄金旧料
	ProductOldClassKGold                                // K金旧料
	ProductOldClassPlatinum                             // 铂金旧料
	ProductOldClassSilver                               // 银旧料
	ProductOldClassInlayGold                            // 足金镶嵌旧料
	ProductOldClassInlay                                // 镶嵌旧料
	ProductOldClassOther                                // 其他
)

var ProductOldClassMap = map[ProductOldClass]string{
	ProductOldClassGold:      "黄金旧料",
	ProductOldClassKGold:     "K金旧料",
	ProductOldClassPlatinum:  "铂金旧料",
	ProductOldClassSilver:    "银旧料",
	ProductOldClassInlayGold: "足金镶嵌旧料",
	ProductOldClassInlay:     "镶嵌旧料",
	ProductOldClassOther:     "其他",
}

func (p ProductOldClass) ToMap() any {
	return ProductOldClassMap
}

func (p ProductOldClass) InMap() error {
	if _, ok := ProductOldClassMap[p]; !ok {
		return errors.New("not in enum")
	}
	return nil
}
