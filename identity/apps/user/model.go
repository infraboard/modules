package user

import (
	"encoding/json"
	"fmt"

	"github.com/infraboard/mcube/v2/pb/resource"
	"golang.org/x/crypto/bcrypt"
)

func NewUser(req *CreateUserRequest) *User {
	req.PasswordHash()

	return &User{
		Meta:              resource.NewMeta(),
		CreateUserRequest: req,
	}
}

// 用于存放 存入数据库的对象(PO)
type User struct {
	// 通用信息
	*resource.Meta
	// 用户传递过来的请求
	*CreateUserRequest
}

func (u *User) String() string {
	dj, _ := json.Marshal(u)
	return string(dj)
}

// 判断该用户的密码是否正确
func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

// 声明你这个对象存储在users表里面
// orm 负责调用TableName() 来动态获取你这个对象要存储的表的名称
func (u *User) TableName() string {
	return "users"
}

func NewCreateUserRequest() *CreateUserRequest {
	return &CreateUserRequest{
		Role:  ROLE_AUDITOR,
		Label: map[string]string{},
	}
}

// VO
type CreateUserRequest struct {
	//
	Username string `json:"username" gorm:"column:username"`
	//
	Password string `json:"password"`
	// 用户的角色(写死)，
	// @Role("admin")
	// CreateBlog( check user Role  )
	Role Role `json:"role"`
	// 对象标签, Dep:部门A
	// Label 没法存入数据库，不是一个结构化的数据
	// 比如就存储在数据里面 ，存储为Json, 需要ORM来帮我们完成 json的序列化和存储
	// 直接序列化为Json存储到 lable字段
	Label map[string]string `json:"label" gorm:"serializer:json"`

	isHashed bool
}

func (req *CreateUserRequest) SetIsHashed() {
	req.isHashed = true
}

func (req *CreateUserRequest) Validate() error {
	if req.Username == "" || req.Password == "" {
		return fmt.Errorf("用户名或者密码需要填写")
	}
	return nil
}

func (req *CreateUserRequest) PasswordHash() {
	if req.isHashed {
		return
	}

	b, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	req.Password = string(b)
	req.isHashed = true
}
