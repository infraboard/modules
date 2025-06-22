package label

import (
	"context"

	"github.com/infraboard/mcube/v2/http/request"
	request1 "github.com/infraboard/mcube/v2/pb/request"
	"github.com/infraboard/mcube/v2/pb/resource"
	"github.com/infraboard/mcube/v2/types"
)

type Service interface {
	// 创建标签
	CreateLabel(context.Context, *CreateLabelRequest) (*Label, error)
	// 修改标签
	UpdateLabel(context.Context, *UpdateLabelRequest) (*Label, error)
	// 删除标签
	DeleteLabel(context.Context, *DeleteLabelRequest) (*Label, error)
	// 查询标签列表
	QueryLabel(context.Context, *QueryLabelRequest) (*types.Set[*Label], error)
	// 查询标签列表
	DescribeLabel(context.Context, *DescribeLabelRequest) (*Label, error)
}

type UpdateLabelRequest struct {
	// 更新模式
	UpdateMode request1.UpdateMode `json:"update_mode"`
	// 标签Id
	Id string `json:"id"`
	// 更新人
	UpdateBy string `json:"update_by"`
	// 标签信息
	Spec *CreateLabelRequest `json:"spec"`
}

type DeleteLabelRequest struct {
	// 标签Id
	Id string `json:"id"`
}

type QueryLabelRequest struct {
	// 资源范围
	Scope *resource.Scope `json:"scope"`
	// 分页请求
	Page *request.PageRequest `json:"page"`
	// key
	Keys []string `json:"keys"`
}

type DescribeLabelRequest struct {
	// 标签Id
	Id string `json:"id"`
}
