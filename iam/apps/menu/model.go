package menu

import "github.com/infraboard/modules/iam/apps"

func NewMenu() *Menu {
	return &Menu{
		Meta: *apps.NewMeta(),
	}
}

type Menu struct {
	// 基础数据
	apps.Meta
	// 菜单定义
	CreateMenuRequest
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
	// 父Menu Id
	ParentId uint64 `json:"parent_id" bson:"parent_id" gorm:"column:parent_id;type:uint;index" description:"父Menu Id" optional:"true"`
	// 菜单路径
	Path string `json:"path" bson:"path" gorm:"column:path" description:"菜单路径" unique:"true"`
	// 菜单名称
	Name string `json:"name" bson:"name" gorm:"column:name" description:"菜单名称"`
	// 图标
	Icon string `json:"icon" bson:"icon" gorm:"column:icon" description:"图标" optional:"true"`
	// 标签
	Label string `json:"label" gorm:"column:label;type:varchar(200);index" description:"标签" optional:"true"`
	// 其他扩展信息
	Extras map[string]string `json:"extras" gorm:"column:extras;serializer:json;type:json" description:"其他扩展信息" optional:"true"`
}
