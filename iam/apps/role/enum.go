package role

const (
	ADMIN = "admin"
)

// EFFECT_TYPE 授权效力包括两种：允许（Allow）和拒绝（Deny）
type EFFECT_TYPE int32

const (
	// 允许访问
	EFFECT_TYPE_ALLOW EFFECT_TYPE = iota
	// 拒绝访问
	EFFECT_TYPE_DENY
)
