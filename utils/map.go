package utils

func IfExistElseCreate(obj map[string]interface{}, key string) map[string]interface{} {
	if obj == nil {
		obj = make(map[string]interface{})
	}
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

func RecursiveGetMap(obj map[string]interface{}, keys ...string) (map[string]interface{}, bool) {
	if len(keys) == 0 {
		return obj, true
	}
	if val, ok := obj[keys[0]]; ok {
		if m, ok := val.(map[string]interface{}); ok {
			return RecursiveGetMap(m, keys[1:]...)
		}
	}
	return nil, false
}

func RecursiveCreateMap(obj map[string]interface{}, keys ...string) map[string]interface{} {
	if len(keys) == 0 {
		return obj
	}
	if val, ok := obj[keys[0]]; ok {
		if m, ok := val.(map[string]interface{}); ok {
			return RecursiveCreateMap(m, keys[1:]...)
		}
	}
	m := make(map[string]interface{})
	obj[keys[0]] = m
	return RecursiveCreateMap(m, keys[1:]...)
}
