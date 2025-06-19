package setting

import (
	"errors"
	"jdy/logic"
	"jdy/model"
	"jdy/types"
	"jdy/utils"
	"log"
)

type PrintLogic struct {
	logic.BaseLogic
}

func (PrintLogic) Create(req *types.PrintReq) error {
	print := model.Print{
		StoreId: req.StoreId,
		Name:    req.Name,
		Type:    req.Type,
		Config:  req.Config,
	}

	if err := model.DB.Create(&print).Error; err != nil {
		log.Printf("创建失败: %v", err)
		return errors.New("创建失败")
	}

	return nil
}

func (PrintLogic) List(req *types.PrintListReq) (*types.PageRes[model.Print], error) {
	var (
		print model.Print

		res types.PageRes[model.Print]
	)

	db := model.DB.Model(&print)
	db = print.WhereCondition(db, &req.Where)

	// 获取总数
	if err := db.Count(&res.Total).Error; err != nil {
		return nil, errors.New("获取总数失败")
	}

	// 获取列表
	db = db.Order("created_at desc")
	db = model.PageCondition(db, req.Page, req.Limit)
	if err := db.Find(&res.List).Error; err != nil {
		return nil, errors.New("获取列表失败")
	}

	res.List = append([]model.Print{model.Print{}.Default(req.Where.Type)}, res.List...)
	res.Total++

	return &res, nil
}

func (PrintLogic) Info(req *types.PrintInfoReq) (*model.Print, error) {
	var (
		print model.Print
	)

	db := model.DB.Model(&print)
	db = print.WhereCondition(db, &types.PrintWhere{
		Id:      req.Id,
		StoreId: req.StoreId,
		Type:    req.Type,
	})

	db = db.Order("updated_at desc")
	if err := db.Attrs(model.Print{}.Default(req.Type)).FirstOrInit(&print).Error; err != nil {
		return nil, errors.New("获取信息失败")
	}

	return &print, nil
}

func (PrintLogic) Update(req *types.PrintUpdateReq) error {
	print, err := utils.StructToStruct[model.Print](req)
	if err != nil {
		return err
	}

	if err := model.DB.Model(&model.Print{}).
		Where("id = ?", req.Id).
		Where(&model.Print{
			IsDefault: false,
		}).
		Updates(&print).Error; err != nil {
		return errors.New("更新失败")
	}

	return nil
}

func (PrintLogic) Delete(req *types.PrintDeleteReq) error {
	if err := model.DB.Where("id = ?", req.Id).
		Where(&model.Print{
			IsDefault: false,
		}).
		Delete(&model.Print{}).Error; err != nil {
		return errors.New("删除失败")
	}

	return nil
}

func (PrintLogic) Copy(req *types.PrintCopyReq) error {
	var (
		print model.Print
	)

	if err := model.DB.Where("id = ?", req.Id).First(&print).Error; err != nil {
		return errors.New("获取信息失败")
	}

	print.Id = ""
	print.StoreId = req.StoreId
	if req.Name != "" {
		print.Name = req.Name
	}

	def := model.Print{}.Default(0)
	if print.Name == def.Name {
		print.Name = print.Name + "副本"
	}

	if err := model.DB.Create(&print).Error; err != nil {
		log.Printf("复制失败: %v", err)
		return errors.New("复制失败")
	}

	return nil
}
