package mysql

import (
	"context"

	"github.com/infraboard/mcube/v2/types"
	"github.com/infraboard/modules/iam/apps/menu"
)

// 创建菜单
func (i *MenuServiceImpl) CreateMenu(ctx context.Context, in *menu.CreateMenuRequest) (*menu.Menu, error) {
	return nil, nil
}

// 查询列表
func (i *MenuServiceImpl) QueryMenu(ctx context.Context, in *menu.QueryMenuRequest) (*types.Set[*menu.Menu], error) {
	return nil, nil
}

// 查询详情
func (i *MenuServiceImpl) DescribeMenu(ctx context.Context, in *menu.DescribeMenuRequest) (*menu.Menu, error) {
	return nil, nil
}

// 更新菜单
func (i *MenuServiceImpl) UpdateMenu(ctx context.Context, in *menu.UpdateMenuRequest) (*menu.Menu, error) {
	return nil, nil
}

// 删除菜单
func (i *MenuServiceImpl) DeleteMenu(ctx context.Context, in *menu.DeleteMenuRequest) (*menu.Menu, error) {
	return nil, nil
}
