package utils

func InArray(needle any, array any) bool {
	switch key := needle.(type) {
	case int:
		for _, item := range array.([]int) {
			if key == item {
				return true
			}
		}
	case string:
		for _, item := range array.([]string) {
			if key == item {
				return true
			}
		}
	case int64:
		for _, item := range array.([]int64) {
			if key == item {
				return true
			}
		}

	default:
		return false
	}
	return false
}
