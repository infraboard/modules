package user

import (
	"context"
	"fmt"
)

const (
	AppName = "user"
)

// 定义User包的能力 就是接口定义
// 站在使用放的角度来定义的   userSvc.Create(ctx, req), userSvc.DeleteUser(id)
// 接口定义好了，不要试图 随意修改接口， 要保证接口的兼容性
type Service interface {
	// 创建用户
	CreateUser(context.Context, *CreateUserRequest) (*User, error)
	// 删除用户
	DeleteUser(context.Context, *DeleteUserRequest) error

	// 查询用户  User.CheckPassword(xxx)
	DescribeUserRequest(context.Context, *DescribeUserRequest) (*User, error)
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

// 删除用户的请求
type DeleteUserRequest struct {
	Id int `json:"id"`
}

func (req *DeleteUserRequest) IdString() string {
	return fmt.Sprintf("%d", req.Id)
}
