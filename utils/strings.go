package utils

import "github.com/gosimple/slug"

func Slug(val string) string {
	slug.CustomSub = map[string]string{
		" ": "_",
	}
	return slug.Make(val)
}
