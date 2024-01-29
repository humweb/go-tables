package utils

import "github.com/gosimple/slug"

// Slug generates snake_case slug string
func Slug(val string) string {
	slug.CustomSub = map[string]string{
		" ": "_",
	}
	return slug.Make(val)
}
