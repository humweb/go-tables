package tables

import (
	"gorm.io/gorm"
)

type ITable interface {
	GetModel() Model
	GetFields() []*Field
	GetFilters() []*Filter
	WithGlobalSearch(db *gorm.DB, val string)
	ApplyFilter(db *gorm.DB)
}
