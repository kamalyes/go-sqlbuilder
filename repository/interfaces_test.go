/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-23 08:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-23 08:11:50
 * @FilePath: \go-sqlbuilder\repository\interfaces_test.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package repository

import (
	"testing"

	"github.com/kamalyes/go-sqlbuilder/constants"
	"github.com/stretchr/testify/assert"
)

// TestNewEqFilter 测试等于过滤器创建
func TestNewEqFilter(t *testing.T) {
	filter := NewEqFilter("name", "test")

	assert.NotNil(t, filter, "过滤器不应为空")
	assert.Equal(t, "name", filter.Field, "字段名应为 'name'")
	assert.Equal(t, constants.OP_EQ, filter.Operator, "操作符应为 constants.OP_EQ")
	assert.Equal(t, "test", filter.Value, "值应为 'test'")
}

// TestNewGtFilter 测试大于过滤器创建
func TestNewGtFilter(t *testing.T) {
	filter := NewGtFilter("age", 18)

	assert.NotNil(t, filter, "过滤器不应为空")
	assert.Equal(t, "age", filter.Field, "字段名应为 'age'")
	assert.Equal(t, constants.OP_GT, filter.Operator, "操作符应为 constants.OP_GT")
	assert.Equal(t, 18, filter.Value, "值应为 18")
}

// TestNewLtFilter 测试小于过滤器创建
func TestNewLtFilter(t *testing.T) {
	filter := NewLtFilter("price", 100.0)

	assert.NotNil(t, filter, "过滤器不应为空")
	assert.Equal(t, "price", filter.Field, "字段名应为 'price'")
	assert.Equal(t, constants.OP_LT, filter.Operator, "操作符应为 constants.OP_LT")
	assert.Equal(t, 100.0, filter.Value, "值应为 100.0")
}

// TestNewGteFilter 测试大于等于过滤器创建
func TestNewGteFilter(t *testing.T) {
	filter := NewGteFilter("score", 60)

	assert.NotNil(t, filter, "过滤器不应为空")
	assert.Equal(t, "score", filter.Field, "字段名应为 'score'")
	assert.Equal(t, constants.OP_GTE, filter.Operator, "操作符应为 constants.OP_GTE")
	assert.Equal(t, 60, filter.Value, "值应为 60")
}

// TestNewLteFilter 测试小于等于过滤器创建
func TestNewLteFilter(t *testing.T) {
	filter := NewLteFilter("level", 5)

	assert.NotNil(t, filter, "过滤器不应为空")
	assert.Equal(t, "level", filter.Field, "字段名应为 'level'")
	assert.Equal(t, constants.OP_LTE, filter.Operator, "操作符应为 constants.OP_LTE")
	assert.Equal(t, 5, filter.Value, "值应为 5")
}

// TestNewInFilter 测试IN过滤器创建
func TestNewInFilter(t *testing.T) {
	filter := NewInFilter("status", "active", "pending", "inactive")

	assert.NotNil(t, filter, "过滤器不应为空")
	assert.Equal(t, "status", filter.Field, "字段名应为 'status'")
	assert.Equal(t, constants.OP_IN, filter.Operator, "操作符应为 constants.OP_IN")

	values, ok := filter.Value.([]interface{})
	assert.True(t, ok, "值应为 []interface{} 类型")
	assert.Len(t, values, 3, "值列表应包含 3 个元素")
	assert.Contains(t, values, "active", "值列表应包含 'active'")
	assert.Contains(t, values, "pending", "值列表应包含 'pending'")
	assert.Contains(t, values, "inactive", "值列表应包含 'inactive'")
}

// TestNewInFilterEmpty 测试空IN过滤器创建
func TestNewInFilterEmpty(t *testing.T) {
	filter := NewInFilter("status")

	assert.NotNil(t, filter, "过滤器不应为空")
	assert.Equal(t, "status", filter.Field, "字段名应为 'status'")
	assert.Equal(t, constants.OP_IN, filter.Operator, "操作符应为 constants.OP_IN")

	values, ok := filter.Value.([]interface{})
	assert.True(t, ok, "值应为 []interface{} 类型")
	assert.Len(t, values, 0, "值列表应为空")
}

// TestNewLikeFilter 测试LIKE过滤器创建
func TestNewLikeFilter(t *testing.T) {
	filter := NewLikeFilter("title", "test")

	assert.NotNil(t, filter, "过滤器不应为空")
	assert.Equal(t, "title", filter.Field, "字段名应为 'title'")
	assert.Equal(t, constants.OP_LIKE, filter.Operator, "操作符应为 constants.OP_LIKE")
	assert.Equal(t, "%test%", filter.Value, "值应为 '%test%'")
}

// TestNewNeqFilter 测试不等于过滤器创建
func TestNewNeqFilter(t *testing.T) {
	filter := NewNeqFilter("id", 0)

	assert.NotNil(t, filter, "过滤器不应为空")
	assert.Equal(t, "id", filter.Field, "字段名应为 'id'")
	assert.Equal(t, constants.OP_NEQ, filter.Operator, "操作符应为 constants.OP_NEQ")
	assert.Equal(t, 0, filter.Value, "值应为 0")
}

// TestNewBetweenFilter 测试BETWEEN过滤器创建
func TestNewBetweenFilter(t *testing.T) {
	filter := NewBetweenFilter("created_at", "2023-01-01", "2023-12-31")

	assert.NotNil(t, filter, "过滤器不应为空")
	assert.Equal(t, "created_at", filter.Field, "字段名应为 'created_at'")
	assert.Equal(t, constants.OP_BETWEEN, filter.Operator, "操作符应为 constants.OP_BETWEEN")

	values, ok := filter.Value.([]interface{})
	assert.True(t, ok, "值应为 []interface{} 类型")
	assert.Len(t, values, 2, "值列表应包含 2 个元素")
	assert.Equal(t, "2023-01-01", values[0], "第一个值应为 '2023-01-01'")
	assert.Equal(t, "2023-12-31", values[1], "第二个值应为 '2023-12-31'")
}

// TestNewNotInFilter 测试NOT IN过滤器创建
func TestNewNotInFilter(t *testing.T) {
	filter := NewNotInFilter("category", "archived", "deleted")

	assert.NotNil(t, filter, "过滤器不应为空")
	assert.Equal(t, "category", filter.Field, "字段名应为 'category'")
	assert.Equal(t, constants.OP_NOT_IN, filter.Operator, "操作符应为 constants.OP_NOT_IN")

	values, ok := filter.Value.([]interface{})
	assert.True(t, ok, "值应为 []interface{} 类型")
	assert.Len(t, values, 2, "值列表应包含 2 个元素")
	assert.Contains(t, values, "archived", "值列表应包含 'archived'")
	assert.Contains(t, values, "deleted", "值列表应包含 'deleted'")
}

// TestNewIsNullFilter 测试IS NULL过滤器创建
func TestNewIsNullFilter(t *testing.T) {
	filter := NewIsNullFilter("deleted_at")

	assert.NotNil(t, filter, "过滤器不应为空")
	assert.Equal(t, "deleted_at", filter.Field, "字段名应为 'deleted_at'")
	assert.Equal(t, constants.OP_IS_NULL, filter.Operator, "操作符应为 constants.OP_IS_NULL")
	assert.Nil(t, filter.Value, "值应为 nil")
}

// TestNewIsNotNullFilter 测试IS NOT NULL过滤器创建
func TestNewIsNotNullFilter(t *testing.T) {
	filter := NewIsNotNullFilter("updated_at")

	assert.NotNil(t, filter, "过滤器不应为空")
	assert.Equal(t, "updated_at", filter.Field, "字段名应为 'updated_at'")
	assert.Equal(t, constants.OP_IS_NOT_NULL, filter.Operator, "操作符应为 constants.OP_IS_NOT_NULL")
	assert.Nil(t, filter.Value, "值应为 nil")
}

// TestNewStartsWithFilter 测试前缀匹配过滤器创建
func TestNewStartsWithFilter(t *testing.T) {
	filter := NewStartsWithFilter("name", "user")

	assert.NotNil(t, filter, "过滤器不应为空")
	assert.Equal(t, "name", filter.Field, "字段名应为 'name'")
	assert.Equal(t, constants.OP_LIKE, filter.Operator, "操作符应为 constants.OP_LIKE")
	assert.Equal(t, "user%", filter.Value, "值应为 'user%'")
}

// TestNewEndsWithFilter 测试后缀匹配过滤器创建
func TestNewEndsWithFilter(t *testing.T) {
	filter := NewEndsWithFilter("email", "@example.com")

	assert.NotNil(t, filter, "过滤器不应为空")
	assert.Equal(t, "email", filter.Field, "字段名应为 'email'")
	assert.Equal(t, constants.OP_LIKE, filter.Operator, "操作符应为 constants.OP_LIKE")
	assert.Equal(t, "%@example.com", filter.Value, "值应为 '%@example.com'")
}

// TestNewContainsFilter 测试包含匹配过滤器创建
func TestNewContainsFilter(t *testing.T) {
	filter := NewContainsFilter("description", "keyword")

	assert.NotNil(t, filter, "过滤器不应为空")
	assert.Equal(t, "description", filter.Field, "字段名应为 'description'")
	assert.Equal(t, constants.OP_LIKE, filter.Operator, "操作符应为 constants.OP_LIKE")
	assert.Equal(t, "%keyword%", filter.Value, "值应为 '%keyword%'")
}

// TestNewNotLikeFilter 测试NOT LIKE过滤器创建
func TestNewNotLikeFilter(t *testing.T) {
	filter := NewNotLikeFilter("content", "spam")

	assert.NotNil(t, filter, "过滤器不应为空")
	assert.Equal(t, "content", filter.Field, "字段名应为 'content'")
	assert.Equal(t, constants.OP_NOT_LIKE, filter.Operator, "操作符应为 constants.OP_NOT_LIKE")
	assert.Equal(t, "%spam%", filter.Value, "值应为 '%spam%'")
}

// TestNewFindInSetFilter 测试FIND_IN_SET过滤器创建
func TestNewFindInSetFilter(t *testing.T) {
	filter := NewFindInSetFilter("tags", "important")

	assert.NotNil(t, filter, "过滤器不应为空")
	assert.Equal(t, "tags", filter.Field, "字段名应为 'tags'")
	assert.Equal(t, constants.OP_FIND_IN_SET, filter.Operator, "操作符应为 constants.OP_FIND_IN_SET")
	assert.Equal(t, "important", filter.Value, "值应为 'important'")
}

// TestNewFilterGroup 测试过滤器组创建
func TestNewFilterGroup(t *testing.T) {
	// 测试 AND 逻辑组
	andGroup := NewFilterGroup(constants.LOGIC_AND)
	assert.NotNil(t, andGroup, "AND 过滤器组不应为空")
	assert.Equal(t, constants.LOGIC_AND, andGroup.LogicOp, "逻辑操作符应为 AND")
	assert.Empty(t, andGroup.Filters, "过滤器列表应为空")
	assert.Empty(t, andGroup.Groups, "子组列表应为空")
	assert.True(t, andGroup.IsEmpty(), "过滤器组应为空")
	assert.Equal(t, 0, andGroup.Count(), "过滤器计数应为 0")

	// 测试 OR 逻辑组
	orGroup := NewFilterGroup(constants.LOGIC_OR)
	assert.NotNil(t, orGroup, "OR 过滤器组不应为空")
	assert.Equal(t, constants.LOGIC_OR, orGroup.LogicOp, "逻辑操作符应为 OR")

	// 测试无效逻辑操作符（应该默认为 AND）
	invalidGroup := NewFilterGroup("INVALID")
	assert.NotNil(t, invalidGroup, "无效逻辑组不应为空")
	assert.Equal(t, constants.LOGIC_AND, invalidGroup.LogicOp, "无效逻辑操作符应默认为 AND")
}

// TestFilterGroupAddFilter 测试向过滤器组添加过滤器
func TestFilterGroupAddFilter(t *testing.T) {
	group := NewFilterGroup(constants.LOGIC_AND)
	filter1 := NewEqFilter("name", "test")
	filter2 := NewGtFilter("age", 18)

	// 添加单个过滤器
	result := group.AddFilter(filter1)
	assert.Same(t, group, result, "应返回自身以支持链式调用")
	assert.Len(t, group.Filters, 1, "应包含 1 个过滤器")
	assert.Equal(t, filter1, group.Filters[0], "应包含添加的过滤器")
	assert.False(t, group.IsEmpty(), "过滤器组不应为空")
	assert.Equal(t, 1, group.Count(), "过滤器计数应为 1")

	// 添加另一个过滤器
	group.AddFilter(filter2)
	assert.Len(t, group.Filters, 2, "应包含 2 个过滤器")
	assert.Equal(t, 2, group.Count(), "过滤器计数应为 2")

	// 添加 nil 过滤器（应忽略）
	group.AddFilter(nil)
	assert.Len(t, group.Filters, 2, "nil 过滤器应被忽略")
}

// TestFilterGroupAddFilters 测试向过滤器组批量添加过滤器
func TestFilterGroupAddFilters(t *testing.T) {
	group := NewFilterGroup(constants.LOGIC_OR)
	filter1 := NewEqFilter("status", "active")
	filter2 := NewGtFilter("score", 80)
	filter3 := NewLikeFilter("title", "important")

	// 批量添加过滤器
	result := group.AddFilters(filter1, filter2, filter3)
	assert.Same(t, group, result, "应返回自身以支持链式调用")
	assert.Len(t, group.Filters, 3, "应包含 3 个过滤器")
	assert.Equal(t, 3, group.Count(), "过滤器计数应为 3")

	// 批量添加包含 nil 的过滤器
	group.AddFilters(nil, NewInFilter("type", "A", "B"))
	assert.Len(t, group.Filters, 4, "应正确处理 nil 过滤器")
}

// TestFilterGroupAddGroup 测试向过滤器组添加子组
func TestFilterGroupAddGroup(t *testing.T) {
	parentGroup := NewFilterGroup(constants.LOGIC_AND)
	childGroup := NewFilterGroup(constants.LOGIC_OR)
	childGroup.AddFilter(NewEqFilter("category", "tech"))
	childGroup.AddFilter(NewEqFilter("category", "science"))

	// 添加子组
	result := parentGroup.AddGroup(childGroup)
	assert.Same(t, parentGroup, result, "应返回自身以支持链式调用")
	assert.Len(t, parentGroup.Groups, 1, "应包含 1 个子组")
	assert.Equal(t, childGroup, parentGroup.Groups[0], "应包含添加的子组")
	assert.Equal(t, 2, parentGroup.Count(), "计数应包含子组中的过滤器")

	// 添加 nil 子组（应忽略）
	parentGroup.AddGroup(nil)
	assert.Len(t, parentGroup.Groups, 1, "nil 子组应被忽略")
}

// TestNewQuery 测试查询对象创建
func TestNewQuery(t *testing.T) {
	query := NewQuery()

	assert.NotNil(t, query, "查询对象不应为空")
	assert.NotNil(t, query.Filters, "过滤器列表不应为空")
	assert.NotNil(t, query.Orders, "排序列表不应为空")
	assert.Empty(t, query.Filters, "过滤器列表应为空")
	assert.Empty(t, query.Orders, "排序列表应为空")
	assert.Nil(t, query.Pagination, "分页信息应为空")
	assert.Nil(t, query.LimitValue, "限制数量应为空")
	assert.Nil(t, query.OffsetValue, "偏移量应为空")
	assert.False(t, query.Distinct, "去重标记应为 false")
	assert.Nil(t, query.GroupBy, "分组字段应为空")
	assert.Nil(t, query.Having, "HAVING 条件应为空")
	assert.False(t, query.HasFilters(), "应没有过滤条件")
}

// TestQueryAddFilter 测试向查询添加过滤器
func TestQueryAddFilter(t *testing.T) {
	query := NewQuery()
	filter := NewEqFilter("status", "published")

	// 添加过滤器
	result := query.AddFilter(filter)
	assert.Same(t, query, result, "应返回自身以支持链式调用")
	assert.Len(t, query.Filters, 1, "应包含 1 个过滤器")
	assert.Equal(t, filter, query.Filters[0], "应包含添加的过滤器")
	assert.True(t, query.HasFilters(), "应有过滤条件")

	// 添加 nil 过滤器
	query.AddFilter(nil)
	assert.Len(t, query.Filters, 1, "nil 过滤器应被忽略")
}

// TestQueryAddFilters 测试向查询批量添加过滤器
func TestQueryAddFilters(t *testing.T) {
	query := NewQuery()
	filter1 := NewGtFilter("views", 1000)
	filter2 := NewLtFilter("created_at", "2024-01-01")

	// 批量添加过滤器
	result := query.AddFilters(filter1, filter2)
	assert.Same(t, query, result, "应返回自身以支持链式调用")
	assert.Len(t, query.Filters, 2, "应包含 2 个过滤器")
	assert.True(t, query.HasFilters(), "应有过滤条件")

	allFilters := query.GetAllFilters()
	assert.Len(t, allFilters, 2, "扁平化后应包含 2 个过滤器")
}

// TestQueryAddOrder 测试向查询添加排序
func TestQueryAddOrder(t *testing.T) {
	query := NewQuery()

	// 添加排序
	result := query.AddOrder("created_at", "DESC")
	assert.Same(t, query, result, "应返回自身以支持链式调用")
	assert.Len(t, query.Orders, 1, "应包含 1 个排序条件")
	assert.Equal(t, "created_at", query.Orders[0].Field, "排序字段应为 'created_at'")
	assert.Equal(t, "DESC", query.Orders[0].Direction, "排序方向应为 'DESC'")

	// 添加另一个排序
	query.AddOrder("title", "ASC")
	assert.Len(t, query.Orders, 2, "应包含 2 个排序条件")
}

// TestQueryWithPaging 测试查询分页设置
func TestQueryWithPaging(t *testing.T) {
	query := NewQuery()

	// 设置正常分页
	result := query.WithPaging(2, 20)
	assert.Same(t, query, result, "应返回自身以支持链式调用")
	assert.NotNil(t, query.Pagination, "分页信息不应为空")
	assert.Equal(t, int32(2), query.Pagination.Page, "页码应为 2")
	assert.Equal(t, int32(20), query.Pagination.PageSize, "页大小应为 20")

	// 设置无效分页（应使用默认值）
	query.WithPaging(0, -5)
	assert.Equal(t, int32(1), query.Pagination.Page, "无效页码应默认为 1")
	assert.Equal(t, int32(10), query.Pagination.PageSize, "无效页大小应默认为 10")
}

// TestQueryLimit 测试查询限制设置
func TestQueryLimit(t *testing.T) {
	query := NewQuery()

	// 设置限制
	result := query.Limit(50)
	assert.Same(t, query, result, "应返回自身以支持链式调用")
	assert.NotNil(t, query.LimitValue, "限制值不应为空")
	assert.Equal(t, 50, *query.LimitValue, "限制值应为 50")
}

// TestQueryOffset 测试查询偏移设置
func TestQueryOffset(t *testing.T) {
	query := NewQuery()

	// 设置偏移
	result := query.Offset(100)
	assert.Same(t, query, result, "应返回自身以支持链式调用")
	assert.NotNil(t, query.OffsetValue, "偏移值不应为空")
	assert.Equal(t, 100, *query.OffsetValue, "偏移值应为 100")
}

// TestQueryWithFilterGroup 测试查询设置过滤器组
func TestQueryWithFilterGroup(t *testing.T) {
	query := NewQuery()
	group := NewFilterGroup(constants.LOGIC_OR)
	group.AddFilter(NewEqFilter("type", "A"))
	group.AddFilter(NewEqFilter("type", "B"))

	// 设置过滤器组
	result := query.WithFilterGroup(group)
	assert.Same(t, query, result, "应返回自身以支持链式调用")
	assert.Equal(t, group, query.FilterGroup, "应设置指定的过滤器组")
	assert.True(t, query.HasFilters(), "应有过滤条件")

	allFilters := query.GetAllFilters()
	assert.Len(t, allFilters, 2, "扁平化后应包含组中的过滤器")
}

// TestQueryWithDistinct 测试查询去重设置
func TestQueryWithDistinct(t *testing.T) {
	query := NewQuery()

	// 设置去重
	result := query.WithDistinct(true)
	assert.Same(t, query, result, "应返回自身以支持链式调用")
	assert.True(t, query.Distinct, "去重标记应为 true")

	// 取消去重
	query.WithDistinct(false)
	assert.False(t, query.Distinct, "去重标记应为 false")
}

// TestQueryAddGroupBy 测试查询添加分组字段
func TestQueryAddGroupBy(t *testing.T) {
	query := NewQuery()

	// 添加分组字段
	result := query.AddGroupBy("category", "status")
	assert.Same(t, query, result, "应返回自身以支持链式调用")
	assert.NotNil(t, query.GroupBy, "分组字段不应为空")
	assert.Len(t, query.GroupBy, 2, "应包含 2 个分组字段")
	assert.Contains(t, query.GroupBy, "category", "应包含 'category' 字段")
	assert.Contains(t, query.GroupBy, "status", "应包含 'status' 字段")

	// 继续添加分组字段
	query.AddGroupBy("region")
	assert.Len(t, query.GroupBy, 3, "应包含 3 个分组字段")
}

// TestQueryAddHaving 测试查询添加HAVING条件
func TestQueryAddHaving(t *testing.T) {
	query := NewQuery()
	havingFilter := NewGtFilter("COUNT(*)", 5)

	// 添加HAVING条件
	result := query.AddHaving(havingFilter)
	assert.Same(t, query, result, "应返回自身以支持链式调用")
	assert.NotNil(t, query.Having, "HAVING 条件不应为空")
	assert.Len(t, query.Having, 1, "应包含 1 个HAVING条件")
	assert.Equal(t, havingFilter, query.Having[0], "应包含添加的HAVING条件")

	// 添加 nil HAVING条件
	query.AddHaving(nil)
	assert.Len(t, query.Having, 1, "nil HAVING条件应被忽略")
}

// TestQueryComplexFiltering 测试复杂过滤条件组合
func TestQueryComplexFiltering(t *testing.T) {
	query := NewQuery()

	// 添加简单过滤器
	query.AddFilter(NewEqFilter("published", true))

	// 创建复合过滤器组
	categoryGroup := NewFilterGroup(constants.LOGIC_OR)
	categoryGroup.AddFilter(NewEqFilter("category", "tech"))
	categoryGroup.AddFilter(NewEqFilter("category", "science"))

	dateGroup := NewFilterGroup(constants.LOGIC_AND)
	dateGroup.AddFilter(NewGteFilter("created_at", "2023-01-01"))
	dateGroup.AddFilter(NewLteFilter("created_at", "2023-12-31"))

	// 创建主过滤器组
	mainGroup := NewFilterGroup(constants.LOGIC_AND)
	mainGroup.AddGroup(categoryGroup)
	mainGroup.AddGroup(dateGroup)

	query.WithFilterGroup(mainGroup)

	// 验证过滤器结构
	assert.True(t, query.HasFilters(), "应有过滤条件")
	assert.Len(t, query.Filters, 1, "应有 1 个简单过滤器")
	assert.NotNil(t, query.FilterGroup, "应有复合过滤器组")
	assert.Equal(t, 4, query.FilterGroup.Count(), "复合过滤器组应包含 4 个条件")

	allFilters := query.GetAllFilters()
	assert.Len(t, allFilters, 5, "扁平化后应包含 5 个过滤器")
}

// TestFindOptions 测试兼容旧API的查询选项
func TestFindOptions(t *testing.T) {
	options := &FindOptions{
		Conditions: []Condition{
			{Field: "status", Op: constants.OP_EQ, Value: "active"},
			{Field: "priority", Op: constants.OP_GT, Value: 3},
		},
		Orders: []OrderBy{
			{Field: "created_at", Direction: "DESC"},
			{Field: "title", Direction: "ASC"},
		},
		Limit:  10,
		Offset: 20,
	}

	assert.Len(t, options.Conditions, 2, "应包含 2 个查询条件")
	assert.Len(t, options.Orders, 2, "应包含 2 个排序条件")
	assert.Equal(t, 10, options.Limit, "限制数量应为 10")
	assert.Equal(t, 20, options.Offset, "偏移量应为 20")

	// 验证条件内容
	assert.Equal(t, "status", options.Conditions[0].Field, "第一个条件字段应为 'status'")
	assert.Equal(t, constants.OP_EQ, options.Conditions[0].Op, "第一个条件操作符应为 constants.OP_EQ")
	assert.Equal(t, "active", options.Conditions[0].Value, "第一个条件值应为 'active'")

	// 验证排序内容
	assert.Equal(t, "created_at", options.Orders[0].Field, "第一个排序字段应为 'created_at'")
	assert.Equal(t, "DESC", options.Orders[0].Direction, "第一个排序方向应为 'DESC'")
}

// TestMetaPaging 测试分页元数据
func TestMetaPaging(t *testing.T) {
	paging := &Pagination{
		Page:     2,
		PageSize: 15,
		Total:    150,
	}

	assert.Equal(t, int32(2), paging.Page, "页码应为 2")
	assert.Equal(t, int32(15), paging.PageSize, "页大小应为 15")
	assert.Equal(t, int64(150), paging.Total, "总数应为 150")
}

// TestFilterGroupNesting 测试过滤器组嵌套
func TestFilterGroupNesting(t *testing.T) {
	// 创建三级嵌套的过滤器组
	level3Group := NewFilterGroup(constants.LOGIC_OR)
	level3Group.AddFilter(NewEqFilter("tag", "urgent"))
	level3Group.AddFilter(NewEqFilter("tag", "important"))

	level2Group := NewFilterGroup(constants.LOGIC_AND)
	level2Group.AddFilter(NewGtFilter("priority", 5))
	level2Group.AddGroup(level3Group)

	level1Group := NewFilterGroup(constants.LOGIC_OR)
	level1Group.AddFilter(NewEqFilter("status", "active"))
	level1Group.AddGroup(level2Group)

	// 验证嵌套结构
	assert.False(t, level1Group.IsEmpty(), "顶级组不应为空")
	assert.Equal(t, 4, level1Group.Count(), "顶级组应包含 4 个条件（递归计算）")
	assert.Len(t, level1Group.Filters, 1, "顶级组直接包含 1 个过滤器")
	assert.Len(t, level1Group.Groups, 1, "顶级组包含 1 个子组")

	assert.Equal(t, 3, level2Group.Count(), "二级组应包含 3 个条件")
	assert.Len(t, level2Group.Filters, 1, "二级组直接包含 1 个过滤器")
	assert.Len(t, level2Group.Groups, 1, "二级组包含 1 个子组")

	assert.Equal(t, 2, level3Group.Count(), "三级组应包含 2 个条件")
	assert.Len(t, level3Group.Filters, 2, "三级组直接包含 2 个过滤器")
	assert.Len(t, level3Group.Groups, 0, "三级组不包含子组")
}

// TestQueryFlattenFilters 测试查询过滤器扁平化
func TestQueryFlattenFilters(t *testing.T) {
	query := NewQuery()

	// 添加简单过滤器
	query.AddFilter(NewEqFilter("published", true))
	query.AddFilter(NewGtFilter("views", 100))

	// 创建嵌套过滤器组
	subGroup1 := NewFilterGroup(constants.LOGIC_OR)
	subGroup1.AddFilter(NewEqFilter("category", "A"))
	subGroup1.AddFilter(NewEqFilter("category", "B"))

	subGroup2 := NewFilterGroup(constants.LOGIC_AND)
	subGroup2.AddFilter(NewGteFilter("rating", 4.0))

	mainGroup := NewFilterGroup(constants.LOGIC_AND)
	mainGroup.AddGroup(subGroup1)
	mainGroup.AddGroup(subGroup2)

	query.WithFilterGroup(mainGroup)

	// 验证扁平化结果
	allFilters := query.GetAllFilters()
	assert.Len(t, allFilters, 5, "扁平化后应包含 5 个过滤器")

	// 验证包含的过滤器类型
	hasPublishedFilter := false
	hasViewsFilter := false
	hasCategoryFilters := 0
	hasRatingFilter := false

	for _, filter := range allFilters {
		switch filter.Field {
		case "published":
			hasPublishedFilter = true
		case "views":
			hasViewsFilter = true
		case "category":
			hasCategoryFilters++
		case "rating":
			hasRatingFilter = true
		}
	}

	assert.True(t, hasPublishedFilter, "应包含 published 过滤器")
	assert.True(t, hasViewsFilter, "应包含 views 过滤器")
	assert.Equal(t, 2, hasCategoryFilters, "应包含 2 个 category 过滤器")
	assert.True(t, hasRatingFilter, "应包含 rating 过滤器")
}
