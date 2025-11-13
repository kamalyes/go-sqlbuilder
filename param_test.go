/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-11 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-13 23:21:08
 * @FilePath: \go-sqlbuilder\param_test.go
 * @Description: 缓存和高级查询测试（重构为assert校验）
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package sqlbuilder

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// ==================== 缓存测试 ====================

func TestMockCacheStore(t *testing.T) {
	ctx := context.Background()
	cache := NewMockCacheStore()

	// 测试Set和Get
	err := cache.Set(ctx, "test_key", "test_value", 1*time.Hour)
	assert.NoError(t, err, "Set should not return error")

	value, err := cache.Get(ctx, "test_key")
	assert.NoError(t, err, "Get should not return error for existing key")
	assert.Equal(t, "test_value", value, "Get should return the stored value")

	// 测试不存在的键
	_, err = cache.Get(ctx, "non_exist_key")
	assert.Error(t, err, "Get should return error for non-existent key")

	// 测试Exists
	exists, err := cache.Exists(ctx, "test_key")
	assert.NoError(t, err, "Exists should not return error")
	assert.True(t, exists, "Exists should return true for existing key")

	// 测试Delete
	err = cache.Delete(ctx, "test_key")
	assert.NoError(t, err, "Delete should not return error")

	exists, err = cache.Exists(ctx, "test_key")
	assert.NoError(t, err, "Exists should not return error after delete")
	assert.False(t, exists, "Exists should return false after delete")
}

func TestCacheExpiration(t *testing.T) {
	ctx := context.Background()
	cache := NewMockCacheStore()

	// 设置短TTL
	err := cache.Set(ctx, "exp_key", "exp_value", 100*time.Millisecond)
	assert.NoError(t, err, "Set should not return error")

	// 立即获取应该成功
	value, err := cache.Get(ctx, "exp_key")
	assert.NoError(t, err, "Get should not return error before expiration")
	assert.Equal(t, "exp_value", value, "Get should return value before expiration")

	// 等待过期
	time.Sleep(150 * time.Millisecond)

	// 获取应该失败
	_, err = cache.Get(ctx, "exp_key")
	assert.Error(t, err, "Get should return error for expired key")
}

func TestCacheClear(t *testing.T) {
	ctx := context.Background()
	cache := NewMockCacheStore()

	// 设置多个键
	err1 := cache.Set(ctx, "sqlbuilder:key1", "value1", 1*time.Hour)
	err2 := cache.Set(ctx, "sqlbuilder:key2", "value2", 1*time.Hour)
	err3 := cache.Set(ctx, "other:key3", "value3", 1*time.Hour)
	assert.NoError(t, err1, "Set sqlbuilder:key1 should not return error")
	assert.NoError(t, err2, "Set sqlbuilder:key2 should not return error")
	assert.NoError(t, err3, "Set other:key3 should not return error")

	// 清除sqlbuilder前缀的缓存
	err := cache.Clear(ctx, "sqlbuilder:")
	assert.NoError(t, err, "Clear should not return error")

	// 验证sqlbuilder:前缀的缓存已删除
	_, err = cache.Get(ctx, "sqlbuilder:key1")
	assert.Error(t, err, "sqlbuilder:key1 should be deleted after Clear")

	// other:前缀的缓存应该还在
	value, err := cache.Get(ctx, "other:key3")
	assert.NoError(t, err, "Get other:key3 should not return error")
	assert.Equal(t, "value3", value, "other:key3 should still exist after Clear")
}

func TestCacheStats(t *testing.T) {
	cache := NewMockCacheStore()
	ctx := context.Background()

	// 初始状态
	stats := cache.GetStats()
	assert.Equal(t, 0, stats["total"], "Initial cache should be empty")

	// 添加项目
	err1 := cache.Set(ctx, "key1", "value1", 1*time.Hour)
	err2 := cache.Set(ctx, "key2", "value2", 1*time.Hour)
	assert.NoError(t, err1, "Set key1 should not return error")
	assert.NoError(t, err2, "Set key2 should not return error")

	stats = cache.GetStats()
	assert.Equal(t, 2, stats["total"], "Cache should have 2 items after adding key1 and key2")

	// 添加过期项目
	err3 := cache.Set(ctx, "key3", "value3", 10*time.Millisecond)
	assert.NoError(t, err3, "Set key3 should not return error")
	time.Sleep(50 * time.Millisecond)

	stats = cache.GetStats()
	assert.Equal(t, 2, stats["valid"], "Cache should have 2 valid items after key3 expires")
}

// ==================== 高级查询测试 ====================

func TestAdvancedQueryParam_Filters(t *testing.T) {
	aq := NewAdvancedQueryParam()

	// 测试添加过滤
	aq.AddEQ("status", "active").
		AddLike("name", "test").
		AddIn("id", 1, 2, 3)

	assert.Len(t, aq.Filters, 3, "AdvancedQueryParam should have 3 filters after adding")

	// 验证过滤内容
	assert.Equal(t, "status", aq.Filters[0].Field, "First filter field should be 'status'")
	assert.Equal(t, OP_EQ, aq.Filters[0].Operator, "First filter operator should be OP_EQ")

	assert.Equal(t, "name", aq.Filters[1].Field, "Second filter field should be 'name'")
	assert.Equal(t, OP_LIKE, aq.Filters[1].Operator, "Second filter operator should be OP_LIKE")

	assert.Equal(t, "id", aq.Filters[2].Field, "Third filter field should be 'id'")
	assert.Equal(t, OP_IN, aq.Filters[2].Operator, "Third filter operator should be OP_IN")
}

func TestAdvancedQueryParam_TimeRange(t *testing.T) {
	aq := NewAdvancedQueryParam()

	startTime := time.Now().Add(-24 * time.Hour)
	endTime := time.Now()
	aq.AddTimeRange("created_at", startTime, endTime)

	assert.Len(t, aq.TimeRanges, 1, "AdvancedQueryParam should have 1 time range after adding")
	_, exists := aq.TimeRanges["created_at"]
	assert.True(t, exists, "Time range for 'created_at' should exist")
}

func TestAdvancedQueryParam_Ordering(t *testing.T) {
	aq := NewAdvancedQueryParam()

	aq.AddOrder("created_at", "DESC").
		AddOrderAsc("name")

	assert.Len(t, aq.Orders, 2, "AdvancedQueryParam should have 2 orders after adding")

	assert.Equal(t, "created_at", aq.Orders[0].Field, "First order field should be 'created_at'")
	assert.Equal(t, "DESC", aq.Orders[0].Order, "First order direction should be 'DESC'")

	assert.Equal(t, "name", aq.Orders[1].Field, "Second order field should be 'name'")
	assert.Equal(t, "ASC", aq.Orders[1].Order, "Second order direction should be 'ASC'")
}

func TestAdvancedQueryParam_Pagination(t *testing.T) {
	aq := NewAdvancedQueryParam()
	aq.SetPage(2, 20)

	assert.Equal(t, 2, aq.Page, "Page should be set to 2")
	assert.Equal(t, 20, aq.PageSize, "PageSize should be set to 20")
	assert.Equal(t, 20, aq.Offset, "Offset should be 20 for page 2 with pageSize 20")
}

func TestAdvancedQueryParam_BuildWhereClause(t *testing.T) {
	aq := NewAdvancedQueryParam()

	aq.AddEQ("status", "active").
		AddLike("name", "test").
		AddIn("id", 1, 2, 3)

	where, args := aq.BuildWhereClause()

	assert.NotEmpty(t, where, "WHERE clause should not be empty")
	assert.Len(t, args, 5, "BuildWhereClause should return 5 arguments (active, %test%, 1, 2, 3)")
}

func TestAdvancedQueryParam_FindInSet(t *testing.T) {
	aq := NewAdvancedQueryParam()
	aq.AddFindInSet("tags", "golang", "database", "cache")

	assert.Len(t, aq.FindInSets, 1, "AdvancedQueryParam should have 1 FIND_IN_SET condition after adding")
	values, exists := aq.FindInSets["tags"]
	assert.True(t, exists, "FIND_IN_SET condition for 'tags' should exist")
	assert.Len(t, values, 3, "FIND_IN_SET for 'tags' should have 3 values")
}

func TestAdvancedQueryParam_Distinct(t *testing.T) {
	aq := NewAdvancedQueryParam()
	aq.SetDistinct(true)

	assert.True(t, aq.Distinct, "Distinct should be set to true")
}

func TestAdvancedQueryParam_SelectFields(t *testing.T) {
	aq := NewAdvancedQueryParam()
	aq.SetSelectFields("id", "name", "email")

	assert.Len(t, aq.SelectFields, 3, "AdvancedQueryParam should have 3 select fields after setting")
	assert.Equal(t, []string{"id", "name", "email"}, aq.SelectFields, "Select fields should be [id, name, email]")
}

func TestAdvancedQueryParam_Having(t *testing.T) {
	aq := NewAdvancedQueryParam()

	aq.AddHaving("COUNT(*) > 5").
		AddHaving("SUM(amount) > 100")

	assert.Len(t, aq.HavingClauses, 2, "AdvancedQueryParam should have 2 HAVING clauses after adding")
}

// ==================== 分页响应测试 ====================

func TestPageBean_Initialization(t *testing.T) {
	rows := []string{"row1", "row2", "row3"}
	page := NewPageBean(2, 20, 100, rows)

	assert.Equal(t, 2, page.Page, "PageBean Page should be 2")
	assert.Equal(t, 20, page.PageSize, "PageBean PageSize should be 20")
	assert.Equal(t, 100, page.Total, "PageBean Total should be 100")
	assert.Equal(t, 5, page.Pages, "PageBean Pages should be 5 (100/20)")
}

func TestPageBean_EdgeCases(t *testing.T) {
	// 测试无数据情况
	page := NewPageBean(0, 10, 0, []string{})
	assert.Equal(t, 1, page.Page, "PageBean should default to page 1 for invalid input")
	assert.Equal(t, 1, page.Pages, "PageBean should have 1 page when total is 0")

	// 测试未满一页的情况
	page = NewPageBean(1, 10, 5, []string{"a", "b", "c", "d", "e"})
	assert.Equal(t, 1, page.Pages, "PageBean should have 1 page for 5 items with pageSize 10")

	// 测试需要多页的情况
	page = NewPageBean(1, 10, 25, []string{})
	assert.Equal(t, 3, page.Pages, "PageBean should have 3 pages for 25 items with pageSize 10")
}

// ==================== 查询选项测试 ====================

func TestFindOption_Builder(t *testing.T) {
	opt := NewFindOption().
		WithBusinessId("biz123").
		WithShopId("shop456").
		WithTablePrefix("user_").
		WithCacheTTL(30 * time.Minute).
		WithNoCache()

	assert.Equal(t, "biz123", opt.BusinessId, "FindOption BusinessId should be 'biz123'")
	assert.Equal(t, "shop456", opt.ShopId, "FindOption ShopId should be 'shop456'")
	assert.Equal(t, "user_", opt.TablePrefix, "FindOption TablePrefix should be 'user_'")
	assert.True(t, opt.NoCache, "FindOption NoCache should be true")
	assert.Equal(t, 30*time.Minute, opt.CacheTTL, "FindOption CacheTTL should be 30 minutes")
}

// ==================== BaseInfoFilter测试 ====================

func TestBaseInfoFilter(t *testing.T) {
	filter := BaseInfoFilter{
		DBField:    "status",
		Values:     []interface{}{"active", "pending"},
		ExactMatch: true,
		AllRegex:   false,
	}

	assert.Equal(t, "status", filter.DBField, "BaseInfoFilter DBField should be 'status'")
	assert.Len(t, filter.Values, 2, "BaseInfoFilter should have 2 values")
	assert.True(t, filter.ExactMatch, "BaseInfoFilter ExactMatch should be true")
	assert.False(t, filter.AllRegex, "BaseInfoFilter AllRegex should be false")
}

// ==================== FilterOperator测试 ====================

func TestFilterOperators(t *testing.T) {
	operators := []FilterOperator{
		OP_EQ, OP_NEQ, OP_GT, OP_GTE, OP_LT, OP_LTE,
		OP_LIKE, OP_IN, OP_BETWEEN, OP_IS_NULL, OP_FIND_IN_SET,
	}

	assert.Len(t, operators, 11, "Should have 11 filter operators")
	assert.Equal(t, "=", OP_EQ, "OP_EQ should be '='")
	assert.Equal(t, "LIKE", OP_LIKE, "OP_LIKE should be 'LIKE'")
	assert.Equal(t, "IN", OP_IN, "OP_IN should be 'IN'")
}

// ==================== 高级查询链式调用测试 ====================

func TestAdvancedQueryParam_FluentAPI(t *testing.T) {
	aq := NewAdvancedQueryParam().
		AddEQ("status", "active").
		AddLike("name", "test").
		SetPage(1, 20).
		AddOrder("created_at", "DESC").
		SetDistinct(true).
		SetSelectFields("id", "name", "email")

	assert.Len(t, aq.Filters, 2, "AdvancedQueryParam should have 2 filters after fluent API calls")
	assert.Equal(t, 1, aq.Page, "Page should be 1")
	assert.Equal(t, 20, aq.PageSize, "PageSize should be 20")
	assert.Len(t, aq.Orders, 1, "Should have 1 order after fluent API call")
	assert.True(t, aq.Distinct, "Distinct should be true")
	assert.Len(t, aq.SelectFields, 3, "Should have 3 select fields after fluent API call")
}

// ==================== 批量测试 ====================

func TestCacheAndQueryIntegration(t *testing.T) {
	cache := NewMockCacheStore()
	ctx := context.Background()

	// 模拟高级查询参数的缓存
	aq := NewAdvancedQueryParam().
		AddEQ("status", "active").
		AddLike("name", "test")

	// 生成缓存键
	cacheKey := "query:" + aq.Filters[0].Field
	cacheValue := `{"page":1,"pageSize":10,"total":100,"rows":[]}`

	// 存储缓存
	err := cache.Set(ctx, cacheKey, cacheValue, 1*time.Hour)
	assert.NoError(t, err, "Set cache should not return error")

	// 检索缓存
	value, err := cache.Get(ctx, cacheKey)
	assert.NoError(t, err, "Get cache should not return error")
	assert.Equal(t, cacheValue, value, "Cache value should match the stored value")
}

// ==================== 性能相关测试 ====================

func BenchmarkAdvancedQueryParam_BuildWhereClause(b *testing.B) {
	aq := NewAdvancedQueryParam()
	for i := 0; i < 10; i++ {
		aq.AddEQ("field"+string(rune(i)), "value"+string(rune(i)))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		aq.BuildWhereClause()
	}
}

func BenchmarkMockCacheStore_SetGet(b *testing.B) {
	cache := NewMockCacheStore()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Set(ctx, "key", "value", 1*time.Hour)
		cache.Get(ctx, "key")
	}
}

// ==================== 覆盖率相关测试 ====================

func TestAdvancedQueryParam_LikeStartFilter(t *testing.T) {
	aq := NewAdvancedQueryParam()
	aq.AddStartsWith("name", "test")

	assert.Len(t, aq.Filters, 1, "AdvancedQueryParam should have 1 filter after AddStartsWith")

	// 验证通配符前缀
	value, ok := aq.Filters[0].Value.(string)
	assert.True(t, ok, "Filter value should be string type")
	assert.Equal(t, "test%", value, "LIKE filter should have % suffix for prefix matching")
}

func TestAdvancedQueryParam_AddOrFilter(t *testing.T) {
	aq := NewAdvancedQueryParam()
	aq.AddEQ("status", "active")
	// 使用新增的便捷方法 AddOrEQ
	aq.AddOrEQ("status", "pending")

	// 验证OR逻辑（使用 assert 校验）
	assert.Equal(t, "OR", aq.Filters[0].Logic, "First filter should have OR logic after AddOrEQ")
}

func TestFilter_BuildSQL(t *testing.T) {
	aq := NewAdvancedQueryParam()

	// 测试IN操作符
	aq.AddFilter("id", OP_IN, []interface{}{1, 2, 3})
	where, args := aq.BuildWhereClause()

	assert.Contains(t, where, "IN", "WHERE clause should contain IN operator")
	assert.Len(t, args, 3, "IN filter should return 3 arguments")
}

// ==================== 便捷方法测试 - 展示简化的API ====================

func TestAdvancedQueryParam_ConvenienceMethods(t *testing.T) {
	// 演示新增便捷方法的用法
	aq := NewAdvancedQueryParam().
		AddEQ("status", "active").        // 简化的等于
		AddGT("age", 18).                 // 大于
		AddLT("price", 1000).             // 小于
		AddGTE("score", 80).              // 大于等于
		AddLTE("quantity", 100).          // 小于等于
		AddNEQ("deleted_at", nil).        // 不等于
		AddLike("name", "test").          // 全模糊匹配
		AddStartsWith("code", "ABC").     // 前缀匹配
		AddEndsWith("suffix", "000").     // 后缀匹配
		AddIn("category", "A", "B", "C"). // IN 列表
		SetPage(1, 20).
		AddOrderDesc("created_at")

	// 使用 assert 验证
	assert.Len(t, aq.Filters, 10, "Should have 10 filter conditions")
	assert.Len(t, aq.Orders, 1, "Should have 1 order condition")
	assert.Equal(t, "DESC", aq.Orders[0].Order, "Sort order should be DESC")
	assert.Equal(t, 1, aq.Page, "Page number should be 1")
	assert.Equal(t, 20, aq.PageSize, "Page size should be 20")
}

func TestAdvancedQueryParam_ConvenienceComparisonMethods(t *testing.T) {
	aq := NewAdvancedQueryParam().
		AddGT("score", 100).
		AddGTE("rating", 4.5).
		AddLT("age", 65).
		AddLTE("balance", 500000)

	assert.Len(t, aq.Filters, 4, "Should have 4 comparison filter conditions")

	// 验证第一个过滤器是GT
	assert.Equal(t, "score", aq.Filters[0].Field)
	assert.Equal(t, OP_GT, aq.Filters[0].Operator)
	assert.Equal(t, 100, aq.Filters[0].Value)

	// 验证第二个过滤器是GTE
	assert.Equal(t, "rating", aq.Filters[1].Field)
	assert.Equal(t, OP_GTE, aq.Filters[1].Operator)
	assert.Equal(t, 4.5, aq.Filters[1].Value)
}

func TestAdvancedQueryParam_ConvenientStringMatching(t *testing.T) {
	aq := NewAdvancedQueryParam().
		AddLike("name", "john").           // %john%
		AddStartsWith("username", "user"). // user%
		AddEndsWith("email", "@qq.com")    // %@qq.com

	assert.Len(t, aq.Filters, 3, "Should have 3 string matching filters")

	// 验证全模糊
	assert.Equal(t, "%john%", aq.Filters[0].Value)

	// 验证前缀
	assert.Equal(t, "user%", aq.Filters[1].Value)

	// 验证后缀
	assert.Equal(t, "%@qq.com", aq.Filters[2].Value)
}

func TestAdvancedQueryParam_OrMethods(t *testing.T) {
	aq := NewAdvancedQueryParam().
		AddEQ("status", "active").
		AddOrEQ("status", "pending").
		AddOrGT("priority", 5).
		AddOrLike("title", "urgent")

	assert.Len(t, aq.Filters, 4, "Should have 4 filter conditions")

	// AddOrFilter 会修改前一个过滤器的Logic为OR
	assert.Equal(t, "OR", aq.Filters[0].Logic, "First filter's Logic should be OR after AddOrEQ")
	assert.Equal(t, "OR", aq.Filters[1].Logic, "Second filter's Logic should be OR after AddOrGT")
	assert.Equal(t, "OR", aq.Filters[2].Logic, "Third filter's Logic should be OR after AddOrLike")
	assert.Equal(t, "AND", aq.Filters[3].Logic, "Fourth filter should keep default AND Logic")
}

func TestAdvancedQueryParam_InMethods(t *testing.T) {
	aq := NewAdvancedQueryParam().
		AddIn("status", "active", "pending", "processing").
		AddOrIn("type", 1, 2, 3)

	assert.Len(t, aq.Filters, 2, "Should have 2 IN filter conditions")
	assert.Equal(t, OP_IN, aq.Filters[0].Operator)
	assert.Equal(t, OP_IN, aq.Filters[1].Operator)

	values1, ok1 := aq.Filters[0].Value.([]interface{})
	assert.True(t, ok1, "First IN filter value should be []interface{}")
	assert.Len(t, values1, 3, "First IN filter should have 3 values")

	values2, ok2 := aq.Filters[1].Value.([]interface{})
	assert.True(t, ok2, "Second IN filter value should be []interface{}")
	assert.Len(t, values2, 3, "Second IN filter should have 3 values")
}

// 移除原有自定义的 testStringContains，直接使用 assert.Contains 替代
