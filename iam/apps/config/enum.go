package config

type DESCRIBE_BY int

const (
	DESCRIBE_BY_ID DESCRIBE_BY = iota
	DESCRIBE_BY_KEY
)

type FORMAT int

const (
	FORMAT_JSON FORMAT = iota
)
