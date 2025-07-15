package types

import (
	"errors"
	"jdy/enums"
	"time"

	"github.com/shopspring/decimal"
)

type OrderRepairWhere struct {
	Id       string `json:"id" label:"订单编号" find:"true" sort:"1" type:"string" input:"text"`                        // 订单编号
	StoreId  string `json:"store_id" label:"门店" find:"false" sort:"2" type:"string" input:"search" required:"true"` // 门店
	MemberId string `json:"member_id" label:"会员" find:"true" create:"true" sort:"3" type:"string" input:"search"`   // 会员

	Status        enums.OrderRepairStatus  `json:"status" label:"订单状态" find:"true" sort:"4" type:"number" input:"select" preset:"typeMap"`                                      // 订单状态
	PaymentMethod enums.OrderPaymentMethod `json:"payment_method" label:"支付方式" find:"false" create:"true" update:"true" sort:"5" type:"number" input:"select" preset:"typeMap"` // 支付方式

	ReceptionistId string `json:"receptionist_id" label:"接待人" create:"true" find:"true" sort:"6" type:"string" input:"search"` // 接待人
	CashierId      string `json:"cashier_id" label:"收银员" create:"true" find:"true" sort:"7" type:"string" input:"search"`      // 收银员

	Name string `json:"name" label:"维修项目" find:"true" create:"true" update:"true" sort:"7" input:"text" type:"string"`      // 维修项目
	Desc string `json:"desc" label:"问题描述" find:"false" create:"true" update:"true" sort:"8" input:"textarea" type:"string"` // 问题描述

	DeliveryMethod enums.DeliveryMethod `json:"delivery_method" label:"取货方式" find:"true" create:"true" update:"true" sort:"9" input:"select" type:"number" preset:"typeMap"`
	Province       string               `json:"province" label:"省" find:"true" create:"true" update:"true" sort:"12" input:"text" type:"string"` // 省份
	City           string               `json:"city" label:"市" find:"true" create:"true" update:"true" sort:"13" input:"text" type:"string"`     // 城市
	Area           string               `json:"area" label:"区" find:"true" create:"true" update:"true" sort:"14" input:"text" type:"string"`     // 区
	Address        string               `json:"address" label:"地址" find:"true" create:"true" update:"true" sort:"15" input:"text" type:"string"` // 地址

	Images string `json:"images" label:"图片" find:"false" create:"true" update:"true" sort:"16" input:"string[]" type:"string"` // 图片

	StartDate *time.Time `json:"start_date" label:"开始日期" find:"true" sort:"10" type:"string" input:"date"` // 开始日期
	EndDate   *time.Time `json:"end_date" label:"结束日期" find:"true" sort:"11" type:"string" input:"date"`   // 结束日期
}

type OrderRepairWhereProduct struct {
	IsOur       bool                  `json:"is_our" label:"是否本司货品" create:"true" sort:"1" input:"switch" type:"boolean"`                     // 是否本司货品
	Code        string                `json:"code" label:"条码" create:"true" sort:"2" input:"text" type:"string"`                              // 条码
	Material    enums.ProductMaterial `json:"material" label:"材质" create:"true" sort:"3" input:"select" type:"number" preset:"typeMap"`       // 材质
	Quality     enums.ProductQuality  `json:"quality" label:"成色" create:"true" sort:"4" input:"select" type:"number" preset:"typeMap"`        // 成色
	Gem         enums.ProductGem      `json:"gem" label:"主石" create:"true" sort:"5" input:"select" type:"number" preset:"typeMap"`            // 主石
	Category    enums.ProductCategory `json:"category" label:"品类" create:"true" sort:"6" input:"select" type:"number" preset:"typeMap"`       // 品类
	Craft       enums.ProductCraft    `json:"craft" label:"工艺" create:"true" sort:"7" input:"select" type:"number" preset:"typeMap"`          // 工艺
	WeightMetal decimal.Decimal       `json:"weight_metal" label:"金重" create:"true" sort:"8" input:"number" type:"string"`                    // 金重
	LabelPrice  decimal.Decimal       `json:"label_price" label:"标签价" create:"true" sort:"9" input:"number" type:"string"`                    // 标签价
	Brand       enums.ProductBrand    `json:"brand" label:"品牌" create:"true" sort:"10" input:"select" type:"number" preset:"typeMap"`         // 品牌
	WeightGem   decimal.Decimal       `json:"weight_gem" label:"主石重" create:"true" sort:"11" input:"number" type:"string"`                    // 主石重
	ColorGem    enums.ProductColor    `json:"color_gem" label:"主石颜色" create:"true" sort:"12" input:"select" type:"number" preset:"typeMap"`   // 主石颜色
	ClarityGem  enums.ProductClarity  `json:"clarity_gem" label:"主石净度" create:"true" sort:"13" input:"select" type:"number" preset:"typeMap"` // 主石净度
	Cut         enums.ProductCut      `json:"cut" label:"主石切工" create:"true" sort:"14" input:"select" type:"number" preset:"typeMap"`         // 主石切工
	WeightTotal decimal.Decimal       `json:"weight_total" label:"总重" create:"true" sort:"15" input:"number" type:"string"`                   // 总重
	Remark      string                `json:"remark" label:"备注" create:"true" sort:"16" input:"text" type:"string"`                           // 备注
	Name        string                `json:"name" label:"名称" create:"true" sort:"17" input:"text" type:"string"`                             // 名称
}

type OrderRepairCreateReq struct {
	StoreId string `json:"store_id" binding:"required"` // 门店ID

	ReceptionistId string `json:"receptionist_id" binding:"required"` // 接待人ID
	MemberId       string `json:"member_id" binding:"required"`       // 会员ID
	CashierId      string `json:"cashier_id" binding:"required"`      // 收银员ID

	Name string `json:"name" binding:"required"` // 维修项目
	Desc string `json:"desc" binding:"required"` // 问题描述

	DeliveryMethod enums.DeliveryMethod `json:"delivery_method" binding:"required"` // 取货方式
	Province       string               `json:"province"`                           // 省份
	City           string               `json:"city"`                               // 城市
	Area           string               `json:"area"`                               // 区
	Address        string               `json:"address"`                            // 地址

	Products []OrderRepairCreateReqProduct `json:"products" binding:"required"` // 商品
	Images   []string                      `json:"images" binding:"required"`   // 图片

	Expense decimal.Decimal `json:"expense" binding:"required"` // 费用
	Cost    decimal.Decimal `json:"cost" binding:"required"`    // 成本

	Payments []OrderPaymentMethods `json:"payments" binding:"required"` // 支付方式
}

func (req *OrderRepairCreateReq) Validate() error {
	// 检查商品
	for _, p := range req.Products {
		if p.IsOur && p.ProductId == "" {
			return errors.New("商品ID不能为空")
		}
	}
	if len(req.Products) == 0 {
		return errors.New("商品不能为空")
	}

	// 检查支付方式
	if len(req.Payments) == 0 {
		return errors.New("支付方式不能为空")
	}

	// 检查取货方式
	if req.DeliveryMethod == enums.DeliveryMethodMail {
		if req.Province == "" || req.City == "" || req.Area == "" || req.Address == "" {
			return errors.New("取货方式为邮寄时，省市区地址不能为空")
		}
	}

	return nil
}

type OrderRepairCreateReqProduct struct {
	IsOur     bool   `json:"is_our" binding:"required"` // 是否本司货品
	ProductId string `json:"product_id"`                // 商品ID

	Code        string                `json:"code"`                    // 条码
	Material    enums.ProductMaterial `json:"material"`                // 材质
	Quality     enums.ProductQuality  `json:"quality"`                 // 成色
	Gem         enums.ProductGem      `json:"gem"`                     // 主石
	Category    enums.ProductCategory `json:"category"`                // 品类
	Craft       enums.ProductCraft    `json:"craft"`                   // 工艺
	WeightMetal decimal.Decimal       `json:"weight_metal"`            // 金重
	LabelPrice  decimal.Decimal       `json:"label_price"`             // 标签价
	Brand       enums.ProductBrand    `json:"brand"`                   // 品牌
	WeightGem   decimal.Decimal       `json:"weight_gem"`              // 主石重
	ColorGem    enums.ProductColor    `json:"color_gem"`               // 主石颜色
	ClarityGem  enums.ProductClarity  `json:"clarity_gem"`             // 主石净度
	Cut         enums.ProductCut      `json:"cut"`                     // 主石切工
	WeightTotal decimal.Decimal       `json:"weight_total"`            // 总重
	Remark      string                `json:"remark"`                  // 备注
	Name        string                `json:"name" binding:"required"` // 名称
}

type OrderRepairListReq struct {
	PageReq
	Where OrderRepairWhere `json:"where"`
}

type OrderRepairInfoReq struct {
	Id string `json:"id" binding:"required"`
}

type OrderRepairUpdateReq struct {
	Id string `json:"id" binding:"required"` // 订单ID

	Name string `json:"name"` // 维修项目
	Desc string `json:"desc"` // 问题描述

	DeliveryMethod enums.DeliveryMethod `json:"delivery_method"` // 取货方式
	Province       string               `json:"province"`        // 省份
	City           string               `json:"city"`            // 城市
	Area           string               `json:"area"`            // 区
	Address        string               `json:"address"`         // 地址

	Images []string `json:"images"` // 图片
}

type OrderRepairOperationReq struct {
	Id        string                  `json:"id" binding:"required"`
	Operation enums.OrderRepairStatus `json:"operation" binding:"required"`
}

type OrderRepairRevokedReq struct {
	Id string `json:"id" binding:"required"`
}

type OrderRepairPayReq struct {
	Id string `json:"id" binding:"required"`
}

type OrderRepairRefundReq struct {
	Id     string `json:"id" binding:"required"`
	Remark string `json:"remark"`
}
