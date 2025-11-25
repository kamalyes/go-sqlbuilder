/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-26 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-26 00:00:00
 * @FilePath: \go-sqlbuilder\repository\query_filter_test.go
 * @Description: Query 泛型过滤方法测试 - 确保100%测试覆盖率
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package repository

import (
	"github.com/kamalyes/go-sqlbuilder/constants"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// 测试用枚举类型
type TestStatus int32

const (
	TestStatusPending TestStatus = 1
	TestStatusActive  TestStatus = 2
	TestStatusClosed  TestStatus = 3
)

type TestPriority int32

const (
	TestPriorityLow    TestPriority = 1
	TestPriorityMedium TestPriority = 2
	TestPriorityHigh   TestPriority = 3
)

// TestAddFilterIfNotEmpty_String 测试字符串过滤
func TestAddFilterIfNotEmpty_String(t *testing.T) {
	// 非空字符串
	query := NewQuery()
	query.AddFilterIfNotEmpty("name", "test")
	assert.Equal(t, 1, len(query.Filters))
	assert.Equal(t, "name", query.Filters[0].Field)
	assert.Equal(t, constants.OP_EQ, query.Filters[0].Operator)
	assert.Equal(t, "test", query.Filters[0].Value)

	// 空字符串 - 不应添加过滤条件
	query = NewQuery()
	query.AddFilterIfNotEmpty("name", "")
	assert.Equal(t, 0, len(query.Filters))

	// nil 值 - 不应添加过滤条件
	query = NewQuery()
	query.AddFilterIfNotEmpty("name", nil)
	assert.Equal(t, 0, len(query.Filters))
}

// TestAddFilterIfNotEmpty_StringPointer 测试字符串指针过滤
func TestAddFilterIfNotEmpty_StringPointer(t *testing.T) {
	// 非空字符串指针
	str := "test"
	query := NewQuery()
	query.AddFilterIfNotEmpty("name", &str)
	assert.Equal(t, 1, len(query.Filters))
	assert.Equal(t, "test", query.Filters[0].Value)

	// 空字符串指针
	emptyStr := ""
	query = NewQuery()
	query.AddFilterIfNotEmpty("name", &emptyStr)
	assert.Equal(t, 0, len(query.Filters))

	// nil 指针
	var nilStr *string
	query = NewQuery()
	query.AddFilterIfNotEmpty("name", nilStr)
	assert.Equal(t, 0, len(query.Filters))
}

// TestAddFilterIfNotEmpty_StringSlice 测试字符串切片过滤
func TestAddFilterIfNotEmpty_StringSlice(t *testing.T) {
	// 非空字符串切片
	query := NewQuery()
	query.AddFilterIfNotEmpty("status", []string{"active", "pending"})
	assert.Equal(t, 1, len(query.Filters))
	assert.Equal(t, "status", query.Filters[0].Field)
	assert.Equal(t, constants.OP_IN, query.Filters[0].Operator)
	values := query.Filters[0].Value.([]interface{})
	assert.Equal(t, 2, len(values))
	assert.Equal(t, "active", values[0])
	assert.Equal(t, "pending", values[1])

	// 空切片 - 不应添加过滤条件
	query = NewQuery()
	query.AddFilterIfNotEmpty("status", []string{})
	assert.Equal(t, 0, len(query.Filters))

	// nil 切片 - 不应添加过滤条件
	query = NewQuery()
	var nilSlice []string
	query.AddFilterIfNotEmpty("status", nilSlice)
	assert.Equal(t, 0, len(query.Filters))
}

// TestAddFilterIfNotEmpty_IntSlice 测试 int 切片过滤
func TestAddFilterIfNotEmpty_IntSlice(t *testing.T) {
	// 非空 int 切片
	query := NewQuery()
	query.AddFilterIfNotEmpty("age", []int{20, 30, 40})
	assert.Equal(t, 1, len(query.Filters))
	values := query.Filters[0].Value.([]interface{})
	assert.Equal(t, 3, len(values))
	assert.Equal(t, 20, values[0])
	assert.Equal(t, 30, values[1])
	assert.Equal(t, 40, values[2])

	// 空 int 切片
	query = NewQuery()
	query.AddFilterIfNotEmpty("age", []int{})
	assert.Equal(t, 0, len(query.Filters))
}

// TestAddFilterIfNotEmpty_Int32Slice 测试 int32 切片过滤
func TestAddFilterIfNotEmpty_Int32Slice(t *testing.T) {
	// 非空 int32 切片
	query := NewQuery()
	query.AddFilterIfNotEmpty("count", []int32{100, 200, 300})
	assert.Equal(t, 1, len(query.Filters))
	values := query.Filters[0].Value.([]interface{})
	assert.Equal(t, 3, len(values))
	assert.Equal(t, int32(100), values[0])
	assert.Equal(t, int32(200), values[1])
	assert.Equal(t, int32(300), values[2])

	// 空 int32 切片
	query = NewQuery()
	query.AddFilterIfNotEmpty("count", []int32{})
	assert.Equal(t, 0, len(query.Filters))
}

// TestAddFilterIfNotEmpty_Int64Slice 测试 int64 切片过滤
func TestAddFilterIfNotEmpty_Int64Slice(t *testing.T) {
	// 非空 int64 切片
	query := NewQuery()
	query.AddFilterIfNotEmpty("id", []int64{1000, 2000, 3000})
	assert.Equal(t, 1, len(query.Filters))
	values := query.Filters[0].Value.([]interface{})
	assert.Equal(t, 3, len(values))
	assert.Equal(t, int64(1000), values[0])
	assert.Equal(t, int64(2000), values[1])
	assert.Equal(t, int64(3000), values[2])

	// 空 int64 切片
	query = NewQuery()
	query.AddFilterIfNotEmpty("id", []int64{})
	assert.Equal(t, 0, len(query.Filters))
}

// TestAddFilterIfNotEmpty_Int 测试单个 int 值过滤
func TestAddFilterIfNotEmpty_Int(t *testing.T) {
	query := NewQuery()
	query.AddFilterIfNotEmpty("age", 25)
	assert.Equal(t, 1, len(query.Filters))
	assert.Equal(t, "age", query.Filters[0].Field)
	assert.Equal(t, constants.OP_EQ, query.Filters[0].Operator)
	assert.Equal(t, 25, query.Filters[0].Value)
}

// TestAddFilterIfNotEmpty_Int32 测试单个 int32 值过滤
func TestAddFilterIfNotEmpty_Int32(t *testing.T) {
	query := NewQuery()
	query.AddFilterIfNotEmpty("count", int32(100))
	assert.Equal(t, 1, len(query.Filters))
	assert.Equal(t, int32(100), query.Filters[0].Value)
}

// TestAddFilterIfNotEmpty_Int64 测试单个 int64 值过滤
func TestAddFilterIfNotEmpty_Int64(t *testing.T) {
	query := NewQuery()
	query.AddFilterIfNotEmpty("id", int64(1000))
	assert.Equal(t, 1, len(query.Filters))
	assert.Equal(t, int64(1000), query.Filters[0].Value)
}

// TestAddFilterIfNotEmpty_Uint 测试 uint 类型过滤
func TestAddFilterIfNotEmpty_Uint(t *testing.T) {
	query := NewQuery()
	query.AddFilterIfNotEmpty("count", uint(100))
	assert.Equal(t, 1, len(query.Filters))
	assert.Equal(t, uint(100), query.Filters[0].Value)
}

// TestAddFilterIfNotEmpty_Uint32 测试 uint32 类型过滤
func TestAddFilterIfNotEmpty_Uint32(t *testing.T) {
	query := NewQuery()
	query.AddFilterIfNotEmpty("count", uint32(100))
	assert.Equal(t, 1, len(query.Filters))
	assert.Equal(t, uint32(100), query.Filters[0].Value)
}

// TestAddFilterIfNotEmpty_Uint64 测试 uint64 类型过滤
func TestAddFilterIfNotEmpty_Uint64(t *testing.T) {
	query := NewQuery()
	query.AddFilterIfNotEmpty("id", uint64(1000))
	assert.Equal(t, 1, len(query.Filters))
	assert.Equal(t, uint64(1000), query.Filters[0].Value)
}

// TestAddFilterIfNotEmpty_Bool 测试布尔值过滤
func TestAddFilterIfNotEmpty_Bool(t *testing.T) {
	// true
	query := NewQuery()
	query.AddFilterIfNotEmpty("is_active", true)
	assert.Equal(t, 1, len(query.Filters))
	assert.Equal(t, true, query.Filters[0].Value)

	// false (false 也应该添加，因为它是有效值)
	query = NewQuery()
	query.AddFilterIfNotEmpty("is_deleted", false)
	assert.Equal(t, 1, len(query.Filters))
	assert.Equal(t, false, query.Filters[0].Value)
}

// TestAddFilterIfNotEmpty_EnumSlice 测试枚举切片过滤
func TestAddFilterIfNotEmpty_EnumSlice(t *testing.T) {
	// 非空枚举切片
	query := NewQuery()
	statuses := []TestStatus{TestStatusPending, TestStatusActive}
	query.AddFilterIfNotEmpty("status", statuses)
	assert.Equal(t, 1, len(query.Filters))
	values := query.Filters[0].Value.([]interface{})
	assert.Equal(t, 2, len(values))
	assert.Equal(t, TestStatusPending, values[0])
	assert.Equal(t, TestStatusActive, values[1])

	// 空枚举切片
	query = NewQuery()
	query.AddFilterIfNotEmpty("status", []TestStatus{})
	assert.Equal(t, 0, len(query.Filters))
}

// TestAddFilterIfNotEmpty_EnumValue 测试单个枚举值过滤
func TestAddFilterIfNotEmpty_EnumValue(t *testing.T) {
	query := NewQuery()
	query.AddFilterIfNotEmpty("status", TestStatusActive)
	assert.Equal(t, 1, len(query.Filters))
	assert.Equal(t, "status", query.Filters[0].Field)
	assert.Equal(t, constants.OP_EQ, query.Filters[0].Operator)
	assert.Equal(t, TestStatusActive, query.Filters[0].Value)
}

// TestAddFilterIfNotEmpty_ChainCall 测试链式调用
func TestAddFilterIfNotEmpty_ChainCall(t *testing.T) {
	query := NewQuery()
	result := query.
		AddFilterIfNotEmpty("name", "test").
		AddFilterIfNotEmpty("age", 25).
		AddFilterIfNotEmpty("status", []string{"active", "pending"}).
		AddFilterIfNotEmpty("empty", "").           // 应该被忽略
		AddFilterIfNotEmpty("nil", nil).            // 应该被忽略
		AddFilterIfNotEmpty("empty_slice", []int{}) // 应该被忽略

	assert.Equal(t, query, result) // 验证返回的是同一个对象
	assert.Equal(t, 3, len(query.Filters))
}

// TestAddLikeFilterIfNotEmpty 测试 LIKE 过滤
func TestAddLikeFilterIfNotEmpty(t *testing.T) {
	// 非空关键词
	query := NewQuery()
	query.AddLikeFilterIfNotEmpty("name", "test")
	assert.Equal(t, 1, len(query.Filters))
	assert.Equal(t, "name", query.Filters[0].Field)
	assert.Equal(t, constants.OP_LIKE, query.Filters[0].Operator)
	assert.Equal(t, "%test%", query.Filters[0].Value)

	// 空关键词 - 不应添加过滤条件
	query = NewQuery()
	query.AddLikeFilterIfNotEmpty("name", "")
	assert.Equal(t, 0, len(query.Filters))

	// 链式调用
	query = NewQuery()
	result := query.AddLikeFilterIfNotEmpty("title", "golang")
	assert.Equal(t, query, result)
	assert.Equal(t, 1, len(query.Filters))
	assert.Equal(t, "%golang%", query.Filters[0].Value)
}

// TestAddTimeRangeFilter 测试时间范围过滤
func TestAddTimeRangeFilter(t *testing.T) {
	startTime := time.Now().Add(-24 * time.Hour)
	endTime := time.Now()

	// 添加完整的时间范围
	query := NewQuery()
	query.AddTimeRangeFilter("created_at", startTime, endTime)
	assert.Equal(t, 2, len(query.Filters))

	// 第一个过滤条件：开始时间 (>=)
	assert.Equal(t, "created_at", query.Filters[0].Field)
	assert.Equal(t, constants.OP_GTE, query.Filters[0].Operator)
	assert.Equal(t, startTime, query.Filters[0].Value)

	// 第二个过滤条件：结束时间 (<=)
	assert.Equal(t, "created_at", query.Filters[1].Field)
	assert.Equal(t, constants.OP_LTE, query.Filters[1].Operator)
	assert.Equal(t, endTime, query.Filters[1].Value)

	// 只有开始时间
	query = NewQuery()
	query.AddTimeRangeFilter("created_at", startTime, nil)
	assert.Equal(t, 1, len(query.Filters))
	assert.Equal(t, constants.OP_GTE, query.Filters[0].Operator)

	// 只有结束时间
	query = NewQuery()
	query.AddTimeRangeFilter("created_at", nil, endTime)
	assert.Equal(t, 1, len(query.Filters))
	assert.Equal(t, constants.OP_LTE, query.Filters[0].Operator)

	// 都为 nil - 不应添加过滤条件
	query = NewQuery()
	query.AddTimeRangeFilter("created_at", nil, nil)
	assert.Equal(t, 0, len(query.Filters))

	// 链式调用
	query = NewQuery()
	result := query.AddTimeRangeFilter("updated_at", startTime, endTime)
	assert.Equal(t, query, result)
	assert.Equal(t, 2, len(query.Filters))
}

// TestAddTimeRangeFilter_TimePointer 测试时间指针
func TestAddTimeRangeFilter_TimePointer(t *testing.T) {
	startTime := time.Now().Add(-24 * time.Hour)
	endTime := time.Now()

	query := NewQuery()
	query.AddTimeRangeFilter("created_at", &startTime, &endTime)
	assert.Equal(t, 2, len(query.Filters))
	assert.Equal(t, &startTime, query.Filters[0].Value)
	assert.Equal(t, &endTime, query.Filters[1].Value)
}

// TestAddInFilterIfNotEmpty 测试 IN 过滤条件
func TestAddInFilterIfNotEmpty(t *testing.T) {
	// 字符串切片
	query := NewQuery()
	query.AddInFilterIfNotEmpty("status", []string{"active", "pending"})
	assert.Equal(t, 1, len(query.Filters))
	assert.Equal(t, constants.OP_IN, query.Filters[0].Operator)
	values := query.Filters[0].Value.([]interface{})
	assert.Equal(t, 2, len(values))

	// int 切片
	query = NewQuery()
	query.AddInFilterIfNotEmpty("id", []int{1, 2, 3})
	assert.Equal(t, 1, len(query.Filters))
	values = query.Filters[0].Value.([]interface{})
	assert.Equal(t, 3, len(values))

	// 枚举切片
	query = NewQuery()
	priorities := []TestPriority{TestPriorityLow, TestPriorityHigh}
	query.AddInFilterIfNotEmpty("priority", priorities)
	assert.Equal(t, 1, len(query.Filters))
	values = query.Filters[0].Value.([]interface{})
	assert.Equal(t, 2, len(values))

	// 空切片 - 不应添加过滤条件
	query = NewQuery()
	query.AddInFilterIfNotEmpty("status", []string{})
	assert.Equal(t, 0, len(query.Filters))

	// nil 值 - 不应添加过滤条件
	query = NewQuery()
	query.AddInFilterIfNotEmpty("status", nil)
	assert.Equal(t, 0, len(query.Filters))

	// 链式调用
	query = NewQuery()
	result := query.AddInFilterIfNotEmpty("tags", []string{"go", "rust"})
	assert.Equal(t, query, result)
	assert.Equal(t, 1, len(query.Filters))
}

// TestAddInFilterIfNotEmpty_NonSliceValue 测试非切片值
func TestAddInFilterIfNotEmpty_NonSliceValue(t *testing.T) {
	// 传入非切片值 - 不应添加过滤条件
	query := NewQuery()
	query.AddInFilterIfNotEmpty("status", "active")
	assert.Equal(t, 0, len(query.Filters))

	query = NewQuery()
	query.AddInFilterIfNotEmpty("count", 100)
	assert.Equal(t, 0, len(query.Filters))
}

// TestCombinedFilterMethods 测试组合使用多个过滤方法
func TestCombinedFilterMethods(t *testing.T) {
	startTime := time.Now().Add(-7 * 24 * time.Hour)
	endTime := time.Now()

	query := NewQuery()
	query.
		AddFilterIfNotEmpty("category", "tech").
		AddFilterIfNotEmpty("status", []TestStatus{TestStatusActive, TestStatusPending}).
		AddLikeFilterIfNotEmpty("title", "golang").
		AddTimeRangeFilter("created_at", startTime, endTime).
		AddInFilterIfNotEmpty("tags", []string{"backend", "database"})

	// 验证所有过滤条件都正确添加
	assert.Equal(t, 6, len(query.Filters))

	// category = 'tech'
	assert.Equal(t, "category", query.Filters[0].Field)
	assert.Equal(t, constants.OP_EQ, query.Filters[0].Operator)

	// status IN (...)
	assert.Equal(t, "status", query.Filters[1].Field)
	assert.Equal(t, constants.OP_IN, query.Filters[1].Operator)

	// title LIKE '%golang%'
	assert.Equal(t, "title", query.Filters[2].Field)
	assert.Equal(t, constants.OP_LIKE, query.Filters[2].Operator)

	// created_at >= startTime
	assert.Equal(t, "created_at", query.Filters[3].Field)
	assert.Equal(t, constants.OP_GTE, query.Filters[3].Operator)

	// created_at <= endTime
	assert.Equal(t, "created_at", query.Filters[4].Field)
	assert.Equal(t, constants.OP_LTE, query.Filters[4].Operator)

	// tags IN (...)
	assert.Equal(t, "tags", query.Filters[5].Field)
	assert.Equal(t, constants.OP_IN, query.Filters[5].Operator)
}

// TestEmptyValueScenarios 测试各种空值场景
func TestEmptyValueScenarios(t *testing.T) {
	query := NewQuery()

	// 添加各种空值 - 都不应该添加过滤条件
	query.
		AddFilterIfNotEmpty("str", "").
		AddFilterIfNotEmpty("nil", nil).
		AddFilterIfNotEmpty("slice", []string{}).
		AddFilterIfNotEmpty("int_slice", []int{}).
		AddFilterIfNotEmpty("int32_slice", []int32{}).
		AddFilterIfNotEmpty("int64_slice", []int64{}).
		AddLikeFilterIfNotEmpty("keyword", "").
		AddTimeRangeFilter("time", nil, nil).
		AddInFilterIfNotEmpty("in", []string{}).
		AddInFilterIfNotEmpty("in_nil", nil)

	// 验证没有添加任何过滤条件
	assert.Equal(t, 0, len(query.Filters))
}

// TestValidValueScenarios 测试各种有效值场景
func TestValidValueScenarios(t *testing.T) {
	query := NewQuery()

	// 添加各种有效值
	query.
		AddFilterIfNotEmpty("str", "value").                                                  // 1
		AddFilterIfNotEmpty("int", 100).                                                      // 2
		AddFilterIfNotEmpty("int32", int32(200)).                                             // 3
		AddFilterIfNotEmpty("int64", int64(300)).                                             // 4
		AddFilterIfNotEmpty("uint", uint(400)).                                               // 5
		AddFilterIfNotEmpty("uint32", uint32(500)).                                           // 6
		AddFilterIfNotEmpty("uint64", uint64(600)).                                           // 7
		AddFilterIfNotEmpty("bool_true", true).                                               // 8
		AddFilterIfNotEmpty("bool_false", false).                                             // 9
		AddFilterIfNotEmpty("str_slice", []string{"a", "b"}).                                 // 10
		AddFilterIfNotEmpty("int_slice", []int{1, 2}).                                        // 11
		AddFilterIfNotEmpty("enum", TestStatusActive).                                        // 12
		AddFilterIfNotEmpty("enum_slice", []TestPriority{TestPriorityLow, TestPriorityHigh}). // 13
		AddLikeFilterIfNotEmpty("keyword", "search").                                         // 14
		AddTimeRangeFilter("time", time.Now(), time.Now()).                                   // 15-16 (2个)
		AddInFilterIfNotEmpty("in", []string{"x", "y"})                                       // 17

	// 验证所有有效值都正确添加
	// AddTimeRangeFilter 添加2个过滤条件 (start >= 和 end <=)
	// 总计: 13 + 1 + 2 + 1 = 17个过滤条件
	assert.Equal(t, 17, len(query.Filters))
}
