package cmd

import (
	"context"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/infraboard/modules/iam/apps/namespace"
	"github.com/infraboard/modules/iam/apps/user"
	"github.com/spf13/cobra"

	"github.com/infraboard/mcube/v2/ioc/server/cmd"
)

func init() {
	cmd.Root.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "初始化Admin用户名密码",
	Run: func(cmd *cobra.Command, args []string) {
		req := user.NewCreateUserRequest()
		req.IsAdmin = true
		req.EnabledApi = true

		cobra.CheckErr(survey.AskOne(
			&survey.Input{
				Message: "请输入管理员用户名称:",
				Default: "admin",
			},
			&req.UserName,
			survey.WithValidator(survey.Required),
		))
		cobra.CheckErr(survey.AskOne(
			&survey.Password{
				Message: "请输入管理员密码:",
			},
			&req.Password,
			survey.WithValidator(survey.Required),
		))

		var repeatPass string
		cobra.CheckErr(survey.AskOne(
			&survey.Password{
				Message: "再次输入管理员密码:",
			},
			&repeatPass,
			survey.WithValidator(survey.Required),
			survey.WithValidator(func(ans any) error {
				if ans.(string) != req.Password {
					return fmt.Errorf("两次输入的密码不一致")
				}
				return nil
			}),
		))

		// 添加管理员用户
		u, err := user.GetService().CreateUser(context.Background(), req)
		cobra.CheckErr(err)
		fmt.Println(u)

		// 添加默认空间
		defaultNS := namespace.NewCreateNamespaceRequest()
		defaultNS.Name = namespace.DEFAULT_NS_NAME
		defaultNS.Description = "默认空间"
		defaultNS.OwnerUserId = u.Id
		ns, err := namespace.GetService().CreateNamespace(context.Background(), defaultNS)
		cobra.CheckErr(err)
		fmt.Println(ns)
	},
}
