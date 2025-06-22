package label

import (
	"context"

	"github.com/infraboard/mcube/v2/http/request"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/types"
)

const (
	APP_NAME = "lable"
)

func GetService() Service {
	return ioc.Controller().Get(APP_NAME).(Service)
}

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
	DescribeLabelRequest
	// 更新人
	UpdateBy string `json:"update_by"`
	// 标签信息
	Spec *CreateLabelRequest `json:"spec"`
}

type DeleteLabelRequest struct {
	DescribeLabelRequest
}

func NewQueryLabelRequest() *QueryLabelRequest {
	return &QueryLabelRequest{
		PageRequest: request.NewDefaultPageRequest(),
	}
}

type QueryLabelRequest struct {
	*request.PageRequest
}

type DescribeLabelRequest struct {
	// 标签Id
	Id string `json:"id"`
}
