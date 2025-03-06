package model

import (
	"jdy/enums"
	"jdy/types"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type GoldPrice struct {
	SoftDelete

	StoreId         string                 `json:"store_id" gorm:"type:varchar(255);comment:店铺ID;"`                // 店铺ID
	Price           decimal.Decimal        `json:"price" gorm:"type:decimal(10,2);comment:金价;"`                    // 金价
	ProductMaterial enums.ProductMaterial  `json:"product_material" gorm:"type:tinyint(1);comment:产品材质;"`          // 产品材质
	ProductType     enums.ProductType      `json:"product_type" gorm:"type:tinyint(1);comment:产品类型;"`              // 产品类型
	ProductBrand    []enums.ProductBrand   `json:"product_brand" gorm:"type:text;serializer:json;comment:产品品牌;"`   // 产品品牌
	ProductQuality  []enums.ProductQuality `json:"product_quality" gorm:"type:text;serializer:json;comment:产品成色;"` // 产品成色
}

func (GoldPrice) WhereCondition(db *gorm.DB, req *types.GoldPriceOptions) *gorm.DB {
	if req.StoreId != "" {
		db = db.Where("store_id = ?", req.StoreId)
	}
	if req.ProductMaterial != 0 {
		db = db.Where("product_material = ?", req.ProductMaterial)
	}
	if req.ProductType != 0 {
		db = db.Where("product_type = ?", req.ProductType)
	}
	if len(req.ProductBrand) > 0 {
		db = db.Where("product_brand in (?)", req.ProductBrand)
	}
	if len(req.ProductQuality) > 0 {
		db = db.Where("product_quality in (?)", req.ProductQuality)
	}

	return db
}

func GetGoldPrice(req *types.GoldPriceOptions) (decimal.Decimal, error) {
	var goldPrice GoldPrice

	db := DB.Order("updated_at desc")
	db = goldPrice.WhereCondition(db, req)
	if err := db.First(&goldPrice).Error; err != nil {
		return decimal.Zero, err
	}

	return goldPrice.Price, nil
}

func init() {
	// 注册模型
	RegisterModels(
		&GoldPrice{},
	)
	// 重置表
	RegisterRefreshModels(
	// &GoldPrice{},
	)
}
