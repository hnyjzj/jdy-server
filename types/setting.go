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
