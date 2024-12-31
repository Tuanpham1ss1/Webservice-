package utils

import "strconv"

func GetPattern(id uint) string {
	return strconv.Itoa(int(id)) + "-:--"
}
