package role

const (
	ADMIN = "admin"
)

type MATCH_BY int32

const (
	MATCH_BY_ID = iota
	MATCH_BY_LABLE
	MATCH_BY_RESOURCE_ACTION
	MATCH_BY_RESOURCE_ACCESS_MODE
)
