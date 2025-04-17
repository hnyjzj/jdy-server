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

type OpenOrder struct {
	SoftDelete

	StoreId string `json:"store_id" gorm:"type:varchar(255);comment:店铺ID;"`                     // 店铺ID
	Store   Store  `json:"store,omitempty" gorm:"foreignKey:StoreId;references:Id;comment:店铺;"` // 店铺

	DiscountRate decimal.Decimal    `json:"discount_rate" gorm:"type:decimal(10,2);comment:积分抵扣比例;"` // 积分抵扣比例
	DecimalPoint enums.DecimalPoint `json:"decimal_point" gorm:"type:int(11);comment:金额小数点控制;"`      // 金额小数点控制
	Rounding     enums.Rounding     `json:"rounding" gorm:"type:int(11);comment:金额进位控制;"`            // 金额进位控制
	UseConfirm   bool               `json:"use_confirm" gorm:"type:tinyint(1);comment:积分使用二次确认;"`    // 积分使用二次确认
}

func (OpenOrder) Default() *OpenOrder {
	return &OpenOrder{
		DiscountRate: decimal.NewFromFloat(0),
		DecimalPoint: enums.DecimalPointNone,
		Rounding:     enums.RoundingRound,
		UseConfirm:   false,
	}
}

func init() {
	// 注册模型
	RegisterModels(
		&GoldPrice{},
		&OpenOrder{},
	)
	// 重置表
	RegisterRefreshModels(
	// &GoldPrice{},
	// &OpenOrder{},
	)
}
