package enums

import "errors"

/* 主石 */
// 全部、素金、钻石、蓝宝石、碧玺、珍珠、翡翠、和田玉、彩宝、其他、红宝石、水晶、祖母绿、玉髓、玛瑙、石榴石、锆石、孔雀石、贝母
type ProductGem int

const (
	ProductGemAll         ProductGem = iota // 全部
	ProductGemGold                          // 素金
	ProductGemDiamond                       // 钻石
	ProductGemSapphire                      // 蓝宝石
	ProductGemTurquoise                     // 碧玺
	ProductGemPearl                         // 珍珠
	ProductGemJade                          // 翡翠
	ProductGemJadeite                       // 和田玉
	ProductGemCoral                         // 彩宝
	ProductGemOther                         // 其他
	ProductGemRuby                          // 红宝石
	ProductGemCrystal                       // 水晶
	ProductGemEmerald                       // 祖母绿
	ProductGemOpal                          // 玉髓
	ProductGemJasper                        // 玛瑙
	ProductGemGarnet                        // 石榴石
	ProductGemZircon                        // 锆石
	ProductGemMalachite                     // 孔雀石
	ProductGemPearlMother                   // 贝母
)

var ProductGemMap = map[ProductGem]string{
	ProductGemAll:         "全部",
	ProductGemGold:        "素金",
	ProductGemDiamond:     "钻石",
	ProductGemSapphire:    "蓝宝石",
	ProductGemTurquoise:   "碧玺",
	ProductGemPearl:       "珍珠",
	ProductGemJade:        "翡翠",
	ProductGemJadeite:     "和田玉",
	ProductGemCoral:       "彩宝",
	ProductGemOther:       "其他",
	ProductGemRuby:        "红宝石",
	ProductGemCrystal:     "水晶",
	ProductGemEmerald:     "祖母绿",
	ProductGemOpal:        "玉髓",
	ProductGemJasper:      "玛瑙",
	ProductGemGarnet:      "石榴石",
	ProductGemZircon:      "锆石",
	ProductGemMalachite:   "孔雀石",
	ProductGemPearlMother: "贝母",
}

func (p ProductGem) ToMap() any {
	return ProductGemMap
}

func (p ProductGem) InMap() error {
	if _, ok := ProductGemMap[p]; !ok {
		return errors.New("not in enum")
	}
	return nil
}
