package utils

// InList 检查是否在list中
func InList(key string, list []string) bool {
	for _, value := range list {
		if key == value {
			return true
		}
	}
	return false
}
