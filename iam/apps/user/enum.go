package user

// 使用这个ROLE类型来表现枚举类型
type Role string

func (r Role) String() string {
	return string(r)
}

const (
	// 普通成员
	ROLE_MEMBER Role = "member"
	// 系统管理员
	ROLE_ADMIN Role = "admin"
)

type DescribeBy int

const (
	DESCRIBE_BY_ID DescribeBy = iota
	DESCRIBE_BY_USERNAME
)
