package enums

import "errors"

/* 产品品牌 */
// 全部、金美福、周大生、老庙、周六福、金大福、潮宏基、中国珠宝、老庙推广
type ProductBrand int

const (
	ProductBrandAll  ProductBrand = iota // 全部
	ProductBrandJMF                      // 金美福
	ProductBrandZDS                      // 周大生
	ProductBrandLM                       // 老庙
	ProductBrandZLF                      // 周六福
	ProductBrandJDF                      // 金大福
	ProductBrandCHJ                      // 潮宏基
	ProductBrandZGJB                     // 中国珠宝
	ProductBrandLMTG                     // 老庙推广
)

var ProductBrandMap = map[ProductBrand]string{
	ProductBrandAll:  "全部",
	ProductBrandJMF:  "金美福",
	ProductBrandZDS:  "周大生",
	ProductBrandLM:   "老庙",
	ProductBrandZLF:  "周六福",
	ProductBrandJDF:  "金大福",
	ProductBrandCHJ:  "潮宏基",
	ProductBrandZGJB: "中国珠宝",
	ProductBrandLMTG: "老庙推广",
}

func (p ProductBrand) ToMap() any {
	return ProductBrandMap
}

func (p ProductBrand) InMap() error {
	if _, ok := ProductBrandMap[p]; !ok {
		return errors.New("not in enum")
	}
	return nil
}
