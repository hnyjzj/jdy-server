package enums

import "errors"

/* 产品大类 */
// 足金（克）、足金（件）、金 750、金 916、铂金、银饰、足金镶嵌、裸钻、钻石、彩宝、玉石、珍珠、其他
type ProductClass int

const (
	ProductClassGoldKg    ProductClass = iota + 1 // 足金（克）
	ProductClassGoldPiece                         // 足金（件）
	ProductClassGold750                           // 金 750
	ProductClassGold916                           // 金 916
	ProductClassPlatinum                          // 铂金
	ProductClassSilver                            // 银饰
	ProductClassGoldInlay                         // 足金镶嵌
	ProductClassGem                               // 裸石
	ProductClassDiamond                           // 钻石
	ProductClassJade                              // 玉石
	ProductClassPearl                             // 珍珠
	ProductClassOther                             // 其他

)

var ProductClassMap = map[ProductClass]string{
	ProductClassGoldKg:    "足金（克）",
	ProductClassGoldPiece: "足金（件）",
	ProductClassGold750:   "金 750",
	ProductClassGold916:   "金 916",
	ProductClassPlatinum:  "铂金",
	ProductClassSilver:    "银饰",
	ProductClassGoldInlay: "足金镶嵌",
	ProductClassGem:       "裸石",
	ProductClassDiamond:   "钻石",
	ProductClassJade:      "玉石",
	ProductClassPearl:     "珍珠",
	ProductClassOther:     "其他",
}

func (p ProductClass) ToMap() any {
	return ProductClassMap
}

func (p ProductClass) InMap() error {
	if _, ok := ProductClassMap[p]; !ok {
		return errors.New("not in enum")
	}
	return nil
}
