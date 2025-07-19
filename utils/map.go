package utils

func IfExistElseCreate(obj map[string]interface{}, key string) map[string]interface{} {
	if val, ok := obj[key]; ok {
		if m, ok := val.(map[string]interface{}); ok {
			return m
		}
	}
	m := make(map[string]interface{})
	obj[key] = m
	return m
}

func GetMap(obj map[string]interface{}, key string) (map[string]interface{}, bool) {
	val, ok := obj[key]
	if !ok {
		return nil, false
	}
	m, ok := val.(map[string]interface{})
	return m, ok
}
