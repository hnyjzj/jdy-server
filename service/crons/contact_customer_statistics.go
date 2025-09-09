package crons

import (
	"context"
	"jdy/config"
	"jdy/enums"
	"jdy/message"
	"jdy/model"
	"log"
	"time"

	"github.com/ArtisanCloud/PowerLibs/v3/object"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/work/externalContact/statistics/request"
	"gorm.io/gorm/clause"
)

func init() {
	RegisterCrons(
		Crons{
			Spec: "0 0 1 * * *",      // 每天凌晨1点执行
			Func: CustomerStatistics, // 客户统计信息
		},
		Crons{
			Spec: "0 0 8 * * *",                  // 每天 8 点执行
			Func: SendCustomerStatisticsPersonal, // 发送个人客户统计信息
		},
	)
}

// 个人客户统计信息
func CustomerStatistics() {
	// 查询所有员工
	var staffs []model.Staff
	if err := model.DB.Find(&staffs).Error; err != nil {
		log.Printf("查询员工信息失败：%s", err.Error())
		return
	}

	var (
		ctx    = context.Background()
		wxwork = config.NewWechatService().JdyWork
	)

	start, end, err := enums.DurationYesterday.GetTime(time.Now())
	if err != nil {
		log.Printf("获取时间失败：%s", err.Error())
		return
	}

	for _, staff := range staffs {
		// 查询个人客户统计信息
		res, err := wxwork.ExternalContactStatistics.GetUserBehaviorData(ctx, &request.RequestGetUserBehaviorData{
			StartTime: start.Unix(),
			EndTime:   end.Add(time.Nanosecond).Unix(),
			UserID:    []string{staff.Username},
		})
		if err != nil || (res != nil && res.ErrCode != 0) {
			continue
		}

		for _, item := range res.MomentList {
			var data model.StaffCustomerStatistics
			if err := object.HashMapToStructure(item.ToHashMap(), &data); err != nil {
				log.Printf("解析个人客户统计信息失败：%s", err.Error())
				continue
			}
			data.StaffId = staff.Id

			// 使用数据库级 UPSERT，显式允许零值更新，避免并发竞态
			if err := model.DB.Clauses(clause.OnConflict{
				Columns: []clause.Column{ // 索引判定列
					{Name: "staff_id"},
					{Name: "stat_time"},
				},
				DoUpdates: clause.AssignmentColumns([]string{ // 显式允许 0 值覆盖
					"chat_cnt",
					"message_cnt",
					"reply_percentage",
					"avg_reply_time",
					"negative_feedback_cnt",
					"new_apply_cnt",
					"new_contact_cnt",
				}),
			}).Create(&data).Error; err != nil {
				log.Printf("保存/更新个人客户统计信息失败 staff_id=%s stat_time=%d：%v", data.StaffId, data.StatTime, err.Error())
				continue
			}
		}
	}
}

// 发送个人客户统计信息
func SendCustomerStatisticsPersonal() {
	// 查询所有员工
	var staffs []model.Staff
	if err := model.DB.Find(&staffs).Error; err != nil {
		log.Printf("查询员工信息失败：%s", err.Error())
		return
	}

	ctx := context.Background()
	strat, _, err := enums.DurationYesterday.GetTime(time.Now())
	if err != nil {
		log.Printf("获取时间失败：%s", err.Error())
		return
	}

	for _, staff := range staffs {
		// 查询个人客户统计信息
		var data model.StaffCustomerStatistics
		if err := model.DB.Where(&model.StaffCustomerStatistics{
			StaffId:  staff.Id,
			StatTime: strat.Unix(),
		}).First(&data).Error; err != nil {
			log.Printf("查询个人客户统计信息失败：%s", err.Error())
			continue
		}

		// if data.ChatCnt == 0 && data.MessageCnt == 0 && data.NewApplyCnt == 0 && data.NewContactCnt == 0 {
		// 	continue
		// }

		req := message.CustomerStatisticsPersonal{
			ToUser:              staff.Username,
			StatTime:            time.Unix(data.StatTime, 0),
			ChatCnt:             data.ChatCnt,
			MessageCnt:          data.MessageCnt,
			ReplyPercentage:     data.ReplyPercentage,
			AvgReplyTime:        data.AvgReplyTime,
			NegativeFeedbackCnt: data.NegativeFeedbackCnt,
			NewApplyCnt:         data.NewApplyCnt,
			NewContactCnt:       data.NewContactCnt,
		}

		// 发送个人客户统计信息
		msg := message.NewMessage(ctx)
		msg.SendCustomerStatisticsPersonal(&req)
	}

}
