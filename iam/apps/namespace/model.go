package namespace

type Nmespace struct {
	// 对象Id
	Id string `json:"id" bson:"_id"`
	// 创建时间
	CreateAt int64 `json:"create_at" bson:"create_at"`
	// 更新时间
	UpdateAt int64 `json:"update_at" bson:"update_at"`
	// 更新人
	UpdateBy string `json:"update_by" bson:"update_by"`
	// 空间属性
	*CreateNamespaceRequest
}

type CreateNamespaceRequest struct {
	// 父Namespace Id
	ParentId string `json:"parent_id" bson:"parent_id" gorm:"column:parent_id"`
	// 空间名称, 不允许修改
	Name string `json:"name" bson:"name" validate:"required" gorm:"column:name"`
	// 空间负责人
	Owner string `json:"owner" bson:"owner" gorm:"column:owner"`
	// 空间负责人助理
	Assistants []string `json:"assistants" bson:"assistants" gorm:"column:assistants"`
	// 禁用项目, 该项目所有人暂时都无法访问
	Enabled bool `json:"enabled" bson:"enabled" gorm:"column:enabled"`
	// 项目描述图片
	Picture string `json:"picture" bson:"picture" gorm:"column:picture"`
	// 项目描述
	Description string `json:"description" bson:"description" gorm:"column:description"`
	// 空间标签
	Label map[string]string `json:"label" bson:"label" gorm:"column:label;serializer:json"`
	// 扩展信息
	Extras map[string]string `json:"extras" bson:"extras" gorm:"column:extras;serializer:json"`
}
