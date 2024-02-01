package tables

import (
	"gorm.io/gorm"
)

type ITable interface {
	GetFields() []*Field
	GetFilters() []*Filter
	WithGlobalSearch(db *gorm.DB, val string)
	ApplyFilter(db *gorm.DB)
}
