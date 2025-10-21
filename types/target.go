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

	StartTime *time.Time `json:"start_time" label:"开始时间" sort:"4" find:"true" create:"true" update:"true" list:"true" info:"true" input:"date" type:"date" required:"true"` // 开始时间
	EndTime   *time.Time `json:"end_time" label:"结束时间" sort:"5" find:"true" create:"true" update:"true" list:"true" info:"true" input:"date" type:"date" required:"true"`   // 结束时间

	Method enums.TargetMethod `json:"method" label:"统计方式" sort:"6" find:"true" create:"true" update:"true" list:"true" info:"true" input:"radio" type:"number" required:"true" preset:"typeMap"` // 统计方式

	Scope    enums.TargetScope          `json:"scope" label:"统计范围" sort:"7" find:"true" create:"true" update:"true" list:"false" info:"true" input:"radio" type:"number" required:"true" preset:"typeMap"`                                                                         // 统计范围
	Class    enums.ProductClassFinished `json:"class" label:"产品大类" sort:"8" find:"true" create:"true" update:"true" list:"false" info:"true" input:"multiple" type:"number" required:"true" preset:"typeMap" condition:"[{\"key\":\"scope\",\"operator\":\"=\",\"value\":1}]"`     // 产品大类
	Material enums.ProductMaterial      `json:"material" label:"产品材质" sort:"9" find:"true" create:"true" update:"true" list:"false" info:"true" input:"multiple" type:"number" required:"true" preset:"typeMap" condition:"[{\"key\":\"scope\",\"operator\":\"=\",\"value\":2}]"`  // 产品材质
	Quality  enums.ProductQuality       `json:"quality" label:"产品成色" sort:"10" find:"true" create:"true" update:"true" list:"false" info:"true" input:"multiple" type:"number" required:"true" preset:"typeMap" condition:"[{\"key\":\"scope\",\"operator\":\"=\",\"value\":2}]"`  // 产品成色
	Category enums.ProductCategory      `json:"category" label:"产品品类" sort:"11" find:"true" create:"true" update:"true" list:"false" info:"true" input:"multiple" type:"number" required:"true" preset:"typeMap" condition:"[{\"key\":\"scope\",\"operator\":\"=\",\"value\":2}]"` // 产品品类
	Gem      enums.ProductGem           `json:"gem" label:"产品主石" sort:"12" find:"true" create:"true" update:"true" list:"false" info:"true" input:"multiple" type:"number" required:"true" preset:"typeMap" condition:"[{\"key\":\"scope\",\"operator\":\"=\",\"value\":2}]"`      // 产品主石
	Craft    enums.ProductCraft         `json:"craft" label:"产品工艺" sort:"13" find:"true" create:"true" update:"true" list:"false" info:"true" input:"multiple" type:"number" required:"true" preset:"typeMap" condition:"[{\"key\":\"scope\",\"operator\":\"=\",\"value\":2}]"`    // 产品工艺

	Object enums.TargetObject `json:"object" label:"统计对象" sort:"14" find:"true" create:"true" update:"true" list:"true" info:"true" input:"radio" type:"number" required:"true" preset:"typeMap"` // 统计对象
}

type TargetWhereGroup struct {
	Id   string `json:"id" label:"编号" sort:"1" find:"false" create:"false" update:"false" list:"true" info:"true" input:"text" type:"string" required:"false"` // 编号
	Name string `json:"name" label:"名称" sort:"2" find:"true" create:"true" update:"true" list:"true" info:"true" input:"text" type:"string" required:"true"`   // 名称
}

type TargetWherePersonal struct {
	StaffId  string           `json:"staff_id" label:"员工编号" sort:"1" find:"true" create:"true" update:"true" list:"true" info:"false" input:"search" type:"string" required:"true"`   // 员工编号
	GroupId  string           `json:"group_id" label:"组别编号" sort:"2" find:"false" create:"false" update:"true" list:"true" info:"false" input:"search" type:"string" required:"true"` // 组别编号
	IsLeader bool             `json:"is_leader" label:"是否组长" sort:"3" find:"true" create:"true" update:"true" list:"true" info:"false" input:"switch" type:"boolean" required:"true"` // 是否组长
	Purpose  *decimal.Decimal `json:"purpose" label:"目标" sort:"4" find:"true" create:"true" update:"true" list:"true" info:"false" input:"text" type:"number" required:"true"`        // 目标
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

	Groups    []TargetWhereGroup    `json:"groups"`    // 组别
	Personals []TargetWherePersonal `json:"personals"` // 员工
}

func (req *TargetCreateReq) Validate() error { // 验证时间范围
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
	switch req.Object {
	case enums.TargetObjectGroup:
		{
			if len(req.Groups) == 0 {
				return errors.New("组别不能为空")
			}
			if len(req.Personals) == 0 {
				return errors.New("员工不能为空")
			} else {
				for _, personal := range req.Personals {
					if personal.GroupId == "" {
						return errors.New("员工组别不能为空")
					}
					if personal.StaffId == "" {
						return errors.New("员工编号不能为空")
					}
				}
			}
		}
	case enums.TargetObjectPersonal:
		{
			if len(req.Personals) == 0 {
				return errors.New("员工不能为空")
			} else {
				for _, personal := range req.Personals {
					if personal.StaffId == "" {
						return errors.New("员工编号不能为空")
					}
				}
			}
		}
	}

	// 验证组是否重复
	groupNames := make(map[string]bool)
	for _, group := range req.Groups {
		if group.Name == "" {
			return errors.New("组别名称不能为空")
		}
		if _, ok := groupNames[group.Name]; ok {
			return errors.New("组别编号不能重复")
		}
		groupNames[group.Name] = true
	}

	// 验证员工是否重复
	staffIds := make(map[string]bool)
	// 验证组长是否重复
	isLeaders := make(map[string]map[string]bool)
	for _, personal := range req.Personals {
		if personal.StaffId == "" {
			return errors.New("员工编号不能为空")
		}
		if _, ok := staffIds[personal.StaffId]; ok {
			return errors.New("员工不能重复")
		}
		staffIds[personal.StaffId] = true

		if _, ok := isLeaders[personal.GroupId]; !ok {
			isLeaders[personal.GroupId] = make(map[string]bool)
		}
		if _, ok := isLeaders[personal.GroupId][personal.StaffId]; ok {
			return errors.New("组长不能重复")
		} else if personal.IsLeader {
			isLeaders[personal.GroupId][personal.StaffId] = true
		}
	}

	return nil
}

type TargetListReq struct {
	PageReq

	Where TargetWhere `json:"where" binding:"required"`
}
