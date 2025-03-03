package types

// 模块
type UploadModel string

func (s UploadModel) String() string {
	return string(s)
}

const (
	UploadModelAvatar    UploadModel = "avatar"    // 头像
	UploadModelWorkbench UploadModel = "workbench" // 工作台
	UploadModelStore     UploadModel = "store"     // 门店
	UploadModelProduct   UploadModel = "product"   // 商品
)

// 类型
type UploadType string

func (s UploadType) String() string {
	return string(s)
}

const (
	UploadTypeImage UploadType = "image"
)
