/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-11 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-11 14:38:09
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

	"github.com/stretchr/testify/assert"
)

// TestBuilderSelectBasic 测试基本SELECT生成
func TestBuilderSelectBasic(t *testing.T) {
	builder := &Builder{
		adapter: nil,
		ctx:     context.Background(),
	}

	sql, args := builder.
		Table("users").
		Select("id", "name", "email").
		ToSQL()

	assert.Contains(t, sql, "SELECT id, name, email", "Expected SELECT clause.")
	assert.Contains(t, sql, "FROM users", "Expected FROM clause.")
	assert.Empty(t, args, "Expected 0 args.")

	t.Logf("✓ SELECT SQL: %s", sql)
}

// TestBuilderSelectDistinct 测试DISTINCTs
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

	assert.Contains(t, sql, "DISTINCT", "Expected DISTINCT keyword.")

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

	assert.Contains(t, sql, "WHERE", "Expected WHERE clause.")
	assert.Contains(t, sql, "age", "Expected age column.")
	assert.Len(t, args, 1)
	assert.Equal(t, 25, args[0], "Expected args to be [25].")

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

	assert.Contains(t, sql, "IN", "Expected IN clause.")
	assert.Len(t, args, 3, "Expected 3 args.")

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

	assert.Contains(t, sql, "BETWEEN", "Expected BETWEEN clause.")
	assert.Len(t, args, 2, "Expected 2 args.")

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

	assert.Contains(t, sql, "IS NULL", "Expected IS NULL clause.")
	assert.Empty(t, args, "Expected 0 args for IS NULL.")

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

	assert.Contains(t, sql, "INSERT INTO", "Expected INSERT INTO.")
	assert.Contains(t, sql, "VALUES", "Expected VALUES.")
	assert.Len(t, args, 3, "Expected 3 args.")

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

	assert.Contains(t, sql, "UPDATE", "Expected UPDATE.")
	assert.Contains(t, sql, "SET", "Expected SET.")
	assert.Contains(t, sql, "WHERE", "Expected WHERE.")
	assert.Len(t, args, 3, "Expected 3 args.")

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

	assert.Contains(t, sql, "DELETE FROM", "Expected DELETE FROM.")
	assert.Contains(t, sql, "WHERE", "Expected WHERE.")
	assert.Len(t, args, 1, "Expected 1 arg.")

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

	assert.Contains(t, sql, "LEFT JOIN", "Expected LEFT JOIN clause.")
	assert.Contains(t, sql, "ON", "Expected ON clause.")

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

	assert.Contains(t, sql, "GROUP BY", "Expected GROUP BY.")
	assert.Contains(t, sql, "HAVING", "Expected HAVING.")
	assert.Len(t, args, 1, "Expected 1 arg for HAVING.")

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

	assert.Contains(t, sql, "ORDER BY", "Expected ORDER BY.")
	assert.Contains(t, sql, "ASC", "Expected ASC.")
	assert.Contains(t, sql, "DESC", "Expected DESC.")

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

	assert.Contains(t, sql, "LIMIT 10", "Expected LIMIT 10.")
	assert.Contains(t, sql, "OFFSET 20", "Expected OFFSET 20.")

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

	assert.Contains(t, sql, "LIMIT 25", "Expected LIMIT 25.")
	assert.Contains(t, sql, "OFFSET 25", "Expected OFFSET 25 (page2 * 25).")

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
		assert.Contains(t, sql, clause, "Expected '%s' in complex query.", clause)
	}

	assert.Len(t, args, 3, "Expected 3 args.")

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

	assert.Same(t, builder, result, "Method chaining did not return the same builder instance")

	sql, _ := builder.ToSQL()
	assert.NotEmpty(t, sql, "Expected non-empty SQL.")

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

	assert.Equal(t, ctx, builder.ctx, "Expected context to be set.")
	assert.Equal(t, 10*time.Second, builder.timeout, "Expected timeout to be 10 seconds.")

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

	assert.Contains(t, sql, "AND", "Expected AND.")
	assert.Contains(t, sql, "OR", "Expected OR.")
	assert.Len(t, args, 3, "Expected 3 args.")

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

	assert.Contains(t, sql, "UPDATE", "Expected UPDATE.")
	assert.Len(t, args, 3, "Expected 3 args.")

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

	assert.Contains(t, sql, "AS u", "Expected 'AS u'.")

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

	assert.Contains(t, sql, "COUNT(*)", "Expected COUNT(*).")

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

	assert.Len(t, args, 3, "Expected 3 args.")
	countWhere := strings.Count(sql, "WHERE")
	countAnd := strings.Count(sql, "AND")

	assert.Equal(t, 1, countWhere, "Expected 1 WHERE.")
	assert.GreaterOrEqual(t, countAnd, 2, "Expected at least 2 AND.")

	t.Logf("✓ Multiple WHERE SQL: %s", sql)
}
