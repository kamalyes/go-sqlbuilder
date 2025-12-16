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

// TestAddFilterIfNotEmptyString 测试字符串过滤
func TestAddFilterIfNotEmptyString(t *testing.T) {
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

// TestAddFilterIfNotEmptyStringPointer 测试字符串指针过滤
func TestAddFilterIfNotEmptyStringPointer(t *testing.T) {
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

// TestAddFilterIfNotEmptyStringSlice 测试字符串切片过滤
func TestAddFilterIfNotEmptyStringSlice(t *testing.T) {
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

// TestAddFilterIfNotEmptyIntSlice 测试 int 切片过滤
func TestAddFilterIfNotEmptyIntSlice(t *testing.T) {
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

// TestAddFilterIfNotEmptyInt32Slice 测试 int32 切片过滤
func TestAddFilterIfNotEmptyInt32Slice(t *testing.T) {
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

// TestAddFilterIfNotEmptyInt64Slice 测试 int64 切片过滤
func TestAddFilterIfNotEmptyInt64Slice(t *testing.T) {
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

// TestAddFilterIfNotEmptyInt 测试单个 int 值过滤
func TestAddFilterIfNotEmptyInt(t *testing.T) {
	query := NewQuery()
	query.AddFilterIfNotEmpty("age", 25)
	assert.Equal(t, 1, len(query.Filters))
	assert.Equal(t, "age", query.Filters[0].Field)
	assert.Equal(t, constants.OP_EQ, query.Filters[0].Operator)
	assert.Equal(t, 25, query.Filters[0].Value)
}

// TestAddFilterIfNotEmptyInt32 测试单个 int32 值过滤
func TestAddFilterIfNotEmptyInt32(t *testing.T) {
	query := NewQuery()
	query.AddFilterIfNotEmpty("count", int32(100))
	assert.Equal(t, 1, len(query.Filters))
	assert.Equal(t, int32(100), query.Filters[0].Value)
}

// TestAddFilterIfNotEmptyInt64 测试单个 int64 值过滤
func TestAddFilterIfNotEmptyInt64(t *testing.T) {
	query := NewQuery()
	query.AddFilterIfNotEmpty("id", int64(1000))
	assert.Equal(t, 1, len(query.Filters))
	assert.Equal(t, int64(1000), query.Filters[0].Value)
}

// TestAddFilterIfNotEmptyUint 测试 uint 类型过滤
func TestAddFilterIfNotEmptyUint(t *testing.T) {
	query := NewQuery()
	query.AddFilterIfNotEmpty("count", uint(100))
	assert.Equal(t, 1, len(query.Filters))
	assert.Equal(t, uint(100), query.Filters[0].Value)
}

// TestAddFilterIfNotEmptyUint32 测试 uint32 类型过滤
func TestAddFilterIfNotEmptyUint32(t *testing.T) {
	query := NewQuery()
	query.AddFilterIfNotEmpty("count", uint32(100))
	assert.Equal(t, 1, len(query.Filters))
	assert.Equal(t, uint32(100), query.Filters[0].Value)
}

// TestAddFilterIfNotEmptyUint64 测试 uint64 类型过滤
func TestAddFilterIfNotEmptyUint64(t *testing.T) {
	query := NewQuery()
	query.AddFilterIfNotEmpty("id", uint64(1000))
	assert.Equal(t, 1, len(query.Filters))
	assert.Equal(t, uint64(1000), query.Filters[0].Value)
}

// TestAddFilterIfNotEmptyBool 测试布尔值过滤
func TestAddFilterIfNotEmptyBool(t *testing.T) {
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

// TestAddFilterIfNotEmptyEnumSlice 测试枚举切片过滤
func TestAddFilterIfNotEmptyEnumSlice(t *testing.T) {
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

// TestAddFilterIfNotEmptyEnumValue 测试单个枚举值过滤
func TestAddFilterIfNotEmptyEnumValue(t *testing.T) {
	query := NewQuery()
	query.AddFilterIfNotEmpty("status", TestStatusActive)
	assert.Equal(t, 1, len(query.Filters))
	assert.Equal(t, "status", query.Filters[0].Field)
	assert.Equal(t, constants.OP_EQ, query.Filters[0].Operator)
	assert.Equal(t, TestStatusActive, query.Filters[0].Value)
}

// TestAddFilterIfNotEmptyChainCall 测试链式调用
func TestAddFilterIfNotEmptyChainCall(t *testing.T) {
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

// TestAddTimeRangeFilterTimePointer 测试时间指针
func TestAddTimeRangeFilterTimePointer(t *testing.T) {
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

// TestAddInFilterIfNotEmptyNonSliceValue 测试非切片值
func TestAddInFilterIfNotEmptyNonSliceValue(t *testing.T) {
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

// TestAddSafeOrderDefaultValues 测试使用默认值
func TestAddSafeOrderDefaultValues(t *testing.T) {
	query := NewQuery()
	query.AddSafeOrder("", "", "created_at", "DESC")

	assert.Equal(t, 1, len(query.Orders))
	assert.Equal(t, "created_at", query.Orders[0].Field)
	assert.Equal(t, "DESC", query.Orders[0].Direction)
}

// TestAddSafeOrderCustomValues 测试自定义排序值
func TestAddSafeOrderCustomValues(t *testing.T) {
	query := NewQuery()
	query.AddSafeOrder("updated_at", "ASC", "created_at", "DESC")

	assert.Equal(t, 1, len(query.Orders))
	assert.Equal(t, "updated_at", query.Orders[0].Field)
	assert.Equal(t, "ASC", query.Orders[0].Direction)
}

// TestAddSafeOrderWhitelistValidField 测试白名单 - 有效字段
func TestAddSafeOrderWhitelistValidField(t *testing.T) {
	allowedFields := []string{"id", "created_at", "updated_at", "name"}
	query := NewQuery()
	query.AddSafeOrder("name", "ASC", "created_at", "DESC", allowedFields)

	assert.Equal(t, 1, len(query.Orders))
	assert.Equal(t, "name", query.Orders[0].Field)
	assert.Equal(t, "ASC", query.Orders[0].Direction)
}

// TestAddSafeOrderWhitelistInvalidField 测试白名单 - 无效字段(使用默认值)
func TestAddSafeOrderWhitelistInvalidField(t *testing.T) {
	allowedFields := []string{"id", "created_at", "updated_at"}
	query := NewQuery()
	query.AddSafeOrder("malicious_field", "ASC", "created_at", "DESC", allowedFields)

	assert.Equal(t, 1, len(query.Orders))
	assert.Equal(t, "created_at", query.Orders[0].Field) // 回退到默认字段
	assert.Equal(t, "DESC", query.Orders[0].Direction)   // 使用默认方向
}

// TestAddSafeOrderSQLInjectionAttempt 测试SQL注入攻击防护
func TestAddSafeOrderSQLInjectionAttempt(t *testing.T) {
	query := NewQuery()
	// 尝试注入恶意SQL
	query.AddSafeOrder("id; DROP TABLE users--", "DESC", "created_at", "DESC")

	assert.Equal(t, 1, len(query.Orders))
	assert.Equal(t, "created_at", query.Orders[0].Field) // 回退到默认字段
}

// TestAddSafeOrderInvalidDirection 测试无效排序方向(使用默认值)
func TestAddSafeOrderInvalidDirection(t *testing.T) {
	query := NewQuery()
	query.AddSafeOrder("id", "INVALID", "created_at", "DESC")

	assert.Equal(t, 1, len(query.Orders))
	assert.Equal(t, "id", query.Orders[0].Field)
	assert.Equal(t, "DESC", query.Orders[0].Direction) // 无效方向,使用默认值
}

// TestAddSafeOrderLowercaseDirection 测试小写排序方向(自动转大写)
func TestAddSafeOrderLowercaseDirection(t *testing.T) {
	query := NewQuery()
	query.AddSafeOrder("id", "asc", "created_at", "DESC")

	assert.Equal(t, 1, len(query.Orders))
	assert.Equal(t, "ASC", query.Orders[0].Direction) // 自动转为大写
}

// TestAddSafeOrderMixedCaseDirection 测试混合大小写排序方向
func TestAddSafeOrderMixedCaseDirection(t *testing.T) {
	query := NewQuery()
	query.AddSafeOrder("id", "DeSc", "created_at", "ASC")

	assert.Equal(t, 1, len(query.Orders))
	assert.Equal(t, "DESC", query.Orders[0].Direction) // 标准化为大写
}

// TestAddSafeOrderEmptyWhitelist 测试空白名单(使用字段名安全检查)
func TestAddSafeOrderEmptyWhitelist(t *testing.T) {
	query := NewQuery()
	query.AddSafeOrder("valid_field_123", "ASC", "created_at", "DESC", []string{})

	assert.Equal(t, 1, len(query.Orders))
	assert.Equal(t, "valid_field_123", query.Orders[0].Field) // 安全字段名,允许通过
}

// TestAddSafeOrderDotNotation 测试点号表示法(表名.字段名)
func TestAddSafeOrderDotNotation(t *testing.T) {
	query := NewQuery()
	query.AddSafeOrder("users.created_at", "ASC", "id", "DESC")

	assert.Equal(t, 1, len(query.Orders))
	assert.Equal(t, "users.created_at", query.Orders[0].Field) // 允许表名.字段名格式
}

// TestAddSafeOrderChainCalls 测试链式调用
func TestAddSafeOrderChainCalls(t *testing.T) {
	query := NewQuery()
	query.AddSafeOrder("name", "ASC", "id", "DESC").
		AddSafeOrder("created_at", "DESC", "updated_at", "ASC")

	assert.Equal(t, 2, len(query.Orders))
	assert.Equal(t, "name", query.Orders[0].Field)
	assert.Equal(t, "ASC", query.Orders[0].Direction)
	assert.Equal(t, "created_at", query.Orders[1].Field)
	assert.Equal(t, "DESC", query.Orders[1].Direction)
}

// TestAddSafeOrderSpecialCharactersBlocked 测试特殊字符被阻止
func TestAddSafeOrderSpecialCharactersBlocked(t *testing.T) {
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

func TestSetPaginationDefaultValues(t *testing.T) {
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

func TestAddRawOrderEmptyExpression(t *testing.T) {
	query := NewQuery()

	// 测试空表达式不会添加排序
	result := query.AddRawOrder("")

	assert.Len(t, result.Orders, 0, "Empty order expression should not add any order")
}

func TestAddRawOrderMultipleOrders(t *testing.T) {
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

func TestSetPaginationIntegrationWithOtherMethods(t *testing.T) {
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

// TestBuildWhereClause 测试 BuildWhereClause 方法
func TestBuildWhereClause(t *testing.T) {
	t.Run("空查询", func(t *testing.T) {
		var query *Query
		whereClause, args := query.BuildWhereClause()
		assert.Empty(t, whereClause)
		assert.Empty(t, args)
	})

	t.Run("只有简单过滤条件", func(t *testing.T) {
		query := NewQuery().
			AddFilter(NewEqFilter("agent_id", "user123")).
			AddFilter(NewEqFilter("work_status", 2))

		whereClause, args := query.BuildWhereClause()
		expected := "agent_id = ? AND work_status = ?"
		assert.Equal(t, expected, whereClause)
		assert.Equal(t, []interface{}{"user123", 2}, args)
	})

	t.Run("包含BETWEEN操作", func(t *testing.T) {
		query := NewQuery().
			AddFilter(NewEqFilter("agent_id", "user123")).
			AddFilter(&Filter{
				Field:    "created_at",
				Operator: constants.OP_BETWEEN,
				Value:    []interface{}{"2023-01-01", "2023-12-31"},
			})

		whereClause, args := query.BuildWhereClause()
		expected := "agent_id = ? AND created_at BETWEEN ? AND ?"
		assert.Equal(t, expected, whereClause)
		assert.Equal(t, []interface{}{"user123", "2023-01-01", "2023-12-31"}, args)
	})

	t.Run("包含LIKE操作", func(t *testing.T) {
		query := NewQuery().
			AddFilter(&Filter{
				Field:    "name",
				Operator: constants.OP_STARTS_WITH,
				Value:    "test",
			}).
			AddFilter(&Filter{
				Field:    "email",
				Operator: constants.OP_CONTAINS,
				Value:    "example",
			})

		whereClause, args := query.BuildWhereClause()
		expected := "name LIKE ? AND email LIKE ?"
		assert.Equal(t, expected, whereClause)
		assert.Equal(t, []interface{}{"test%", "%example%"}, args)
	})

	t.Run("包含NULL操作", func(t *testing.T) {
		query := NewQuery().
			AddFilter(&Filter{
				Field:    "deleted_at",
				Operator: constants.OP_IS_NULL,
				Value:    nil,
			}).
			AddFilter(&Filter{
				Field:    "updated_at",
				Operator: constants.OP_IS_NOT_NULL,
				Value:    nil,
			})

		whereClause, args := query.BuildWhereClause()
		expected := "deleted_at IS NULL AND updated_at IS NOT NULL"
		assert.Equal(t, expected, whereClause)
		assert.Empty(t, args)
	})
}

// TestBuildWhereClauseFilterGroup 测试过滤条件组
func TestBuildWhereClauseFilterGroup(t *testing.T) {
	t.Run("AND条件组", func(t *testing.T) {
		group := NewFilterGroup(constants.LOGIC_AND).
			AddFilter(NewEqFilter("status", "active")).
			AddFilter(NewGtFilter("age", 18))

		query := NewQuery().WithFilterGroup(group)
		whereClause, args := query.BuildWhereClause()
		expected := "(status = ? AND age > ?)"
		assert.Equal(t, expected, whereClause)
		assert.Equal(t, []interface{}{"active", 18}, args)
	})

	t.Run("OR条件组", func(t *testing.T) {
		group := NewFilterGroup(constants.LOGIC_OR).
			AddFilter(NewEqFilter("type", "admin")).
			AddFilter(NewEqFilter("type", "manager"))

		query := NewQuery().WithFilterGroup(group)
		whereClause, args := query.BuildWhereClause()
		expected := "(type = ? OR type = ?)"
		assert.Equal(t, expected, whereClause)
		assert.Equal(t, []interface{}{"admin", "manager"}, args)
	})

	t.Run("嵌套条件组", func(t *testing.T) {
		subGroup1 := NewFilterGroup(constants.LOGIC_OR).
			AddFilter(NewEqFilter("status", "active")).
			AddFilter(NewEqFilter("status", "pending"))

		subGroup2 := NewFilterGroup(constants.LOGIC_AND).
			AddFilter(NewGtFilter("age", 21)).
			AddFilter(NewLtFilter("age", 65))

		mainGroup := NewFilterGroup(constants.LOGIC_AND).
			AddGroup(subGroup1).
			AddGroup(subGroup2)

		query := NewQuery().WithFilterGroup(mainGroup)
		whereClause, args := query.BuildWhereClause()
		expected := "((status = ? OR status = ?) AND (age > ? AND age < ?))"
		assert.Equal(t, expected, whereClause)
		assert.Equal(t, []interface{}{"active", "pending", 21, 65}, args)
	})
}

// TestBuildWhereClauseMixedConditions 测试混合条件
func TestBuildWhereClauseMixedConditions(t *testing.T) {
	t.Run("简单过滤条件和条件组", func(t *testing.T) {
		group := NewFilterGroup(constants.LOGIC_OR).
			AddFilter(NewEqFilter("role", "admin")).
			AddFilter(NewEqFilter("role", "manager"))

		query := NewQuery().
			AddFilter(NewEqFilter("company_id", "123")).
			AddFilter(NewEqFilter("active", true)).
			WithFilterGroup(group)

		whereClause, args := query.BuildWhereClause()
		expected := "company_id = ? AND active = ? AND (role = ? OR role = ?)"
		assert.Equal(t, expected, whereClause)
		assert.Equal(t, []interface{}{"123", true, "admin", "manager"}, args)
	})
}

// TestSpecialOperators 测试特殊操作符处理
func TestSpecialOperators(t *testing.T) {
	t.Run("STARTS_WITH", func(t *testing.T) {
		query := NewQuery().AddFilter(&Filter{
			Field:    "name",
			Operator: constants.OP_STARTS_WITH,
			Value:    "John",
		})

		whereClause, args := query.BuildWhereClause()
		expected := "name LIKE ?"
		assert.Equal(t, expected, whereClause)
		assert.Equal(t, []interface{}{"John%"}, args)
	})

	t.Run("ENDS_WITH", func(t *testing.T) {
		query := NewQuery().AddFilter(&Filter{
			Field:    "email",
			Operator: constants.OP_ENDS_WITH,
			Value:    "@example.com",
		})

		whereClause, args := query.BuildWhereClause()
		expected := "email LIKE ?"
		assert.Equal(t, expected, whereClause)
		assert.Equal(t, []interface{}{"%@example.com"}, args)
	})

	t.Run("CONTAINS", func(t *testing.T) {
		query := NewQuery().AddFilter(&Filter{
			Field:    "description",
			Operator: constants.OP_CONTAINS,
			Value:    "test",
		})

		whereClause, args := query.BuildWhereClause()
		expected := "description LIKE ?"
		assert.Equal(t, expected, whereClause)
		assert.Equal(t, []interface{}{"%test%"}, args)
	})

	t.Run("FIND_IN_SET", func(t *testing.T) {
		query := NewQuery().AddFilter(&Filter{
			Field:    "tags",
			Operator: constants.OP_FIND_IN_SET,
			Value:    "important",
		})

		whereClause, args := query.BuildWhereClause()
		expected := "FIND_IN_SET(?, tags) > 0"
		assert.Equal(t, expected, whereClause)
		assert.Equal(t, []interface{}{"important"}, args)
	})
}

// TestQueryAddEqFilterIfNotEmpty 测试 AddEqFilterIfNotEmpty 方法
func TestQueryAddEqFilterIfNotEmpty(t *testing.T) {
	t.Run("非空值", func(t *testing.T) {
		query := NewQuery().AddEqFilterIfNotEmpty("name", "John")
		assert.Equal(t, 1, len(query.Filters))
		assert.Equal(t, constants.OP_EQ, query.Filters[0].Operator)
	})

	t.Run("空字符串", func(t *testing.T) {
		query := NewQuery().AddEqFilterIfNotEmpty("name", "")
		assert.Equal(t, 0, len(query.Filters))
	})

	t.Run("nil值", func(t *testing.T) {
		query := NewQuery().AddEqFilterIfNotEmpty("name", nil)
		assert.Equal(t, 0, len(query.Filters))
	})
}

// TestQueryAddNeqFilterIfNotEmpty 测试 AddNeqFilterIfNotEmpty 方法
func TestQueryAddNeqFilterIfNotEmpty(t *testing.T) {
	t.Run("非空值", func(t *testing.T) {
		query := NewQuery().AddNeqFilterIfNotEmpty("status", "deleted")
		assert.Equal(t, 1, len(query.Filters))
		assert.Equal(t, constants.OP_NEQ, query.Filters[0].Operator)
	})

	t.Run("空值", func(t *testing.T) {
		query := NewQuery().AddNeqFilterIfNotEmpty("status", "")
		assert.Equal(t, 0, len(query.Filters))
	})
}

// TestQueryAddGtFilterIfNotEmpty 测试 AddGtFilterIfNotEmpty 方法
func TestQueryAddGtFilterIfNotEmpty(t *testing.T) {
	t.Run("非空值", func(t *testing.T) {
		query := NewQuery().AddGtFilterIfNotEmpty("age", 18)
		assert.Equal(t, 1, len(query.Filters))
		assert.Equal(t, constants.OP_GT, query.Filters[0].Operator)
	})

	t.Run("零值整数", func(t *testing.T) {
		query := NewQuery().AddGtFilterIfNotEmpty("age", 0)
		assert.Equal(t, 0, len(query.Filters))
	})
}

// TestQueryAddGteFilterIfNotEmpty 测试 AddGteFilterIfNotEmpty 方法
func TestQueryAddGteFilterIfNotEmpty(t *testing.T) {
	t.Run("非空值", func(t *testing.T) {
		query := NewQuery().AddGteFilterIfNotEmpty("score", 60)
		assert.Equal(t, 1, len(query.Filters))
		assert.Equal(t, constants.OP_GTE, query.Filters[0].Operator)
	})

	t.Run("空值", func(t *testing.T) {
		query := NewQuery().AddGteFilterIfNotEmpty("score", nil)
		assert.Equal(t, 0, len(query.Filters))
	})
}

// TestQueryAddLtFilterIfNotEmpty 测试 AddLtFilterIfNotEmpty 方法
func TestQueryAddLtFilterIfNotEmpty(t *testing.T) {
	t.Run("非空值", func(t *testing.T) {
		query := NewQuery().AddLtFilterIfNotEmpty("age", 65)
		assert.Equal(t, 1, len(query.Filters))
		assert.Equal(t, constants.OP_LT, query.Filters[0].Operator)
	})

	t.Run("空值", func(t *testing.T) {
		query := NewQuery().AddLtFilterIfNotEmpty("age", nil)
		assert.Equal(t, 0, len(query.Filters))
	})
}

// TestQueryAddLteFilterIfNotEmpty 测试 AddLteFilterIfNotEmpty 方法
func TestQueryAddLteFilterIfNotEmpty(t *testing.T) {
	t.Run("非空值", func(t *testing.T) {
		query := NewQuery().AddLteFilterIfNotEmpty("price", 100.0)
		assert.Equal(t, 1, len(query.Filters))
		assert.Equal(t, constants.OP_LTE, query.Filters[0].Operator)
	})

	t.Run("空值", func(t *testing.T) {
		query := NewQuery().AddLteFilterIfNotEmpty("price", nil)
		assert.Equal(t, 0, len(query.Filters))
	})
}

// TestQueryAddNotInFilterIfNotEmpty 测试 AddNotInFilterIfNotEmpty 方法
func TestQueryAddNotInFilterIfNotEmpty(t *testing.T) {
	t.Run("非空切片", func(t *testing.T) {
		query := NewQuery().AddNotInFilterIfNotEmpty("status", []string{"deleted", "banned"})
		assert.Equal(t, 1, len(query.Filters))
		assert.Equal(t, constants.OP_NOT_IN, query.Filters[0].Operator)
	})

	t.Run("空切片", func(t *testing.T) {
		query := NewQuery().AddNotInFilterIfNotEmpty("status", []string{})
		assert.Equal(t, 0, len(query.Filters))
	})

	t.Run("nil切片", func(t *testing.T) {
		var nilSlice []string
		query := NewQuery().AddNotInFilterIfNotEmpty("status", nilSlice)
		assert.Equal(t, 0, len(query.Filters))
	})
}

// TestQueryAddBetweenFilterIfNotEmpty 测试 AddBetweenFilterIfNotEmpty 方法
func TestQueryAddBetweenFilterIfNotEmpty(t *testing.T) {
	t.Run("非空值", func(t *testing.T) {
		query := NewQuery().AddBetweenFilterIfNotEmpty("age", 18, 65)
		assert.Equal(t, 1, len(query.Filters))
		assert.Equal(t, constants.OP_BETWEEN, query.Filters[0].Operator)
	})

	t.Run("第一个值为空", func(t *testing.T) {
		query := NewQuery().AddBetweenFilterIfNotEmpty("age", nil, 65)
		assert.Equal(t, 0, len(query.Filters))
	})

	t.Run("第二个值为空", func(t *testing.T) {
		query := NewQuery().AddBetweenFilterIfNotEmpty("age", 18, nil)
		assert.Equal(t, 0, len(query.Filters))
	})
}

// TestQueryAddStartsWithFilterIfNotEmpty 测试 AddStartsWithFilterIfNotEmpty 方法
func TestQueryAddStartsWithFilterIfNotEmpty(t *testing.T) {
	t.Run("非空值", func(t *testing.T) {
		query := NewQuery().AddStartsWithFilterIfNotEmpty("name", "John")
		assert.Equal(t, 1, len(query.Filters))
		assert.Equal(t, constants.OP_STARTS_WITH, query.Filters[0].Operator)
	})

	t.Run("空字符串", func(t *testing.T) {
		query := NewQuery().AddStartsWithFilterIfNotEmpty("name", "")
		assert.Equal(t, 0, len(query.Filters))
	})
}

// TestQueryAddEndsWithFilterIfNotEmpty 测试 AddEndsWithFilterIfNotEmpty 方法
func TestQueryAddEndsWithFilterIfNotEmpty(t *testing.T) {
	t.Run("非空值", func(t *testing.T) {
		query := NewQuery().AddEndsWithFilterIfNotEmpty("email", "@example.com")
		assert.Equal(t, 1, len(query.Filters))
		assert.Equal(t, constants.OP_ENDS_WITH, query.Filters[0].Operator)
	})

	t.Run("空字符串", func(t *testing.T) {
		query := NewQuery().AddEndsWithFilterIfNotEmpty("email", "")
		assert.Equal(t, 0, len(query.Filters))
	})
}

// TestQueryAddContainsFilterIfNotEmpty 测试 AddContainsFilterIfNotEmpty 方法
func TestQueryAddContainsFilterIfNotEmpty(t *testing.T) {
	t.Run("非空值", func(t *testing.T) {
		query := NewQuery().AddContainsFilterIfNotEmpty("description", "test")
		assert.Equal(t, 1, len(query.Filters))
		assert.Equal(t, constants.OP_CONTAINS, query.Filters[0].Operator)
	})

	t.Run("空字符串", func(t *testing.T) {
		query := NewQuery().AddContainsFilterIfNotEmpty("description", "")
		assert.Equal(t, 0, len(query.Filters))
	})
}

// TestQueryAddNotLikeFilterIfNotEmpty 测试 AddNotLikeFilterIfNotEmpty 方法
func TestQueryAddNotLikeFilterIfNotEmpty(t *testing.T) {
	t.Run("非空值", func(t *testing.T) {
		query := NewQuery().AddNotLikeFilterIfNotEmpty("name", "%test%")
		assert.Equal(t, 1, len(query.Filters))
		assert.Equal(t, constants.OP_NOT_LIKE, query.Filters[0].Operator)
	})

	t.Run("空字符串", func(t *testing.T) {
		query := NewQuery().AddNotLikeFilterIfNotEmpty("name", "")
		assert.Equal(t, 0, len(query.Filters))
	})
}

// TestQueryAddFindInSetFilterIfNotEmpty 测试 AddFindInSetFilterIfNotEmpty 方法
func TestQueryAddFindInSetFilterIfNotEmpty(t *testing.T) {
	t.Run("非空值", func(t *testing.T) {
		query := NewQuery().AddFindInSetFilterIfNotEmpty("tags", "important")
		assert.Equal(t, 1, len(query.Filters))
		assert.Equal(t, constants.OP_FIND_IN_SET, query.Filters[0].Operator)
	})

	t.Run("空字符串", func(t *testing.T) {
		query := NewQuery().AddFindInSetFilterIfNotEmpty("tags", "")
		assert.Equal(t, 0, len(query.Filters))
	})
}

// TestFilterGroupConditionalMethods 测试 FilterGroup 的条件方法
func TestFilterGroupConditionalMethods(t *testing.T) {
	t.Run("AddFilterIf - 条件为true", func(t *testing.T) {
		group := NewFilterGroup(constants.LOGIC_AND)
		filter := NewEqFilter("name", "John")
		group.AddFilterIf(true, filter)
		assert.Equal(t, 1, len(group.Filters))
	})

	t.Run("AddFilterIf - 条件为false", func(t *testing.T) {
		group := NewFilterGroup(constants.LOGIC_AND)
		filter := NewEqFilter("name", "John")
		group.AddFilterIf(false, filter)
		assert.Equal(t, 0, len(group.Filters))
	})

	t.Run("AddFilterIfNotEmpty - 非空值", func(t *testing.T) {
		group := NewFilterGroup(constants.LOGIC_AND)
		group.AddFilterIfNotEmpty("name", constants.OP_EQ, "John")
		assert.Equal(t, 1, len(group.Filters))
	})

	t.Run("AddFilterIfNotEmpty - 空字符串", func(t *testing.T) {
		group := NewFilterGroup(constants.LOGIC_AND)
		group.AddFilterIfNotEmpty("name", constants.OP_EQ, "")
		assert.Equal(t, 0, len(group.Filters))
	})

	t.Run("AddFilterIfNotEmpty - 空切片", func(t *testing.T) {
		group := NewFilterGroup(constants.LOGIC_AND)
		group.AddFilterIfNotEmpty("status", constants.OP_IN, []string{})
		assert.Equal(t, 0, len(group.Filters))
	})
}

// TestFilterGroupHelperMethods 测试 FilterGroup 的辅助方法
func TestFilterGroupHelperMethods(t *testing.T) {
	t.Run("AddEqFilterIfNotEmpty", func(t *testing.T) {
		group := NewFilterGroup(constants.LOGIC_AND)
		group.AddEqFilterIfNotEmpty("name", "John")
		assert.Equal(t, 1, len(group.Filters))
		assert.Equal(t, constants.OP_EQ, group.Filters[0].Operator)
	})

	t.Run("AddNeqFilterIfNotEmpty", func(t *testing.T) {
		group := NewFilterGroup(constants.LOGIC_AND)
		group.AddNeqFilterIfNotEmpty("status", "deleted")
		assert.Equal(t, 1, len(group.Filters))
		assert.Equal(t, constants.OP_NEQ, group.Filters[0].Operator)
	})

	t.Run("AddGtFilterIfNotEmpty", func(t *testing.T) {
		group := NewFilterGroup(constants.LOGIC_AND)
		group.AddGtFilterIfNotEmpty("age", 18)
		assert.Equal(t, 1, len(group.Filters))
		assert.Equal(t, constants.OP_GT, group.Filters[0].Operator)
	})

	t.Run("AddGteFilterIfNotEmpty", func(t *testing.T) {
		group := NewFilterGroup(constants.LOGIC_AND)
		group.AddGteFilterIfNotEmpty("score", 60)
		assert.Equal(t, 1, len(group.Filters))
		assert.Equal(t, constants.OP_GTE, group.Filters[0].Operator)
	})

	t.Run("AddLtFilterIfNotEmpty", func(t *testing.T) {
		group := NewFilterGroup(constants.LOGIC_AND)
		group.AddLtFilterIfNotEmpty("age", 65)
		assert.Equal(t, 1, len(group.Filters))
		assert.Equal(t, constants.OP_LT, group.Filters[0].Operator)
	})

	t.Run("AddLteFilterIfNotEmpty", func(t *testing.T) {
		group := NewFilterGroup(constants.LOGIC_AND)
		group.AddLteFilterIfNotEmpty("price", 100.0)
		assert.Equal(t, 1, len(group.Filters))
		assert.Equal(t, constants.OP_LTE, group.Filters[0].Operator)
	})

	t.Run("AddLikeFilterIfNotEmpty", func(t *testing.T) {
		group := NewFilterGroup(constants.LOGIC_AND)
		group.AddLikeFilterIfNotEmpty("name", "%test%")
		assert.Equal(t, 1, len(group.Filters))
		assert.Equal(t, constants.OP_LIKE, group.Filters[0].Operator)
	})

	t.Run("AddInFilterIfNotEmpty", func(t *testing.T) {
		group := NewFilterGroup(constants.LOGIC_AND)
		group.AddInFilterIfNotEmpty("status", []interface{}{"active", "pending"})
		assert.Equal(t, 1, len(group.Filters))
		assert.Equal(t, constants.OP_IN, group.Filters[0].Operator)
	})

	t.Run("AddNotInFilterIfNotEmpty", func(t *testing.T) {
		group := NewFilterGroup(constants.LOGIC_AND)
		group.AddNotInFilterIfNotEmpty("status", []interface{}{"deleted"})
		assert.Equal(t, 1, len(group.Filters))
		assert.Equal(t, constants.OP_NOT_IN, group.Filters[0].Operator)
	})

	t.Run("AddBetweenFilterIfNotEmpty", func(t *testing.T) {
		group := NewFilterGroup(constants.LOGIC_AND)
		group.AddBetweenFilterIfNotEmpty("age", 18, 65)
		assert.Equal(t, 1, len(group.Filters))
		assert.Equal(t, constants.OP_BETWEEN, group.Filters[0].Operator)
	})

	t.Run("AddStartsWithFilterIfNotEmpty", func(t *testing.T) {
		group := NewFilterGroup(constants.LOGIC_AND)
		group.AddStartsWithFilterIfNotEmpty("name", "John")
		assert.Equal(t, 1, len(group.Filters))
		assert.Equal(t, constants.OP_STARTS_WITH, group.Filters[0].Operator)
	})

	t.Run("AddEndsWithFilterIfNotEmpty", func(t *testing.T) {
		group := NewFilterGroup(constants.LOGIC_AND)
		group.AddEndsWithFilterIfNotEmpty("email", "@example.com")
		assert.Equal(t, 1, len(group.Filters))
		assert.Equal(t, constants.OP_ENDS_WITH, group.Filters[0].Operator)
	})

	t.Run("AddContainsFilterIfNotEmpty", func(t *testing.T) {
		group := NewFilterGroup(constants.LOGIC_AND)
		group.AddContainsFilterIfNotEmpty("description", "test")
		assert.Equal(t, 1, len(group.Filters))
		assert.Equal(t, constants.OP_CONTAINS, group.Filters[0].Operator)
	})

	t.Run("AddNotLikeFilterIfNotEmpty", func(t *testing.T) {
		group := NewFilterGroup(constants.LOGIC_AND)
		group.AddNotLikeFilterIfNotEmpty("name", "%test%")
		assert.Equal(t, 1, len(group.Filters))
		assert.Equal(t, constants.OP_NOT_LIKE, group.Filters[0].Operator)
	})

	t.Run("AddFindInSetFilterIfNotEmpty", func(t *testing.T) {
		group := NewFilterGroup(constants.LOGIC_AND)
		group.AddFindInSetFilterIfNotEmpty("tags", "important")
		assert.Equal(t, 1, len(group.Filters))
		assert.Equal(t, constants.OP_FIND_IN_SET, group.Filters[0].Operator)
	})
}

// TestFilterGroupAddGroupIf 测试 AddGroupIf 方法
func TestFilterGroupAddGroupIf(t *testing.T) {
	t.Run("条件为true", func(t *testing.T) {
		mainGroup := NewFilterGroup(constants.LOGIC_AND)
		subGroup := NewFilterGroup(constants.LOGIC_OR).AddFilter(NewEqFilter("status", "active"))
		mainGroup.AddGroupIf(true, subGroup)
		assert.Equal(t, 1, len(mainGroup.Groups))
	})

	t.Run("条件为false", func(t *testing.T) {
		mainGroup := NewFilterGroup(constants.LOGIC_AND)
		subGroup := NewFilterGroup(constants.LOGIC_OR).AddFilter(NewEqFilter("status", "active"))
		mainGroup.AddGroupIf(false, subGroup)
		assert.Equal(t, 0, len(mainGroup.Groups))
	})
}

// TestFilterGroupAddGroupIfNotEmpty 测试 AddGroupIfNotEmpty 方法
func TestFilterGroupAddGroupIfNotEmpty(t *testing.T) {
	t.Run("非空组", func(t *testing.T) {
		mainGroup := NewFilterGroup(constants.LOGIC_AND)
		subGroup := NewFilterGroup(constants.LOGIC_OR).AddFilter(NewEqFilter("status", "active"))
		mainGroup.AddGroupIfNotEmpty(subGroup)
		assert.Equal(t, 1, len(mainGroup.Groups))
	})

	t.Run("空组", func(t *testing.T) {
		mainGroup := NewFilterGroup(constants.LOGIC_AND)
		emptyGroup := NewFilterGroup(constants.LOGIC_OR)
		mainGroup.AddGroupIfNotEmpty(emptyGroup)
		assert.Equal(t, 0, len(mainGroup.Groups))
	})
}

// TestFilterGroupClear 测试 Clear 方法
func TestFilterGroupClear(t *testing.T) {
	group := NewFilterGroup(constants.LOGIC_AND).
		AddFilter(NewEqFilter("name", "John")).
		AddFilter(NewEqFilter("age", 30)).
		AddGroup(NewFilterGroup(constants.LOGIC_OR))

	assert.Equal(t, 2, len(group.Filters))
	assert.Equal(t, 1, len(group.Groups))

	group.Clear()
	assert.Equal(t, 0, len(group.Filters))
	assert.Equal(t, 0, len(group.Groups))
}

// TestFilterGroupClone 测试 Clone 方法
func TestFilterGroupClone(t *testing.T) {
	original := NewFilterGroup(constants.LOGIC_AND).
		AddFilter(NewEqFilter("name", "John")).
		AddFilter(NewEqFilter("age", 30))

	cloned := original.Clone()

	// 验证克隆的内容相同
	assert.Equal(t, original.LogicOp, cloned.LogicOp)
	assert.Equal(t, len(original.Filters), len(cloned.Filters))
	assert.Equal(t, original.Filters[0].Field, cloned.Filters[0].Field)

	// 修改克隆的对象，不应影响原对象
	cloned.AddFilter(NewEqFilter("status", "active"))
	assert.Equal(t, 2, len(original.Filters))
	assert.Equal(t, 3, len(cloned.Filters))
}

// TestNewInFilterSliceEmptySlice 测试空切片情况
func TestNewInFilterSliceEmptySlice(t *testing.T) {
	filter := NewInFilterSlice("status", []interface{}{})
	assert.Equal(t, "status", filter.Field)
	assert.Equal(t, constants.OP_IN, filter.Operator)
	values := filter.Value.([]interface{})
	assert.Equal(t, 0, len(values))
}

// TestNewNotInFilterSlice 测试 NewNotInFilterSlice 函数
func TestNewNotInFilterSlice(t *testing.T) {
	t.Run("非空切片", func(t *testing.T) {
		filter := NewNotInFilterSlice("status", []interface{}{"deleted", "banned"})
		assert.Equal(t, "status", filter.Field)
		assert.Equal(t, constants.OP_NOT_IN, filter.Operator)
		values := filter.Value.([]interface{})
		assert.Equal(t, 2, len(values))
	})

	t.Run("空切片", func(t *testing.T) {
		filter := NewNotInFilterSlice("status", []interface{}{})
		assert.Equal(t, "status", filter.Field)
		values := filter.Value.([]interface{})
		assert.Equal(t, 0, len(values))
	})
}

// TestIsTimeValidEdgeCases 测试 isTimeValid 的边界情况
func TestIsTimeValidEdgeCases(t *testing.T) {
	t.Run("零值time.Time", func(t *testing.T) {
		var zeroTime time.Time
		result := isTimeValid(&zeroTime)
		assert.False(t, result)
	})

	t.Run("非零time.Time", func(t *testing.T) {
		now := time.Now()
		result := isTimeValid(&now)
		assert.True(t, result)
	})

	t.Run("nil指针", func(t *testing.T) {
		result := isTimeValid(nil)
		assert.False(t, result)
	})
}

// TestNewInFilterSliceNilValues 测试 NewInFilterSlice 处理 nil 值
func TestNewInFilterSliceNilValues(t *testing.T) {
	filter := NewInFilterSlice("status", nil)
	assert.Equal(t, "status", filter.Field)
	assert.Equal(t, constants.OP_IN, filter.Operator)
	values := filter.Value.([]interface{})
	assert.NotNil(t, values)
	assert.Equal(t, 0, len(values))
}

// TestNewNotInFilterNilValues 测试 NewNotInFilter 处理 nil 值
func TestNewNotInFilterNilValues(t *testing.T) {
	filter := NewNotInFilter("status", nil)
	assert.Equal(t, "status", filter.Field)
	assert.Equal(t, constants.OP_NOT_IN, filter.Operator)
	assert.Nil(t, filter.Value)
}

// TestNewNotInFilterSliceNilValues 测试 NewNotInFilterSlice 处理 nil 值
func TestNewNotInFilterSliceNilValues(t *testing.T) {
	filter := NewNotInFilterSlice("status", nil)
	assert.Equal(t, "status", filter.Field)
	assert.Equal(t, constants.OP_NOT_IN, filter.Operator)
	values := filter.Value.([]interface{})
	assert.NotNil(t, values)
	assert.Equal(t, 0, len(values))
}

// TestFilterGroupCloneWithGroups 测试 Clone 方法包含嵌套组
func TestFilterGroupCloneWithGroups(t *testing.T) {
	subGroup := NewFilterGroup(constants.LOGIC_OR).
		AddFilter(NewEqFilter("status", "active"))

	original := NewFilterGroup(constants.LOGIC_AND).
		AddFilter(NewEqFilter("name", "John")).
		AddGroup(subGroup)

	cloned := original.Clone()

	// 验证克隆的内容相同
	assert.Equal(t, len(original.Groups), len(cloned.Groups))
	assert.Equal(t, original.Groups[0].LogicOp, cloned.Groups[0].LogicOp)

	// 修改克隆的嵌套组，不应影响原对象
	cloned.Groups[0].AddFilter(NewEqFilter("age", 30))
	assert.Equal(t, 1, len(original.Groups[0].Filters))
	assert.Equal(t, 2, len(cloned.Groups[0].Filters))
}

// TestQueryAddNotInFilterIfNotEmptyEdgeCases 测试 AddNotInFilterIfNotEmpty 边界情况
func TestQueryAddNotInFilterIfNotEmptyEdgeCases(t *testing.T) {
	t.Run("非空数组", func(t *testing.T) {
		values := [3]string{"a", "b", "c"}
		query := NewQuery().AddNotInFilterIfNotEmpty("status", values)
		assert.Equal(t, 1, len(query.Filters))
		assert.Equal(t, constants.OP_NOT_IN, query.Filters[0].Operator)
	})

	t.Run("空数组", func(t *testing.T) {
		var emptyArray [0]string
		query := NewQuery().AddNotInFilterIfNotEmpty("status", emptyArray)
		assert.Equal(t, 0, len(query.Filters))
	})
}

// TestIsTimeValidWithPointer 测试 isTimeValid 处理指针类型
func TestIsTimeValidWithPointer(t *testing.T) {
	t.Run("time.Time指针", func(t *testing.T) {
		now := time.Now()
		result := isTimeValid(&now)
		assert.True(t, result)
	})

	t.Run("time.Time值", func(t *testing.T) {
		now := time.Now()
		result := isTimeValid(now)
		assert.True(t, result)
	})
}

// TestHandleSpecialOperatorsWithSubquery 测试特殊操作符处理子查询
func TestHandleSpecialOperatorsWithSubquery(t *testing.T) {
	t.Run("IS_NULL操作符", func(t *testing.T) {
		query := NewQuery().AddFilter(&Filter{
			Field:    "deleted_at",
			Operator: constants.OP_IS_NULL,
			Value:    nil,
		})

		whereClause, args := query.BuildWhereClause()
		assert.Equal(t, "deleted_at IS NULL", whereClause)
		assert.Equal(t, 0, len(args))
	})

	t.Run("IS_NOT_NULL操作符", func(t *testing.T) {
		query := NewQuery().AddFilter(&Filter{
			Field:    "updated_at",
			Operator: constants.OP_IS_NOT_NULL,
			Value:    nil,
		})

		whereClause, args := query.BuildWhereClause()
		assert.Equal(t, "updated_at IS NOT NULL", whereClause)
		assert.Equal(t, 0, len(args))
	})

	t.Run("BETWEEN操作符带数组值", func(t *testing.T) {
		query := NewQuery().AddFilter(&Filter{
			Field:    "age",
			Operator: constants.OP_BETWEEN,
			Value:    []interface{}{18, 65},
		})

		whereClause, args := query.BuildWhereClause()
		assert.Equal(t, "age BETWEEN ? AND ?", whereClause)
		assert.Equal(t, []interface{}{18, 65}, args)
	})

	t.Run("FIND_IN_SET操作符", func(t *testing.T) {
		query := NewQuery().AddFilter(&Filter{
			Field:    "tags",
			Operator: constants.OP_FIND_IN_SET,
			Value:    "urgent",
		})

		whereClause, args := query.BuildWhereClause()
		assert.Equal(t, "FIND_IN_SET(?, tags) > 0", whereClause)
		assert.Equal(t, []interface{}{"urgent"}, args)
	})
}

// TestBuildWhereClauseWithEmptyFilters 测试空过滤条件
func TestBuildWhereClauseWithEmptyFilters(t *testing.T) {
	query := NewQuery()
	whereClause, args := query.BuildWhereClause()
	assert.Equal(t, "", whereClause)
	assert.Equal(t, 0, len(args))
}

// TestBuildGroupConditionEdgeCases 测试 buildGroupCondition 边界情况
func TestBuildGroupConditionEdgeCases(t *testing.T) {
	t.Run("空组", func(t *testing.T) {
		emptyGroup := NewFilterGroup(constants.LOGIC_AND)
		query := NewQuery().WithFilterGroup(emptyGroup)
		whereClause, args := query.BuildWhereClause()
		assert.Equal(t, "", whereClause)
		assert.Equal(t, 0, len(args))
	})

	t.Run("仅有过滤条件的组", func(t *testing.T) {
		group := NewFilterGroup(constants.LOGIC_AND).
			AddFilter(NewEqFilter("name", "John"))
		query := NewQuery().WithFilterGroup(group)
		whereClause, args := query.BuildWhereClause()
		assert.Equal(t, "name = ?", whereClause)
		assert.Equal(t, []interface{}{"John"}, args)
	})

	t.Run("仅有嵌套组的组", func(t *testing.T) {
		subGroup := NewFilterGroup(constants.LOGIC_OR).
			AddFilter(NewEqFilter("status", "active"))
		mainGroup := NewFilterGroup(constants.LOGIC_AND).
			AddGroup(subGroup)
		query := NewQuery().WithFilterGroup(mainGroup)
		whereClause, args := query.BuildWhereClause()
		assert.Equal(t, "status = ?", whereClause)
		assert.Equal(t, []interface{}{"active"}, args)
	})
}

// TestNewNotInFilterVariadic 测试 NewNotInFilter 可变参数
func TestNewNotInFilterVariadic(t *testing.T) {
	t.Run("多个参数", func(t *testing.T) {
		filter := NewNotInFilter("status", "deleted", "banned", "archived")
		assert.Equal(t, "status", filter.Field)
		assert.Equal(t, constants.OP_NOT_IN, filter.Operator)
		values := filter.Value.([]interface{})
		assert.Equal(t, 3, len(values))
		assert.Equal(t, "deleted", values[0])
		assert.Equal(t, "banned", values[1])
		assert.Equal(t, "archived", values[2])
	})

	t.Run("单个参数", func(t *testing.T) {
		filter := NewNotInFilter("status", "deleted")
		values := filter.Value.([]interface{})
		assert.Equal(t, 1, len(values))
	})

	t.Run("无参数", func(t *testing.T) {
		filter := NewNotInFilter("status")
		values := filter.Value.([]interface{})
		assert.Equal(t, 0, len(values))
	})
}

// TestIsTimeValidUnixZero 测试 isTimeValid 处理 Unix 零点前的时间
func TestIsTimeValidUnixZero(t *testing.T) {
	t.Run("Unix零点前的时间", func(t *testing.T) {
		beforeUnix := time.Date(1960, 1, 1, 0, 0, 0, 0, time.UTC)
		result := isTimeValid(beforeUnix)
		assert.False(t, result)
	})

	t.Run("Unix零点后的时间", func(t *testing.T) {
		afterUnix := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
		result := isTimeValid(afterUnix)
		assert.True(t, result)
	})

	t.Run("Unix零点前的时间指针", func(t *testing.T) {
		beforeUnix := time.Date(1960, 1, 1, 0, 0, 0, 0, time.UTC)
		result := isTimeValid(&beforeUnix)
		assert.False(t, result)
	})
}

// TestAddNotInFilterIfNotEmptyNilValue 测试 AddNotInFilterIfNotEmpty 处理 nil
func TestAddNotInFilterIfNotEmptyNilValue(t *testing.T) {
	query := NewQuery().AddNotInFilterIfNotEmpty("status", nil)
	assert.Equal(t, 0, len(query.Filters))
}

// TestAddThisWeekOnSunday 测试 AddThisWeek 在周日的情况
func TestAddThisWeekOnSunday(t *testing.T) {
	// 这个测试会覆盖 weekday == 0 的分支
	query := NewQuery().AddThisWeek("created_at")
	assert.Equal(t, 1, len(query.Filters))
	assert.Equal(t, constants.OP_BETWEEN, query.Filters[0].Operator)
}

// TestHandleOperatorsWithInvalidTypes 测试各种 handle 函数处理无效类型
func TestHandleOperatorsWithInvalidTypes(t *testing.T) {
	query := NewQuery()

	t.Run("handleNullOperators返回空字符串", func(t *testing.T) {
		filter := &Filter{
			Field:    "test",
			Operator: "INVALID_OP",
			Value:    nil,
		}
		sql, arg := query.handleNullOperators(filter)
		assert.Equal(t, "", sql)
		assert.Nil(t, arg)
	})

	t.Run("handleBetweenOperator处理非切片值", func(t *testing.T) {
		filter := &Filter{
			Field:    "age",
			Operator: constants.OP_BETWEEN,
			Value:    "not a slice",
		}
		sql, arg := query.handleBetweenOperator(filter)
		assert.Equal(t, "", sql)
		assert.Nil(t, arg)
	})

	t.Run("handleBetweenOperator处理长度不为2的切片", func(t *testing.T) {
		filter := &Filter{
			Field:    "age",
			Operator: constants.OP_BETWEEN,
			Value:    []interface{}{18},
		}
		sql, arg := query.handleBetweenOperator(filter)
		assert.Equal(t, "", sql)
		assert.Nil(t, arg)
	})

	t.Run("handleStartsWithOperator处理非字符串值", func(t *testing.T) {
		filter := &Filter{
			Field:    "name",
			Operator: constants.OP_STARTS_WITH,
			Value:    123,
		}
		sql, arg := query.handleStartsWithOperator(filter)
		assert.Equal(t, "", sql)
		assert.Nil(t, arg)
	})

	t.Run("handleEndsWithOperator处理非字符串值", func(t *testing.T) {
		filter := &Filter{
			Field:    "email",
			Operator: constants.OP_ENDS_WITH,
			Value:    []byte("test"),
		}
		sql, arg := query.handleEndsWithOperator(filter)
		assert.Equal(t, "", sql)
		assert.Nil(t, arg)
	})

	t.Run("handleContainsOperator处理非字符串值", func(t *testing.T) {
		filter := &Filter{
			Field:    "description",
			Operator: constants.OP_CONTAINS,
			Value:    42,
		}
		sql, arg := query.handleContainsOperator(filter)
		assert.Equal(t, "", sql)
		assert.Nil(t, arg)
	})

	t.Run("handleFindInSetOperator返回空字符串对于无效操作符", func(t *testing.T) {
		filter := &Filter{
			Field:    "tags",
			Operator: "INVALID_FIND_IN_SET",
			Value:    "test",
		}
		sql, arg := query.handleFindInSetOperator(filter)
		assert.Equal(t, "", sql)
		assert.Nil(t, arg)
	})
}

// TestBuildFilterConditionNilFilter 测试 buildFilterCondition 处理 nil
func TestBuildFilterConditionNilFilter(t *testing.T) {
	query := NewQuery()
	sql, arg := query.buildFilterCondition(nil)
	assert.Equal(t, "", sql)
	assert.Nil(t, arg)
}

// TestBuildFilterConditionUnknownOperator 测试 buildFilterCondition 处理未知操作符
func TestBuildFilterConditionUnknownOperator(t *testing.T) {
	query := NewQuery()
	filter := &Filter{
		Field:    "test",
		Operator: "UNKNOWN_OPERATOR",
		Value:    "value",
	}
	sql, arg := query.buildFilterCondition(filter)
	assert.Equal(t, "", sql)
	assert.Nil(t, arg)
}

// TestBuildGroupConditionNilGroup 测试 buildGroupCondition 处理 nil
func TestBuildGroupConditionNilGroup(t *testing.T) {
	query := NewQuery()
	sql, args := query.buildGroupCondition(nil)
	assert.Equal(t, "", sql)
	assert.Nil(t, args)
}

// TestBuildGroupConditionWithNilFilters 测试 buildGroupCondition 处理包含 nil 过滤条件的组
func TestBuildGroupConditionWithNilFilters(t *testing.T) {
	group := NewFilterGroup(constants.LOGIC_AND)
	group.Filters = append(group.Filters, nil)
	group.Filters = append(group.Filters, NewEqFilter("name", "John"))
	group.Filters = append(group.Filters, nil)

	query := NewQuery()
	sql, args := query.buildGroupCondition(group)
	assert.Equal(t, "name = ?", sql)
	assert.Equal(t, []interface{}{"John"}, args)
}

// TestProcessSubGroupsWithNilGroup 测试 processSubGroups 处理 nil 子组
func TestProcessSubGroupsWithNilGroup(t *testing.T) {
	mainGroup := NewFilterGroup(constants.LOGIC_AND)
	mainGroup.Groups = append(mainGroup.Groups, nil)
	subGroup := NewFilterGroup(constants.LOGIC_OR).AddFilter(NewEqFilter("status", "active"))
	mainGroup.Groups = append(mainGroup.Groups, subGroup)
	mainGroup.Groups = append(mainGroup.Groups, nil)

	query := NewQuery().WithFilterGroup(mainGroup)
	whereClause, args := query.BuildWhereClause()
	assert.Equal(t, "status = ?", whereClause)
	assert.Equal(t, []interface{}{"active"}, args)
}

// TestProcessSubGroupsWithEmptyGroup 测试 processSubGroups 处理空子组
func TestProcessSubGroupsWithEmptyGroup(t *testing.T) {
	mainGroup := NewFilterGroup(constants.LOGIC_AND)
	emptyGroup := NewFilterGroup(constants.LOGIC_OR)
	mainGroup.Groups = append(mainGroup.Groups, emptyGroup)
	validGroup := NewFilterGroup(constants.LOGIC_OR).AddFilter(NewEqFilter("name", "John"))
	mainGroup.Groups = append(mainGroup.Groups, validGroup)

	query := NewQuery().WithFilterGroup(mainGroup)
	whereClause, args := query.BuildWhereClause()
	assert.Equal(t, "name = ?", whereClause)
	assert.Equal(t, []interface{}{"John"}, args)
}

// TestIsTimeValidWithOtherTypes 测试 isTimeValid 处理其他类型（返回 true）
func TestIsTimeValidWithOtherTypes(t *testing.T) {
	t.Run("字符串类型", func(t *testing.T) {
		result := isTimeValid("2025-01-01")
		assert.True(t, result)
	})

	t.Run("整数类型", func(t *testing.T) {
		result := isTimeValid(123456789)
		assert.True(t, result)
	})

	t.Run("布尔类型", func(t *testing.T) {
		result := isTimeValid(true)
		assert.True(t, result)
	})
}

// TestAddThisWeekCoverage 测试 AddThisWeek 完整覆盖
func TestAddThisWeekCoverage(t *testing.T) {
	// 通过多次调用确保覆盖不同的星期几
	for i := 0; i < 7; i++ {
		query := NewQuery().AddThisWeek("created_at")
		assert.Equal(t, 1, len(query.Filters))
	}
}

// TestBuildGroupConditionAllBranches 测试 buildGroupCondition 所有分支
func TestBuildGroupConditionAllBranches(t *testing.T) {
	query := NewQuery()

	t.Run("空条件列表", func(t *testing.T) {
		group := NewFilterGroup(constants.LOGIC_AND)
		// 添加一个返回空条件的过滤器
		group.Filters = append(group.Filters, &Filter{
			Field:    "test",
			Operator: "UNKNOWN_OP",
			Value:    "value",
		})
		sql, args := query.buildGroupCondition(group)
		assert.Equal(t, "", sql)
		assert.Nil(t, args)
	})

	t.Run("混合有效和无效过滤器", func(t *testing.T) {
		group := NewFilterGroup(constants.LOGIC_AND)
		group.Filters = append(group.Filters, &Filter{
			Field:    "test",
			Operator: "UNKNOWN_OP",
			Value:    "value",
		})
		group.Filters = append(group.Filters, NewEqFilter("name", "John"))
		group.Filters = append(group.Filters, &Filter{
			Field:    "test2",
			Operator: "UNKNOWN_OP",
			Value:    "value2",
		})
		sql, args := query.buildGroupCondition(group)
		assert.Equal(t, "name = ?", sql)
		assert.Equal(t, []interface{}{"John"}, args)
	})
}
