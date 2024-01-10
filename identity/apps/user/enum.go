package user

// 使用这个ROLE类型来表现枚举类型
type Role int

const (
	// 创建者, 负责博客创作
	ROLE_AUTHOR Role = iota
	// 审核员
	ROLE_AUDITOR
	// 系统管理员
	ROLE_ADMIN
)

type DescribeBy int

const (
	DESCRIBE_BY_ID DescribeBy = iota
	DESCRIBE_BY_USERNAME
)
