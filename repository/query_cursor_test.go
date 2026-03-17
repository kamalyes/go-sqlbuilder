/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-03-17 09:21:33
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-03-17 09:55:07
 * @FilePath: \go-sqlbuilder\repository\query_cursor_test.go
 * @Description: AddCursorFilter 游标分页过滤方法测试
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */
package repository

import (
	"testing"

	"github.com/kamalyes/go-sqlbuilder/constants"
	"github.com/stretchr/testify/assert"
)

// TestAddCursorFilterEmptyString 空字符串游标不应添加过滤条件
func TestAddCursorFilterEmptyString(t *testing.T) {
	query := NewQuery()
	result := query.AddCursorFilter("message_id", "", false)

	assert.Equal(t, query, result)
	assert.Equal(t, 0, len(query.Filters))
}

// TestAddCursorFilterNilCursor nil 游标不应添加过滤条件
func TestAddCursorFilterNilCursor(t *testing.T) {
	query := NewQuery()
	result := query.AddCursorFilter("message_id", nil, false)

	assert.Equal(t, query, result)
	assert.Equal(t, 0, len(query.Filters))
}

// TestAddCursorFilterZeroInt 零值 int 游标不应添加过滤条件
func TestAddCursorFilterZeroInt(t *testing.T) {
	query := NewQuery()
	result := query.AddCursorFilter("id", 0, false)

	assert.Equal(t, query, result)
	assert.Equal(t, 0, len(query.Filters))
}

// TestAddCursorFilterNextPage isPrev=false 时使用 > 向后翻页
func TestAddCursorFilterNextPage(t *testing.T) {
	query := NewQuery()
	result := query.AddCursorFilter("message_id", "msg_100", false)

	assert.Equal(t, query, result)
	assert.Equal(t, 1, len(query.Filters))
	assert.Equal(t, "message_id", query.Filters[0].Field)
	assert.Equal(t, constants.OP_GT, query.Filters[0].Operator)
	assert.Equal(t, "msg_100", query.Filters[0].Value)
}

// TestAddCursorFilterPrevPage isPrev=true 时使用 < 向前翻页
func TestAddCursorFilterPrevPage(t *testing.T) {
	query := NewQuery()
	result := query.AddCursorFilter("message_id", "msg_100", true)

	assert.Equal(t, query, result)
	assert.Equal(t, 1, len(query.Filters))
	assert.Equal(t, "message_id", query.Filters[0].Field)
	assert.Equal(t, constants.OP_LT, query.Filters[0].Operator)
	assert.Equal(t, "msg_100", query.Filters[0].Value)
}

// TestAddCursorFilterInt64Cursor int64 类型游标
func TestAddCursorFilterInt64Cursor(t *testing.T) {
	query := NewQuery()
	query.AddCursorFilter("id", int64(9999), false)

	assert.Equal(t, 1, len(query.Filters))
	assert.Equal(t, constants.OP_GT, query.Filters[0].Operator)
	assert.Equal(t, int64(9999), query.Filters[0].Value)
}

// TestAddCursorFilterChainWithOtherFilters 与其他过滤条件链式组合
func TestAddCursorFilterChainWithOtherFilters(t *testing.T) {
	tests := []struct {
		name        string
		cursor      any
		isPrev      bool
		wantFilters int
		wantOp      constants.Operator
	}{
		{"向后翻页-字符串游标", "msg_50", false, 2, constants.OP_GT},
		{"向前翻页-字符串游标", "msg_50", true, 2, constants.OP_LT},
		{"空游标-不添加", "", false, 1, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := NewQuery().AddEqual("session_id", "sess_abc")
			query.AddCursorFilter("message_id", tt.cursor, tt.isPrev)

			assert.Equal(t, tt.wantFilters, len(query.Filters))
			if tt.wantFilters == 2 {
				assert.Equal(t, tt.wantOp, query.Filters[1].Operator)
				assert.Equal(t, tt.cursor, query.Filters[1].Value)
			}
		})
	}
}
