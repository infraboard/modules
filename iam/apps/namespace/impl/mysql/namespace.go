package mysql

import (
	"context"

	"github.com/infraboard/mcube/v2/types"
	"github.com/infraboard/modules/iam/apps/namespace"
)

// 创建空间
func (i *NameSpaceServiceImpl) CreateNamespace(ctx context.Context, in *namespace.CreateNamespaceRequest) (*namespace.Namespace, error) {
	return nil, nil
}

// 查询空间
func (i *NameSpaceServiceImpl) QueryNamespace(ctx context.Context, in *namespace.QueryNamespaceRequest) (*types.Set[*namespace.Namespace], error) {
	return nil, nil
}

// 查询空间详情
func (i *NameSpaceServiceImpl) DescribeNamespace(ctx context.Context, in *namespace.DescribeNamespaceRequest) (*namespace.Namespace, error) {
	return nil, nil
}

// 更新空间
func (i *NameSpaceServiceImpl) UpdateNamespace(ctx context.Context, in *namespace.UpdateNamespaceRequest) (*namespace.Namespace, error) {
	return nil, nil
}

// 删除空间
func (i *NameSpaceServiceImpl) DeleteNamespace(ctx context.Context, in *namespace.DeleteNamespaceRequest) (*namespace.Namespace, error) {
	return nil, nil
}
