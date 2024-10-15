package role

type Role struct {
	// 在添加数据需要村的定义
	Id int64 `json:"id" gorm:"column:id"`
	// 创建时间
	CreatedAt int64 `json:"created_at" gorm:"column:created_at"`
	// 更新时间
	UpdatedAt int64 `json:"updated_at" gorm:"column:updated_at"`
	// 角色创建信息
	*CreateRoleRequest
}

func NewCreateRoleRequest() *CreateRoleRequest {
	return &CreateRoleRequest{
		Lables: map[string]string{},
	}
}

type CreateRoleRequest struct {
	// 创建者ID
	CreateBy string `json:"create_by" gorm:"column:create_by"`
	// 角色名称
	Name string `json:"name" gorm:"column:name" bson:"name"`
	// 角色描述
	Description string `json:"description" gorm:"column:description" bson:"description"`
	// 是否启用
	Enabled bool `json:"enabled" gorm:"column:enabled" bson:"enabled"`
	// 角色关联的其他信息，比如展示的视图
	Lables map[string]string `json:"lables" gorm:"column:lables" bson:"lables"`
}
