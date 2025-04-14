package model

import (
	"jdy/enums"
	"jdy/types"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// 产品盘点单
type ProductInventory struct {
	SoftDelete

	StoreId string `json:"store_id" gorm:"type:varchar(255);not NULL;comment:门店ID;"`  // 门店ID
	Store   Store  `json:"store" gorm:"foreignKey:StoreId;references:Id;comment:门店;"` // 门店

	InventoryPersonId string `json:"inventory_person_id" gorm:"type:varchar(255);not NULL;comment:盘点人ID;"`            // 盘点人ID
	InventoryPerson   Staff  `json:"inventory_person" gorm:"foreignKey:InventoryPersonId;references:Id;comment:盘点人;"` // 盘点人

	InspectorId string `json:"inspector_id" gorm:"type:varchar(255);not NULL;comment:监盘人ID;"`      // 监盘人ID
	Inspector   Staff  `json:"inspector" gorm:"foreignKey:InspectorId;references:Id;comment:监盘人;"` // 监盘人

	CreatorId string `json:"creator_id" gorm:"type:varchar(255);not NULL;comment:创建人ID;"`    // 创建人ID
	Creator   Staff  `json:"creator" gorm:"foreignKey:CreatorId;references:Id;comment:创建人;"` // 创建人

	Type  enums.ProductType           `json:"type" gorm:"type:tinyint(2);comment:产品类型;"`  // 仓库类型
	Range enums.ProductInventoryRange `json:"range" gorm:"type:tinyint(2);comment:盘点范围;"` // 盘点范围

	Brand         []enums.ProductBrand         `json:"brand" gorm:"type:text;serializer:json;comment:产品品牌;"`          // 产品品牌
	ClassFinished []enums.ProductClassFinished `json:"class_finished" gorm:"type:text;serializer:json;comment:成品大类;"` // 成品大类
	ClassOld      []enums.ProductClassOld      `json:"class_old" gorm:"type:text;serializer:json;comment:旧料大类;"`      // 旧料大类
	Category      []enums.ProductCategory      `json:"category" gorm:"type:text;serializer:json;comment:产品品类;"`       // 产品品类
	Craft         []enums.ProductCraft         `json:"craft" gorm:"type:text;serializer:json;comment:产品工艺;"`          // 产品工艺
	Material      []enums.ProductMaterial      `json:"material" gorm:"type:text;serializer:json;comment:产品材质;"`       // 产品材质
	Quality       []enums.ProductQuality       `json:"quality" gorm:"type:text;serializer:json;comment:产品成色;"`        // 产品成色
	Gem           []enums.ProductGem           `json:"gem" gorm:"type:text;serializer:json;comment:宝石种类;"`            // 宝石种类

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

// 产品盘点产品
type ProductInventoryProduct struct {
	Model

	ProductInventoryId string           `json:"product_inventory_id" gorm:"type:varchar(255);not NULL;comment:盘点ID;"` // 盘点ID
	ProductInventory   ProductInventory `json:"-" gorm:"foreignKey:ProductInventoryId;references:Id;comment:盘点;"`

	ProductId       string            `json:"product_id" gorm:"type:varchar(255);not NULL;comment:产品ID;"`            // 产品ID
	ProductType     enums.ProductType `json:"product_type" gorm:"type:tinyint(2);not NULL;comment:产品类型;"`            // 产品类型
	ProductFinished ProductFinished   `json:"product_finished" gorm:"foreignKey:ProductId;references:Id;comment:成品"` // 成品
	ProductOld      ProductOld        `json:"product_old"  gorm:"foreignKey:ProductId;references:Id;comment:旧料"`     // 旧料

	Status enums.ProductInventoryProductStatus `json:"status" gorm:"type:tinyint(2);comment:盘点状态;"` // 盘点状态

	InventoryTime *time.Time `json:"inventory_time" gorm:"type:datetime;comment:盘点时间;"` // 盘点时间
}

func (ProductInventory) WhereCondition(db *gorm.DB, req *types.ProductInventoryWhere) *gorm.DB {
	if req.Id != "" {
		db = db.Where("id = ?", req.Id)
	}
	if req.StoreId != "" {
		db = db.Where("store_id = ?", req.StoreId)
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

// 产品盘点创建条件
func CreateProductInventoryCondition(db *gorm.DB, req *types.ProductInventoryCreateReq) *gorm.DB {
	if len(req.Brand) > 0 {
		db = db.Where("brand in (?)", req.Brand)
	}
	if len(req.ClassFinished) > 0 {
		db = db.Where("class in (?)", req.ClassFinished)
	}
	if len(req.ClassOld) > 0 {
		db = db.Where("class in (?)", req.ClassOld)
	}
	if len(req.Category) > 0 {
		db = db.Where("category in (?)", req.Category)
	}
	if len(req.Craft) > 0 {
		db = db.Where("craft in (?)", req.Craft)
	}
	if len(req.Material) > 0 {
		db = db.Where("material in (?)", req.Material)
	}
	if len(req.Quality) > 0 {
		db = db.Where("quality in (?)", req.Quality)
	}
	if len(req.Gem) > 0 {
		db = db.Where("gem in (?)", req.Gem)
	}

	return db
}

// 产品盘点关联条件
func (ProductInventory) Preloads(db *gorm.DB, req *types.ProductInventoryWhere) *gorm.DB {
	db = db.Preload("Store")
	db = db.Preload("InventoryPerson")
	db = db.Preload("Inspector")
	db = db.Preload("Products", func(tx *gorm.DB) *gorm.DB {
		pdb := tx
		pdb = pdb.Preload("ProductFinished")
		pdb = pdb.Preload("ProductOld")
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

func init() {
	// 注册模型
	RegisterModels(
		&ProductInventory{},
		&ProductInventoryProduct{},
	)
	// 重置表
	RegisterRefreshModels(
	// &ProductInventory{},
	// &ProductInventoryProduct{},
	)
}
