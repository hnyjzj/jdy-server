package types

import (
	"errors"
	"jdy/enums"
	"time"

	"github.com/shopspring/decimal"
)

type TargetWhere struct {
	StoreId string `json:"store_id" label:"门店" sort:"1" find:"false" create:"false" update:"false" list:"true" info:"true" input:"text" type:"string" required:"true"` // 门店

	Id   string `json:"id" label:"编号" sort:"1" find:"true" create:"false" update:"false" list:"true" info:"true" input:"text" type:"string" required:"false"` // 编号
	Name string `json:"name" label:"名称" sort:"2" find:"true" create:"true" update:"true" list:"true" info:"true" input:"text" type:"string" required:"true"`  // 名称

	IsDefault *bool `json:"is_default" label:"是否默认" sort:"3" find:"true" create:"true" update:"true" list:"true" info:"true" input:"switch" type:"boolean" required:"true"` // 是否默认

	StartTime *time.Time `json:"start_time" label:"开始时间" sort:"4" find:"true" create:"true" update:"true" list:"true" info:"true" input:"date" type:"date" required:"true"`   // 开始时间
	EndTime   *time.Time `json:"end_time" label:"结束时间" sort:"5" find:"true" create:"true" update:"true" list:"true" info:"true" input:"datetime" type:"date" required:"true"` // 结束时间

	Method enums.TargetMethod `json:"method" label:"统计方式" sort:"6" find:"true" create:"true" update:"false" list:"true" info:"true" input:"radio" type:"number" required:"true" preset:"typeMap"` // 统计方式

	Scope    enums.TargetScope          `json:"scope" label:"统计范围" sort:"7" find:"true" create:"true" update:"false" list:"false" info:"true" input:"radio" type:"number" required:"true" preset:"typeMap"`                                                                         // 统计范围
	Class    enums.ProductClassFinished `json:"class" label:"产品大类" sort:"8" find:"true" create:"true" update:"false" list:"false" info:"true" input:"multiple" type:"number" required:"true" preset:"typeMap" condition:"[{\"key\":\"scope\",\"operator\":\"=\",\"value\":1}]"`     // 产品大类
	Material enums.ProductMaterial      `json:"material" label:"产品材质" sort:"9" find:"true" create:"true" update:"false" list:"false" info:"true" input:"multiple" type:"number" required:"true" preset:"typeMap" condition:"[{\"key\":\"scope\",\"operator\":\"=\",\"value\":2}]"`  // 产品材质
	Quality  enums.ProductQuality       `json:"quality" label:"产品成色" sort:"10" find:"true" create:"true" update:"false" list:"false" info:"true" input:"multiple" type:"number" required:"true" preset:"typeMap" condition:"[{\"key\":\"scope\",\"operator\":\"=\",\"value\":2}]"`  // 产品成色
	Category enums.ProductCategory      `json:"category" label:"产品品类" sort:"11" find:"true" create:"true" update:"false" list:"false" info:"true" input:"multiple" type:"number" required:"true" preset:"typeMap" condition:"[{\"key\":\"scope\",\"operator\":\"=\",\"value\":2}]"` // 产品品类
	Gem      enums.ProductGem           `json:"gem" label:"产品主石" sort:"12" find:"true" create:"true" update:"false" list:"false" info:"true" input:"multiple" type:"number" required:"true" preset:"typeMap" condition:"[{\"key\":\"scope\",\"operator\":\"=\",\"value\":2}]"`      // 产品主石
	Craft    enums.ProductCraft         `json:"craft" label:"产品工艺" sort:"13" find:"true" create:"true" update:"false" list:"false" info:"true" input:"multiple" type:"number" required:"true" preset:"typeMap" condition:"[{\"key\":\"scope\",\"operator\":\"=\",\"value\":2}]"`    // 产品工艺

	Object enums.TargetObject `json:"object" label:"统计对象" sort:"14" find:"true" create:"true" update:"false" list:"true" info:"true" input:"radio" type:"number" required:"true" preset:"typeMap"` // 统计对象
}

type TargetWhereGroup struct {
	Id   string `json:"id" label:"编号" sort:"1" find:"false" create:"false" update:"false" list:"true" info:"true" input:"text" type:"string" required:"false"` // 编号
	Name string `json:"name" label:"名称" sort:"2" find:"true" create:"true" update:"true" list:"true" info:"true" input:"text" type:"string" required:"true"`   // 名称
}

type TargetWherePersonal struct {
	Id       string           `json:"id" label:"编号" sort:"1" find:"false" create:"false" update:"false" list:"false" info:"true" input:"text" type:"string" required:"false"`       // 编号
	StaffId  string           `json:"staff_id" label:"员工" sort:"2" find:"true" create:"true" update:"true" list:"true" info:"false" input:"search" type:"string" required:"true"`   // 员工
	GroupId  string           `json:"group_id" label:"组别" sort:"3" find:"false" create:"false" update:"true" list:"true" info:"false" input:"search" type:"string" required:"true"` // 组别
	IsLeader bool             `json:"is_leader" label:"组长" sort:"4" find:"true" create:"true" update:"true" list:"true" info:"false" input:"switch" type:"boolean" required:"true"` // 是否组长
	Purpose  *decimal.Decimal `json:"purpose" label:"目标" sort:"5" find:"true" create:"true" update:"true" list:"true" info:"false" input:"text" type:"number" required:"true"`      // 目标
}

type TargetCreateReq struct {
	StoreId   string             `json:"store_id" binding:"required"`   // 店铺编号
	Name      string             `json:"name" binding:"required"`       // 名称
	IsDefault bool               `json:"is_default"`                    // 是否默认
	StartTime *time.Time         `json:"start_time" binding:"required"` // 开始时间
	EndTime   *time.Time         `json:"end_time" binding:"required"`   // 结束时间
	Method    enums.TargetMethod `json:"method" binding:"required"`     // 统计方式
	Scope     enums.TargetScope  `json:"scope" binding:"required"`      // 统计范围
	Object    enums.TargetObject `json:"object"`                        // 统计对象

	Class    []enums.ProductClassFinished `json:"class"`    // 产品大类
	Material []enums.ProductMaterial      `json:"material"` // 产品材质
	Quality  []enums.ProductQuality       `json:"quality"`  // 产品成色
	Category []enums.ProductCategory      `json:"category"` // 产品品类
	Gem      []enums.ProductGem           `json:"gem"`      // 产品主石
	Craft    []enums.ProductCraft         `json:"craft"`    // 产品工艺
}

func (req *TargetCreateReq) Validate() error {
	// 验证时间范围
	if req.StartTime != nil && req.EndTime != nil {
		if req.EndTime.Before(*req.StartTime) {
			return errors.New("结束时间必须晚于开始时间")
		}
	}

	switch req.Scope {
	case enums.TargetScopeClass:
		{
			if len(req.Class) == 0 {
				return errors.New("产品大类不能为空")
			}
		}
	case enums.TargetScopeOther:
		{
			if len(req.Material) == 0 {
				return errors.New("产品材质不能为空")
			}
			if len(req.Quality) == 0 {
				return errors.New("产品成色不能为空")
			}
			if len(req.Category) == 0 {
				return errors.New("产品品类不能为空")
			}
			if len(req.Gem) == 0 {
				return errors.New("产品主石不能为空")
			}
			if len(req.Craft) == 0 {
				return errors.New("产品工艺不能为空")
			}
		}
	}

	return nil
}

type TargetGroupCreateReq struct {
	TargetId string `json:"target_id" binding:"required"` // 目标编号

	Name string `json:"name" binding:"required"` // 名称
}

type TargetPersonalCreateReq struct {
	TargetId string `json:"target_id" binding:"required"` // 目标编号

	StaffId  string           `json:"staff_id" binding:"required"` // 员工编号
	GroupId  string           `json:"group_id"`                    // 组别编号
	IsLeader bool             `json:"is_leader"`                   // 是否组长
	Purpose  *decimal.Decimal `json:"purpose" binding:"required"`  // 目标
}

type TargetListReq struct {
	PageReq

	Where TargetWhere `json:"where" binding:"required"`
}

type TargetInfoReq struct {
	Id string `json:"id" binding:"required"` // 目标编号
}

type TargetUpdateReq struct {
	Id string `json:"id" binding:"required"` // 目标编号

	Name      string     `json:"name"`       // 名称
	IsDefault bool       `json:"is_default"` // 是否默认
	StartTime *time.Time `json:"start_time"` // 开始时间
	EndTime   *time.Time `json:"end_time"`   // 结束时间

	Groups    []TargetWhereGroup    `json:"groups"`    // 组别
	Personals []TargetWherePersonal `json:"personals"` // 员工
}

type TargetUpdateGroupReq struct {
	Id string `json:"id" binding:"required"` // 目标编号

	Name string `json:"name"` // 名称
}

type TargetUpdatePersonalReq struct {
	Id string `json:"id" binding:"required"` // 目标编号

	IsLeader bool             `json:"is_leader"`                  // 是否组长
	Purpose  *decimal.Decimal `json:"purpose" binding:"required"` // 目标
}

type TargetDeleteReq struct {
	Id string `json:"id" binding:"required"` // 目标编号
}
type TargetDeleteGroupReq struct {
	Id string `json:"id" binding:"required"` // 目标编号
}
type TargetDeletePersonalReq struct {
	Id string `json:"id" binding:"required"` // 目标编号
}
