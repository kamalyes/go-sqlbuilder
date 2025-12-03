/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-26 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-03 16:36:21
 * @FilePath: \go-sqlbuilder\repository\query_filter_test.go
 * @Description: Query 泛型过滤方法测试 - 确保100%测试覆盖率
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package repository

import (
	"testing"
	"time"

	"github.com/kamalyes/go-sqlbuilder/constants"
	"github.com/stretchr/testify/assert"
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

// =============== AddSafeOrder 排序方法测试 ===============

// TestAddSafeOrder_DefaultValues 测试使用默认值
func TestAddSafeOrder_DefaultValues(t *testing.T) {
	query := NewQuery()
	query.AddSafeOrder("", "", "created_at", "DESC")

	assert.Equal(t, 1, len(query.Orders))
	assert.Equal(t, "created_at", query.Orders[0].Field)
	assert.Equal(t, "DESC", query.Orders[0].Direction)
}

// TestAddSafeOrder_CustomValues 测试自定义排序值
func TestAddSafeOrder_CustomValues(t *testing.T) {
	query := NewQuery()
	query.AddSafeOrder("updated_at", "ASC", "created_at", "DESC")

	assert.Equal(t, 1, len(query.Orders))
	assert.Equal(t, "updated_at", query.Orders[0].Field)
	assert.Equal(t, "ASC", query.Orders[0].Direction)
}

// TestAddSafeOrder_WhitelistValidField 测试白名单 - 有效字段
func TestAddSafeOrder_WhitelistValidField(t *testing.T) {
	allowedFields := []string{"id", "created_at", "updated_at", "name"}
	query := NewQuery()
	query.AddSafeOrder("name", "ASC", "created_at", "DESC", allowedFields)

	assert.Equal(t, 1, len(query.Orders))
	assert.Equal(t, "name", query.Orders[0].Field)
	assert.Equal(t, "ASC", query.Orders[0].Direction)
}

// TestAddSafeOrder_WhitelistInvalidField 测试白名单 - 无效字段(使用默认值)
func TestAddSafeOrder_WhitelistInvalidField(t *testing.T) {
	allowedFields := []string{"id", "created_at", "updated_at"}
	query := NewQuery()
	query.AddSafeOrder("malicious_field", "ASC", "created_at", "DESC", allowedFields)

	assert.Equal(t, 1, len(query.Orders))
	assert.Equal(t, "created_at", query.Orders[0].Field) // 回退到默认字段
	assert.Equal(t, "DESC", query.Orders[0].Direction)   // 使用默认方向
}

// TestAddSafeOrder_SQLInjectionAttempt 测试SQL注入攻击防护
func TestAddSafeOrder_SQLInjectionAttempt(t *testing.T) {
	query := NewQuery()
	// 尝试注入恶意SQL
	query.AddSafeOrder("id; DROP TABLE users--", "DESC", "created_at", "DESC")

	assert.Equal(t, 1, len(query.Orders))
	assert.Equal(t, "created_at", query.Orders[0].Field) // 回退到默认字段
}

// TestAddSafeOrder_InvalidDirection 测试无效排序方向(使用默认值)
func TestAddSafeOrder_InvalidDirection(t *testing.T) {
	query := NewQuery()
	query.AddSafeOrder("id", "INVALID", "created_at", "DESC")

	assert.Equal(t, 1, len(query.Orders))
	assert.Equal(t, "id", query.Orders[0].Field)
	assert.Equal(t, "DESC", query.Orders[0].Direction) // 无效方向,使用默认值
}

// TestAddSafeOrder_LowercaseDirection 测试小写排序方向(自动转大写)
func TestAddSafeOrder_LowercaseDirection(t *testing.T) {
	query := NewQuery()
	query.AddSafeOrder("id", "asc", "created_at", "DESC")

	assert.Equal(t, 1, len(query.Orders))
	assert.Equal(t, "ASC", query.Orders[0].Direction) // 自动转为大写
}

// TestAddSafeOrder_MixedCaseDirection 测试混合大小写排序方向
func TestAddSafeOrder_MixedCaseDirection(t *testing.T) {
	query := NewQuery()
	query.AddSafeOrder("id", "DeSc", "created_at", "ASC")

	assert.Equal(t, 1, len(query.Orders))
	assert.Equal(t, "DESC", query.Orders[0].Direction) // 标准化为大写
}

// TestAddSafeOrder_EmptyWhitelist 测试空白名单(使用字段名安全检查)
func TestAddSafeOrder_EmptyWhitelist(t *testing.T) {
	query := NewQuery()
	query.AddSafeOrder("valid_field_123", "ASC", "created_at", "DESC", []string{})

	assert.Equal(t, 1, len(query.Orders))
	assert.Equal(t, "valid_field_123", query.Orders[0].Field) // 安全字段名,允许通过
}

// TestAddSafeOrder_DotNotation 测试点号表示法(表名.字段名)
func TestAddSafeOrder_DotNotation(t *testing.T) {
	query := NewQuery()
	query.AddSafeOrder("users.created_at", "ASC", "id", "DESC")

	assert.Equal(t, 1, len(query.Orders))
	assert.Equal(t, "users.created_at", query.Orders[0].Field) // 允许表名.字段名格式
}

// TestAddSafeOrder_ChainCalls 测试链式调用
func TestAddSafeOrder_ChainCalls(t *testing.T) {
	query := NewQuery()
	query.AddSafeOrder("name", "ASC", "id", "DESC").
		AddSafeOrder("created_at", "DESC", "updated_at", "ASC")

	assert.Equal(t, 2, len(query.Orders))
	assert.Equal(t, "name", query.Orders[0].Field)
	assert.Equal(t, "ASC", query.Orders[0].Direction)
	assert.Equal(t, "created_at", query.Orders[1].Field)
	assert.Equal(t, "DESC", query.Orders[1].Direction)
}

// TestAddSafeOrder_SpecialCharactersBlocked 测试特殊字符被阻止
func TestAddSafeOrder_SpecialCharactersBlocked(t *testing.T) {
	testCases := []struct {
		name      string
		sortBy    string
		expectUse string // "default" 表示应使用默认字段
	}{
		{"空格", "created at", "default"},
		{"单引号", "id'OR'1'='1", "default"},
		{"双引号", "id\"OR\"1\"=\"1", "default"},
		{"反引号", "id`", "default"},
		{"括号", "id()", "default"},
		{"星号", "id*", "default"},
		{"逗号", "id,name", "default"},
		{"分号", "id;DROP TABLE", "default"},
		{"减号", "id-name", "default"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			query := NewQuery()
			query.AddSafeOrder(tc.sortBy, "ASC", "created_at", "DESC")

			assert.Equal(t, 1, len(query.Orders))
			if tc.expectUse == "default" {
				assert.Equal(t, "created_at", query.Orders[0].Field)
			}
		})
	}
}

// Test_isSafeFieldName 测试字段名安全检查函数
func Test_isSafeFieldName(t *testing.T) {
	testCases := []struct {
		name     string
		field    string
		expected bool
	}{
		{"空字符串", "", false},
		{"简单字段名", "id", true},
		{"下划线字段名", "user_id", true},
		{"数字结尾", "field123", true},
		{"大写字母", "UserId", true},
		{"点号表示法", "users.id", true},
		{"包含空格", "user id", false},
		{"包含单引号", "user'id", false},
		{"包含分号", "id;DROP", false},
		{"包含星号", "id*", false},
		{"包含减号", "user-id", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := isSafeFieldName(tc.field)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestAddTimeRangeValidFilter(t *testing.T) {
	t.Log("测试 AddTimeRangeFilter 的零值处理...")

	// 测试1: nil 值
	t.Run("测试 nil 值", func(t *testing.T) {
		query := NewQuery()
		query.AddTimeRangeFilter("created_at", nil, nil)
		assert.Equal(t, 0, len(query.Filters), "nil值不应该添加任何过滤条件")
	})

	// 测试2: 零值时间
	t.Run("测试零值时间", func(t *testing.T) {
		query := NewQuery()
		zeroTime := time.Time{}
		query.AddTimeRangeFilter("created_at", zeroTime, zeroTime)
		assert.Equal(t, 0, len(query.Filters), "零值时间不应该添加任何过滤条件")
	})

	// 测试3: 零值时间指针
	t.Run("测试零值时间指针", func(t *testing.T) {
		query := NewQuery()
		zeroTimePtr := &time.Time{}
		query.AddTimeRangeFilter("created_at", zeroTimePtr, zeroTimePtr)
		assert.Equal(t, 0, len(query.Filters), "零值时间指针不应该添加任何过滤条件")
	})

	// 测试4: 有效时间
	t.Run("测试有效时间", func(t *testing.T) {
		query := NewQuery()
		now := time.Now()
		query.AddTimeRangeFilter("created_at", &now, &now)
		assert.Equal(t, 2, len(query.Filters), "有效时间应该添加2个过滤条件(>= 和 <=)")
	})

	// 测试5: 混合情况 - 只有开始时间
	t.Run("测试只有开始时间", func(t *testing.T) {
		query := NewQuery()
		now := time.Now()
		query.AddTimeRangeFilter("created_at", &now, nil)
		assert.Equal(t, 1, len(query.Filters), "只有开始时间应该添加1个过滤条件(>=)")
	})

	// 测试6: 混合情况 - 只有结束时间
	t.Run("测试只有结束时间", func(t *testing.T) {
		query := NewQuery()
		now := time.Now()
		query.AddTimeRangeFilter("created_at", nil, &now)
		assert.Equal(t, 1, len(query.Filters), "只有结束时间应该添加1个过滤条件(<=)")
	})

	// 测试7: 验证过滤条件的操作符
	t.Run("测试过滤条件的操作符", func(t *testing.T) {
		query := NewQuery()
		now := time.Now()
		query.AddTimeRangeFilter("created_at", &now, &now)
		if len(query.Filters) == 2 {
			// 第一个应该是 >=，第二个应该是 <=
			assert.Equal(t, constants.OP_GTE, query.Filters[0].Operator, "第一个过滤条件应该是 >=")
			assert.Equal(t, constants.OP_LTE, query.Filters[1].Operator, "第二个过滤条件应该是 <=")
			assert.Equal(t, "created_at", query.Filters[0].Field, "过滤条件字段应该是 created_at")
			assert.Equal(t, "created_at", query.Filters[1].Field, "过滤条件字段应该是 created_at")
		} else {
			t.Errorf("❌ FAIL: 过滤条件数量不正确，无法验证操作符")
		}
	})

	t.Log("测试完成！")
}

func TestSetPagination(t *testing.T) {
	query := NewQuery()

	// 测试正常分页
	result := query.SetPagination(2, 20)

	assert.NotNil(t, result.Pagination)
	assert.Equal(t, int32(2), result.Pagination.Page)
	assert.Equal(t, int32(20), result.Pagination.PageSize)

	// 测试链式调用
	assert.Equal(t, query, result, "SetPagination should return the same Query instance for chaining")
}

func TestSetPagination_DefaultValues(t *testing.T) {
	query := NewQuery()

	// 测试零值或负值的处理
	query.SetPagination(0, 0)

	assert.NotNil(t, query.Pagination)
	assert.Equal(t, int32(1), query.Pagination.Page, "Page should default to 1 when <= 0")
	assert.Equal(t, int32(10), query.Pagination.PageSize, "PageSize should default to 10 when <= 0")

	// 测试负值
	query2 := NewQuery()
	query2.SetPagination(-5, -10)

	assert.Equal(t, int32(1), query2.Pagination.Page, "Page should default to 1 when negative")
	assert.Equal(t, int32(10), query2.Pagination.PageSize, "PageSize should default to 10 when negative")
}

func TestAddRawOrder(t *testing.T) {
	query := NewQuery()

	// 测试添加原始排序表达式
	result := query.AddRawOrder("created_at DESC, updated_at ASC")

	assert.Len(t, result.Orders, 1)
	assert.Equal(t, "created_at DESC, updated_at ASC", result.Orders[0].Field)
	assert.Equal(t, "", result.Orders[0].Direction, "Direction should be empty for raw order")

	// 测试链式调用
	assert.Equal(t, query, result, "AddRawOrder should return the same Query instance for chaining")
}

func TestAddRawOrder_EmptyExpression(t *testing.T) {
	query := NewQuery()

	// 测试空表达式不会添加排序
	result := query.AddRawOrder("")

	assert.Len(t, result.Orders, 0, "Empty order expression should not add any order")
}

func TestAddRawOrder_MultipleOrders(t *testing.T) {
	query := NewQuery()

	// 测试多个原始排序
	query.AddRawOrder("priority DESC").
		AddRawOrder("CASE WHEN status = 'urgent' THEN 1 ELSE 2 END").
		AddOrder("id", "ASC")

	assert.Len(t, query.Orders, 3, "Should have 3 order conditions")

	// 验证第一个原始排序
	assert.Equal(t, "priority DESC", query.Orders[0].Field)
	assert.Equal(t, "", query.Orders[0].Direction)

	// 验证第二个原始排序（复杂表达式）
	assert.Equal(t, "CASE WHEN status = 'urgent' THEN 1 ELSE 2 END", query.Orders[1].Field)
	assert.Equal(t, "", query.Orders[1].Direction)

	// 验证普通排序仍然正常工作
	assert.Equal(t, "id", query.Orders[2].Field)
	assert.Equal(t, "ASC", query.Orders[2].Direction)
}

func TestSetPagination_Integration_WithOtherMethods(t *testing.T) {
	query := NewQuery()

	// 测试与其他方法的集成
	result := query.
		AddFilter(NewEqFilter("status", "active")).
		SetPagination(3, 25).
		AddRawOrder("created_at DESC").
		AddOrder("name", "ASC")

	// 验证过滤器
	assert.Len(t, result.Filters, 1)
	assert.Equal(t, "status", result.Filters[0].Field)

	// 验证分页
	assert.NotNil(t, result.Pagination)
	assert.Equal(t, int32(3), result.Pagination.Page)
	assert.Equal(t, int32(25), result.Pagination.PageSize)

	// 验证排序
	assert.Len(t, result.Orders, 2)
	assert.Equal(t, "created_at DESC", result.Orders[0].Field)
	assert.Equal(t, "name", result.Orders[1].Field)
	assert.Equal(t, "ASC", result.Orders[1].Direction)
}
