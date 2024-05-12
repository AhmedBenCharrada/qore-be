package utils

import "strconv"

// StringToInt ..
func StringToInt(str string, def int) int {
	num, err := strconv.Atoi(str)
	if err != nil {
		return def
	}
	return num
}
