package view

import "github.com/infraboard/modules/iam/apps"

func NewPage() *Page {
	return &Page{
		Meta: *apps.NewMeta(),
	}
}

type Page struct {
	// 基础数据
	apps.Meta
	// 菜单定义
	CreatePageRequest
	// 用户是否有权限访问该页面, 只有在策略模块查询时，才会计算出该字段
	HasPermission *bool `json:"has_permission,omitempty" gorm:"column:has_permission;type:tinyint(1)" optional:"true" description:"用户是否有权限访问该页面"`
}

func (p *Page) TableName() string {
	return "pages"
}

func NewCreatePageRequest() *CreatePageRequest {
	return &CreatePageRequest{
		Extras: map[string]string{},
	}
}

type CreatePageRequest struct {
	// 菜单Id
	MenuId uint64 `json:"menu_id" bson:"menu_id" gorm:"column:menu_id;type:uint;index" description:"菜单Id"`
	// 页面路径
	Path string `json:"path" bson:"path" gorm:"column:path" description:"页面路径" unique:"true"`
	// 页面名称
	Name string `json:"name" bson:"name" gorm:"column:name" description:"页面名称"`
	// 标签
	Label string `json:"label" gorm:"column:label;type:varchar(200);index" description:"标签" optional:"true"`
	// 关联的Api接口
	RefEndpointId []uint64 `json:"ref_endpoints" gorm:"column:ref_endpoints;serializer:json;type:json" description:"该页面管理的Api接口关联的接口" optional:"true"`
	// 其他扩展信息
	Extras map[string]string `json:"extras" gorm:"column:extras;serializer:json;type:json" description:"其他扩展信息" optional:"true"`
}
