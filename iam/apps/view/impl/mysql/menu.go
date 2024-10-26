package mysql

import (
	"context"

	"github.com/infraboard/mcube/v2/types"
	"github.com/infraboard/modules/iam/apps/view"
)

// 创建菜单
func (i *MenuServiceImpl) CreateMenu(ctx context.Context, in *view.CreateMenuRequest) (*view.Menu, error) {
	return nil, nil
}

// 查询列表
func (i *MenuServiceImpl) QueryMenu(ctx context.Context, in *view.QueryMenuRequest) (*types.Set[*view.Menu], error) {
	return nil, nil
}

// 查询详情
func (i *MenuServiceImpl) DescribeMenu(ctx context.Context, in *view.DescribeMenuRequest) (*view.Menu, error) {
	return nil, nil
}

// 更新菜单
func (i *MenuServiceImpl) UpdateMenu(ctx context.Context, in *view.UpdateMenuRequest) (*view.Menu, error) {
	return nil, nil
}

// 删除菜单
func (i *MenuServiceImpl) DeleteMenu(ctx context.Context, in *view.DeleteMenuRequest) (*view.Menu, error) {
	return nil, nil
}
