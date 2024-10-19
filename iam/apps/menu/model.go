package menu

type Menu struct {
	// 端点名称
	Id string `json:"id" bson:"_id" validate:"required,lte=64" gorm:"column:id"`
	// 创建时间
	CreateAt int64 `json:"create_at" bson:"create_at" gorm:"column:create_at"`
	// 更新时间
	UpdateAt int64 `json:"update_at" bson:"update_at" gorm:"column:update_at"`
	// 菜单定义
	*CreateMenuRequest
}

type CreateMenuRequest struct {
	// 菜单路径
	Path string `json:"path" bson:"path" gorm:"column:path"`
	// 菜单名称
	Name string `json:"name" bson:"name" gorm:"column:name"`
	// 图标
	Icon string `json:"icon" bson:"icon" gorm:"column:icon"`
	// 策略标签
	Label map[string]string `json:"label" bson:"label" gorm:"column:label;serializer:json"`
	// 扩展信息
	Extras map[string]string `json:"extras" bson:"extras" gorm:"column:extras;serializer:json"`
}
