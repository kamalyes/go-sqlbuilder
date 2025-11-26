/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-26 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-26 00:00:00
 * @FilePath: \go-sqlbuilder\repository\query_extensions_test.go
 * @Description: Query 扩展方法测试 - 测试新增的便捷查询构建方法
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

// =============== 基础条件构建方法测试 ===============

// TestAddEqual 测试等于条件
func TestAddEqual(t *testing.T) {
	query := NewQuery()
	result := query.AddEqual("status", 1)

	assert.Equal(t, query, result, "应该返回同一个查询对象")
	assert.Equal(t, 1, len(query.Filters))
	assert.Equal(t, "status", query.Filters[0].Field)
	assert.Equal(t, constants.OP_EQ, query.Filters[0].Operator)
	assert.Equal(t, 1, query.Filters[0].Value)
}

// TestAddNotEqual 测试不等于条件
func TestAddNotEqual(t *testing.T) {
	query := NewQuery()
	query.AddNotEqual("status", 0)

	assert.Equal(t, 1, len(query.Filters))
	assert.Equal(t, "status", query.Filters[0].Field)
	assert.Equal(t, constants.OP_NEQ, query.Filters[0].Operator)
	assert.Equal(t, 0, query.Filters[0].Value)
}

// TestAddLike 测试LIKE条件
func TestAddLike(t *testing.T) {
	// 非空关键词
	query := NewQuery()
	query.AddLike("name", "test")

	assert.Equal(t, 1, len(query.Filters))
	assert.Equal(t, "name", query.Filters[0].Field)
	assert.Equal(t, constants.OP_LIKE, query.Filters[0].Operator)
	assert.Equal(t, "%test%", query.Filters[0].Value)

	// 空关键词 - 不应添加过滤条件
	query = NewQuery()
	query.AddLike("name", "")
	assert.Equal(t, 0, len(query.Filters))
}

// TestAddStartsWith 测试前缀匹配条件
func TestAddStartsWith(t *testing.T) {
	// 非空前缀
	query := NewQuery()
	query.AddStartsWith("username", "admin")

	assert.Equal(t, 1, len(query.Filters))
	assert.Equal(t, "username", query.Filters[0].Field)
	assert.Equal(t, constants.OP_LIKE, query.Filters[0].Operator)
	assert.Equal(t, "admin%", query.Filters[0].Value)

	// 空前缀 - 不应添加过滤条件
	query = NewQuery()
	query.AddStartsWith("username", "")
	assert.Equal(t, 0, len(query.Filters))
}

// TestAddEndsWith 测试后缀匹配条件
func TestAddEndsWith(t *testing.T) {
	// 非空后缀
	query := NewQuery()
	query.AddEndsWith("email", "@example.com")

	assert.Equal(t, 1, len(query.Filters))
	assert.Equal(t, "email", query.Filters[0].Field)
	assert.Equal(t, constants.OP_LIKE, query.Filters[0].Operator)
	assert.Equal(t, "%@example.com", query.Filters[0].Value)

	// 空后缀 - 不应添加过滤条件
	query = NewQuery()
	query.AddEndsWith("email", "")
	assert.Equal(t, 0, len(query.Filters))
}

// TestAddIn 测试IN条件
func TestAddIn(t *testing.T) {
	// 有值的IN条件
	query := NewQuery()
	query.AddIn("status", 1, 2, 3)

	assert.Equal(t, 1, len(query.Filters))
	assert.Equal(t, "status", query.Filters[0].Field)
	assert.Equal(t, constants.OP_IN, query.Filters[0].Operator)
	values := query.Filters[0].Value.([]interface{})
	assert.Equal(t, 3, len(values))
	assert.Equal(t, 1, values[0])
	assert.Equal(t, 2, values[1])
	assert.Equal(t, 3, values[2])

	// 空值 - 不应添加过滤条件
	query = NewQuery()
	query.AddIn("status")
	assert.Equal(t, 0, len(query.Filters))
}

// TestAddNotIn 测试NOT IN条件
func TestAddNotIn(t *testing.T) {
	// 有值的NOT IN条件
	query := NewQuery()
	query.AddNotIn("status", "deleted", "disabled")

	assert.Equal(t, 1, len(query.Filters))
	assert.Equal(t, "status", query.Filters[0].Field)
	assert.Equal(t, constants.OP_NOT_IN, query.Filters[0].Operator)
	values := query.Filters[0].Value.([]interface{})
	assert.Equal(t, 2, len(values))
	assert.Equal(t, "deleted", values[0])
	assert.Equal(t, "disabled", values[1])

	// 空值 - 不应添加过滤条件
	query = NewQuery()
	query.AddNotIn("status")
	assert.Equal(t, 0, len(query.Filters))
}

// TestAddGreaterThan 测试大于条件
func TestAddGreaterThan(t *testing.T) {
	query := NewQuery()
	query.AddGreaterThan("age", 18)

	assert.Equal(t, 1, len(query.Filters))
	assert.Equal(t, "age", query.Filters[0].Field)
	assert.Equal(t, constants.OP_GT, query.Filters[0].Operator)
	assert.Equal(t, 18, query.Filters[0].Value)
}

// TestAddGreaterEqual 测试大于等于条件
func TestAddGreaterEqual(t *testing.T) {
	query := NewQuery()
	query.AddGreaterEqual("score", 60)

	assert.Equal(t, 1, len(query.Filters))
	assert.Equal(t, "score", query.Filters[0].Field)
	assert.Equal(t, constants.OP_GTE, query.Filters[0].Operator)
	assert.Equal(t, 60, query.Filters[0].Value)
}

// TestAddLessThan 测试小于条件
func TestAddLessThan(t *testing.T) {
	query := NewQuery()
	query.AddLessThan("price", 100.0)

	assert.Equal(t, 1, len(query.Filters))
	assert.Equal(t, "price", query.Filters[0].Field)
	assert.Equal(t, constants.OP_LT, query.Filters[0].Operator)
	assert.Equal(t, 100.0, query.Filters[0].Value)
}

// TestAddLessEqual 测试小于等于条件
func TestAddLessEqual(t *testing.T) {
	query := NewQuery()
	query.AddLessEqual("count", 50)

	assert.Equal(t, 1, len(query.Filters))
	assert.Equal(t, "count", query.Filters[0].Field)
	assert.Equal(t, constants.OP_LTE, query.Filters[0].Operator)
	assert.Equal(t, 50, query.Filters[0].Value)
}

// TestAddBetween 测试BETWEEN条件
func TestAddBetween(t *testing.T) {
	query := NewQuery()
	query.AddBetween("age", 18, 65)

	assert.Equal(t, 1, len(query.Filters))
	assert.Equal(t, "age", query.Filters[0].Field)
	assert.Equal(t, constants.OP_BETWEEN, query.Filters[0].Operator)
	values := query.Filters[0].Value.([]interface{})
	assert.Equal(t, 2, len(values))
	assert.Equal(t, 18, values[0])
	assert.Equal(t, 65, values[1])
}

// TestAddIsNull 测试IS NULL条件
func TestAddIsNull(t *testing.T) {
	query := NewQuery()
	query.AddIsNull("deleted_at")

	assert.Equal(t, 1, len(query.Filters))
	assert.Equal(t, "deleted_at", query.Filters[0].Field)
	assert.Equal(t, constants.OP_IS_NULL, query.Filters[0].Operator)
	assert.Nil(t, query.Filters[0].Value)
}

// TestAddIsNotNull 测试IS NOT NULL条件
func TestAddIsNotNull(t *testing.T) {
	query := NewQuery()
	query.AddIsNotNull("email")

	assert.Equal(t, 1, len(query.Filters))
	assert.Equal(t, "email", query.Filters[0].Field)
	assert.Equal(t, constants.OP_IS_NOT_NULL, query.Filters[0].Operator)
	assert.Nil(t, query.Filters[0].Value)
}

// =============== 排序方法测试 ===============

// TestAddOrderAsc 测试升序排序
func TestAddOrderAsc(t *testing.T) {
	query := NewQuery()
	result := query.AddOrderAsc("name")

	assert.Equal(t, query, result, "应该返回同一个查询对象")
	assert.Equal(t, 1, len(query.Orders))
	assert.Equal(t, "name", query.Orders[0].Field)
	assert.Equal(t, constants.Asc, query.Orders[0].Direction)
}

// TestAddOrderDesc 测试降序排序
func TestAddOrderDesc(t *testing.T) {
	query := NewQuery()
	query.AddOrderDesc("created_at")

	assert.Equal(t, 1, len(query.Orders))
	assert.Equal(t, "created_at", query.Orders[0].Field)
	assert.Equal(t, constants.Desc, query.Orders[0].Direction)
}

// =============== 时间相关方法测试 ===============

// TestAddTimeAfter 测试时间晚于条件
func TestAddTimeAfter(t *testing.T) {
	// 非零时间
	testTime := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
	query := NewQuery()
	query.AddTimeAfter("created_at", testTime)

	assert.Equal(t, 1, len(query.Filters))
	assert.Equal(t, "created_at", query.Filters[0].Field)
	assert.Equal(t, constants.OP_GT, query.Filters[0].Operator)
	assert.Equal(t, testTime, query.Filters[0].Value)

	// 零时间 - 不应添加过滤条件
	query = NewQuery()
	query.AddTimeAfter("created_at", time.Time{})
	assert.Equal(t, 0, len(query.Filters))
}

// TestAddTimeBefore 测试时间早于条件
func TestAddTimeBefore(t *testing.T) {
	// 非零时间
	testTime := time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC)
	query := NewQuery()
	query.AddTimeBefore("updated_at", testTime)

	assert.Equal(t, 1, len(query.Filters))
	assert.Equal(t, "updated_at", query.Filters[0].Field)
	assert.Equal(t, constants.OP_LT, query.Filters[0].Operator)
	assert.Equal(t, testTime, query.Filters[0].Value)

	// 零时间 - 不应添加过滤条件
	query = NewQuery()
	query.AddTimeBefore("updated_at", time.Time{})
	assert.Equal(t, 0, len(query.Filters))
}

// TestAddTimeBetween 测试时间范围条件
func TestAddTimeBetween(t *testing.T) {
	start := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC)

	// 完整时间范围
	query := NewQuery()
	query.AddTimeBetween("created_at", start, end)

	assert.Equal(t, 2, len(query.Filters))
	// 第一个过滤条件：>= start
	assert.Equal(t, "created_at", query.Filters[0].Field)
	assert.Equal(t, constants.OP_GTE, query.Filters[0].Operator)
	assert.Equal(t, start, query.Filters[0].Value)
	// 第二个过滤条件：<= end
	assert.Equal(t, "created_at", query.Filters[1].Field)
	assert.Equal(t, constants.OP_LTE, query.Filters[1].Operator)
	assert.Equal(t, end, query.Filters[1].Value)

	// 只有开始时间
	query = NewQuery()
	query.AddTimeBetween("created_at", start, time.Time{})
	assert.Equal(t, 1, len(query.Filters))
	assert.Equal(t, constants.OP_GTE, query.Filters[0].Operator)

	// 只有结束时间
	query = NewQuery()
	query.AddTimeBetween("created_at", time.Time{}, end)
	assert.Equal(t, 1, len(query.Filters))
	assert.Equal(t, constants.OP_LTE, query.Filters[0].Operator)

	// 两个都是零时间
	query = NewQuery()
	query.AddTimeBetween("created_at", time.Time{}, time.Time{})
	assert.Equal(t, 0, len(query.Filters))
}

// TestAddToday 测试今天条件
func TestAddToday(t *testing.T) {
	query := NewQuery()
	query.AddToday("created_at")

	assert.Equal(t, 1, len(query.Filters))
	assert.Equal(t, "created_at", query.Filters[0].Field)
	assert.Equal(t, constants.OP_BETWEEN, query.Filters[0].Operator)

	// 验证时间范围是今天的开始到结束
	values := query.Filters[0].Value.([]interface{})
	assert.Equal(t, 2, len(values))

	startTime := values[0].(time.Time)
	endTime := values[1].(time.Time)
	now := time.Now()

	// 验证是今天的日期
	assert.Equal(t, now.Year(), startTime.Year())
	assert.Equal(t, now.Month(), startTime.Month())
	assert.Equal(t, now.Day(), startTime.Day())
	assert.Equal(t, 0, startTime.Hour())
	assert.Equal(t, 0, startTime.Minute())
	assert.Equal(t, 0, startTime.Second())

	// 验证结束时间是明天的前一纳秒
	assert.True(t, endTime.After(startTime))
	assert.True(t, endTime.Sub(startTime) < 24*time.Hour)
}

// TestAddThisWeek 测试本周条件
func TestAddThisWeek(t *testing.T) {
	query := NewQuery()
	query.AddThisWeek("created_at")

	assert.Equal(t, 1, len(query.Filters))
	assert.Equal(t, "created_at", query.Filters[0].Field)
	assert.Equal(t, constants.OP_BETWEEN, query.Filters[0].Operator)

	values := query.Filters[0].Value.([]interface{})
	assert.Equal(t, 2, len(values))

	startTime := values[0].(time.Time)
	endTime := values[1].(time.Time)

	// 验证是一周的时间范围
	duration := endTime.Sub(startTime)
	assert.True(t, duration >= 6*24*time.Hour && duration < 7*24*time.Hour+time.Second)
}

// TestAddThisMonth 测试本月条件
func TestAddThisMonth(t *testing.T) {
	query := NewQuery()
	query.AddThisMonth("created_at")

	assert.Equal(t, 1, len(query.Filters))
	assert.Equal(t, "created_at", query.Filters[0].Field)
	assert.Equal(t, constants.OP_BETWEEN, query.Filters[0].Operator)

	values := query.Filters[0].Value.([]interface{})
	assert.Equal(t, 2, len(values))

	startTime := values[0].(time.Time)
	endTime := values[1].(time.Time)
	now := time.Now()

	// 验证开始时间是本月1号
	assert.Equal(t, now.Year(), startTime.Year())
	assert.Equal(t, now.Month(), startTime.Month())
	assert.Equal(t, 1, startTime.Day())

	// 验证结束时间在下个月1号之前
	assert.True(t, endTime.After(startTime))
	nextMonth := startTime.AddDate(0, 1, 0)
	assert.True(t, endTime.Before(nextMonth) || endTime.Equal(nextMonth.Add(-time.Nanosecond)))
}

// TestAddThisYear 测试今年条件
func TestAddThisYear(t *testing.T) {
	query := NewQuery()
	query.AddThisYear("created_at")

	assert.Equal(t, 1, len(query.Filters))
	assert.Equal(t, "created_at", query.Filters[0].Field)
	assert.Equal(t, constants.OP_BETWEEN, query.Filters[0].Operator)

	values := query.Filters[0].Value.([]interface{})
	assert.Equal(t, 2, len(values))

	startTime := values[0].(time.Time)
	endTime := values[1].(time.Time)
	now := time.Now()

	// 验证开始时间是今年1月1号
	assert.Equal(t, now.Year(), startTime.Year())
	assert.Equal(t, time.January, startTime.Month())
	assert.Equal(t, 1, startTime.Day())

	// 验证结束时间在明年1月1号之前
	assert.True(t, endTime.After(startTime))
	nextYear := time.Date(now.Year()+1, 1, 1, 0, 0, 0, 0, now.Location())
	assert.True(t, endTime.Before(nextYear))
}

// =============== 简化方法测试 ===============

// TestSetDistinct 测试设置去重
func TestSetDistinct(t *testing.T) {
	query := NewQuery()
	result := query.SetDistinct()

	assert.Equal(t, query, result, "应该返回同一个查询对象")
	assert.True(t, query.Distinct, "应该设置为去重")
}

// TestPage 测试设置分页
func TestPage(t *testing.T) {
	query := NewQuery()
	result := query.Page(2, 20)

	assert.Equal(t, query, result, "应该返回同一个查询对象")
	assert.NotNil(t, query.Pagination)
	assert.Equal(t, int32(2), query.Pagination.Page)
	assert.Equal(t, int32(20), query.Pagination.PageSize)

	// 测试边界情况
	query = NewQuery()
	query.Page(0, -5) // 无效值应该被修正
	assert.Equal(t, int32(1), query.Pagination.Page)
	assert.Equal(t, int32(10), query.Pagination.PageSize)
}

// TestTake 测试设置限制数量
func TestTake(t *testing.T) {
	query := NewQuery()
	result := query.Take(10)

	assert.Equal(t, query, result, "应该返回同一个查询对象")
	assert.NotNil(t, query.LimitValue)
	assert.Equal(t, 10, *query.LimitValue)
}

// TestSkip 测试设置跳过数量
func TestSkip(t *testing.T) {
	query := NewQuery()
	result := query.Skip(20)

	assert.Equal(t, query, result, "应该返回同一个查询对象")
	assert.NotNil(t, query.OffsetValue)
	assert.Equal(t, 20, *query.OffsetValue)
}

// =============== 链式调用测试 ===============

// TestChainedCalls_BasicConditions 测试基础条件的链式调用
func TestChainedCalls_BasicConditions(t *testing.T) {
	query := NewQuery()
	result := query.
		AddEqual("status", 1).
		AddNotEqual("type", "deleted").
		AddLike("name", "test").
		AddStartsWith("username", "admin").
		AddEndsWith("email", "@example.com").
		AddIn("category", "tech", "business").
		AddNotIn("tag", "spam", "adult").
		AddGreaterThan("age", 18).
		AddLessEqual("score", 100)

	assert.Equal(t, query, result, "链式调用应该返回同一个对象")
	assert.Equal(t, 9, len(query.Filters), "应该添加9个过滤条件")
}

// TestChainedCalls_TimeConditions 测试时间条件的链式调用
func TestChainedCalls_TimeConditions(t *testing.T) {
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)

	query := NewQuery()
	query.
		AddTimeAfter("created_at", yesterday).
		AddTimeBefore("updated_at", now).
		AddTimeBetween("last_login", yesterday, now).
		AddToday("register_date").
		AddThisWeek("activity_date").
		AddThisMonth("payment_date").
		AddThisYear("join_date")

	// AddTimeBetween 添加2个条件，其他时间方法各添加1个条件
	// 总共: 1 + 1 + 2 + 1 + 1 + 1 + 1 = 8个条件
	assert.Equal(t, 8, len(query.Filters), "应该添加8个时间过滤条件")
}

// TestChainedCalls_SortingAndPaging 测试排序和分页的链式调用
func TestChainedCalls_SortingAndPaging(t *testing.T) {
	query := NewQuery()
	query.
		AddEqual("status", 1).
		AddOrderAsc("name").
		AddOrderDesc("created_at").
		SetDistinct().
		Page(1, 20).
		Take(100).
		Skip(50)

	assert.Equal(t, 1, len(query.Filters))
	assert.Equal(t, 2, len(query.Orders))
	assert.True(t, query.Distinct)
	assert.NotNil(t, query.Pagination)
	assert.NotNil(t, query.LimitValue)
	assert.NotNil(t, query.OffsetValue)
}

// TestChainedCalls_ComplexQuery 测试复杂查询的链式调用
func TestChainedCalls_ComplexQuery(t *testing.T) {
	now := time.Now()
	lastMonth := now.AddDate(0, -1, 0)

	query := NewQuery()
	query.
		AddEqual("tenant_id", 123).
		AddIn("status", "active", "pending", "processing").
		AddLike("description", "important").
		AddStartsWith("code", "ORD").
		AddBetween("amount", 100, 10000).
		AddTimeAfter("created_at", lastMonth).
		AddIsNotNull("email").
		AddOrderDesc("created_at").
		AddOrderAsc("amount").
		Page(1, 50)

	assert.Equal(t, 7, len(query.Filters), "应该添加7个过滤条件")
	assert.Equal(t, 2, len(query.Orders), "应该添加2个排序条件")
	assert.NotNil(t, query.Pagination, "应该设置分页")
}

// =============== 边界条件测试 ===============

// TestEmptyStringFilters 测试空字符串过滤
func TestEmptyStringFilters(t *testing.T) {
	query := NewQuery()
	query.
		AddLike("name", "").       // 应该被忽略
		AddStartsWith("code", ""). // 应该被忽略
		AddEndsWith("email", "").  // 应该被忽略
		AddEqual("status", "").    // 不应该被忽略（空字符串是有效值）
		AddNotEqual("type", "")    // 不应该被忽略

	assert.Equal(t, 2, len(query.Filters), "只有Equal和NotEqual应该添加空字符串条件")
}

// TestEmptySliceFilters 测试空切片过滤
func TestEmptySliceFilters(t *testing.T) {
	query := NewQuery()
	query.
		AddIn("status").     // 空参数，应该被忽略
		AddNotIn("type").    // 空参数，应该被忽略
		AddEqual("name", "") // 有效的空字符串

	assert.Equal(t, 1, len(query.Filters), "只有Equal应该添加条件")
}

// TestZeroTimeFilters 测试零时间过滤
func TestZeroTimeFilters(t *testing.T) {
	zeroTime := time.Time{}
	validTime := time.Now()

	query := NewQuery()
	query.
		AddTimeAfter("created_at", zeroTime).             // 应该被忽略
		AddTimeBefore("updated_at", zeroTime).            // 应该被忽略
		AddTimeBetween("deleted_at", zeroTime, zeroTime). // 应该被忽略
		AddTimeAfter("active_at", validTime)              // 应该被添加

	assert.Equal(t, 1, len(query.Filters), "只有有效时间条件应该被添加")
}

// TestMethodReturnValues 测试方法返回值一致性
func TestMethodReturnValues(t *testing.T) {
	query := NewQuery()

	// 所有方法都应该返回同一个查询对象以支持链式调用
	assert.Equal(t, query, query.AddEqual("a", 1))
	assert.Equal(t, query, query.AddNotEqual("b", 2))
	assert.Equal(t, query, query.AddLike("c", "test"))
	assert.Equal(t, query, query.AddStartsWith("d", "pre"))
	assert.Equal(t, query, query.AddEndsWith("e", "suf"))
	assert.Equal(t, query, query.AddIn("f", 1, 2))
	assert.Equal(t, query, query.AddNotIn("g", 3, 4))
	assert.Equal(t, query, query.AddGreaterThan("h", 10))
	assert.Equal(t, query, query.AddGreaterEqual("i", 20))
	assert.Equal(t, query, query.AddLessThan("j", 30))
	assert.Equal(t, query, query.AddLessEqual("k", 40))
	assert.Equal(t, query, query.AddBetween("l", 1, 10))
	assert.Equal(t, query, query.AddIsNull("m"))
	assert.Equal(t, query, query.AddIsNotNull("n"))
	assert.Equal(t, query, query.AddOrderAsc("o"))
	assert.Equal(t, query, query.AddOrderDesc("p"))
	assert.Equal(t, query, query.AddTimeAfter("q", time.Now()))
	assert.Equal(t, query, query.AddTimeBefore("r", time.Now()))
	assert.Equal(t, query, query.AddTimeBetween("s", time.Now(), time.Now()))
	assert.Equal(t, query, query.AddToday("t"))
	assert.Equal(t, query, query.AddThisWeek("u"))
	assert.Equal(t, query, query.AddThisMonth("v"))
	assert.Equal(t, query, query.AddThisYear("w"))
	assert.Equal(t, query, query.SetDistinct())
	assert.Equal(t, query, query.Page(1, 10))
	assert.Equal(t, query, query.Take(50))
	assert.Equal(t, query, query.Skip(100))
}

// =============== 性能测试 ===============

// TestLargeChainedQuery 测试大量链式调用的性能
func TestLargeChainedQuery(t *testing.T) {
	query := NewQuery()

	// 构建一个包含大量条件的查询
	for i := 0; i < 100; i++ {
		query.AddEqual("field_"+string(rune(i)), i)
	}

	assert.Equal(t, 100, len(query.Filters), "应该添加100个过滤条件")

	// 验证第一个和最后一个条件
	assert.Equal(t, "field_\x00", query.Filters[0].Field)
	assert.Equal(t, 0, query.Filters[0].Value)
	assert.Equal(t, "field_c", query.Filters[99].Field)
	assert.Equal(t, 99, query.Filters[99].Value)
}

// =============== 集成测试 ===============

// TestRealWorldScenario 测试真实世界场景
func TestRealWorldScenario(t *testing.T) {
	// 模拟一个电商订单查询场景
	startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC)

	query := NewQuery()
	query.
		AddEqual("tenant_id", 1001).                             // 租户过滤
		AddIn("status", "pending", "processing", "shipped").     // 状态过滤
		AddGreaterEqual("amount", 100).                          // 最小金额
		AddLessEqual("amount", 50000).                           // 最大金额
		AddStartsWith("order_no", "ORD2025").                    // 订单号前缀
		AddIsNotNull("customer_email").                          // 必须有邮箱
		AddTimeBetween("created_at", startDate, endDate).        // 时间范围
		AddLike("shipping_address", "北京").                       // 地址包含
		AddNotIn("payment_method", "cash_on_delivery", "check"). // 排除支付方式
		AddOrderDesc("created_at").                              // 按创建时间降序
		AddOrderAsc("amount").                                   // 按金额升序
		Page(1, 20)                                              // 分页

	// 验证查询条件数量
	expectedFilters := 10 // AddTimeBetween 添加2个条件，其他8个方法各1个：8+2=10
	assert.Equal(t, expectedFilters, len(query.Filters))

	// 验证排序条件
	assert.Equal(t, 2, len(query.Orders))
	assert.Equal(t, "created_at", query.Orders[0].Field)
	assert.Equal(t, constants.Desc, query.Orders[0].Direction)
	assert.Equal(t, "amount", query.Orders[1].Field)
	assert.Equal(t, constants.Asc, query.Orders[1].Direction)

	// 验证分页
	assert.NotNil(t, query.Pagination)
	assert.Equal(t, int32(1), query.Pagination.Page)
	assert.Equal(t, int32(20), query.Pagination.PageSize)
}
