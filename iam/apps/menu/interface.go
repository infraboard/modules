package menu

import (
	"context"

	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/types"
)

const (
	AppName = "menu"
)

func GetService() Service {
	return ioc.Controller().Get(AppName).(Service)
}

type Service interface {
	// 创建菜单
	CreateMenu(context.Context, *CreateMenuRequest) (*Menu, error)
	// 查询列表
	QueryMenu(context.Context, *QueryMenuRequest) (*types.Set[*Menu], error)
	// 查询详情
	DescribeMenu(context.Context, *DescribeMenuRequest) (*Menu, error)
	// 更新菜单
	UpdateMenu(context.Context, *UpdateMenuRequest) (*Menu, error)
	// 删除菜单
	DeleteMenu(context.Context, *DeleteMenuRequest) (*Menu, error)
}

type QueryMenuRequest struct {
}

type DescribeMenuRequest struct {
}

type UpdateMenuRequest struct {
}

type DeleteMenuRequest struct {
}
