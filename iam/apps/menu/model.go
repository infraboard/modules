package menu

import "github.com/infraboard/modules/iam/apps"

type Menu struct {
	// 基础数据
	*apps.Meta
	// 菜单定义
	*CreateMenuRequest
}

func (u *Menu) TableName() string {
	return "menus"
}

func NewCreateMenuRequest() *CreateMenuRequest {
	return &CreateMenuRequest{
		Extras: map[string]string{},
	}
}

type CreateMenuRequest struct {
	// 父Namespace Id
	ParentId uint64 `json:"parent_id" bson:"parent_id" gorm:"column:parent_id;type:uint;index"`
	// 菜单路径
	Path string `json:"path" bson:"path" gorm:"column:path"`
	// 菜单名称
	Name string `json:"name" bson:"name" gorm:"column:name"`
	// 图标
	Icon string `json:"icon" bson:"icon" gorm:"column:icon"`
	// 标签
	Label string `json:"label" gorm:"column:label;type:varchar(200);index"`
	// 其他扩展信息
	Extras map[string]string `json:"extras" gorm:"column:extras;serializer:json;type:json"`
}
