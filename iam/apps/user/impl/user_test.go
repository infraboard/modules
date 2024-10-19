package impl_test

import (
	"testing"

	"github.com/infraboard/modules/iam/apps/user"
)

func TestQueryUser(t *testing.T) {
	req := user.NewQueryUserRequest()
	set, err := impl.QueryUser(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(set)
}

func TestCreateAdminUser(t *testing.T) {
	req := user.NewCreateUserRequest()
	req.UserName = "admin"
	req.Password = "123456"
	req.IsAdmin = true
	u, err := impl.CreateUser(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(u)
}

func TestCreateAuthor2(t *testing.T) {
	req := user.NewCreateUserRequest()
	req.UserName = "张三"
	req.Password = "123456"
	u, err := impl.CreateUser(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(u)
}

func TestCreateAuditorUser(t *testing.T) {
	req := user.NewCreateUserRequest()
	req.UserName = "auditor"
	req.Password = "123456"
	u, err := impl.CreateUser(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(u)
}

func TestDeleteUser(t *testing.T) {
	_, err := impl.DeleteUser(ctx, &user.DeleteUserRequest{
		Id: "9",
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestDescribeUserRequestById(t *testing.T) {
	req := user.NewDescribeUserRequestById("0192a3ca-08a9-76d6-8ceb-ab89087eaf6c")
	ins, err := impl.DescribeUser(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ins)
}

// SELECT * FROM `users` WHERE username = 'admin' ORDER BY `users`.`id` LIMIT 1
func TestDescribeUserRequestByName(t *testing.T) {
	req := user.NewDescribeUserRequestByUsername("admin")
	ins, err := impl.DescribeUser(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ins)

	err = ins.CheckPassword("123456")
	if err != nil {
		t.Fatal(err)
	}
}

func TestUserJson(t *testing.T) {
	u := user.NewUser(user.NewCreateUserRequest())
	t.Log(u)
}
