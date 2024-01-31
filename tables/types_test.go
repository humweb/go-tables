package tables

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPaginationDefaults(t *testing.T) {
	is := assert.New(t)
	p := Pagination{}

	// assert equality
	is.Equal(10, p.GetLimit())
	is.Equal(0, p.GetOffset())
	is.Equal(1, p.GetPage())
	is.Equal("id desc", p.GetSort())
}
