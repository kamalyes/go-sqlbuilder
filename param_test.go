/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-11 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-11 15:51:31
 * @FilePath: \go-sqlbuilder\param_test.go
 * @Description: 缓存和高级查询测试
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
	if err != nil {
		t.Errorf("Set failed: %v", err)
	}

	value, err := cache.Get(ctx, "test_key")
	if err != nil || value != "test_value" {
		t.Errorf("Get failed: got %s, want test_value", value)
	}

	// 测试不存在的键
	_, err = cache.Get(ctx, "non_exist_key")
	if err == nil {
		t.Error("Expected error for non-existent key")
	}

	// 测试Exists
	exists, err := cache.Exists(ctx, "test_key")
	if err != nil || !exists {
		t.Error("Exists should return true for existing key")
	}

	// 测试Delete
	err = cache.Delete(ctx, "test_key")
	if err != nil {
		t.Errorf("Delete failed: %v", err)
	}

	exists, err = cache.Exists(ctx, "test_key")
	if err != nil || exists {
		t.Error("Exists should return false after delete")
	}
}

func TestCacheExpiration(t *testing.T) {
	ctx := context.Background()
	cache := NewMockCacheStore()

	// 设置短TTL
	err := cache.Set(ctx, "exp_key", "exp_value", 100*time.Millisecond)
	if err != nil {
		t.Errorf("Set failed: %v", err)
	}

	// 立即获取应该成功
	value, err := cache.Get(ctx, "exp_key")
	if err != nil || value != "exp_value" {
		t.Error("Should get value before expiration")
	}

	// 等待过期
	time.Sleep(150 * time.Millisecond)

	// 获取应该失败
	_, err = cache.Get(ctx, "exp_key")
	if err == nil {
		t.Error("Should fail to get expired value")
	}
}

func TestCacheClear(t *testing.T) {
	ctx := context.Background()
	cache := NewMockCacheStore()

	// 设置多个键
	cache.Set(ctx, "sqlbuilder:key1", "value1", 1*time.Hour)
	cache.Set(ctx, "sqlbuilder:key2", "value2", 1*time.Hour)
	cache.Set(ctx, "other:key3", "value3", 1*time.Hour)

	// 清除sqlbuilder前缀的缓存
	err := cache.Clear(ctx, "sqlbuilder:")
	if err != nil {
		t.Errorf("Clear failed: %v", err)
	}

	// 验证sqlbuilder:前缀的缓存已删除
	_, err = cache.Get(ctx, "sqlbuilder:key1")
	if err == nil {
		t.Error("sqlbuilder:key1 should be deleted")
	}

	// other:前缀的缓存应该还在
	value, err := cache.Get(ctx, "other:key3")
	if err != nil || value != "value3" {
		t.Error("other:key3 should still exist")
	}
}

func TestCacheStats(t *testing.T) {
	cache := NewMockCacheStore()
	ctx := context.Background()

	// 初始状态
	stats := cache.GetStats()
	if stats["total"] != 0 {
		t.Error("Initial cache should be empty")
	}

	// 添加项目
	cache.Set(ctx, "key1", "value1", 1*time.Hour)
	cache.Set(ctx, "key2", "value2", 1*time.Hour)

	stats = cache.GetStats()
	if stats["total"] != 2 {
		t.Errorf("Expected 2 items, got %d", stats["total"])
	}

	// 添加过期项目
	cache.Set(ctx, "key3", "value3", 10*time.Millisecond)
	time.Sleep(50 * time.Millisecond)

	stats = cache.GetStats()
	if stats["valid"] != 2 {
		t.Errorf("Expected 2 valid items, got %d", stats["valid"])
	}
}

// ==================== 高级查询测试 ====================

func TestAdvancedQueryParam_Filters(t *testing.T) {
	aq := NewAdvancedQueryParam()

	// 测试添加过滤
	aq.AddEQ("status", "active").
		AddLike("name", "test").
		AddIn("id", 1, 2, 3)

	if len(aq.Filters) != 3 {
		t.Errorf("Expected 3 filters, got %d", len(aq.Filters))
	}

	// 验证过滤内容
	if aq.Filters[0].Field != "status" || aq.Filters[0].Operator != OP_EQ {
		t.Error("First filter is incorrect")
	}

	if aq.Filters[1].Field != "name" || aq.Filters[1].Operator != OP_LIKE {
		t.Error("Second filter is incorrect")
	}

	if aq.Filters[2].Field != "id" || aq.Filters[2].Operator != OP_IN {
		t.Error("Third filter is incorrect")
	}
}

func TestAdvancedQueryParam_TimeRange(t *testing.T) {
	aq := NewAdvancedQueryParam()

	startTime := time.Now().Add(-24 * time.Hour)
	endTime := time.Now()

	aq.AddTimeRange("created_at", startTime, endTime)

	if len(aq.TimeRanges) != 1 {
		t.Error("Expected 1 time range")
	}

	if _, exists := aq.TimeRanges["created_at"]; !exists {
		t.Error("Time range for created_at not found")
	}
}

func TestAdvancedQueryParam_Ordering(t *testing.T) {
	aq := NewAdvancedQueryParam()

	aq.AddOrder("created_at", "DESC").
		AddOrderAsc("name")

	if len(aq.Orders) != 2 {
		t.Errorf("Expected 2 orders, got %d", len(aq.Orders))
	}

	if aq.Orders[0].Field != "created_at" || aq.Orders[0].Order != "DESC" {
		t.Error("First order is incorrect")
	}

	if aq.Orders[1].Field != "name" || aq.Orders[1].Order != "ASC" {
		t.Error("Second order is incorrect")
	}
}

func TestAdvancedQueryParam_Pagination(t *testing.T) {
	aq := NewAdvancedQueryParam()

	aq.SetPage(2, 20)

	if aq.Page != 2 || aq.PageSize != 20 {
		t.Errorf("Pagination not set correctly: page=%d, pageSize=%d", aq.Page, aq.PageSize)
	}

	if aq.Offset != 20 {
		t.Errorf("Offset should be 20 for page 2 with pageSize 20, got %d", aq.Offset)
	}
}

func TestAdvancedQueryParam_BuildWhereClause(t *testing.T) {
	aq := NewAdvancedQueryParam()

	aq.AddEQ("status", "active").
		AddLike("name", "test").
		AddIn("id", 1, 2, 3)

	where, args := aq.BuildWhereClause()

	if where == "" {
		t.Error("WHERE clause should not be empty")
	}

	if len(args) != 5 { // "active", "%test%", 1, 2, 3
		t.Errorf("Expected 5 arguments, got %d", len(args))
	}
}

func TestAdvancedQueryParam_FindInSet(t *testing.T) {
	aq := NewAdvancedQueryParam()

	aq.AddFindInSet("tags", "golang", "database", "cache")

	if len(aq.FindInSets) != 1 {
		t.Error("Expected 1 FIND_IN_SET condition")
	}

	if values, exists := aq.FindInSets["tags"]; !exists || len(values) != 3 {
		t.Error("FIND_IN_SET values not stored correctly")
	}
}

func TestAdvancedQueryParam_Distinct(t *testing.T) {
	aq := NewAdvancedQueryParam()

	aq.SetDistinct(true)

	if !aq.Distinct {
		t.Error("Distinct should be true")
	}
}

func TestAdvancedQueryParam_SelectFields(t *testing.T) {
	aq := NewAdvancedQueryParam()

	aq.SetSelectFields("id", "name", "email")

	if len(aq.SelectFields) != 3 {
		t.Errorf("Expected 3 select fields, got %d", len(aq.SelectFields))
	}

	if aq.SelectFields[0] != "id" || aq.SelectFields[1] != "name" || aq.SelectFields[2] != "email" {
		t.Error("Select fields not set correctly")
	}
}

func TestAdvancedQueryParam_Having(t *testing.T) {
	aq := NewAdvancedQueryParam()

	aq.AddHaving("COUNT(*) > 5").
		AddHaving("SUM(amount) > 100")

	if len(aq.HavingClauses) != 2 {
		t.Errorf("Expected 2 HAVING clauses, got %d", len(aq.HavingClauses))
	}
}

// ==================== 分页响应测试 ====================

func TestPageBean_Initialization(t *testing.T) {
	rows := []string{"row1", "row2", "row3"}
	page := NewPageBean(2, 20, 100, rows)

	if page.Page != 2 || page.PageSize != 20 || page.Total != 100 {
		t.Error("PageBean initialization failed")
	}

	if page.Pages != 5 {
		t.Errorf("Expected 5 pages, got %d", page.Pages)
	}
}

func TestPageBean_EdgeCases(t *testing.T) {
	// 测试无数据情况
	page := NewPageBean(0, 10, 0, []string{})
	if page.Page != 1 || page.Pages != 1 {
		t.Error("PageBean should default to page 1 for invalid input")
	}

	// 测试未满一页的情况
	page = NewPageBean(1, 10, 5, []string{"a", "b", "c", "d", "e"})
	if page.Pages != 1 {
		t.Errorf("Expected 1 page for 5 items with pageSize 10, got %d", page.Pages)
	}

	// 测试需要多页的情况
	page = NewPageBean(1, 10, 25, []string{})
	if page.Pages != 3 {
		t.Errorf("Expected 3 pages for 25 items with pageSize 10, got %d", page.Pages)
	}
}

// ==================== 查询选项测试 ====================

func TestFindOption_Builder(t *testing.T) {
	opt := NewFindOption().
		WithBusinessId("biz123").
		WithShopId("shop456").
		WithTablePrefix("user_").
		WithCacheTTL(30 * time.Minute).
		WithNoCache()

	if opt.BusinessId != "biz123" {
		t.Errorf("BusinessId should be biz123, got %s", opt.BusinessId)
	}

	if opt.ShopId != "shop456" {
		t.Errorf("ShopId should be shop456, got %s", opt.ShopId)
	}

	if opt.TablePrefix != "user_" {
		t.Errorf("TablePrefix should be user_, got %s", opt.TablePrefix)
	}

	if !opt.NoCache {
		t.Error("NoCache should be true")
	}

	if opt.CacheTTL != 30*time.Minute {
		t.Errorf("CacheTTL should be 30m, got %v", opt.CacheTTL)
	}
}

// ==================== BaseInfoFilter测试 ====================

func TestBaseInfoFilter(t *testing.T) {
	filter := BaseInfoFilter{
		DBField:    "status",
		Values:     []interface{}{"active", "pending"},
		ExactMatch: true,
		AllRegex:   false,
	}

	if filter.DBField != "status" {
		t.Error("DBField not set correctly")
	}

	if len(filter.Values) != 2 {
		t.Errorf("Expected 2 values, got %d", len(filter.Values))
	}

	if !filter.ExactMatch {
		t.Error("ExactMatch should be true")
	}
}

// ==================== FilterOperator测试 ====================

func TestFilterOperators(t *testing.T) {
	operators := []FilterOperator{
		OP_EQ, OP_NEQ, OP_GT, OP_GTE, OP_LT, OP_LTE,
		OP_LIKE, OP_IN, OP_BETWEEN, OP_IS_NULL, OP_FIND_IN_SET,
	}

	if len(operators) != 11 {
		t.Errorf("Expected 11 operators, got %d", len(operators))
	}

	// 验证操作符值
	if OP_EQ != "=" {
		t.Error("OP_EQ should be '='")
	}

	if OP_LIKE != "LIKE" {
		t.Error("OP_LIKE should be 'LIKE'")
	}

	if OP_IN != "IN" {
		t.Error("OP_IN should be 'IN'")
	}
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

	if len(aq.Filters) != 2 {
		t.Errorf("Expected 2 filters, got %d", len(aq.Filters))
	}

	if aq.Page != 1 || aq.PageSize != 20 {
		t.Error("Pagination not set correctly")
	}

	if len(aq.Orders) != 1 {
		t.Errorf("Expected 1 order, got %d", len(aq.Orders))
	}

	if !aq.Distinct {
		t.Error("Distinct should be true")
	}

	if len(aq.SelectFields) != 3 {
		t.Errorf("Expected 3 select fields, got %d", len(aq.SelectFields))
	}
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
	if err != nil {
		t.Errorf("Failed to set cache: %v", err)
	}

	// 检索缓存
	value, err := cache.Get(ctx, cacheKey)
	if err != nil {
		t.Errorf("Failed to get cache: %v", err)
	}

	if value != cacheValue {
		t.Error("Cache value mismatch")
	}
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

	if len(aq.Filters) != 1 {
		t.Error("Filter not added")
	}

	// 验证通配符前缀
	if !testStringContains(aq.Filters[0].Value.(string), "test%") {
		t.Error("LIKE filter should have % suffix for prefix matching")
	}
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

	if !testStringContains(where, "IN") {
		t.Error("WHERE clause should contain IN operator")
	}

	if len(args) != 3 {
		t.Errorf("Expected 3 args for IN filter, got %d", len(args))
	}
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
	assert.Equal(t, 10, len(aq.Filters), "应该有10个过滤条件")
	assert.Equal(t, 1, len(aq.Orders), "应该有1个排序条件")
	assert.Equal(t, "DESC", aq.Orders[0].Order, "排序顺序应为DESC")
	assert.Equal(t, 1, aq.Page, "页码应为1")
	assert.Equal(t, 20, aq.PageSize, "每页大小应为20")
}

func TestAdvancedQueryParam_ConvenienceComparisonMethods(t *testing.T) {
	aq := NewAdvancedQueryParam().
		AddGT("score", 100).
		AddGTE("rating", 4.5).
		AddLT("age", 65).
		AddLTE("balance", 500000)

	assert.Equal(t, 4, len(aq.Filters), "应该有4个比较过滤条件")

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

	assert.Equal(t, 3, len(aq.Filters), "应该有3个字符串匹配过滤")

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

	assert.Equal(t, 4, len(aq.Filters), "应该有4个过滤条件")

	// AddOrFilter 会修改前一个过滤器的Logic为OR
	// 因此：Filters[0].Logic = OR（被AddOrEQ改为OR）
	//      Filters[1].Logic = OR（被AddOrGT改为OR）
	//      Filters[2].Logic = OR（被AddOrLike改为OR）
	//      Filters[3].Logic = AND（新增，保持默认AND）
	assert.Equal(t, "OR", aq.Filters[0].Logic, "第一个过滤的Logic被AddOrEQ改为OR")
	assert.Equal(t, "OR", aq.Filters[1].Logic, "第二个过滤的Logic被AddOrGT改为OR")
	assert.Equal(t, "OR", aq.Filters[2].Logic, "第三个过滤的Logic被AddOrLike改为OR")
	assert.Equal(t, "AND", aq.Filters[3].Logic, "第四个过滤保持默认AND")
}

func TestAdvancedQueryParam_InMethods(t *testing.T) {
	aq := NewAdvancedQueryParam().
		AddIn("status", "active", "pending", "processing").
		AddOrIn("type", 1, 2, 3)

	assert.Equal(t, 2, len(aq.Filters))
	assert.Equal(t, OP_IN, aq.Filters[0].Operator)
	assert.Equal(t, OP_IN, aq.Filters[1].Operator)
	assert.Equal(t, 3, len(aq.Filters[0].Value.([]interface{})))
	assert.Equal(t, 3, len(aq.Filters[1].Value.([]interface{})))
}

func testStringContains(str string, substr string) bool {
	for i := 0; i <= len(str)-len(substr); i++ {
		if str[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
