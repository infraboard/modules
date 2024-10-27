package view

import (
	"context"

	"github.com/infraboard/mcube/v2/http/request"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/types"
	"github.com/infraboard/modules/iam/apps"
)

const (
	AppName = "view"
)

func GetService() Service {
	return ioc.Controller().Get(AppName).(Service)
}

type Service interface {
	// 菜单管理
	MenuService
	// 页面管理
	PageService
}

type MenuService interface {
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

func NewQueryMenuRequest() *QueryMenuRequest {
	return &QueryMenuRequest{
		PageRequest: request.NewDefaultPageRequest(),
	}
}

type QueryMenuRequest struct {
	*request.PageRequest
}

func NewDescribeMenuRequest() *DescribeMenuRequest {
	return &DescribeMenuRequest{}
}

type DescribeMenuRequest struct {
	apps.GetRequest
}

type UpdateMenuRequest struct {
}

func NewDeleteMenuRequest() *DeleteMenuRequest {
	return &DeleteMenuRequest{}
}

type DeleteMenuRequest struct {
	apps.GetRequest
}

type PageService interface {
	// 页面列表
	QueryPage(context.Context, *QueryPageRequest) (*types.Set[*Page], error)
	// 页面详情
	DescribePage(context.Context, *DescribePageRequest) (*Page, error)
	// 添加页面
	CreatePage(context.Context, *CreatePageRequest) (*Page, error)
	// 移除页面
	DeletePage(context.Context, *DeletePageRequest) (*Page, error)
	// 更新页面
	UpdatePage(context.Context, *UpdatePageRequest) (*Page, error)
}

func NewQueryPageRequest() *QueryPageRequest {
	return &QueryPageRequest{
		PageRequest: request.NewDefaultPageRequest(),
	}
}

type QueryPageRequest struct {
	*request.PageRequest
	// 菜单Id
	MenuId uint64 `json:"menu_id"`
}

func NewDescribePageRequest() *DescribePageRequest {
	return &DescribePageRequest{}
}

type DescribePageRequest struct {
	apps.GetRequest
}

func NewDeletePageRequest() *DeletePageRequest {
	return &DeletePageRequest{}
}

type DeletePageRequest struct {
	apps.GetRequest
}

type UpdatePageRequest struct {
}
