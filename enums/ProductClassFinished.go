package enums

import "errors"

/* 产品大类 */
// 足金（克）、足金（件）、金 750、金 916、铂金、银饰、足金镶嵌、钻石、裸钻、彩宝、玉石、珍珠、其他
type ProductClassFinished int

const (
	ProductClassFinishedGoldKg       ProductClassFinished = iota + 1 // 足金（克）
	ProductClassFinishedGoldPiece                                    // 足金（件）
	ProductClassFinishedGold750                                      // 金 750
	ProductClassFinishedGold916                                      // 金 916
	ProductClassFinishedPlatinum                                     // 铂金
	ProductClassFinishedSilver                                       // 银饰
	ProductClassFinishedGoldInlay                                    // 足金镶嵌
	ProductClassFinishedDiamond                                      // 钻石
	ProductClassFinishedDiamondNaked                                 // 裸钻
	ProductClassFinishedCoral                                        // 彩宝
	ProductClassFinishedJade                                         // 玉石
	ProductClassFinishedPearl                                        // 珍珠
	ProductClassFinishedOther                                        // 其他
)

var ProductClassFinishedMap = map[ProductClassFinished]string{
	ProductClassFinishedGoldKg:       "足金（克）",
	ProductClassFinishedGoldPiece:    "足金（件）",
	ProductClassFinishedGold750:      "金 750",
	ProductClassFinishedGold916:      "金 916",
	ProductClassFinishedPlatinum:     "铂金",
	ProductClassFinishedSilver:       "银饰",
	ProductClassFinishedGoldInlay:    "足金镶嵌",
	ProductClassFinishedDiamond:      "钻石",
	ProductClassFinishedDiamondNaked: "裸钻",
	ProductClassFinishedCoral:        "彩宝",
	ProductClassFinishedJade:         "玉石",
	ProductClassFinishedPearl:        "珍珠",
	ProductClassFinishedOther:        "其他",
}

func (p ProductClassFinished) ToMap() any {
	return ProductClassFinishedMap
}

func (p ProductClassFinished) InMap() error {
	if _, ok := ProductClassFinishedMap[p]; !ok {
		return errors.New("not in enum")
	}
	return nil
}
