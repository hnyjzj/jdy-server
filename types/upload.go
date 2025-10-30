package types

// 模块
type UploadModel string

func (s UploadModel) String() string {
	return string(s)
}

const (
	UploadModelCommon    UploadModel = "common"    // 通用
	UploadModelAvatar    UploadModel = "avatar"    // 头像
	UploadModelWorkbench UploadModel = "workbench" // 工作台
	UploadModelStore     UploadModel = "store"     // 门店
	UploadModelProduct   UploadModel = "product"   // 商品
	UploadModelOrder     UploadModel = "order"     // 订单
)

// 类型
type UploadType string

func (s UploadType) String() string {
	return string(s)
}

const (
	UploadTypeImage UploadType = "image"
)
