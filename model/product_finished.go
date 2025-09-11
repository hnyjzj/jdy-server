package model

import (
	"jdy/enums"
	"jdy/types"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// 成品
type ProductFinished struct {
	SoftDelete

	Code        string                     `json:"code" gorm:"uniqueIndex;type:varchar(255);comment:条码;"`                            // 条码
	Name        string                     `json:"name" gorm:"type:varchar(255);comment:名称;"`                                        // 名称
	Status      enums.ProductStatus        `json:"status" gorm:"index:idx_pf_store_status_time,priority:2;type:int(11);comment:状态;"` // 状态
	Images      []string                   `json:"images" gorm:"type:text;serializer:json;comment:图片;"`                              // 图片
	Class       enums.ProductClassFinished `json:"class" gorm:"type:int(11);comment:大类;"`                                            // 大类
	AccessFee   decimal.Decimal            `json:"access_fee" gorm:"type:decimal(10,2);not NULL;comment:入网费;"`                       // 入网费
	RetailType  enums.ProductRetailType    `json:"retail_type" gorm:"type:int(11);not NULL;comment:零售方式;"`                           // 零售方式
	LabelPrice  decimal.Decimal            `json:"label_price" gorm:"type:decimal(10,2);not NULL;comment:标签价;"`                      // 标签价
	LaborFee    decimal.Decimal            `json:"labor_fee" gorm:"type:decimal(10,2);not NULL;comment:工费;"`                         // 工费
	Style       string                     `json:"style" gorm:"type:varchar(255);comment:款式;"`                                       // 款式
	Supplier    enums.ProductSupplier      `json:"supplier" gorm:"type:int(11);not NULL;comment:供应商;"`                               // 供应商
	Brand       enums.ProductBrand         `json:"brand" gorm:"type:int(11);comment:品牌;"`                                            // 品牌
	Material    enums.ProductMaterial      `json:"material" gorm:"type:int(11);not NULL;comment:材质;"`                                // 材质
	Quality     enums.ProductQuality       `json:"quality" gorm:"type:int(11);not NULL;comment:成色;"`                                 // 成色
	Gem         enums.ProductGem           `json:"gem" gorm:"type:int(11);not NULL;comment:主石;"`                                     // 主石
	Category    enums.ProductCategory      `json:"category" gorm:"type:int(11);not NULL;comment:品类;"`                                // 品类
	Craft       enums.ProductCraft         `json:"craft" gorm:"type:int(11);comment:工艺;"`                                            // 工艺
	WeightMetal decimal.Decimal            `json:"weight_metal" gorm:"type:decimal(15,4);comment:金重;"`                               // 金重
	WeightTotal decimal.Decimal            `json:"weight_total" gorm:"type:decimal(15,4);comment:总重;"`                               // 总重
	Size        string                     `json:"size" gorm:"type:varchar(255);comment:手寸;"`                                        // 手寸
	ColorMetal  string                     `json:"color_metal" gorm:"type:varchar(255);comment:贵金属颜色;"`                              // 贵金属颜色
	WeightGem   decimal.Decimal            `json:"weight_gem" gorm:"type:decimal(15,4);comment:主石重;"`                                // 主石重
	NumGem      int                        `json:"num_gem" gorm:"type:int(11);comment:主石数;"`                                         // 主石数
	WeightOther decimal.Decimal            `json:"weight_other" gorm:"type:decimal(15,4);comment:杂料重;"`                              // 杂料重
	NumOther    int                        `json:"num_other" gorm:"type:int(11);comment:杂料数;"`                                       // 杂料数
	ColorGem    enums.ProductColor         `json:"color_gem" gorm:"type:int(11);comment:颜色;"`                                        // 颜色
	Clarity     enums.ProductClarity       `json:"clarity" gorm:"type:int(11);comment:主石净度;"`                                        // 净度
	Certificate []string                   `json:"certificate" gorm:"type:text;serializer:json;comment:证书;"`                         // 证书
	Series      string                     `json:"series" gorm:"type:varchar(255);comment:系列;"`                                      // 系列
	Remark      string                     `json:"remark" gorm:"type:text;comment:备注;"`                                              // 备注

	StoreId        string `json:"store_id" gorm:"index:idx_pf_store_status_time,priority:1;type:varchar(255);comment:门店ID;"` // 门店ID
	Store          Store  `json:"store,omitempty" gorm:"foreignKey:StoreId;references:Id;comment:门店;"`                       // 门店
	IsSpecialOffer bool   `json:"is_special_offer" gorm:"comment:是否特价;"`                                                     // 是否特价

	EnterId   string                `json:"enter_id" gorm:"type:varchar(255);not NULL;comment:成品入库单ID;"`                    // 成品入库单ID
	Enter     *ProductFinishedEnter `json:"product_enter,omitempty" gorm:"foreignKey:EnterId;references:Id;comment:成品入库单;"` // 成品入库单
	EnterTime time.Time             `json:"enter_time" gorm:"index:idx_pf_store_status_time,priority:3;comment:入库时间;"`      // 入库时间
}

func (ProductFinished) WhereCondition(db *gorm.DB, query *types.ProductFinishedWhere) *gorm.DB {
	if query.Code != "" {
		db = db.Where("code = ?", query.Code) // 编码
	}
	if query.Name != "" {
		db = db.Where("name LIKE ?", "%"+query.Name+"%") // 名称
	}
	if query.AccessFee != nil {
		db = db.Where("access_fee = ?", query.AccessFee) // 入网费
	}
	if query.LabelPrice != nil {
		db = db.Where("label_price = ?", query.LabelPrice) // 标签价
	}
	if query.LaborFee != nil {
		db = db.Where("labor_fee = ?", query.LaborFee) // 工费
	}
	if query.Style != "" {
		db = db.Where("style LIKE (?)", "%"+query.Style+"%") // 款式
	}
	if query.WeightTotal != nil {
		db = db.Where("weight_total = ?", query.WeightTotal) // 总重
	}
	if query.Size != "" {
		db = db.Where("size LIKE (?)", "%"+query.Size+"%") // 手寸
	}
	if query.ColorMetal != "" {
		db = db.Where("color_metal LIKE (?)", "%"+query.ColorMetal+"%") // 贵金属颜色
	}
	if query.WeightMetal != nil {
		db = db.Where("weight_metal = ?", query.WeightMetal) // 金重
	}
	if query.WeightGem != nil {
		db = db.Where("weight_gem = ?", query.WeightGem) // 主石重
	}
	if query.WeightOther != nil {
		db = db.Where("weight_other = ?", query.WeightOther) // 杂料重
	}
	if query.NumGem != 0 {
		db = db.Where("num_gem = ?", query.NumGem) // 主石数
	}
	if query.NumOther != 0 {
		db = db.Where("num_other = ?", query.NumOther) // 杂料数
	}
	if query.ColorGem != 0 {
		db = db.Where("color_gem = ?", query.ColorGem) // 主石颜色
	}
	if query.Category != 0 {
		db = db.Where("category = ?", query.Category) // 品类
	}
	if query.RetailType != 0 {
		db = db.Where("retail_type = ?", query.RetailType) // 销售方式
	}
	if query.Class != 0 {
		db = db.Where("class = ?", query.Class) // 大类
	}
	if query.Supplier != 0 {
		db = db.Where("supplier = ?", query.Supplier) // 供应商
	}
	if query.Material != 0 {
		db = db.Where("material = ?", query.Material) // 材质
	}
	if query.Quality != 0 {
		db = db.Where("quality = ?", query.Quality) // 成色
	}
	if query.Gem != 0 {
		db = db.Where("gem = ?", query.Gem) // 主石
	}
	if query.Clarity != 0 {
		db = db.Where("clarity = ?", query.Clarity) // 主石净度
	}
	if query.Series != "" {
		db = db.Where("series LIKE (?)", "%"+query.Series+"%") // 系列
	}
	if query.Remark != "" {
		db = db.Where("remark LIKE (?)", "%"+query.Remark+"%") // 备注
	}
	if query.Brand != 0 {
		db = db.Where("brand = ?", query.Brand) // 品牌
	}
	if query.Craft != 0 {
		db = db.Where("craft = ?", query.Craft) // 工艺
	}
	if query.IsSpecialOffer != nil {
		db = db.Where("is_special_offer = ?", query.IsSpecialOffer) // 是否特价
	}
	if query.Status != 0 {
		db = db.Where("status = ?", query.Status) // 状态
	}
	if query.StoreId != "" {
		db = db.Where("store_id = ?", query.StoreId) // 门店ID
	}
	if query.EnterId != "" {
		db = db.Where("enter_id = ?", query.EnterId) // 入库单ID
	}

	return db
}

func (ProductFinished) Preloads(db *gorm.DB) *gorm.DB {
	db = db.Preload("Store")
	db = db.Preload("Enter")

	return db
}

func (p *ProductFinished) GetClass() enums.ProductClassFinished {
	switch {
	case p.RetailType == enums.ProductRetailTypeGoldKg &&
		p.Material == enums.ProductMaterialGold &&
		(p.Quality == enums.ProductQuality99999 || p.Quality == enums.ProductQuality9999 || p.Quality == enums.ProductQuality999) &&
		p.Gem == enums.ProductGemGold:
		{
			return enums.ProductClassFinishedGoldKg
		}
	case p.RetailType == enums.ProductRetailTypeGoldPiece &&
		p.Material == enums.ProductMaterialGold &&
		(p.Quality == enums.ProductQuality99999 || p.Quality == enums.ProductQuality9999 || p.Quality == enums.ProductQuality999) &&
		p.Gem == enums.ProductGemGold:
		{
			return enums.ProductClassFinishedGoldKg
		}
	case p.RetailType == enums.ProductRetailTypePiece &&
		p.Material == enums.ProductMaterialGold &&
		p.Quality == enums.ProductQuality999 &&
		p.Gem == enums.ProductGemGold:
		{
			return enums.ProductClassFinishedGoldPiece
		}
	case p.Material == enums.ProductMaterialGold &&
		p.Quality == enums.ProductQuality750 &&
		p.Gem == enums.ProductGemGold:
		{
			return enums.ProductClassFinishedGold750
		}
	case p.Material == enums.ProductMaterialGold &&
		p.Quality == enums.ProductQuality916 &&
		p.Gem == enums.ProductGemGold:
		{
			return enums.ProductClassFinishedGold916
		}
	case p.Material == enums.ProductMaterialPlatinum &&
		(p.Quality == enums.ProductQuality999 || p.Quality == enums.ProductQuality990 || p.Quality == enums.ProductQuality950) &&
		p.Gem == enums.ProductGemGold:
		{
			return enums.ProductClassFinishedPlatinum
		}
	case p.Material == enums.ProductMaterialSilver &&
		(p.Quality == enums.ProductQuality990 || p.Quality == enums.ProductQuality925) &&
		p.Gem == enums.ProductGemGold:
		{
			return enums.ProductClassFinishedSilver
		}
	case p.Material == enums.ProductMaterialGold &&
		p.Quality == enums.ProductQuality999 &&
		p.Gem != enums.ProductGemGold:
		{
			return enums.ProductClassFinishedGoldInlay
		}
	case p.Material == enums.ProductMaterialGem &&
		p.Quality == enums.ProductQualityGem &&
		p.Gem == enums.ProductGemDiamond:
		{
			return enums.ProductClassFinishedDiamondNaked
		}
	case p.Material == enums.ProductMaterialGold &&
		p.Quality == enums.ProductQuality750 &&
		p.Gem == enums.ProductGemDiamond:
		{
			return enums.ProductClassFinishedDiamond
		}
	case p.Material == enums.ProductMaterialPlatinum &&
		p.Quality == enums.ProductQuality950 &&
		p.Gem == enums.ProductGemDiamond:
		{
			return enums.ProductClassFinishedDiamond
		}
	case p.Material == enums.ProductMaterialGold &&
		p.Quality == enums.ProductQuality750 &&
		(p.Gem != enums.ProductGemGold && p.Gem != enums.ProductGemPearl &&
			p.Gem != enums.ProductGemDiamond && p.Gem != enums.ProductGemPearlMother &&
			p.Gem != enums.ProductGemJade && p.Gem != enums.ProductGemJadeite):
		{
			return enums.ProductClassFinishedCoral
		}
	case p.Material == enums.ProductMaterialGem &&
		(p.Gem == enums.ProductGemPearl || p.Gem == enums.ProductGemPearlMother):
		{
			return enums.ProductClassFinishedPearl
		}
	case p.Material == enums.ProductMaterialGem &&
		(p.Gem == enums.ProductGemJade || p.Gem == enums.ProductGemJadeite ||
			p.Gem == enums.ProductGemOpal || p.Gem == enums.ProductGemJasper ||
			p.Gem == enums.ProductGemEmerald || p.Gem == enums.ProductGemGarnet):
		{
			return enums.ProductClassFinishedJade
		}
	default:
		{
			return enums.ProductClassFinishedOther
		}
	}
}

// 是否滞销（180 天）
func (p *ProductFinished) IsUnsalable(t time.Time) bool {
	const unsalableDays = 180
	cutoff := t.AddDate(0, 0, -unsalableDays)
	return p.EnterTime.Before(cutoff)
}

// 成品入库单
type ProductFinishedEnter struct {
	SoftDelete

	StoreId string `json:"store_id" gorm:"index:idx_pf_store_status_time,priority:1;type:varchar(255);not NULL;comment:门店ID;"` // 门店ID
	Store   *Store `json:"store" gorm:"foreignKey:StoreId;references:Id;comment:门店;"`                                          // 门店

	Remark string                   `json:"remark" gorm:"type:text;comment:备注;"`                                                       // 备注
	Status enums.ProductEnterStatus `json:"status" gorm:"index:idx_pf_store_status_time,priority:2;type:int(11);not NULL;comment:状态;"` // 状态

	Products                []ProductFinished `json:"products" gorm:"foreignKey:EnterId;references:Id;comment:成品;"`                 // 成品
	ProductCount            int64             `json:"product_count" gorm:"type:int(11);not NULL;comment:成品数量;"`                     // 成品数量
	ProductTotalWeightMetal decimal.Decimal   `json:"product_total_weight_metal" gorm:"type:decimal(15,4);not NULL;comment:成品总重;"`  // 成品总重
	ProductTotalLabelPrice  decimal.Decimal   `json:"product_total_label_price" gorm:"type:decimal(10,2);not NULL;comment:成品总标签价;"` // 成品总标签价
	ProductTotalAccessFee   decimal.Decimal   `json:"product_total_access_fee" gorm:"type:decimal(10,2);not NULL;comment:成品总加工费;"`  // 成品总加工费

	OperatorId string `json:"operator_id" gorm:"type:varchar(255);not NULL;comment:操作人ID;"`     // 操作人ID
	Operator   *Staff `json:"operator" gorm:"foreignKey:OperatorId;references:Id;comment:操作人;"` // 操作人
	IP         string `json:"ip" gorm:"type:varchar(255);not NULL;comment:IP;"`                 // IP
}

func (ProductFinishedEnter) WhereCondition(db *gorm.DB, query *types.ProductFinishedEnterWhere) *gorm.DB {
	if query.Id != "" {
		db = db.Where("id = ?", query.Id)
	}
	if query.StoreId != "" {
		db = db.Where("store_id = ?", query.StoreId)
	}
	if query.Status != 0 {
		db = db.Where("status = ?", query.Status)
	}
	if query.Remark != "" {
		db = db.Where("remark LIKE ?", "%"+query.Remark+"%")
	}
	if query.StartAt != nil && query.EndAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", query.StartAt, query.EndAt)
	}
	return db
}

func init() {
	// 注册模型
	RegisterModels(
		&ProductFinished{},
		&ProductFinishedEnter{},
	)
	// 重置表
	RegisterRefreshModels(
	// &ProductFinished{},
	// &ProductFinishedEnter{},
	)
}
