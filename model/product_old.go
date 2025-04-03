package model

import (
	"jdy/enums"
	"jdy/types"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type ProductOld struct {
	SoftDelete

	Code        string                `json:"code" gorm:"uniqueIndex;type:varchar(255);comment:条码;"`       // 条码
	Name        string                `json:"name" gorm:"type:varchar(255);comment:名称;"`                   // 名称
	Status      enums.ProductStatus   `json:"status" gorm:"type:tinyint(2);comment:状态;"`                   // 状态
	LabelPrice  decimal.Decimal       `json:"label_price" gorm:"type:decimal(10,2);not NULL;comment:标签价;"` // 标签价
	Brand       enums.ProductBrand    `json:"brand" gorm:"type:tinyint(2);comment:品牌;"`                    // 品牌
	Material    enums.ProductMaterial `json:"material" gorm:"type:tinyint(2);not NULL;comment:材质;"`        // 材质
	Quality     enums.ProductQuality  `json:"quality" gorm:"type:tinyint(2);not NULL;comment:成色;"`         // 成色
	Gem         enums.ProductGem      `json:"gem" gorm:"type:tinyint(2);not NULL;comment:主石;"`             // 主石
	Category    enums.ProductCategory `json:"category" gorm:"type:tinyint(2);not NULL;comment:品类;"`        // 品类
	Craft       enums.ProductCraft    `json:"craft" gorm:"type:tinyint(2);comment:工艺;"`                    // 工艺
	WeightMetal decimal.Decimal       `json:"weight_metal" gorm:"type:decimal(10,2);comment:金重;"`          // 金重
	WeightTotal decimal.Decimal       `json:"weight_total" gorm:"type:decimal(10,2);comment:总重;"`          // 总重
	ColorGem    enums.ProductColor    `json:"color_gem" gorm:"type:tinyint(2);comment:颜色;"`                // 颜色
	WeightGem   decimal.Decimal       `json:"weight_gem" gorm:"type:decimal(10,2);comment:主石重;"`           // 主石重
	NumGem      int                   `json:"num_gem" gorm:"type:tinyint(3);comment:主石数;"`                 // 主石数
	Clarity     enums.ProductClarity  `json:"clarity" gorm:"type:tinyint(2);comment:主石净度;"`                // 主石净度
	Cut         enums.ProductCut      `json:"cut" gorm:"type:tinyint(2);comment:主石切工;"`                    // 主石切工
	WeightOther decimal.Decimal       `json:"weight_other" gorm:"type:decimal(10,2);comment:杂料重;"`         // 杂料重
	NumOther    int                   `json:"num_other" gorm:"type:tinyint(2);comment:杂料数;"`               // 杂料数
	Remark      string                `json:"remark" gorm:"type:text;comment:备注;"`                         // 备注

	StoreId string `json:"store_id" gorm:"type:varchar(255);comment:门店ID;"`                     // 门店ID
	Store   Store  `json:"store,omitempty" gorm:"foreignKey:StoreId;references:Id;comment:门店;"` // 门店

	IsOur                   bool                       `json:"is_our" gorm:"comment:是否本司货品;"`                                                        // 是否本司货品
	RecycleMethod           enums.ProductRecycleMethod `json:"recycle_method,omitempty" gorm:"type:tinyint(2);comment:回收方式;"`                        // 回收方式
	RecycleType             enums.ProductRecycleType   `json:"recycle_type,omitempty" gorm:"type:tinyint(2);comment:回收类型;"`                          // 回收类型
	RecyclePrice            decimal.Decimal            `json:"recycle_price" gorm:"type:decimal(10,2);comment:回收价格;"`                                // 回收价格
	RecyclePriceGold        decimal.Decimal            `json:"recycle_price_gold" gorm:"type:decimal(10,2);comment:回收金价;"`                           // 回收金价
	RecyclePriceLabor       decimal.Decimal            `json:"recycle_price_labor" gorm:"type:decimal(10,2);comment:回收工费;"`                          // 回收工费
	RecyclePriceLaborMethod enums.ProductRecycleMethod `json:"recycle_price_labor_method,omitempty" gorm:"type:tinyint(2);comment:回收工费方式;"`          // 回收工费方式
	QualityActual           decimal.Decimal            `json:"quality_actual" gorm:"type:decimal(3,2);comment:实际成色;"`                                // 实际成色
	RecycleSource           enums.ProductRecycleSource `json:"recycle_source,omitempty" gorm:"type:tinyint(2);comment:回收来源;"`                        // 回收来源
	RecycleSourceId         string                     `json:"recycle_source_id" gorm:"type:varchar(255);comment:回收来源ID;"`                           // 回收来源ID
	RecycleStoreId          string                     `json:"recycle_store_id" gorm:"type:varchar(255);comment:回收门店ID;"`                            // 回收门店ID
	RecycleStore            Store                      `json:"recycle_store,omitempty" gorm:"foreignKey:RecycleStoreId;references:Id;comment:回收门店;"` // 回收门店
}

func (ProductOld) WhereCondition(db *gorm.DB, query *types.ProductOldWhere) *gorm.DB {
	if query.Code != "" {
		db = db.Where("code LIKE ?", "%"+query.Code+"%")
	}
	if query.Name != "" {
		db = db.Where("name LIKE ?", "%"+query.Name+"%")
	}
	if query.Status != 0 {
		db = db.Where("status = ?", query.Status)
	}
	if query.LabelPrice != nil {
		db = db.Where("label_price = ?", query.LabelPrice)
	}
	if query.Brand != 0 {
		db = db.Where("brand = ?", query.Brand)
	}
	if query.Material != 0 {
		db = db.Where("material = ?", query.Material)
	}
	if query.Quality != 0 {
		db = db.Where("quality = ?", query.Quality)
	}
	if query.Gem != 0 {
		db = db.Where("gem = ?", query.Gem)
	}
	if query.Category != 0 {
		db = db.Where("category = ?", query.Category)
	}
	if query.Craft != 0 {
		db = db.Where("craft = ?", query.Craft)
	}
	if query.WeightMetal != nil {
		db = db.Where("weight_metal = ?", query.WeightMetal)
	}
	if query.WeightTotal != nil {
		db = db.Where("weight_total = ?", query.WeightTotal)
	}
	if query.ColorGem != 0 {
		db = db.Where("color_gem = ?", query.ColorGem)
	}
	if query.WeightGem != nil {
		db = db.Where("weight_gem = ?", query.WeightGem)
	}
	if query.NumGem != 0 {
		db = db.Where("num_gem = ?", query.NumGem)
	}
	if query.Clarity != 0 {
		db = db.Where("clarity = ?", query.Clarity)
	}
	if query.Cut != 0 {
		db = db.Where("cut = ?", query.Cut)
	}
	if query.WeightOther != nil {
		db = db.Where("weight_other = ?", query.WeightOther)
	}
	if query.NumOther != 0 {
		db = db.Where("num_other = ?", query.NumOther)
	}
	if query.Remark != "" {
		db = db.Where("remark LIKE ?", "%"+query.Remark+"%")
	}
	if query.StoreId != "" {
		db = db.Where("store_id = ?", query.StoreId)
	}
	if query.IsOur != nil {
		db = db.Where("is_our = ?", query.IsOur)
	}
	if query.RecycleMethod != 0 {
		db = db.Where("recycle_method = ?", query.RecycleMethod)
	}
	if query.RecycleType != 0 {
		db = db.Where("recycle_type = ?", query.RecycleType)
	}
	if query.RecyclePrice != nil {
		db = db.Where("recycle_price = ?", query.RecyclePrice)
	}
	if query.RecyclePriceGold != nil {
		db = db.Where("recycle_price_gold = ?", query.RecyclePriceGold)
	}
	if query.RecyclePriceLabor != nil {
		db = db.Where("recycle_price_labor = ?", query.RecyclePriceLabor)
	}
	if query.RecyclePriceLaborMethod != 0 {
		db = db.Where("recycle_price_labor_method = ?", query.RecyclePriceLaborMethod)
	}
	if query.QualityActual != nil {
		db = db.Where("quality_actual = ?", query.QualityActual)
	}
	if query.RecycleSource != 0 {
		db = db.Where("recycle_source = ?", query.RecycleSource)
	}
	if query.RecycleSourceId != "" {
		db = db.Where("recycle_source_id = ?", query.RecycleSourceId)
	}
	if query.RecycleStoreId != "" {
		db = db.Where("recycle_store_id = ?", query.RecycleStoreId)
	}

	return db
}

// 产品调拨单
type ProductOldAllocate struct {
	SoftDelete

	Method enums.ProductAllocateMethod `json:"method" gorm:"type:tinyint(2);not NULL;comment:调拨类型;"` // 调拨类型
	Type   enums.ProductType           `json:"type" gorm:"type:tinyint(2);not NULL;comment:产品类型;"`   // 仓库类型
	Reason enums.ProductAllocateReason `json:"reason" gorm:"type:tinyint(2);not NULL;comment:调拨原因;"` // 调拨原因
	Status enums.ProductAllocateStatus `json:"status" gorm:"type:tinyint(2);comment:状态;"`            // 状态
	Remark string                      `json:"remark" gorm:"type:text;comment:备注;"`                  // 备注

	FromStoreId string `json:"from_store_id" gorm:"type:varchar(255);comment:调出门店;"` // 调出门店
	FromStore   *Store `json:"from_store" gorm:"foreignKey:FromStoreId;references:Id;comment:调出门店;"`
	ToStoreId   string `json:"to_store_id" gorm:"type:varchar(255);comment:调入门店;"` // 调入门店
	ToStore     *Store `json:"to_store" gorm:"foreignKey:ToStoreId;references:Id;comment:调入门店;"`

	Products []ProductOld `json:"products" gorm:"many2many:product_old_allocate_products;comment:产品;"`

	OperatorId string `json:"operator_id" gorm:"type:varchar(255);not NULL;comment:操作人ID;"`     // 操作人ID
	Operator   *Staff `json:"operator" gorm:"foreignKey:OperatorId;references:Id;comment:操作人;"` // 操作人
	IP         string `json:"-" gorm:"type:varchar(255);not NULL;comment:IP;"`                  // IP
}

func (ProductOldAllocate) WhereCondition(db *gorm.DB, query *types.ProductOldAllocateWhere) *gorm.DB {
	if query.Method != 0 {
		db = db.Where("method = ?", query.Method)
	}
	if query.Type != 0 {
		db = db.Where("type = ?", query.Type)
	}
	if query.Reason != 0 {
		db = db.Where("reason = ?", query.Reason)
	}
	if query.FromStoreId != "" {
		db = db.Where("from_store_id = ?", query.FromStoreId)
	}
	if query.ToStoreId != "" {
		db = db.Where("to_store_id = ?", query.ToStoreId)
	}
	if query.StartTime != nil && query.EndTime != nil {
		db = db.Where("created_at BETWEEN ? AND ?", query.StartTime, query.EndTime)
	}
	return db
}

func init() {
	// 注册模型
	RegisterModels(
		&ProductOld{},
		&ProductOldAllocate{},
	)
	// 重置表
	RegisterRefreshModels(
	// &ProductOld{},
	// &ProductOldAllocate{},
	)
}
