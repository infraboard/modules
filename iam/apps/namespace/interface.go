package namespace

import (
	"context"

	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/types"
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

type QueryNamespaceRequest struct {
}

type DescribeNamespaceRequest struct {
}

type UpdateNamespaceRequest struct {
}

type DeleteNamespaceRequest struct {
}
