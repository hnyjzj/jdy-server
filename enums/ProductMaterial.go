package enums

import "errors"

/* 产品材质 */
// 全部、黄金、银饰、铂金、钯金、裸石
type ProductMaterial int

const (
	ProductMaterialAll       ProductMaterial = iota // 全部
	ProductMaterialGold                             // 黄金
	ProductMaterialSilver                           // 银饰
	ProductMaterialPlatinum                         // 铂金
	ProductMaterialPalladium                        // 钯金
	ProductMaterialGem                              // 裸石
)

var ProductMaterialMap = map[ProductMaterial]string{
	ProductMaterialAll:       "全部",
	ProductMaterialGold:      "黄金",
	ProductMaterialSilver:    "银饰",
	ProductMaterialPlatinum:  "铂金",
	ProductMaterialPalladium: "钯金",
	ProductMaterialGem:       "裸石",
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
