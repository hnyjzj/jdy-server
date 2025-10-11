package store

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

func (l *StoreLogic) Create(ctx *gin.Context, req *types.StoreCreateReq) error {
	// 验证门店名称是否以"门店"结尾
	if !strings.HasSuffix(req.Name, enums.DepartmentStore.String()) {
		return errors.New("门店名称必须以" + enums.DepartmentStore.String() + "结尾")
	}

	// 转换结构体
	store, err := utils.StructToStruct[model.Store](req)
	if err != nil {
		return errors.New("验证信息失败")
	}

	// 查询区域
	var region model.Region
	if err := model.DB.First(&region, "id = ?", req.RegionId).Error; err != nil {
		return errors.New("区域不存在或已被删除")
	}

	// 创建门店
	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 同步企业微信
		wxwork := config.NewWechatService().ContactsWork.Department
		regionId, err := strconv.Atoi(region.IdWx)
		if err != nil {
			return errors.New("区域ID转换失败")
		}
		params := &request.RequestDepartmentInsert{
			Name:     req.Name,
			NameEn:   req.Alias,
			ParentID: regionId,
			Order:    req.Order,
		}
		// 创建部门
		res, err := wxwork.Create(l.Ctx, params)
		if err != nil || (res != nil && res.ErrCode != 0) {
			log.Printf("创建部门失败: %v, %+v", err, res)
			return errors.New("创建部门失败")
		}

		store.IdWx = strconv.Itoa(res.ID)

		if err := tx.Create(&store).Error; err != nil {
			return errors.New("创建失败")
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (l *StoreLogic) Update(ctx *gin.Context, req *types.StoreUpdateReq) error {
	// 查询门店信息
	var store model.Store
	if err := model.DB.First(&store, "id = ?", req.Id).Error; err != nil {
		return errors.New("门店不存在或已被删除")
	}

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 验证信息
		data, err := utils.StructToStruct[model.Store](req)
		if err != nil {
			return errors.New("验证信息失败")
		}
		if data.Name != "" {
			// 验证门店名称是否以"门店"结尾
			if !strings.HasSuffix(data.Name, enums.DepartmentStore.String()) {
				return errors.New("门店名称必须以" + enums.DepartmentStore.String() + "结尾")
			}
		}

		// 同步企业微信
		wxwork := config.NewWechatService().ContactsWork.Department
		id, err := strconv.Atoi(store.IdWx)
		if err != nil {
			return errors.New("门店ID转换失败")
		}
		params := &request.RequestDepartmentUpdate{
			ID:     id,
			Name:   data.Name,
			NameEn: data.Alias,
		}
		res, err := wxwork.Update(l.Ctx, params)
		if err != nil || (res != nil && res.ErrCode != 0) {
			log.Printf("更新部门失败: %v, %+v", err, res)
			return errors.New("更新部门失败")
		}

		// 更新门店信息
		if err := tx.Model(&model.Store{}).Where("id = ?", store.Id).Updates(data).Error; err != nil {
			return errors.New("更新失败")
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

// 删除门店
func (l *StoreLogic) Delete(ctx *gin.Context, req *types.StoreDeleteReq) error {
	// 查询门店信息
	var (
		store model.Store
	)
	db := model.DB.Where("id = ?", req.Id)
	db = store.Preloads(db)
	if err := db.First(&store).Error; err != nil {
		return errors.New("门店不存在或已被删除")
	}

	if len(store.Staffs) > 0 || len(store.Superiors) > 0 {
		return errors.New("门店下有员工或负责人，无法删除")
	}

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 同步企业微信
		wxwork := config.NewWechatService().ContactsWork.Department
		// 删除部门
		id, err := strconv.Atoi(store.IdWx)
		if err != nil {
			return errors.New("门店ID转换失败")
		}
		res, err := wxwork.Delete(l.Ctx, id)
		if err != nil || (res != nil && res.ErrCode != 0) {
			log.Printf("删除部门失败: %v", err)
			return errors.New("删除部门失败")
		}

		if err := tx.Delete(&store).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
