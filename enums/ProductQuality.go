package enums

import "errors"

/* 产品成色 */
// 全部、999.99、999.9、999、990、950、925、916、900、750、585、375、1、裸石
type ProductQuality int

const (
	ProductQualityAll   ProductQuality = iota // 全部
	ProductQuality99999                       // 999.99
	ProductQuality9999                        // 999.9
	ProductQuality999                         // 999
	ProductQuality990                         // 990
	ProductQuality950                         // 950
	ProductQuality925                         // 925
	ProductQuality916                         // 916
	ProductQuality900                         // 900
	ProductQuality750                         // 750
	ProductQuality585                         // 585
	ProductQuality375                         // 375
	ProductQuality1                           // 1
	ProductQualityGem                         // 裸石
)

var ProductQualityMap = map[ProductQuality]string{
	ProductQualityAll:   "全部",
	ProductQuality99999: "999.99",
	ProductQuality9999:  "999.9",
	ProductQuality999:   "999",
	ProductQuality990:   "990",
	ProductQuality950:   "950",
	ProductQuality925:   "925",
	ProductQuality916:   "916",
	ProductQuality900:   "900",
	ProductQuality750:   "750",
	ProductQuality585:   "585",
	ProductQuality375:   "375",
	ProductQuality1:     "1",
	ProductQualityGem:   "裸石",
}

func (p ProductQuality) ToMap() any {
	return ProductQualityMap
}

func (p ProductQuality) InMap() error {
	if _, ok := ProductQualityMap[p]; !ok {
		return errors.New("not in enum")
	}
	return nil
}
