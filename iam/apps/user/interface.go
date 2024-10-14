package user

import (
	"context"

	"github.com/infraboard/mcube/v2/http/request"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/types"
)

const (
	AppName = "users"
)

func GetService() Service {
	return ioc.Controller().Get(AppName).(Service)
}

// 定义User包的能力 就是接口定义
// 站在使用放的角度来定义的   userSvc.Create(ctx, req), userSvc.DeleteUser(id)
// 接口定义好了，不要试图 随意修改接口， 要保证接口的兼容性
type Service interface {
	// 创建用户
	CreateUser(context.Context, *CreateUserRequest) (*User, error)
	// 删除用户
	DeleteUser(context.Context, *DeleteUserRequest) (*User, error)

	// 查询用户详情
	DescribeUser(context.Context, *DescribeUserRequest) (*User, error)
	// 查询用户列表
	QueryUser(context.Context, *QueryUserRequest) (*types.Set[*User], error)
}

func NewQueryUserRequest() *QueryUserRequest {
	return &QueryUserRequest{
		PageRequest: request.NewDefaultPageRequest(),
	}
}

type QueryUserRequest struct {
	*request.PageRequest
}

func NewDescribeUserRequestById(id string) *DescribeUserRequest {
	return &DescribeUserRequest{
		DescribeValue: id,
	}
}

func NewDescribeUserRequestByUsername(username string) *DescribeUserRequest {
	return &DescribeUserRequest{
		DescribeBy:    DESCRIBE_BY_USERNAME,
		DescribeValue: username,
	}
}

// 同时支持通过Id来查询，也要支持通过username来查询
type DescribeUserRequest struct {
	DescribeBy    DescribeBy `json:"describe_by"`
	DescribeValue string     `json:"describe_value"`
}

func NewDeleteUserRequest(id string) *DeleteUserRequest {
	return &DeleteUserRequest{
		Id: id,
	}
}

// 删除用户的请求
type DeleteUserRequest struct {
	Id string `json:"id"`
}
