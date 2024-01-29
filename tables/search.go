package tables

import (
	"strconv"

	"gorm.io/gorm"
)

type Search struct {
	Label   string `json:"label"`
	Field   string `json:"field"`
	Value   string `json:"value"`
	Enabled bool   `json:"enabled"`
}

func (f *Search) ApplySearch(db *gorm.DB) {
	if v, err := strconv.Atoi(f.Value); err == nil {
		db.Where(f.Field+" = ?", v)
	} else {
		db.Where(f.Field+" ILIKE ?", "%"+f.Value+"%")
	}
}
