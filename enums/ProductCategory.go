package enums

import "errors"

/* 产品品类 */
// 全部、戒指、项链、套链、吊坠、耳饰、手链、手镯、脚链、胸针、珠串、挂件、金条、摆件、饰品、其他
type ProductCategory int

const (
	ProductCategoryAll       ProductCategory = iota // 全部
	ProductCategoryRing                             // 戒指
	ProductCategoryNecklace                         // 项链
	ProductCategoryChoker                           // 套链
	ProductCategoryPendant                          // 吊坠
	ProductCategoryEarring                          // 耳饰
	ProductCategoryBracelet                         // 手链
	ProductCategoryBangle                           // 手镯
	ProductCategoryAnklet                           // 脚链
	ProductCategoryBrooch                           // 胸针
	ProductCategoryBead                             // 珠串
	ProductCategoryAccessory                        // 挂件
	ProductCategoryGoldBar                          // 金条
	ProductCategoryOrnament                         // 摆件
	ProductCategoryJewelry                          // 饰品
	ProductCategoryOther                            // 其他
)

var ProductCategoryMap = map[ProductCategory]string{
	ProductCategoryAll:       "全部",
	ProductCategoryRing:      "戒指",
	ProductCategoryNecklace:  "项链",
	ProductCategoryChoker:    "套链",
	ProductCategoryPendant:   "吊坠",
	ProductCategoryEarring:   "耳饰",
	ProductCategoryBracelet:  "手链",
	ProductCategoryBangle:    "手镯",
	ProductCategoryAnklet:    "脚链",
	ProductCategoryBrooch:    "胸针",
	ProductCategoryBead:      "珠串",
	ProductCategoryAccessory: "挂件",
	ProductCategoryGoldBar:   "金条",
	ProductCategoryOrnament:  "摆件",
	ProductCategoryJewelry:   "饰品",
	ProductCategoryOther:     "其他",
}

func (p ProductCategory) ToMap() any {
	return ProductCategoryMap
}

func (p ProductCategory) InMap() error {
	if _, ok := ProductCategoryMap[p]; !ok {
		return errors.New("not in enum")
	}
	return nil
}
