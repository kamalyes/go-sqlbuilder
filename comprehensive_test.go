/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-11 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-11 00:00:00
 * @FilePath: \go-sqlbuilder\comprehensive_test.go
 * @Description: 完整的SQL生成测试 - 不依赖真实数据库
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package sqlbuilder

import (
	"context"
	"strings"
	"testing"
	"time"
)

// TestBuilderSelectBasic 测试基本SELECT生成
func TestBuilderSelectBasic(t *testing.T) {
	// 使用nil作为adapter以测试SQL生成逻辑
	builder := &Builder{
		adapter: nil,
		ctx:     context.Background(),
	}

	sql, args := builder.
		Table("users").
		Select("id", "name", "email").
		ToSQL()

	if !strings.Contains(sql, "SELECT id, name, email") {
		t.Errorf("Expected SELECT clause. Got: %s", sql)
	}

	if !strings.Contains(sql, "FROM users") {
		t.Errorf("Expected FROM clause. Got: %s", sql)
	}

	if len(args) != 0 {
		t.Errorf("Expected 0 args, got %d", len(args))
	}

	t.Logf("✓ SELECT SQL: %s", sql)
}

// TestBuilderSelectDistinct 测试DISTINCT
func TestBuilderSelectDistinct(t *testing.T) {
	builder := &Builder{
		adapter: nil,
		ctx:     context.Background(),
	}

	sql, _ := builder.
		Table("users").
		Select("age").
		Distinct().
		ToSQL()

	if !strings.Contains(sql, "DISTINCT") {
		t.Errorf("Expected DISTINCT keyword. Got: %s", sql)
	}

	t.Logf("✓ DISTINCT SQL: %s", sql)
}

// TestBuilderWhereEquals 测试WHERE equals
func TestBuilderWhereEquals(t *testing.T) {
	builder := &Builder{
		adapter: nil,
		ctx:     context.Background(),
	}

	sql, args := builder.
		Table("users").
		Select("*").
		Where("age", "=", 25).
		ToSQL()

	if !strings.Contains(sql, "WHERE") {
		t.Errorf("Expected WHERE clause. Got: %s", sql)
	}

	if !strings.Contains(sql, "age") {
		t.Errorf("Expected age column. Got: %s", sql)
	}

	if len(args) != 1 || args[0] != 25 {
		t.Errorf("Expected args [25], got %v", args)
	}

	t.Logf("✓ WHERE SQL: %s with args %v", sql, args)
}

// TestBuilderWhereIn 测试WHERE IN
func TestBuilderWhereIn(t *testing.T) {
	builder := &Builder{
		adapter: nil,
		ctx:     context.Background(),
	}

	sql, args := builder.
		Table("users").
		Select("*").
		WhereIn("id", 1, 2, 3).
		ToSQL()

	if !strings.Contains(sql, "IN") {
		t.Errorf("Expected IN clause. Got: %s", sql)
	}

	if len(args) != 3 {
		t.Errorf("Expected 3 args, got %d", len(args))
	}

	t.Logf("✓ WHERE IN SQL: %s with args %v", sql, args)
}

// TestBuilderWhereBetween 测试WHERE BETWEEN
func TestBuilderWhereBetween(t *testing.T) {
	builder := &Builder{
		adapter: nil,
		ctx:     context.Background(),
	}

	sql, args := builder.
		Table("users").
		Select("*").
		WhereBetween("age", 20, 50).
		ToSQL()

	if !strings.Contains(sql, "BETWEEN") {
		t.Errorf("Expected BETWEEN clause. Got: %s", sql)
	}

	if len(args) != 2 {
		t.Errorf("Expected 2 args, got %d", len(args))
	}

	t.Logf("✓ WHERE BETWEEN SQL: %s with args %v", sql, args)
}

// TestBuilderWhereNull 测试WHERE IS NULL
func TestBuilderWhereNull(t *testing.T) {
	builder := &Builder{
		adapter: nil,
		ctx:     context.Background(),
	}

	sql, args := builder.
		Table("users").
		Select("*").
		WhereNull("deleted_at").
		ToSQL()

	if !strings.Contains(sql, "IS NULL") {
		t.Errorf("Expected IS NULL clause. Got: %s", sql)
	}

	if len(args) != 0 {
		t.Errorf("Expected 0 args for IS NULL, got %d", len(args))
	}

	t.Logf("✓ WHERE IS NULL SQL: %s", sql)
}

// TestBuilderInsert 测试INSERT生成
func TestBuilderInsert(t *testing.T) {
	builder := &Builder{
		adapter: nil,
		ctx:     context.Background(),
	}

	data := map[string]interface{}{
		"name":  "Alice",
		"email": "alice@example.com",
		"age":   30,
	}

	sql, args := builder.
		Table("users").
		Insert(data).
		ToSQL()

	if !strings.Contains(sql, "INSERT INTO") {
		t.Errorf("Expected INSERT INTO. Got: %s", sql)
	}

	if !strings.Contains(sql, "VALUES") {
		t.Errorf("Expected VALUES. Got: %s", sql)
	}

	if len(args) != 3 {
		t.Errorf("Expected 3 args, got %d", len(args))
	}

	t.Logf("✓ INSERT SQL: %s with %d args", sql, len(args))
}

// TestBuilderUpdate 测试UPDATE生成
func TestBuilderUpdate(t *testing.T) {
	builder := &Builder{
		adapter: nil,
		ctx:     context.Background(),
	}

	data := map[string]interface{}{
		"age":     35,
		"balance": 1500.00,
	}

	builder = builder.Table("users")
	b2, _ := builder.Update(data)
	sql, args := b2.
		Where("id", "=", 1).
		ToSQL()

	if !strings.Contains(sql, "UPDATE") {
		t.Errorf("Expected UPDATE. Got: %s", sql)
	}

	if !strings.Contains(sql, "SET") {
		t.Errorf("Expected SET. Got: %s", sql)
	}

	if !strings.Contains(sql, "WHERE") {
		t.Errorf("Expected WHERE. Got: %s", sql)
	}

	if len(args) != 3 {
		t.Errorf("Expected 3 args, got %d", len(args))
	}

	t.Logf("✓ UPDATE SQL: %s with %d args", sql, len(args))
}

// TestBuilderDelete 测试DELETE生成
func TestBuilderDelete(t *testing.T) {
	builder := &Builder{
		adapter: nil,
		ctx:     context.Background(),
	}

	sql, args := builder.
		Table("users").
		Delete().
		Where("id", "=", 1).
		ToSQL()

	if !strings.Contains(sql, "DELETE FROM") {
		t.Errorf("Expected DELETE FROM. Got: %s", sql)
	}

	if !strings.Contains(sql, "WHERE") {
		t.Errorf("Expected WHERE. Got: %s", sql)
	}

	if len(args) != 1 {
		t.Errorf("Expected 1 arg, got %d", len(args))
	}

	t.Logf("✓ DELETE SQL: %s with args %v", sql, args)
}

// TestBuilderJoin 测试JOIN生成
func TestBuilderJoin(t *testing.T) {
	builder := &Builder{
		adapter: nil,
		ctx:     context.Background(),
	}

	sql, _ := builder.
		Table("users").
		Select("*").
		LeftJoin("orders", "users.id = orders.user_id").
		ToSQL()

	if !strings.Contains(sql, "LEFT JOIN") {
		t.Errorf("Expected LEFT JOIN clause. Got: %s", sql)
	}

	if !strings.Contains(sql, "ON") {
		t.Errorf("Expected ON clause. Got: %s", sql)
	}

	t.Logf("✓ JOIN SQL: %s", sql)
}

// TestBuilderGroupByHaving 测试GROUP BY和HAVING
func TestBuilderGroupByHaving(t *testing.T) {
	builder := &Builder{
		adapter: nil,
		ctx:     context.Background(),
	}

	sql, args := builder.
		Table("users").
		Select("age", "COUNT(*) as count").
		GroupBy("age").
		Having("count(*)", ">", 5).
		ToSQL()

	if !strings.Contains(sql, "GROUP BY") {
		t.Errorf("Expected GROUP BY. Got: %s", sql)
	}

	if !strings.Contains(sql, "HAVING") {
		t.Errorf("Expected HAVING. Got: %s", sql)
	}

	if len(args) != 1 {
		t.Errorf("Expected 1 arg for HAVING, got %d", len(args))
	}

	t.Logf("✓ GROUP BY/HAVING SQL: %s", sql)
}

// TestBuilderOrderBy 测试ORDER BY
func TestBuilderOrderBy(t *testing.T) {
	builder := &Builder{
		adapter: nil,
		ctx:     context.Background(),
	}

	sql, _ := builder.
		Table("users").
		Select("*").
		OrderBy("name").
		OrderByDesc("age").
		ToSQL()

	if !strings.Contains(sql, "ORDER BY") {
		t.Errorf("Expected ORDER BY. Got: %s", sql)
	}

	if !strings.Contains(sql, "ASC") {
		t.Errorf("Expected ASC. Got: %s", sql)
	}

	if !strings.Contains(sql, "DESC") {
		t.Errorf("Expected DESC. Got: %s", sql)
	}

	t.Logf("✓ ORDER BY SQL: %s", sql)
}

// TestBuilderLimitOffset 测试LIMIT和OFFSET
func TestBuilderLimitOffset(t *testing.T) {
	builder := &Builder{
		adapter: nil,
		ctx:     context.Background(),
	}

	sql, _ := builder.
		Table("users").
		Select("*").
		Limit(10).
		Offset(20).
		ToSQL()

	if !strings.Contains(sql, "LIMIT 10") {
		t.Errorf("Expected LIMIT 10. Got: %s", sql)
	}

	if !strings.Contains(sql, "OFFSET 20") {
		t.Errorf("Expected OFFSET 20. Got: %s", sql)
	}

	t.Logf("✓ LIMIT/OFFSET SQL: %s", sql)
}

// TestBuilderPaginate 测试分页
func TestBuilderPaginate(t *testing.T) {
	builder := &Builder{
		adapter: nil,
		ctx:     context.Background(),
	}

	// 第2页，每页25条
	sql, _ := builder.
		Table("users").
		Select("*").
		Paginate(2, 25).
		ToSQL()

	if !strings.Contains(sql, "LIMIT 25") {
		t.Errorf("Expected LIMIT 25. Got: %s", sql)
	}

	if !strings.Contains(sql, "OFFSET 25") {
		t.Errorf("Expected OFFSET 25 (page2 * 25). Got: %s", sql)
	}

	t.Logf("✓ Paginate SQL: %s", sql)
}

// TestBuilderComplexQuery 测试复杂查询
func TestBuilderComplexQuery(t *testing.T) {
	builder := &Builder{
		adapter: nil,
		ctx:     context.Background(),
	}

	sql, args := builder.
		Table("users").
		As("u").
		Select("u.id", "u.name", "u.age").
		LeftJoin("orders o", "u.id = o.user_id").
		Where("u.age", ">", 20).
		OrWhere("u.balance", ">", 1000).
		WhereNull("u.deleted_at").
		GroupBy("u.id", "u.name").
		Having("COUNT(*)", ">", 0).
		OrderByDesc("u.created_at").
		Limit(50).
		Offset(10).
		ToSQL()

	expectedClauses := []string{"SELECT", "FROM", "LEFT JOIN", "WHERE", "GROUP BY", "HAVING", "ORDER BY", "LIMIT", "OFFSET"}
	for _, clause := range expectedClauses {
		if !strings.Contains(sql, clause) {
			t.Errorf("Expected '%s' in complex query. Got: %s", clause, sql)
		}
	}

	if len(args) != 3 {
		t.Errorf("Expected 3 args, got %d", len(args))
	}

	t.Logf("✓ Complex SQL: %s", sql)
}

// TestBuilderMethodChaining 测试方法链式调用
func TestBuilderMethodChaining(t *testing.T) {
	builder := &Builder{
		adapter: nil,
		ctx:     context.Background(),
	}

	// 验证所有方法返回相同的实例
	result := builder.
		Table("users").
		Select("*").
		Where("age", ">", 20).
		OrderBy("name").
		Limit(10)

	if result != builder {
		t.Errorf("Method chaining did not return the same builder instance")
	}

	sql, _ := builder.ToSQL()
	if sql == "" {
		t.Errorf("Expected non-empty SQL")
	}

	t.Logf("✓ Method chaining works: %s", sql)
}

// TestBuilderContext 测试上下文支持
func TestBuilderContext(t *testing.T) {
	builder := &Builder{
		adapter: nil,
		ctx:     context.Background(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	builder.WithContext(ctx).WithTimeout(10 * time.Second)

	if builder.ctx != ctx {
		t.Errorf("Expected context to be set")
	}

	if builder.timeout != 10*time.Second {
		t.Errorf("Expected timeout to be 10 seconds, got %v", builder.timeout)
	}

	t.Logf("✓ Context support works")
}

// TestBuilderWhereRaw 测试原始WHERE
func TestBuilderWhereRaw(t *testing.T) {
	builder := &Builder{
		adapter: nil,
		ctx:     context.Background(),
	}

	sql, args := builder.
		Table("users").
		Select("*").
		WhereRaw("(age > ? AND balance > ?) OR name = ?", 25, 500.0, "VIP").
		ToSQL()

	if !strings.Contains(sql, "AND") {
		t.Errorf("Expected AND. Got: %s", sql)
	}

	if !strings.Contains(sql, "OR") {
		t.Errorf("Expected OR. Got: %s", sql)
	}

	if len(args) != 3 {
		t.Errorf("Expected 3 args, got %d", len(args))
	}

	t.Logf("✓ WhereRaw SQL: %s with args %v", sql, args)
}

// TestBuilderSet 测试Set方法
func TestBuilderSet(t *testing.T) {
	builder := &Builder{
		adapter: nil,
		ctx:     context.Background(),
	}

	sql, args := builder.
		Table("users").
		Set("age", 30).
		Set("balance", 1500.0).
		Where("id", "=", 1).
		ToSQL()

	if !strings.Contains(sql, "UPDATE") {
		t.Errorf("Expected UPDATE. Got: %s", sql)
	}

	if len(args) != 3 {
		t.Errorf("Expected 3 args, got %d", len(args))
	}

	t.Logf("✓ Set SQL: %s", sql)
}

// TestBuilderTableAlias 测试表别名
func TestBuilderTableAlias(t *testing.T) {
	builder := &Builder{
		adapter: nil,
		ctx:     context.Background(),
	}

	sql, _ := builder.
		Table("users").
		As("u").
		Select("u.id", "u.name").
		ToSQL()

	if !strings.Contains(sql, "AS u") {
		t.Errorf("Expected 'AS u'. Got: %s", sql)
	}

	t.Logf("✓ Table alias SQL: %s", sql)
}

// TestBuilderSelectRaw 测试SelectRaw
func TestBuilderSelectRaw(t *testing.T) {
	builder := &Builder{
		adapter: nil,
		ctx:     context.Background(),
	}

	sql, _ := builder.
		Table("users").
		SelectRaw("COUNT(*) as total, SUM(balance) as sum_balance").
		ToSQL()

	if !strings.Contains(sql, "COUNT(*)") {
		t.Errorf("Expected COUNT(*). Got: %s", sql)
	}

	t.Logf("✓ SelectRaw SQL: %s", sql)
}

// BenchmarkBuilderSQL 基准测试 - SQL生成性能
func BenchmarkBuilderSQL(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		builder := &Builder{
			adapter: nil,
			ctx:     context.Background(),
		}

		builder.
			Table("users").
			Select("id", "name", "email", "age").
			Where("age", ">", 20).
			OrWhere("status", "=", "active").
			OrderByDesc("created_at").
			Limit(100).
			Paginate(1, 10).
			ToSQL()
	}
}

// TestBuilderMultipleWheres 测试多个WHERE条件
func TestBuilderMultipleWheres(t *testing.T) {
	builder := &Builder{
		adapter: nil,
		ctx:     context.Background(),
	}

	sql, args := builder.
		Table("users").
		Select("*").
		Where("age", ">", 20).
		Where("status", "=", "active").
		Where("balance", ">", 100).
		ToSQL()

	if len(args) != 3 {
		t.Errorf("Expected 3 args, got %d", len(args))
	}

	countWhere := strings.Count(sql, "WHERE")
	countAnd := strings.Count(sql, "AND")

	if countWhere != 1 {
		t.Errorf("Expected 1 WHERE, got %d", countWhere)
	}

	if countAnd < 2 {
		t.Errorf("Expected at least 2 AND, got %d", countAnd)
	}

	t.Logf("✓ Multiple WHERE SQL: %s", sql)
}
