package event

const (
	// DEBUG
	LEVEL_DEBUG LEVEL = iota
	// INFO
	LEVEL_INFO
	// WARN
	LEVEL_WARN
	// ERROR
	LEVEL_ERROR
)

type LEVEL int

const (
	ORDER_BY_DESC = "DESC"
	ORDER_BY_ASC  = "ASC"
)

type ORDER_BY string
