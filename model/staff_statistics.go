package model

type StaffCustomerStatistics struct {
	SoftDelete

	StaffId string `json:"staff_id" gorm:"type:varchar(255);not NULL;uniqueIndex:idx_staff_time;comment:员工ID;"` // 员工ID
	Staff   Staff  `json:"staff"  gorm:"foreignKey:StaffId;references:Id;comment:员工;"`                          // 员工

	StatTime int64 `json:"stat_time" gorm:"type:bigint;size:20;uniqueIndex:idx_staff_time;comment:统计时间(时间戳);"` // 统计时间(时间戳)

	ChatCnt             int64   `json:"chat_cnt" gorm:"type:int;size:11;not null;default:0;comment:聊天总数;"`                     // 聊天总数
	MessageCnt          int64   `json:"message_cnt" gorm:"type:int;size:11;not null;default:0;comment:发送消息数;"`                 // 发送消息数
	ReplyPercentage     float64 `json:"reply_percentage" gorm:"type:decimal(10,2);not null;default:0;comment:已回复聊天占比;"`        // 已回复聊天占比
	AvgReplyTime        int64   `json:"avg_reply_time" gorm:"type:int;size:11;not null;default:0;comment:平均首次回复时长(分钟);"`       // 平均首次回复时长(分钟)
	NegativeFeedbackCnt int64   `json:"negative_feedback_cnt" gorm:"type:int;size:11;not null;default:0;comment:删除/拉黑成员的客户数;"` // 删除/拉黑成员的客户数
	NewApplyCnt         int64   `json:"new_apply_cnt" gorm:"type:int;size:11;not null;default:0;comment:发起申请数;"`               // 发起申请数
	NewContactCnt       int64   `json:"new_contact_cnt" gorm:"type:int;size:11;not null;default:0;comment:新增客户数;"`             // 新增客户数
}

func init() {
	// 注册模型
	RegisterModels(
		&StaffCustomerStatistics{},
	)
	// 重置表
	RegisterRefreshModels(
	// &StaffCustomerStatistics{},
	)
}
