package model

import (
	"fmt"
	"jdy/enums"
	"jdy/types"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// 成品
type Product struct {
	SoftDelete

	Type enums.ProductType `json:"type" gorm:"type:tinyint(2);comment:类型;"` // 类型

	Code           string                  `json:"code" gorm:"uniqueIndex;type:varchar(255);comment:条码;"`       // 条码
	Name           string                  `json:"name" gorm:"type:varchar(255);not NULL;comment:名称;"`          // 名称
	Status         enums.ProductStatus     `json:"status" gorm:"type:tinyint(2);comment:状态;"`                   // 状态
	Images         []string                `json:"images" gorm:"type:text;serializer:json;comment:图片;"`         // 图片
	Class          enums.ProductClass      `json:"class" gorm:"type:tinyint(2);comment:大类;"`                    // 大类
	AccessFee      decimal.Decimal         `json:"access_fee" gorm:"type:decimal(10,2);not NULL;comment:入网费;"`  // 入网费
	RetailType     enums.ProductRetailType `json:"retail_type" gorm:"type:tinyint(2);not NULL;comment:零售方式;"`   // 零售方式
	LabelPrice     decimal.Decimal         `json:"label_price" gorm:"type:decimal(10,2);not NULL;comment:标签价;"` // 标签价
	LaborFee       decimal.Decimal         `json:"labor_fee" gorm:"type:decimal(10,2);not NULL;comment:工费;"`    // 工费
	Style          string                  `json:"style" gorm:"type:varchar(255);comment:款式;"`                  // 款式
	Supplier       enums.ProductSupplier   `json:"supplier" gorm:"type:tinyint(2);not NULL;comment:供应商;"`       // 供应商
	Brand          enums.ProductBrand      `json:"brand" gorm:"type:tinyint(2);comment:品牌;"`                    // 品牌
	Material       enums.ProductMaterial   `json:"material" gorm:"type:tinyint(2);not NULL;comment:材质;"`        // 材质
	Quality        enums.ProductQuality    `json:"quality" gorm:"type:tinyint(2);not NULL;comment:成色;"`         // 成色
	Gem            enums.ProductGem        `json:"gem" gorm:"type:tinyint(2);not NULL;comment:主石;"`             // 主石
	Category       enums.ProductCategory   `json:"category" gorm:"type:tinyint(2);not NULL;comment:品类;"`        // 品类
	Craft          enums.ProductCraft      `json:"craft" gorm:"type:tinyint(2);comment:工艺;"`                    // 工艺
	WeightMetal    decimal.Decimal         `json:"weight_metal" gorm:"type:decimal(10,2);comment:金重;"`          // 金重
	WeightTotal    decimal.Decimal         `json:"weight_total" gorm:"type:decimal(10,2);comment:总重;"`          // 总重
	Size           string                  `json:"size" gorm:"type:varchar(255);comment:手寸;"`                   // 手寸
	ColorMetal     string                  `json:"color_metal" gorm:"type:varchar(255);comment:贵金属颜色;"`         // 贵金属颜色
	WeightGem      decimal.Decimal         `json:"weight_gem" gorm:"type:decimal(10,2);comment:主石重;"`           // 主石重
	NumGem         int                     `json:"num_gem" gorm:"type:tinyint(3);comment:主石数;"`                 // 主石数
	WeightOther    decimal.Decimal         `json:"weight_other" gorm:"type:decimal(10,2);comment:杂料重;"`         // 杂料重
	NumOther       int                     `json:"num_other" gorm:"type:tinyint(2);comment:杂料数;"`               // 杂料数
	ColorGem       enums.ProductColor      `json:"color_gem" gorm:"type:tinyint(2);comment:主石色;"`               // 主石色
	Clarity        enums.ProductClarity    `json:"clarity" gorm:"type:tinyint(2);comment:主石净度;"`                // 净度
	Certificate    []string                `json:"certificate" gorm:"type:text;serializer:json;comment:证书;"`    // 证书
	Series         string                  `json:"series" gorm:"type:varchar(255);comment:系列;"`                 // 系列
	Remark         string                  `json:"remark" gorm:"type:text;comment:备注;"`                         // 备注
	IsSpecialOffer bool                    `json:"is_special_offer" gorm:"comment:是否特价;"`                       // 是否特价

	StoreId string `json:"store_id" gorm:"type:varchar(255);comment:门店ID;"`           // 门店ID
	Store   Store  `json:"store" gorm:"foreignKey:StoreId;references:Id;comment:门店;"` // 门店

	ProductEnterId string        `json:"product_enter_id" gorm:"type:varchar(255);not NULL;comment:产品入库单ID;"`         // 产品入库单ID
	ProductEnter   *ProductEnter `json:"product_enter" gorm:"foreignKey:ProductEnterId;references:Id;comment:产品入库单;"` // 产品入库单
	EnterTime      time.Time     `json:"enter_time" gorm:"comment:入库时间;"`                                             // 入库时间

	IsOur                   bool                       `json:"is_our" gorm:"comment:是否本司货品;"`                                              // 是否本司货品
	RecycleMethod           enums.ProductRecycleMethod `json:"recycle_method" gorm:"type:tinyint(2);comment:回收方式;"`                        // 回收方式
	RecycleType             enums.ProductRecycleType   `json:"recycle_type" gorm:"type:tinyint(2);comment:回收类型;"`                          // 回收类型
	RecyclePrice            decimal.Decimal            `json:"recycle_price" gorm:"type:decimal(10,2);comment:回收价格;"`                      // 回收价格
	RecyclePriceGold        decimal.Decimal            `json:"recycle_price_gold" gorm:"type:decimal(10,2);comment:回收金价;"`                 // 回收金价
	RecyclePriceLabor       decimal.Decimal            `json:"recycle_price_labor" gorm:"type:decimal(10,2);comment:回收工费;"`                // 回收工费
	RecyclePriceLaborMethon enums.ProductRecycleMethod `json:"recycle_price_labor_methon" gorm:"type:tinyint(2);comment:回收工费方式;"`          // 回收工费方式
	QualityActual           decimal.Decimal            `json:"quality_actual" gorm:"type:decimal(3,2);comment:实际成色;"`                      // 实际成色
	RecycleSource           enums.ProductRecycleSource `json:"recycle_source" gorm:"type:tinyint(2);comment:回收来源;"`                        // 回收来源
	RecycleStoreId          string                     `json:"recycle_store_id" gorm:"type:varchar(255);comment:回收门店ID;"`                  // 回收门店ID
	RecycleStore            Store                      `json:"recycle_store" gorm:"foreignKey:RecycleStoreId;references:Id;comment:回收门店;"` // 回收门店

	TypePart enums.ProductTypePart `json:"type_part" gorm:"type:tinyint(2);comment:配件类型;"` // 配件类型
	Stock    int64                 `json:"stock" gorm:"default:1;comment:库存;"`             // 库存
}

func (Product) WhereCondition(db *gorm.DB, query *types.ProductWhere) *gorm.DB {
	if query.Code != "" {
		db = db.Where("code = ?", fmt.Sprint(query.Code))
	}
	if query.Name != "" {
		db = db.Where("name LIKE ?", fmt.Sprintf("%%%s%%", query.Name))
	}
	if !query.AccessFee.IsZero() {
		db = db.Where("access_fee = ?", query.AccessFee)
	}
	if !query.Price.IsZero() {
		db = db.Where("price = ?", query.Price)
	}
	if !query.LaborFee.IsZero() {
		db = db.Where("labor_fee = ?", query.LaborFee)
	}
	if !query.Weight.IsZero() {
		db = db.Where("weight = ?", query.Weight)
	}
	if !query.WeightMetal.IsZero() {
		db = db.Where("weight_metal = ?", query.WeightMetal)
	}
	if !query.WeightGem.IsZero() {
		db = db.Where("weight_gem = ?", query.WeightGem)
	}
	if !query.WeightOther.IsZero() {
		db = db.Where("weight_other = ?", query.WeightOther)
	}
	if query.NumGem != 0 {
		db = db.Where("num_gem = ?", int(query.NumGem))
	}
	if query.NumOther != 0 {
		db = db.Where("num_other = ?", int(query.NumOther))
	}
	if query.ColorMetal != 0 {
		db = db.Where("color_metal = ?", query.ColorMetal)
	}
	if query.ColorGem != 0 {
		db = db.Where("color_gem = ?", query.ColorGem)
	}
	if query.Clarity != 0 {
		db = db.Where("clarity = ?", query.Clarity)
	}
	if query.RetailType != 0 {
		db = db.Where("retail_type = ?", query.RetailType)
	}
	if query.Class != 0 {
		db = db.Where("class = ?", query.Class)
	}
	if query.Supplier != 0 {
		db = db.Where("supplier = ?", query.Supplier)
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
	if query.Brand != 0 {
		db = db.Where("brand = ?", query.Brand)
	}
	if query.Craft != 0 {
		db = db.Where("craft = ?", query.Craft)
	}
	if query.Style != "" {
		db = db.Where("style LIKE ?", fmt.Sprintf("%%%s%%", query.Name))
	}
	if query.Size != "" {
		db = db.Where("size LIKE ?", fmt.Sprintf("%%%s%%", query.Name))
	}
	if query.IsSpecialOffer {
		db = db.Where("is_special_offer = ?", query.IsSpecialOffer)
	}
	if query.Status != 0 {
		db = db.Where("status = ?", query.Status)
	}
	if query.Type != 0 {
		db = db.Where("type = ?", query.Type)
	}
	if query.ProductEnterId != "" {
		db = db.Where("product_enter_id = ?", query.ProductEnterId)
	}
	if query.StoreId != "" {
		db = db.Where("store_id = ?", query.StoreId)
	}

	return db
}

type ProductHistory struct {
	Model

	Action enums.ProductAction `json:"action" gorm:"type:tinyint(2);comment:操作;"` // 操作

	Key      string `json:"key" gorm:"type:varchar(255);comment:键;"`                // 键
	Value    any    `json:"value" gorm:"type:text;serializer:json;comment:值;"`      // 值
	OldValue any    `json:"old_value" gorm:"type:text;serializer:json;comment:旧值;"` // 旧值

	ProductId string   `json:"product_id" gorm:"type:varchar(255);not NULL;comment:产品ID;"`    // 产品ID
	Product   *Product `json:"product" gorm:"foreignKey:ProductId;references:Id;comment:产品;"` // 产品

	StoreId string `json:"store_id" gorm:"type:varchar(255);not NULL;comment:门店ID;"`  // 门店ID
	Store   *Store `json:"store" gorm:"foreignKey:StoreId;references:Id;comment:门店;"` // 门店

	SourceId string `json:"source_id" gorm:"column:source_id;size:255;not NULL;comment:来源id;"` // 来源id

	OperatorId string `json:"operator_id" gorm:"type:varchar(255);not NULL;comment:操作人ID;"`     // 操作人ID
	Operator   *Staff `json:"operator" gorm:"foreignKey:OperatorId;references:Id;comment:操作人;"` // 操作人
	IP         string `json:"ip" gorm:"type:varchar(255);not NULL;comment:IP;"`                 // IP
}

func (ProductHistory) WhereCondition(db *gorm.DB, query *types.ProductHistoryWhere) *gorm.DB {
	if query.ProductId != "" {
		db = db.Where("product_id = ?", query.ProductId)
	}
	if query.StoreId != "" {
		db = db.Where("store_id = ?", query.StoreId)
	}
	if query.Action != 0 {
		db = db.Where("action = ?", query.Action)
	}

	return db
}

// 产品入库单
type ProductEnter struct {
	SoftDelete

	Products []Product `json:"products" gorm:"foreignKey:ProductEnterId;references:Id;comment:产品;"` // 产品

	StoreId string `json:"store_id" gorm:"type:varchar(255);not NULL;comment:门店ID;"`  // 门店ID
	Store   *Store `json:"store" gorm:"foreignKey:StoreId;references:Id;comment:门店;"` // 门店

	OperatorId string `json:"operator_id" gorm:"type:varchar(255);not NULL;comment:操作人ID;"`     // 操作人ID
	Operator   *Staff `json:"operator" gorm:"foreignKey:OperatorId;references:Id;comment:操作人;"` // 操作人

	IP string `json:"ip" gorm:"type:varchar(255);not NULL;comment:IP;"` // IP
}

func (ProductEnter) WhereCondition(db *gorm.DB, query *types.ProductEnterWhere) *gorm.DB {
	if query.Id != "" {
		db = db.Where("id = ?", query.Id)
	}
	if query.StartTime != nil && query.EndTime != nil {
		db = db.Where("created_at BETWEEN ? AND ?", query.StartTime, query.EndTime)
	}
	return db
}

type ProductDamage struct {
	SoftDelete

	ProductId string   `json:"product_id" gorm:"type:varchar(255);not NULL;comment:产品ID;"`
	Product   *Product `json:"product" gorm:"foreignKey:ProductId;references:Id;comment:产品;"`

	OperatorId string `json:"-" gorm:"type:varchar(255);not NULL;comment:操作人ID;"`        // 操作人ID
	Operator   *Staff `json:"-" gorm:"foreignKey:OperatorId;references:Id;comment:操作人;"` // 操作人

	Reason string `json:"reason" gorm:"type:text;not NULL;comment:原因;"`    // 原因
	IP     string `json:"-" gorm:"type:varchar(255);not NULL;comment:IP;"` // IP
}

type ProductAllocate struct {
	SoftDelete

	Method enums.ProductAllocateMethod `json:"method" gorm:"type:tinyint(2);not NULL;comment:调拨方式;"` // 调拨方式
	Type   enums.ProductType           `json:"type" gorm:"type:tinyint(2);not NULL;comment:产品类型;"`   // 仓库类型
	Reason enums.ProductAllocateReason `json:"reason" gorm:"type:tinyint(2);not NULL;comment:调拨原因;"` // 调拨原因
	Status enums.ProductAllocateStatus `json:"status" gorm:"type:tinyint(2);comment:状态;"`            // 状态
	Remark string                      `json:"remark" gorm:"type:text;comment:备注;"`                  // 备注

	FromStoreId string `json:"from_store_id" gorm:"type:varchar(255);comment:调出门店;"` // 调出门店
	FromStore   *Store `json:"from_store" gorm:"foreignKey:FromStoreId;references:Id;comment:调出门店;"`
	ToStoreId   string `json:"to_store_id" gorm:"type:varchar(255);comment:调入门店;"` // 调入门店
	ToStore     *Store `json:"to_store" gorm:"foreignKey:ToStoreId;references:Id;comment:调入门店;"`

	Products []Product `json:"product" gorm:"many2many:product_allocate_list;"` // 产品

	OperatorId string `json:"operator_id" gorm:"type:varchar(255);not NULL;comment:操作人ID;"`     // 操作人ID
	Operator   *Staff `json:"operator" gorm:"foreignKey:OperatorId;references:Id;comment:操作人;"` // 操作人
	IP         string `json:"-" gorm:"type:varchar(255);not NULL;comment:IP;"`                  // IP
}

func (ProductAllocate) WhereCondition(db *gorm.DB, query *types.ProductAllocateWhere) *gorm.DB {
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

// 产品盘点
type ProductInventory struct {
	SoftDelete

	StoreId string `json:"store_id" gorm:"type:varchar(255);not NULL;comment:门店ID;"`  // 门店ID
	Store   Store  `json:"store" gorm:"foreignKey:StoreId;references:Id;comment:门店;"` // 门店

	InventoryPersonId string `json:"inventory_person_id" gorm:"type:varchar(255);not NULL;comment:盘点人ID;"`            // 盘点人ID
	InventoryPerson   Staff  `json:"inventory_person" gorm:"foreignKey:InventoryPersonId;references:Id;comment:盘点人;"` // 盘点人

	InspectorId string `json:"inspector_id" gorm:"type:varchar(255);not NULL;comment:监盘人ID;"`      // 监盘人ID
	Inspector   Staff  `json:"inspector" gorm:"foreignKey:InspectorId;references:Id;comment:监盘人;"` // 监盘人

	Type  enums.ProductType           `json:"type" gorm:"type:tinyint(2);comment:产品类型;"`  // 仓库类型
	Range enums.ProductInventoryRange `json:"range" gorm:"type:tinyint(2);comment:盘点范围;"` // 盘点范围

	Brand    []enums.ProductBrand    `json:"brand" gorm:"type:text;serializer:json;comment:产品品牌;"`    // 产品品牌
	Class    []enums.ProductClass    `json:"class" gorm:"type:text;serializer:json;comment:产品大类;"`    // 产品大类
	Category []enums.ProductCategory `json:"category" gorm:"type:text;serializer:json;comment:产品品类;"` // 产品品类
	Craft    []enums.ProductCraft    `json:"craft" gorm:"type:text;serializer:json;comment:产品工艺;"`    // 产品工艺
	Material []enums.ProductMaterial `json:"material" gorm:"type:text;serializer:json;comment:产品材质;"` // 产品材质
	Quality  []enums.ProductQuality  `json:"quality" gorm:"type:text;serializer:json;comment:产品成色;"`  // 产品成色
	Gem      []enums.ProductGem      `json:"gem" gorm:"type:text;serializer:json;comment:宝石种类;"`      // 宝石种类

	Remark string                       `json:"remark" gorm:"type:text;comment:备注;"`         // 备注
	Status enums.ProductInventoryStatus `json:"status" gorm:"type:tinyint(2);comment:盘点状态;"` // 盘点状态

	Products []ProductInventoryProduct `json:"products" gorm:"foreignKey:ProductInventoryId;references:Id;comment:盘点产品;"`

	CountShould      int64           `json:"count_should" gorm:"type:tinyint(5);comment:应盘数量;"`         // 应盘数量
	CountActual      int64           `json:"count_actual" gorm:"type:tinyint(5);comment:实盘数量;"`         // 实盘数量
	CountExtra       int64           `json:"count_extra" gorm:"type:tinyint(5);comment:盘盈数量;"`          // 盘盈数量
	CountLoss        int64           `json:"count_loss" gorm:"type:tinyint(5);comment:盘亏数量;"`           // 盘亏数量
	CountWeightMetal decimal.Decimal `json:"count_weight_metal" gorm:"type:decimal(10,2);comment:总重量;"` // 总重量
	CountPrice       decimal.Decimal `json:"count_price" gorm:"type:decimal(10,2);comment:总价值;"`        // 总价值
	ContQuantity     int64           `json:"cont_quantity" gorm:"type:tinyint(5);comment:总件数;"`         // 总件数
}

func (ProductInventory) WhereCondition(db *gorm.DB, req *types.ProductInventoryWhere) *gorm.DB {
	if req.Id != "" {
		db = db.Where("id = ?", req.Id)
	}
	if req.Type != 0 {
		db = db.Where("type = ?", req.Type)
	}
	if req.Status != 0 {
		db = db.Where("status = ?", req.Status)
	}
	if req.InventoryPersonId != "" {
		db = db.Where("inventory_person_id = ?", req.InventoryPersonId)
	}
	if req.InspectorId != "" {
		db = db.Where("inspector_id = ?", req.InspectorId)
	}
	if req.StartTime != nil && req.EndTime == nil {
		db = db.Where("created_at >= ?", req.StartTime)
	}
	if req.StartTime == nil && req.EndTime != nil {
		db = db.Where("created_at <= ?", req.EndTime)
	}
	if req.StartTime != nil && req.EndTime != nil {
		db = db.Where("created_at BETWEEN ? AND ?", req.StartTime, req.EndTime)
	}

	return db
}

func (ProductInventory) Preloads(db *gorm.DB, req *types.ProductInventoryWhere) *gorm.DB {
	db = db.Preload("Store")
	db = db.Preload("InventoryPerson")
	db = db.Preload("Inspector")
	db = db.Preload("Products", func(tx *gorm.DB) *gorm.DB {
		pdb := tx.Preload("Product")
		if req != nil && req.ProductStatus != enums.ProductInventoryProductStatusShould {
			pdb = pdb.Where(&ProductInventoryProduct{Status: req.ProductStatus})
		}
		return pdb
	})
	db = db.Preload("InventoryPerson", func(tx *gorm.DB) *gorm.DB {
		pdb := tx.Preload("Account", func(tx *gorm.DB) *gorm.DB {
			pdb := tx.Where(&Account{Platform: enums.PlatformTypeWxWork})
			return pdb
		})
		return pdb
	})
	db = db.Preload("Inspector", func(tx *gorm.DB) *gorm.DB {
		pdb := tx.Preload("Account", func(tx *gorm.DB) *gorm.DB {
			pdb := tx.Where(&Account{Platform: enums.PlatformTypeWxWork})
			return pdb
		})
		return pdb
	})

	return db
}

type ProductInventoryProduct struct {
	Model

	ProductInventoryId string           `json:"product_inventory_id" gorm:"type:varchar(255);not NULL;comment:盘点ID;"` // 盘点ID
	ProductInventory   ProductInventory `json:"-" gorm:"foreignKey:ProductInventoryId;references:Id;comment:盘点;"`

	ProductId string  `json:"product_id" gorm:"type:varchar(255);not NULL;comment:产品ID;"`    // 产品ID
	Product   Product `json:"product" gorm:"foreignKey:ProductId;references:Id;comment:产品;"` // 产品

	Status enums.ProductInventoryProductStatus `json:"status" gorm:"type:tinyint(2);comment:盘点状态;"` // 盘点状态

	InventoryTime *time.Time `json:"inventory_time" gorm:"type:datetime;comment:盘点时间;"` // 盘点时间
}

func init() {
	// 注册模型
	RegisterModels(
		&Product{},
		&ProductEnter{},
		&ProductDamage{},
		&ProductAllocate{},
		&ProductInventory{},
		&ProductInventoryProduct{},
		&ProductHistory{},
	)
	// 重置表
	RegisterRefreshModels(
	// &Product{},
	// &ProductEnter{},
	// &ProductDamage{},
	// &ProductAllocate{},
	// &ProductInventory{},
	// &ProductInventoryProduct{},
	// &ProductHistory{},
	)
}
