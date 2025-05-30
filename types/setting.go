package types

import (
	"errors"
	"jdy/enums"

	"github.com/shopspring/decimal"
)

type GoldPriceListReq struct {
	StoreId string `json:"store_id" binding:"required"` // 门店ID
}

type GoldPriceOptions struct {
	Id              string                 `json:"id"`                                  // 金价ID
	StoreId         string                 `json:"store_id" binding:"required"`         // 门店ID
	Price           decimal.Decimal        `json:"price" binding:"required"`            // 金价
	ProductMaterial enums.ProductMaterial  `json:"product_material" binding:"required"` // 产品材质
	ProductType     enums.ProductType      `json:"product_type" binding:"required"`     // 产品类型
	ProductBrand    []enums.ProductBrand   `json:"product_brand"`                       // 产品品牌
	ProductQuality  []enums.ProductQuality `json:"product_quality" binding:"required"`  // 产品成色
}

type GoldPriceCreateReq struct {
	Options []GoldPriceOptions `json:"options" binding:"required"` // 金价创建参数
	Deletes []string           `json:"deletes"`                    // 金价删除参数
}

func (q *GoldPriceCreateReq) Validate() error {
	for _, v := range q.Options {
		if v.Price.LessThanOrEqual(decimal.Zero) {
			return errors.New("金价不能小于等于0")
		}

		if err := v.ProductMaterial.InMap(); err != nil {
			return err
		}

		if err := v.ProductType.InMap(); err != nil {
			return err
		}

		for _, v := range v.ProductBrand {
			if err := v.InMap(); err != nil {
				return err
			}
		}

		for _, v := range v.ProductQuality {
			if err := v.InMap(); err != nil {
				return err
			}
		}

	}

	return nil
}

type OpenOrderWhere struct {
	StoreId      string             `json:"store_id" label:"门店" input:"text" type:"string" find:"false" sort:"1" required:"true"`                                           // 门店
	DiscountRate decimal.Decimal    `json:"discount_rate" label:"积分抵扣比例" input:"text" type:"decimal" create:"true" find:"true" sort:"2" required:"true"`                    // 积分抵扣比例
	DecimalPoint enums.DecimalPoint `json:"decimal_point" label:"金额小数点控制" input:"select" type:"number" create:"true" find:"true" sort:"3" required:"true" preset:"typeMap"` // 金额小数点控制
	Rounding     enums.Rounding     `json:"rounding" label:"金额进位控制" input:"select" type:"number" create:"true" find:"true" sort:"4" required:"true" preset:"typeMap"`       // 金额进位控制
	UseConfirm   bool               `json:"use_confirm" label:"积分使用二次确认" input:"switch" type:"boolean" create:"true" find:"true" sort:"5" required:"true"`                  // 积分使用二次确认
}

type OpenOrderInfoReq struct {
	StoreId string `json:"store_id" binding:"required"` // 门店ID
}

type OpenOrderUpdateReq struct {
	StoreId      string             `json:"store_id" binding:"required"` // 门店ID
	DiscountRate *decimal.Decimal   `json:"discount_rate"`               // 积分抵扣比例
	DecimalPoint enums.DecimalPoint `json:"decimal_point"`               // 金额小数点控制
	Rounding     enums.Rounding     `json:"rounding"`                    // 金额进位控制
	UseConfirm   bool               `json:"use_confirm"`                 // 积分使用二次确认
}

func (req *OpenOrderUpdateReq) Validate() error {
	// 验证折扣率是否有效（如果提供了）
	if req.DiscountRate != nil && req.DiscountRate.LessThan(decimal.Zero) {
		return errors.New("积分抵扣比例不能为负数")
	}

	// 验证小数点控制枚举值是否有效
	if err := req.DecimalPoint.InMap(); err != nil {
		return errors.New("金额小数点控制枚举值无效")
	}

	// 验证进位控制枚举值是否有效
	if err := req.Rounding.InMap(); err != nil {
		return errors.New("金额进位控制枚举值无效")
	}

	return nil
}
