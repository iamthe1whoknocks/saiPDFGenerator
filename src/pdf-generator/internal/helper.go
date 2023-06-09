package internal

// mapGetValue get value from the map by path, and fill default if not exists
func (is InternalService) mapGetValue(m map[string]interface{}, path, def string) string {
	var result string
	val, ok := m[path]
	if ok {
		result, ok = val.(string)
		if !ok {
			return def
		}
	} else {
		return def
	}

	return result
}
