package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestAddFilterIfNotEmptyBoolPointer 测试布尔指针过滤
func TestAddFilterIfNotEmptyBoolPointer(t *testing.T) {
	// true 指针
	trueVal := true
	query := NewQuery()
	query.AddFilterIfNotEmpty("is_active", &trueVal)
	assert.Equal(t, 1, len(query.Filters))
	assert.Equal(t, "is_active", query.Filters[0].Field)
	assert.Equal(t, true, query.Filters[0].Value)

	// false 指针
	falseVal := false
	query = NewQuery()
	query.AddFilterIfNotEmpty("is_active", &falseVal)
	assert.Equal(t, 1, len(query.Filters))
	assert.Equal(t, "is_active", query.Filters[0].Field)
	assert.Equal(t, false, query.Filters[0].Value)

	// nil 指针 - 不应添加过滤条件
	var nilBool *bool
	query = NewQuery()
	query.AddFilterIfNotEmpty("is_active", nilBool)
	assert.Equal(t, 0, len(query.Filters))
}
