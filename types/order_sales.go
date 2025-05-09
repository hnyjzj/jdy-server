package types

import (
	"errors"
	"jdy/enums"
	"time"

	"github.com/shopspring/decimal"
)

type OrderSalesWhere struct {
	Id       string `json:"id" label:"订单编号" find:"true" sort:"1" type:"string" input:"text"`                                           // 订单编号
	StoreId  string `json:"store_id" label:"门店" find:"false" sort:"2" type:"string" input:"search" required:"true" binding:"required"` // 门店
	MemberId string `json:"member_id" label:"会员" find:"true" create:"true" sort:"3" type:"string" input:"search"`                      // 会员

	Status        enums.OrderSalesStatus   `json:"status" label:"订单状态" find:"true" sort:"4" type:"number" input:"select" preset:"typeMap"`                                     // 订单状态
	Source        enums.OrderSource        `json:"source" label:"订单来源" find:"true" create:"true" update:"true" sort:"6" type:"number" input:"select" preset:"typeMap"`         // 订单来源
	PaymentMethod enums.OrderPaymentMethod `json:"payment_method" label:"支付方式" find:"true" create:"true" update:"true" sort:"6" type:"number" input:"select" preset:"typeMap"` // 支付方式

	CashierId  string `json:"cashier_id" label:"收银员" find:"true" sort:"7" type:"string" input:"search"`  // 收银员
	SalesmanId string `json:"salesman_id" label:"导购员" find:"true" sort:"8" type:"string" input:"search"` // 导购员
	ProductId  string `json:"product_id" label:"商品" find:"true" sort:"9" type:"string" input:"search"`   // 商品

	StartDate *time.Time `json:"start_date" label:"开始日期" find:"true" sort:"10" type:"string" input:"date"` // 开始日期
	EndDate   *time.Time `json:"end_date" label:"结束日期" find:"true" sort:"11" type:"string" input:"date"`   // 结束日期
}

type OrderSalesCreateReq struct {
	StoreId string            `json:"store_id" binding:"required"` // 门店ID
	Source  enums.OrderSource `json:"source" binding:"required"`   // 订单来源

	CashierId string                 `json:"cashier_id" binding:"required"` // 收银员ID
	Clerks    []OrderCreateReqClerks `json:"clerks" binding:"required"`     // 导购员
	MemberId  string                 `json:"member_id" binding:"required"`  // 会员ID

	HasIntegral       bool            `json:"has_integral"`                          // 是否使用积分
	DiscountRate      decimal.Decimal `json:"discount_rate" binding:"required"`      // 整单折扣
	IntegralDeduction decimal.Decimal `json:"integral_deduction" binding:"required"` // 积分抵扣
	RoundOff          decimal.Decimal `json:"round_off" binding:"required"`          // 抹零

	ProductFinisheds   []OrderSalesCreateReqProductFinished   `json:"product_finisheds" binding:"required"`   // 成品
	ProductOlds        []OrderSalesCreateReqProductOld        `json:"product_olds" binding:"required"`        // 旧料
	ProductAccessories []OrderSalesCreateReqProductAccessorie `json:"product_accessories" binding:"required"` // 配件

	OrderDepositIds []string `json:"order_deposit_ids" binding:"required"` // 订金单ID

	Payments []OrderPaymentMethods `json:"payments" binding:"required"` // 支付方式

	Remark string `json:"remark"` // 备注
}

type OrderCreateReqClerks struct {
	SalesmanId      string          `json:"salesman_id" required:"true"`      // 导购员ID
	PerformanceRate decimal.Decimal `json:"performance_rate" required:"true"` // 绩效比例
	IsMain          bool            `json:"is_main" required:"true"`          // 是否主导购员
}

type OrderSalesCreateReqProductFinished struct {
	ProductId string `json:"product_id" binding:"required"` // 商品ID

	PriceGold         decimal.Decimal `json:"price_gold" binding:"required"`         // 金价
	LaborFee          decimal.Decimal `json:"labor_fee" binding:"required"`          // 工费
	DiscountFixed     decimal.Decimal `json:"discount_fixed" binding:"required"`     // 固定折扣
	IntegralDeduction decimal.Decimal `json:"integral_deduction" binding:"required"` // 积分抵扣
	DiscountMember    decimal.Decimal `json:"discount_member" binding:"required"`    // 会员折扣
	RoundOff          decimal.Decimal `json:"round_off" binding:"required"`          // 抹零
	Integral          decimal.Decimal `json:"integral" binding:"required"`           // 积分

	PriceOriginal decimal.Decimal `json:"price_original" binding:"required"` // 原价
	DiscountFinal decimal.Decimal `json:"discount_final" binding:"required"` // 折扣
	Price         decimal.Decimal `json:"price" binding:"required"`          // 应付金额
}

type OrderSalesCreateReqProductOld struct {
	ProductId string `json:"product_id"` // 商品ID

	IsOur                   bool                       `json:"is_our" binding:"required"`         // 是否本司货品
	RecycleMethod           enums.ProductRecycleMethod `json:"recycle_method" binding:"required"` // 回收方式
	RecycleType             enums.ProductRecycleType   `json:"recycle_type"`                      // 回收类型
	Code                    string                     `json:"code"`                              // 条码
	Material                enums.ProductMaterial      `json:"material" binding:"required"`       // 材质
	Quality                 enums.ProductQuality       `json:"quality" binding:"required"`        // 成色
	Gem                     enums.ProductGem           `json:"gem" binding:"required"`            // 主石
	Category                enums.ProductCategory      `json:"category"`                          // 品类
	Craft                   enums.ProductCraft         `json:"craft"`                             // 工艺
	WeightMetal             decimal.Decimal            `json:"weight_metal" binding:"required"`   // 金重
	LabelPrice              decimal.Decimal            `json:"label_price"`                       // 标签价
	RecyclePriceGold        decimal.Decimal            `json:"recycle_price_gold"`                // 回收金价
	RecyclePriceLabor       decimal.Decimal            `json:"recycle_price_labor"`               // 回收工费
	RecyclePriceLaborMethod enums.ProductRecycleMethod `json:"recycle_price_labor_method"`        // 回收工费方式
	Brand                   enums.ProductBrand         `json:"brand"`                             // 品牌
	WeightGem               decimal.Decimal            `json:"weight_gem"`                        // 主石重
	ColorGem                enums.ProductColor         `json:"color_gem"`                         // 主石颜色
	ClarityGem              enums.ProductClarity       `json:"clarity_gem"`                       // 主石净度
	Cut                     enums.ProductCut           `json:"cut"`                               // 主石切工
	NumGem                  int                        `json:"num_gem"`                           // 主石数量
	WeightOther             decimal.Decimal            `json:"weight_other"`                      // 杂料重
	NumOther                int                        `json:"num_other"`                         // 杂料数量
	WeightTotal             decimal.Decimal            `json:"weight_total"`                      // 总重
	QualityActual           decimal.Decimal            `json:"quality_actual" binding:"required"` // 实际成色
	Remark                  string                     `json:"remark"`                            // 备注
	Name                    string                     `json:"name"`                              // 名称
	RecyclePrice            decimal.Decimal            `json:"recycle_price"`                     // 回收金额
	Integral                decimal.Decimal            `json:"integral" binding:"required"`       // 积分
}

type OrderSalesCreateReqProductAccessorie struct {
	ProductId string `json:"product_id" binding:"required"` // 商品ID

	Quantity int64           `json:"quantity" binding:"required"` // 数量
	Price    decimal.Decimal `json:"price" binding:"required"`    // 价格
	Integral decimal.Decimal `json:"integral" binding:"required"` // 积分
}

func (req *OrderSalesCreateReq) Validate() error {
	if !req.DiscountRate.IsZero() {
		if req.DiscountRate.LessThan(decimal.NewFromFloat(0)) || req.DiscountRate.GreaterThan(decimal.NewFromFloat(100)) {
			return errors.New("整单折扣错误")
		}
	} else {
		req.DiscountRate = decimal.NewFromFloat(10)
	}

	// 检查导购员数量
	if len(req.Clerks) == 0 {
		return errors.New("导购员不能为空")
	}
	// 总佣金比例
	var totalPerformanceRate decimal.Decimal
	// 主导购数量
	var mainSalesmanCount int
	// 检查导购员
	for i, salesman := range req.Clerks {
		totalPerformanceRate = totalPerformanceRate.Add(salesman.PerformanceRate)
		if salesman.IsMain {
			mainSalesmanCount++
		}

		if !salesman.PerformanceRate.IsZero() {
			if salesman.PerformanceRate.LessThan(decimal.NewFromFloat(0)) || salesman.PerformanceRate.GreaterThan(decimal.NewFromFloat(100)) {
				return errors.New("佣金比例错误")
			}
		} else {
			req.Clerks[i].PerformanceRate = decimal.NewFromFloat(10)
		}
	}
	// 总佣金比例必须等于100
	if totalPerformanceRate.Cmp(decimal.NewFromFloat(100)) != 0 {
		return errors.New("总佣金比例必须等于100%")
	}
	// 主导购数量必须等于1
	if mainSalesmanCount != 1 {
		return errors.New("必须有且仅有一个主导购员")
	}

	// 检查商品
	for _, old := range req.ProductOlds {
		if old.IsOur && old.ProductId == "" {
			return errors.New("商品ID不能为空")
		}
	}

	// 检查支付方式
	if len(req.Payments) == 0 {
		return errors.New("支付方式不能为空")
	}
	for _, payment := range req.Payments {
		if payment.Amount.LessThan(decimal.NewFromFloat(0)) {
			return errors.New("支付金额错误")
		}
	}

	return nil
}

type OrderSalesListReq struct {
	PageReq
	Where OrderSalesWhere `json:"where"`
}

type OrderSalesInfoReq struct {
	Id string `json:"id" required:"true"`
}

type OrderSalesRevokedReq struct {
	Id string `json:"id" required:"true"`
}

type OrderSalesPayReq struct {
	Id string `json:"id" required:"true"` // 订单ID
}

type OrderSalesRefundReq struct {
	Id          string                `json:"id" required:"true"`           // 订单ID
	Method      enums.ProductTypeUsed `json:"method"`                       // 入库方式
	ProductType enums.ProductType     `json:"product_type" required:"true"` // 货品类型
	ProductId   string                `json:"product_id" required:"true"`   // 商品ID
	Price       decimal.Decimal       `json:"price" required:"true"`        // 退款金额
	Remark      string                `json:"remark" required:"true"`       // 备注

	Payments []OrderPaymentMethods `json:"payments" binding:"required"` // 支付方式
}

func (req *OrderSalesRefundReq) Validate() error {
	if req.ProductType == enums.ProductTypeFinished {
		if req.Method != 0 && req.Method.InMap() != nil {
			return errors.New("入库方式错误")
		}
	}

	// 检查支付方式
	if len(req.Payments) == 0 {
		return errors.New("支付方式不能为空")
	}
	for _, payment := range req.Payments {
		if payment.Amount.LessThan(decimal.NewFromFloat(0)) {
			return errors.New("支付金额错误")
		}
	}

	return nil
}
