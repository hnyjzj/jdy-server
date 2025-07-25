package enums

import "errors"

/* 产品材质 */
// 黄金、银饰、铂金、钯金、裸石类、钻石、金750
type ProductMaterial int

const (
	ProductMaterialGold      ProductMaterial = iota + 1 // 黄金
	ProductMaterialSilver                               // 银饰
	ProductMaterialPlatinum                             // 铂金
	ProductMaterialPalladium                            // 钯金
	ProductMaterialGem                                  // 裸石类
	ProductMaterialDiamond                              // 钻石
	ProductMaterialGold750                              // 金750
)

var ProductMaterialMap = map[ProductMaterial]string{
	ProductMaterialGold:      "黄金",
	ProductMaterialSilver:    "银饰",
	ProductMaterialPlatinum:  "铂金",
	ProductMaterialPalladium: "钯金",
	ProductMaterialGem:       "裸石类",
	ProductMaterialDiamond:   "钻石",
	ProductMaterialGold750:   "金750",
}

func (p ProductMaterial) ToMap() any {
	return ProductMaterialMap
}

func (p ProductMaterial) InMap() error {
	if _, ok := ProductMaterialMap[p]; !ok {
		return errors.New("not in enum")
	}
	return nil
}
