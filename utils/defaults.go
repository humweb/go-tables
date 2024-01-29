package utils

import "fmt"

// DefaultInt Checks if initial int is not empty and returns default
func DefaultInt(val int, def int) int {
	if val == 0 {
		return def
	}
	return val
}

// DefaultSort checks if initial value is not empty and returns a default value
func DefaultSort(val string, def string) string {
	if val == "" {
		val = def
	}
	if val[0:1] == "-" {
		return fmt.Sprintf("%s DESC", val[1:])
	}
	return fmt.Sprintf("%s ASC", val)
}

// DefaultString checks if initial string is empty and returns default
func DefaultString(val string, def string) string {
	if val == "" {
		return def
	}
	return val
}
