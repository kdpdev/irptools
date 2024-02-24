package alg

func ArrToMap[Key comparable, Value any](keys []Key, defaultValue Value) map[Key]Value {
	result := make(map[Key]Value, len(keys))
	for _, k := range keys {
		result[k] = defaultValue
	}
	return result
}

func MapKeys[Key comparable, Value any](m map[Key]Value) []Key {
	result := make([]Key, 0, len(m))
	for k := range m {
		result = append(result, k)
	}
	return result
}

func MapValues[Key comparable, Value any](m map[Key]Value) []Value {
	result := make([]Value, 0, len(m))
	for _, v := range m {
		result = append(result, v)
	}
	return result
}
