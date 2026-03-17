/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-03-17 10:21:05
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-03-17 10:21:05
 * @FilePath: \go-sqlbuilder\repository\query_eq_or_in_test.go
 * @Description: AddEqOrInFilter 单值等于/多值 IN 过滤方法测试
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */
package repository

import (
	"testing"

	"github.com/kamalyes/go-sqlbuilder/constants"
	"github.com/stretchr/testify/assert"
)

// TestAddEqOrInFilterNil nil 不应添加过滤条件
func TestAddEqOrInFilterNil(t *testing.T) {
	query := NewQuery()
	result := query.AddEqOrInFilter("session_id", nil)

	assert.Equal(t, query, result)
	assert.Equal(t, 0, len(query.Filters))
}

// TestAddEqOrInFilterEmptySlice 空切片不应添加过滤条件
func TestAddEqOrInFilterEmptySlice(t *testing.T) {
	query := NewQuery()
	query.AddEqOrInFilter("session_id", []string{})

	assert.Equal(t, 0, len(query.Filters))
}

// TestAddEqOrInFilterSingleValue 单值使用 = 过滤
func TestAddEqOrInFilterSingleValue(t *testing.T) {
	query := NewQuery()
	query.AddEqOrInFilter("session_id", []string{"session_abc"})

	assert.Equal(t, 1, len(query.Filters))
	assert.Equal(t, "session_id", query.Filters[0].Field)
	assert.Equal(t, constants.OP_EQ, query.Filters[0].Operator)
	assert.Equal(t, "session_abc", query.Filters[0].Value)
}

// TestAddEqOrInFilterMultipleValues 多值使用 IN 过滤
func TestAddEqOrInFilterMultipleValues(t *testing.T) {
	query := NewQuery()
	query.AddEqOrInFilter("session_id", []string{"s1", "s2", "s3"})

	assert.Equal(t, 1, len(query.Filters))
	assert.Equal(t, "session_id", query.Filters[0].Field)
	assert.Equal(t, constants.OP_IN, query.Filters[0].Operator)
	values := query.Filters[0].Value.([]any)
	assert.Equal(t, 3, len(values))
	assert.Equal(t, "s1", values[0])
	assert.Equal(t, "s2", values[1])
	assert.Equal(t, "s3", values[2])
}

// TestAddEqOrInFilterInt64Slice int64 切片多值
func TestAddEqOrInFilterInt64Slice(t *testing.T) {
	query := NewQuery()
	query.AddEqOrInFilter("id", []int64{100, 200})

	assert.Equal(t, 1, len(query.Filters))
	assert.Equal(t, constants.OP_IN, query.Filters[0].Operator)
	values := query.Filters[0].Value.([]any)
	assert.Equal(t, 2, len(values))
	assert.Equal(t, int64(100), values[0])
	assert.Equal(t, int64(200), values[1])
}

// TestAddEqOrInFilterChain 链式调用：多值->IN，单值->EQ，空->跳过
func TestAddEqOrInFilterChain(t *testing.T) {
	tests := []struct {
		name        string
		values      any
		wantFilters int
		wantOp      constants.Operator
	}{
		{"多值-IN", []string{"s1", "s2"}, 1, constants.OP_IN},
		{"单值-EQ", []string{"s1"}, 1, constants.OP_EQ},
		{"空切片-跳过", []string{}, 0, ""},
		{"nil-跳过", nil, 0, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := NewQuery()
			query.AddEqOrInFilter("session_id", tt.values)

			assert.Equal(t, tt.wantFilters, len(query.Filters))
			if tt.wantFilters == 1 {
				assert.Equal(t, tt.wantOp, query.Filters[0].Operator)
			}
		})
	}
}
