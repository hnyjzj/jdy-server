package enums

import "errors"

/* 产品工艺 */
// 全部、无、3D、5D、5G、古法、复古、万足金、九五金、珐琅彩、精品A、精品B、精品C、精品D
type ProductCraft int

const (
	ProductCraftAll      ProductCraft = iota // 全部
	ProductCraftNone                         // 无
	ProductCraft3D                           // 3D
	ProductCraft5D                           // 5D
	ProductCraft5G                           // 5G
	ProductCraftAncient                      // 古法
	ProductCraftRetro                        // 复古
	ProductCraftFiveGold                     // 万足金
	ProductCraftNineGold                     // 九五金
	ProductCraftFengLan                      // 珐琅彩
	ProductCraftFineA                        // 精品A
	ProductCraftFineB                        // 精品B
	ProductCraftFineC                        // 精品C
	ProductCraftFineD                        // 精品D
)

var ProductCraftMap = map[ProductCraft]string{
	ProductCraftAll:      "全部",
	ProductCraftNone:     "无",
	ProductCraft3D:       "3D",
	ProductCraft5D:       "5D",
	ProductCraft5G:       "5G",
	ProductCraftAncient:  "古法",
	ProductCraftRetro:    "复古",
	ProductCraftFiveGold: "万足金",
	ProductCraftNineGold: "九五金",
	ProductCraftFengLan:  "珐琅彩",
	ProductCraftFineA:    "精品A",
	ProductCraftFineB:    "精品B",
	ProductCraftFineC:    "精品C",
	ProductCraftFineD:    "精品D",
}

func (p ProductCraft) ToMap() any {
	return ProductCraftMap
}

func (p ProductCraft) InMap() error {
	if _, ok := ProductCraftMap[p]; !ok {
		return errors.New("not in enum")
	}
	return nil
}
