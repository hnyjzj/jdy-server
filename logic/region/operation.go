package region

import (
	"errors"
	"jdy/config"
	"jdy/enums"
	"jdy/model"
	"jdy/types"
	"jdy/utils"
	"log"
	"strconv"
	"strings"

	"github.com/ArtisanCloud/PowerWeChat/v3/src/work/department/request"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (l *RegionLogic) Create(ctx *gin.Context, req *types.RegionCreateReq) error {
	// 验证区域名称是否以"区域"结尾
	if !strings.HasSuffix(req.Name, enums.DepartmentRegion.String()) {
		return errors.New("区域名称必须以" + enums.DepartmentRegion.String() + "结尾")
	}

	region := &model.Region{
		Name:  req.Name,
		Alias: req.Alias,
	}

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 同步企业微信
		wxwork := config.NewWechatService().ContactsWork.Department
		params := &request.RequestDepartmentInsert{
			Name:     req.Name,
			NameEn:   req.Alias,
			ParentID: config.Config.Wechat.Work.MarketingCenter,
		}
		// 创建部门
		res, err := wxwork.Create(l.Ctx, params)
		if err != nil || (res != nil && res.ErrCode != 0) {
			log.Printf("创建部门失败: %v, %+v", err, res)
			return errors.New("创建部门失败")
		}

		region.IdWx = strconv.Itoa(res.ID)
		// 创建区域
		if err := tx.Create(region).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (l *RegionLogic) Update(ctx *gin.Context, req *types.RegionUpdateReq) error {
	// 查询区域信息
	var region model.Region
	if err := model.DB.First(&region, "id = ?", req.Id).Error; err != nil {
		return errors.New("区域不存在或已被删除")
	}

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 验证信息
		data, err := utils.StructToStruct[model.Region](req)
		if err != nil {
			return errors.New("验证信息失败")
		}
		if data.Name != "" {
			// 验证区域名称是否以"区域"结尾
			if !strings.HasSuffix(req.Name, enums.DepartmentRegion.String()) {
				return errors.New("区域名称必须以" + enums.DepartmentRegion.String() + "结尾")
			}
		}

		// 同步企业微信
		wxwork := config.NewWechatService().ContactsWork.Department
		id, err := strconv.Atoi(region.IdWx)
		if err != nil {
			return errors.New("门店ID转换失败")
		}
		params := &request.RequestDepartmentUpdate{
			ID:       id,
			Name:     data.Name,
			NameEn:   data.Alias,
			ParentID: config.Config.Wechat.Work.MarketingCenter,
		}
		res, err := wxwork.Update(l.Ctx, params)
		if err != nil || (res != nil && res.ErrCode != 0) {
			log.Printf("更新部门失败: %v, %+v", err, res)
			return errors.New("更新部门失败")
		}

		// 更新区域
		if err := tx.Model(&model.Region{}).Where("id = ?", region.Id).Updates(data).Error; err != nil {
			return errors.New("更新失败")
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

// 删除区域
func (l *RegionLogic) Delete(ctx *gin.Context, req *types.RegionDeleteReq) error {
	// 查询区域信息
	var (
		region model.Region
	)
	db := model.DB.Where("id = ?", req.Id)
	db = region.Preloads(db)
	if err := db.First(&region).Error; err != nil {
		return errors.New("区域不存在或已被删除")
	}
	if len(region.Stores) > 0 {
		return errors.New("区域下有门店，无法删除")
	}

	if len(region.Staffs) > 0 || len(region.Superiors) > 0 {
		return errors.New("门店下有员工或负责人，无法删除")
	}

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 同步企业微信
		wxwork := config.NewWechatService().ContactsWork.Department
		// 删除部门
		id, err := strconv.Atoi(region.IdWx)
		if err != nil {
			return errors.New("门店ID转换失败")
		}
		res, err := wxwork.Delete(l.Ctx, id)
		if err != nil || (res != nil && res.ErrCode != 0) {
			log.Printf("删除部门失败: %v", err)
			return errors.New("删除部门失败")
		}
		// 删除区域
		if err := tx.Delete(&region).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
