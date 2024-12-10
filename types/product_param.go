package types

/* 产品成色 */
// 全部、999.99、999.9、999、990、950、925、916、900、750、585、375、1、裸石
type ProductQuality int

const (
	ProductQuality99999 ProductQuality = iota + 1 // 999.99
	ProductQuality9999                            // 999.9
	ProductQuality999                             // 999
	ProductQuality990                             // 990
	ProductQuality950                             // 950
	ProductQuality925                             // 925
	ProductQuality916                             // 916
	ProductQuality900                             // 900
	ProductQuality750                             // 750
	ProductQuality585                             // 585
	ProductQuality375                             // 375
	ProductQuality1                               // 1
	ProductQualityGem                             // 裸石
)

var ProductQualityMap = map[ProductQuality]string{
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

/* 产品品牌 */
// 全部、金美福、周大生、老庙、周六福、金大福、潮宏基、中国珠宝、老庙推广
type ProductBrand int

const (
	ProductBrandJMF  ProductBrand = iota + 1 // 金美福
	ProductBrandZDS                          // 周大生
	ProductBrandLM                           // 老庙
	ProductBrandZLF                          // 周六福
	ProductBrandJDF                          // 金大福
	ProductBrandCHJ                          // 潮宏基
	ProductBrandZGJB                         // 中国珠宝
	ProductBrandLMTG                         // 老庙推广
)

var ProductBrandMap = map[ProductBrand]string{
	ProductBrandJMF:  "金美福",
	ProductBrandZDS:  "周大生",
	ProductBrandLM:   "老庙",
	ProductBrandZLF:  "周六福",
	ProductBrandJDF:  "金大福",
	ProductBrandCHJ:  "潮宏基",
	ProductBrandZGJB: "中国珠宝",
	ProductBrandLMTG: "老庙推广",
}

/* 主石 */
// 全部、素金、钻石、蓝宝石、碧玺、珍珠、翡翠、和田玉、彩宝、其他、红宝石、水晶、祖母绿、玉髓、玛瑙、石榴石、锆石、孔雀石、贝母
type ProductGem int

const (
	ProductGemGold        ProductGem = iota + 1 // 素金
	ProductGemDiamond                           // 钻石
	ProductGemSapphire                          // 蓝宝石
	ProductGemTurquoise                         // 碧玺
	ProductGemPearl                             // 珍珠
	ProductGemJade                              // 翡翠
	ProductGemJadeite                           // 和田玉
	ProductGemCoral                             // 彩宝
	ProductGemOther                             // 其他
	ProductGemRuby                              // 红宝石
	ProductGemCrystal                           // 水晶
	ProductGemEmerald                           // 祖母绿
	ProductGemOpal                              // 玉髓
	ProductGemJasper                            // 玛瑙
	ProductGemGarnet                            // 石榴石
	ProductGemZircon                            // 锆石
	ProductGemMalachite                         // 孔雀石
	ProductGemPearlMother                       // 贝母
)

var ProductGemMap = map[ProductGem]string{
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

/* 供应商 */
// 全部、金美福、周大生、老庙、潮宏基
type ProductSupplier int

const (
	ProductSupplierJMF ProductSupplier = iota + 1 // 金美福
	ProductSupplierZDS                            // 周大生
	ProductSupplierLM                             // 老庙
	ProductSupplierZLF                            // 周六福
	ProductSupplierJDF                            // 金大福
)

var ProductSupplierMap = map[ProductSupplier]string{
	ProductSupplierJMF: "金美福",
	ProductSupplierZDS: "周大生",
	ProductSupplierLM:  "老庙",
	ProductSupplierZLF: "周六福",
	ProductSupplierJDF: "金大福",
}

/* 产品品类 */
// 全部、戒指、项链、套链、吊坠、耳饰、手链、手镯、脚链、胸针、珠串、挂件、金条、摆件、饰品、其他
type ProductCategory int

const (
	ProductCategoryRing      ProductCategory = iota + 1 // 戒指
	ProductCategoryNecklace                             // 项链
	ProductCategoryChoker                               // 套链
	ProductCategoryPendant                              // 吊坠
	ProductCategoryEarring                              // 耳饰
	ProductCategoryBracelet                             // 手链
	ProductCategoryBangle                               // 手镯
	ProductCategoryAnklet                               // 脚链
	ProductCategoryBrooch                               // 胸针
	ProductCategoryBead                                 // 珠串
	ProductCategoryAccessory                            // 挂件
	ProductCategoryGoldBar                              // 金条
	ProductCategoryOrnament                             // 摆件
	ProductCategoryJewelry                              // 饰品
	ProductCategoryOther                                // 其他
)

var ProductCategoryMap = map[ProductCategory]string{
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

/* 产品工艺 */
// 全部、无、3D、5D、5G、古法、复古、万足金、九五金、珐琅彩、精品A、精品B、精品C、精品D
type ProductCraft int

const (
	ProductCraftNone     ProductCraft = iota + 1 // 无
	ProductCraft3D                               // 3D
	ProductCraft5D                               // 5D
	ProductCraft5G                               // 5G
	ProductCraftAncient                          // 古法
	ProductCraftRetro                            // 复古
	ProductCraftFiveGold                         // 万足金
	ProductCraftNineGold                         // 九五金
	ProductCraftFengLan                          // 珐琅彩
	ProductCraftFineA                            // 精品A
	ProductCraftFineB                            // 精品B
	ProductCraftFineC                            // 精品C
	ProductCraftFineD                            // 精品D
)

var ProductCraftMap = map[ProductCraft]string{
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

/* 产品颜色 */
// 全部、D、E、D-E、F、G、F-G、H、I、J、I-J、K、L、K-L、M、N、M-N
type ProductColor int

const (
	ProductColorD  ProductColor = iota + 1 // D
	ProductColorE                          // E
	ProductColorDE                         // D-E
	ProductColorF                          // F
	ProductColorG                          // G
	ProductColorFG                         // F-G
	ProductColorH                          // H
	ProductColorI                          // I
	ProductColorJ                          // J
	ProductColorIJ                         // I-J
	ProductColorK                          // K
	ProductColorL                          // L
	ProductColorKL                         // K-L
	ProductColorM                          // M
	ProductColorN                          // N
	ProductColorMN                         // M-N
)

var ProductColorMap = map[ProductColor]string{
	ProductColorD:  "D",
	ProductColorE:  "E",
	ProductColorDE: "D-E",
	ProductColorF:  "F",
	ProductColorG:  "G",
	ProductColorFG: "F-G",
	ProductColorH:  "H",
	ProductColorI:  "I",
	ProductColorJ:  "J",
	ProductColorIJ: "I-J",
	ProductColorK:  "K",
	ProductColorL:  "L",
	ProductColorKL: "K-L",
	ProductColorM:  "M",
	ProductColorN:  "N",
	ProductColorMN: "M-N",
}

/* 产品净度 */
// 全部、IF、VVS、VVS1、VVS2、VS、VS1、VS2、SI、SI1、SI2、SI3、P、FL、LC
type ProductClarity int

const (
	ProductClarityIF   ProductClarity = iota + 1 // IF
	ProductClarityVVS                            // VVS
	ProductClarityVVS1                           // VVS1
	ProductClarityVVS2                           // VVS2
	ProductClarityVS                             // VS
	ProductClarityVS1                            // VS1
	ProductClarityVS2                            // VS2
	ProductClaritySI                             // SI
	ProductClaritySI1                            // SI1
	ProductClaritySI2                            // SI2
	ProductClaritySI3                            // SI3
	ProductClarityP                              // P
	ProductClarityFL                             // FL
	ProductClarityLC                             // LC
)

var ProductClarityMap = map[ProductClarity]string{
	ProductClarityIF:   "IF",
	ProductClarityVVS:  "VVS",
	ProductClarityVVS1: "VVS1",
	ProductClarityVVS2: "VVS2",
	ProductClarityVS:   "VS",
	ProductClarityVS1:  "VS1",
	ProductClarityVS2:  "VS2",
	ProductClaritySI:   "SI",
	ProductClaritySI1:  "SI1",
	ProductClaritySI2:  "SI2",
	ProductClaritySI3:  "SI3",
	ProductClarityP:    "P",
	ProductClarityFL:   "FL",
	ProductClarityLC:   "LC",
}

/* 产品切工 */
// 全部、EX、VG、G、P
type ProductCut int

const (
	ProductCutEX ProductCut = iota + 1 // EX
	ProductCutVG                       // VG
	ProductCutG                        // G
	ProductCutP                        // P
)

var ProductCutMap = map[ProductCut]string{
	ProductCutEX: "EX",
	ProductCutVG: "VG",
	ProductCutG:  "G",
	ProductCutP:  "P",
}

/* 产品材质 */
// 全部、黄金、银饰、铂金、钯金、裸石
type ProductMaterial int

const (
	ProductMaterialGold      ProductMaterial = iota + 1 // 黄金
	ProductMaterialSilver                               // 银饰
	ProductMaterialPlatinum                             // 铂金
	ProductMaterialPalladium                            // 钯金
	ProductMaterialGem                                  // 裸石
)

var ProductMaterialMap = map[ProductMaterial]string{
	ProductMaterialGold:      "黄金",
	ProductMaterialSilver:    "银饰",
	ProductMaterialPlatinum:  "铂金",
	ProductMaterialPalladium: "钯金",
	ProductMaterialGem:       "裸石",
}

/* 产品大类 */
// 全部、足金（克）、足金（件）、金 750、金 916、铂金、银饰、足金镶嵌、裸钻、钻石、彩宝、玉石、珍珠、其他
type ProductClass int

const (
	ProductClassGoldKg    ProductClass = iota + 1 // 足金（克）
	ProductClassGoldPiece                         // 足金（件）
	ProductClassGold750                           // 金 750
	ProductClassGold916                           // 金 916
	ProductClassPlatinum                          // 铂金
	ProductClassSilver                            // 银饰
	ProductClassGoldInlay                         // 足金镶嵌
	ProductClassGem                               // 裸石
	ProductClassDiamond                           // 钻石
	ProductClassJade                              // 玉石
	ProductClassPearl                             // 珍珠
	ProductClassOther                             // 其他

)

var ProductClassMap = map[ProductClass]string{
	ProductClassGoldKg:    "足金（克）",
	ProductClassGoldPiece: "足金（件）",
	ProductClassGold750:   "金 750",
	ProductClassGold916:   "金 916",
	ProductClassPlatinum:  "铂金",
	ProductClassSilver:    "银饰",
	ProductClassGoldInlay: "足金镶嵌",
	ProductClassGem:       "裸石",
	ProductClassDiamond:   "钻石",
	ProductClassJade:      "玉石",
	ProductClassPearl:     "珍珠",
	ProductClassOther:     "其他",
}

/* 零售方式 */
// 全部、计件、计重工费按克、计重工费按件
type ProductRetailType int

const (
	ProductRetailTypePiece     ProductRetailType = iota + 1 // 计件
	ProductRetailTypeGoldKg                                 // 计重工费按克
	ProductRetailTypeGoldPiece                              // 计重工费按件
)

var ProductRetailTypeMap = map[ProductRetailType]string{
	ProductRetailTypePiece:     "计件",
	ProductRetailTypeGoldKg:    "计重工费按克",
	ProductRetailTypeGoldPiece: "计重工费按件",
}

/* 状态 */
// 全部、在库、维修中、调出在途
type ProductStatus int

const (
	ProductStatusInStock  ProductStatus = iota + 1 // 在库
	ProductStatusInRepair                          // 维修中
	ProductStatusOutStock                          // 调出在途
)

var ProductStatusMap = map[ProductStatus]string{
	ProductStatusInStock:  "在库",
	ProductStatusInRepair: "维修中",
	ProductStatusOutStock: "调出在途",
}
