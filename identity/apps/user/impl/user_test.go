package impl_test

import (
	"testing"

	"github.com/infraboard/modules/identity/apps/user"
)

func TestCreateAuth1(t *testing.T) {
	req := user.NewCreateUserRequest()
	req.Username = "admin"
	req.Password = "123456"
	req.Role = user.ROLE_AUTHOR
	u, err := impl.CreateUser(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(u)
}

func TestCreateAuthor2(t *testing.T) {
	req := user.NewCreateUserRequest()
	req.Username = "张三"
	req.Password = "123456"
	req.Role = user.ROLE_AUTHOR
	u, err := impl.CreateUser(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(u)
}

func TestCreateAuditorUser(t *testing.T) {
	req := user.NewCreateUserRequest()
	req.Username = "auditor"
	req.Password = "123456"
	req.Role = user.ROLE_AUDITOR
	u, err := impl.CreateUser(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(u)
}

func TestDeleteUser(t *testing.T) {
	err := impl.DeleteUser(ctx, &user.DeleteUserRequest{
		Id: 9,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestDescribeUserRequestById(t *testing.T) {
	req := user.NewDescribeUserRequestById("9")
	ins, err := impl.DescribeUserRequest(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ins)
}

// SELECT * FROM `users` WHERE username = 'admin' ORDER BY `users`.`id` LIMIT 1
func TestDescribeUserRequestByName(t *testing.T) {
	req := user.NewDescribeUserRequestByUsername("admin")
	ins, err := impl.DescribeUserRequest(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ins)

	err = ins.CheckPassword("123456")
	if err != nil {
		t.Fatal(err)
	}
}