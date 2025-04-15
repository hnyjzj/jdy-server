package enums

import "errors"

/* 产品材质 */
// 黄金、银饰、铂金、钯金、裸石类、合金、其他金属、木质、线材、工艺品、编绳费、加工费、维修费、手表
type ProductAccessorieMaterial int

const (
	ProductAccessorieMaterialGold          ProductAccessorieMaterial = iota + 1 // 黄金
	ProductAccessorieMaterialSilver                                             // 银饰
	ProductAccessorieMaterialPlatinum                                           // 铂金
	ProductAccessorieMaterialPalladium                                          // 钯金
	ProductAccessorieMaterialGem                                                // 裸石类
	ProductAccessorieMaterialAlloy                                              // 合金
	ProductAccessorieMaterialOtherMetal                                         // 其他金属
	ProductAccessorieMaterialWood                                               // 木质
	ProductAccessorieMaterialWire                                               // 线材
	ProductAccessorieMaterialCraftsman                                          // 工艺品
	ProductAccessorieMaterialRopeFee                                            // 编绳费
	ProductAccessorieMaterialProcessingFee                                      // 加工费
	ProductAccessorieMaterialRepairFee                                          // 维修费
	ProductAccessorieMaterialWatch                                              // 手表
)

var ProductAccessorieMaterialMap = map[ProductAccessorieMaterial]string{
	ProductAccessorieMaterialGold:          "黄金",
	ProductAccessorieMaterialSilver:        "银饰",
	ProductAccessorieMaterialPlatinum:      "铂金",
	ProductAccessorieMaterialPalladium:     "钯金",
	ProductAccessorieMaterialGem:           "裸石类",
	ProductAccessorieMaterialAlloy:         "合金",
	ProductAccessorieMaterialOtherMetal:    "其他金属",
	ProductAccessorieMaterialWood:          "木质",
	ProductAccessorieMaterialWire:          "线材",
	ProductAccessorieMaterialCraftsman:     "工艺品",
	ProductAccessorieMaterialRopeFee:       "编绳费",
	ProductAccessorieMaterialProcessingFee: "加工费",
	ProductAccessorieMaterialRepairFee:     "维修费",
	ProductAccessorieMaterialWatch:         "手表",
}

func (p ProductAccessorieMaterial) ToMap() any {
	return ProductAccessorieMaterialMap
}

func (p ProductAccessorieMaterial) InMap() error {
	if _, ok := ProductAccessorieMaterialMap[p]; !ok {
		return errors.New("not in enum")
	}
	return nil
}
