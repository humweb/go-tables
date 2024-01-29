package utils

import "fmt"

func DefaultInt(val int, def int) int {
	if val == 0 {
		return def
	}
	return val
}
func DefaultSort(val string, def string) string {
	if val == "" {
		val = def
	}
	if val[0:1] == "-" {
		return fmt.Sprintf("%s DESC", val[1:len(val)])
	}
	return fmt.Sprintf("%s ASC", val)
}
func DefaultString(val string, def string) string {
	if val == "" {
		return def
	} else {
		return val
	}
}
