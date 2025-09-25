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

	StoreId string `json:"store_id" gorm:"index;type:varchar(255);not NULL;comment:门店ID;"` // 门店ID
	Store   Store  `json:"store" gorm:"foreignKey:StoreId;references:Id;comment:门店;"`      // 门店

	InventoryPersonIds []string `json:"inventory_person_ids" gorm:"-"`                                              // 盘点人员ID
	InventoryPersons   []Staff  `json:"inventory_persons" gorm:"many2many:product_inventory_persons;comment:盘点人员;"` // 盘点人员

	InspectorId string `json:"inspector_id" gorm:"index;type:varchar(255);not NULL;comment:监盘人ID;"` // 监盘人ID
	Inspector   Staff  `json:"inspector" gorm:"foreignKey:InspectorId;references:Id;comment:监盘人;"`  // 监盘人

	CreatorId string `json:"creator_id" gorm:"type:varchar(255);not NULL;comment:创建人ID;"`    // 创建人ID
	Creator   Staff  `json:"creator" gorm:"foreignKey:CreatorId;references:Id;comment:创建人;"` // 创建人

	Type  enums.ProductTypeUsed       `json:"type" gorm:"index;type:int(11);comment:产品类型;"` // 仓库类型
	Range enums.ProductInventoryRange `json:"range" gorm:"type:int(11);comment:盘点范围;"`      // 盘点范围

	Brand         []enums.ProductBrand         `json:"brand" gorm:"type:text;serializer:json;comment:产品品牌;"`          // 产品品牌
	ClassFinished []enums.ProductClassFinished `json:"class_finished" gorm:"type:text;serializer:json;comment:成品大类;"` // 成品大类
	ClassOld      []enums.ProductClassOld      `json:"class_old" gorm:"type:text;serializer:json;comment:旧料大类;"`      // 旧料大类
	Category      []enums.ProductCategory      `json:"category" gorm:"type:text;serializer:json;comment:产品品类;"`       // 产品品类
	Craft         []enums.ProductCraft         `json:"craft" gorm:"type:text;serializer:json;comment:产品工艺;"`          // 产品工艺
	Material      []enums.ProductMaterial      `json:"material" gorm:"type:text;serializer:json;comment:产品材质;"`       // 产品材质
	Quality       []enums.ProductQuality       `json:"quality" gorm:"type:text;serializer:json;comment:产品成色;"`        // 产品成色
	Gem           []enums.ProductGem           `json:"gem" gorm:"type:text;serializer:json;comment:宝石种类;"`            // 宝石种类

	Remark string                       `json:"remark" gorm:"type:text;comment:备注;"`            // 备注
	Status enums.ProductInventoryStatus `json:"status" gorm:"index;type:int(11);comment:盘点状态;"` // 盘点状态

	ShouldCount    int64                     `json:"should_count" gorm:"type:int(11);comment:应盘数量;"`                                   // 应盘数量
	ShouldProducts []ProductInventoryProduct `json:"should_products" gorm:"foreignKey:ProductInventoryId;references:Id;comment:应盘产品;"` // 应盘产品
	ActualCount    int64                     `json:"actual_count" gorm:"type:int(11);comment:实盘数量;"`                                   // 实盘数量
	ActualProducts []ProductInventoryProduct `json:"actual_products" gorm:"foreignKey:ProductInventoryId;references:Id;comment:实盘产品;"` // 实盘产品
	ExtraCount     int64                     `json:"extra_count" gorm:"type:int(11);comment:盘盈数量;"`                                    // 盘盈数量
	ExtraProducts  []ProductInventoryProduct `json:"extra_products" gorm:"foreignKey:ProductInventoryId;references:Id;comment:盘盈产品;"`  // 盘盈产品
	LossCount      int64                     `json:"loss_count" gorm:"type:int(11);comment:盘亏数量;"`                                     // 盘亏数量
	LossProducts   []ProductInventoryProduct `json:"loss_products" gorm:"foreignKey:ProductInventoryId;references:Id;comment:盘亏产品;"`   // 盘亏产品

	CountWeightMetal decimal.Decimal `json:"count_weight_metal" gorm:"type:decimal(15,4);comment:总重量;"` // 总重量
	CountPrice       decimal.Decimal `json:"count_price" gorm:"type:decimal(10,2);comment:总价值;"`        // 总价值
	CountQuantity    int64           `json:"count_quantity" gorm:"type:int(11);comment:总件数;"`           // 总件数
}

// 产品盘点产品
type ProductInventoryProduct struct {
	Model

	ProductInventoryId string           `json:"product_inventory_id" gorm:"uniqueIndex:unique_product;type:varchar(255);not NULL;comment:盘点ID;"` // 盘点ID
	ProductInventory   ProductInventory `json:"-" gorm:"foreignKey:ProductInventoryId;references:Id;comment:盘点;"`

	ProductType     enums.ProductTypeUsed `json:"product_type" gorm:"type:int(11);not NULL;comment:产品类型;"`                                 // 产品类型
	ProductCode     string                `json:"product_code" gorm:"uniqueIndex:unique_product;type:varchar(255);not NULL;comment:产品编码;"` // 产品编码
	ProductFinished ProductFinished       `json:"product_finished" gorm:"foreignKey:ProductCode;references:Code;comment:成品"`               // 成品
	ProductOld      ProductOld            `json:"product_old"  gorm:"foreignKey:ProductCode;references:Code;comment:旧料"`                   // 旧料

	Status enums.ProductInventoryProductStatus `json:"status" gorm:"uniqueIndex:unique_product;type:int(11);comment:盘点状态;"` // 盘点状态

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
	if req.Range != 0 {
		db = db.Where("`range` = ?", req.Range)
	}
	if req.Status != 0 {
		db = db.Where("status = ?", req.Status)
	}
	if req.InspectorId != "" {
		db = db.Where("inspector_id = ?", req.InspectorId)
	}
	if req.StartTime != nil {
		db = db.Where("created_at >= ?", req.StartTime)
	}
	if req.EndTime != nil {
		db = db.Where("created_at <= ?", req.EndTime)
	}
	if req.InventoryPersonIds != nil {
		db = db.Where("id IN (SELECT product_inventory_id FROM product_inventory_persons WHERE staff_id IN (?))", req.InventoryPersonIds)
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
func (ProductInventory) Preloads(db *gorm.DB, req *types.ProductInventoryWhere, isOver bool) *gorm.DB {
	db = db.Preload("Store")
	db = db.Preload("InventoryPersons")
	db = db.Preload("Inspector")

	if isOver {
		// 应盘产品
		db = db.Preload("ShouldProducts", func(tx *gorm.DB) *gorm.DB {
			pdb := tx.Where(&ProductInventoryProduct{Status: enums.ProductInventoryProductStatusShould})
			pdb = pdb.Order("created_at desc")
			if req != nil && req.Page != 0 && req.Limit != 0 {
				pdb = PageCondition(pdb, &types.PageReq{Page: req.Page, Limit: req.Limit})
			}
			pdb = pdb.Preload("ProductFinished", func(finished *gorm.DB) *gorm.DB {
				finished = ProductFinished{}.Preloads(finished)
				return finished
			})
			pdb = pdb.Preload("ProductOld", func(old *gorm.DB) *gorm.DB {
				old = ProductOld{}.Preloads(old)
				return old
			})

			return pdb
		})
		// 盘盈产品
		db = db.Preload("ExtraProducts", func(tx *gorm.DB) *gorm.DB {
			pdb := tx.Where(&ProductInventoryProduct{Status: enums.ProductInventoryProductStatusExtra})
			pdb = pdb.Order("created_at desc")
			if req != nil && req.Page != 0 && req.Limit != 0 {
				pdb = PageCondition(pdb, &types.PageReq{Page: req.Page, Limit: req.Limit})
			}
			pdb = pdb.Preload("ProductFinished", func(finished *gorm.DB) *gorm.DB {
				finished = ProductFinished{}.Preloads(finished)
				return finished
			})
			pdb = pdb.Preload("ProductOld", func(old *gorm.DB) *gorm.DB {
				old = ProductOld{}.Preloads(old)
				return old
			})

			return pdb
		})
		// 盘亏产品
		db = db.Preload("LossProducts", func(tx *gorm.DB) *gorm.DB {
			pdb := tx.Where(&ProductInventoryProduct{Status: enums.ProductInventoryProductStatusLoss})
			pdb = pdb.Order("created_at desc")
			if req != nil && req.Page != 0 && req.Limit != 0 {
				pdb = PageCondition(pdb, &types.PageReq{Page: req.Page, Limit: req.Limit})
			}

			pdb = pdb.Preload("ProductFinished", func(finished *gorm.DB) *gorm.DB {
				finished = ProductFinished{}.Preloads(finished)
				return finished
			})
			pdb = pdb.Preload("ProductOld", func(old *gorm.DB) *gorm.DB {
				old = ProductOld{}.Preloads(old)
				return old
			})
			return pdb
		})
	}

	// 实盘产品
	if req != nil && req.Page != 0 && req.Limit != 0 {
		db = db.Preload("ActualProducts", func(tx *gorm.DB) *gorm.DB {
			pdb := tx.Where(&ProductInventoryProduct{Status: enums.ProductInventoryProductStatusActual})
			pdb = pdb.Order("created_at desc")
			pdb = PageCondition(pdb, &types.PageReq{Page: req.Page, Limit: req.Limit})
			pdb = pdb.Preload("ProductFinished", func(finished *gorm.DB) *gorm.DB {
				finished = ProductFinished{}.Preloads(finished)
				return finished
			})
			pdb = pdb.Preload("ProductOld", func(old *gorm.DB) *gorm.DB {
				old = ProductOld{}.Preloads(old)
				return old
			})

			return pdb
		})
	}

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
