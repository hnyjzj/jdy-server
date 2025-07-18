package enums

import "errors"

/* 产品品牌 */
// 金美福、周大生、老庙、周六福、金大福、中国珠宝、老庙推广价、金大福推广价、金六福推广价、金美福推广价、中国珠宝推广价、潮宏基、谢瑞麟
type ProductBrand int

const (
	ProductBrandJMF   ProductBrand = iota + 1 // 金美福
	ProductBrandZDS                           // 周大生
	ProductBrandLM                            // 老庙
	ProductBrandZLF                           // 周六福
	ProductBrandJDF                           // 金大福
	ProductBrandZJ                            // 中国珠宝
	ProductBrandLMTP                          // 老庙推广价
	ProductBrandJDFTP                         // 金大福推广价
	ProductBrandJLFTP                         // 金六福推广价
	ProductBrandJMFTP                         // 金美福推广价
	ProductBrandZJTP                          // 中国珠宝推广价
	ProductBrandCHJ                           // 潮宏基
	ProductBrandXSL                           // 谢瑞麟
)

var ProductBrandMap = map[ProductBrand]string{
	ProductBrandJMF:   "金美福",
	ProductBrandZDS:   "周大生",
	ProductBrandLM:    "老庙",
	ProductBrandZLF:   "周六福",
	ProductBrandJDF:   "金大福",
	ProductBrandZJ:    "中国珠宝",
	ProductBrandLMTP:  "老庙推广价",
	ProductBrandJDFTP: "金大福推广价",
	ProductBrandJLFTP: "金六福推广价",
	ProductBrandJMFTP: "金美福推广价",
	ProductBrandZJTP:  "中国珠宝推广价",
	ProductBrandCHJ:   "潮宏基",
	ProductBrandXSL:   "谢瑞麟",
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

func (ProductBrand) All() []ProductBrand {
	var all []ProductBrand
	for k := range ProductBrandMap {
		all = append(all, k)
	}

	return all
}
