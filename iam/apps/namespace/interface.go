package namespace

import (
	"context"

	"github.com/infraboard/mcube/v2/http/request"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/types"
	"github.com/infraboard/modules/iam/apps"
)

const (
	AppName = "namespace"
)

func GetService() Service {
	return ioc.Controller().Get(AppName).(Service)
}

type Service interface {
	// 创建空间
	CreateNamespace(context.Context, *CreateNamespaceRequest) (*Namespace, error)
	// 查询空间
	QueryNamespace(context.Context, *QueryNamespaceRequest) (*types.Set[*Namespace], error)
	// 查询空间详情
	DescribeNamespace(context.Context, *DescribeNamespaceRequest) (*Namespace, error)
	// 更新空间
	UpdateNamespace(context.Context, *UpdateNamespaceRequest) (*Namespace, error)
	// 删除空间
	DeleteNamespace(context.Context, *DeleteNamespaceRequest) (*Namespace, error)
}

func NewQueryNamespaceRequest() *QueryNamespaceRequest {
	return &QueryNamespaceRequest{
		PageRequest: request.NewDefaultPageRequest(),
	}
}

type QueryNamespaceRequest struct {
	*request.PageRequest
}

func NewDescribeNamespaceRequest() *DescribeNamespaceRequest {
	return &DescribeNamespaceRequest{}
}

type DescribeNamespaceRequest struct {
	apps.GetRequest
}

func NewUpdateNamespaceRequest() *UpdateNamespaceRequest {
	return &UpdateNamespaceRequest{}
}

type UpdateNamespaceRequest struct {
	apps.GetRequest
	CreateNamespaceRequest
}

func NewDeleteNamespaceRequest() *DeleteNamespaceRequest {
	return &DeleteNamespaceRequest{}
}

type DeleteNamespaceRequest struct {
	apps.GetRequest
}
