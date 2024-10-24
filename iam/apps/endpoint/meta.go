package endpoint

func GetRouteMeta[T any](m map[string]any, key string) T {
	if v, ok := m[key]; ok {
		return v.(T)
	}

	var t T
	return t
}
