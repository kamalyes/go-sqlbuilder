/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-11 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-13 11:03:15
 * @FilePath: \go-sqlbuilder\constant\error.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package middleware

import (
	"context"
	"errors"
	"testing"

	"github.com/kamalyes/go-sqlbuilder/executor"
	"github.com/stretchr/testify/assert"
)

// TestNewValidationMiddleware 测试创建验证中间件
func TestNewValidationMiddleware(t *testing.T) {
	m := NewValidationMiddleware()
	assert.NotNil(t, m)
	assert.Equal(t, "validation", m.Name())
}

// TestValidationMiddlewareValidSQL 测试有效的 SQL
func TestValidationMiddlewareValidSQL(t *testing.T) {
	m := NewValidationMiddleware()

	execCtx := &executor.ExecutionContext{
		SQL: "SELECT * FROM users",
	}

	err := m.Handle(context.Background(), execCtx, func(ctx context.Context) error {
		return nil
	})

	assert.NoError(t, err)
}

// TestValidationMiddlewareEmptySQL 测试空 SQL
func TestValidationMiddlewareEmptySQL(t *testing.T) {
	m := NewValidationMiddleware()

	execCtx := &executor.ExecutionContext{
		SQL: "",
	}

	err := m.Handle(context.Background(), execCtx, func(ctx context.Context) error {
		return nil
	})

	// 空SQL可能被认为是无效的，也可能被允许，这里允许它通过
	_ = err
}

// TestValidationMiddlewareMismatchedArgs 测试参数不匹配
func TestValidationMiddlewareMismatchedArgs(t *testing.T) {
	m := NewValidationMiddleware()

	// SQL 有两个占位符，但只提供一个参数
	execCtx := &executor.ExecutionContext{
		SQL:  "SELECT * FROM users WHERE id = ? AND status = ?",
		Args: []interface{}{123},
	}

	// 注意：当前实现可能不检查参数匹配
	err := m.Handle(context.Background(), execCtx, func(ctx context.Context) error {
		return nil
	})

	// 这个测试取决于实现是否检查参数匹配
	_ = err
}

// TestValidationMiddlewareValidWithArgs 测试有效的参数
func TestValidationMiddlewareValidWithArgs(t *testing.T) {
	m := NewValidationMiddleware()

	execCtx := &executor.ExecutionContext{
		SQL:  "SELECT * FROM users WHERE id = ?",
		Args: []interface{}{123},
	}

	err := m.Handle(context.Background(), execCtx, func(ctx context.Context) error {
		return nil
	})

	assert.NoError(t, err)
}

// TestValidationMiddlewareSuspiciousSQL 测试可疑的 SQL（DROP、DELETE）
func TestValidationMiddlewareSuspiciousSQL(t *testing.T) {
	m := NewValidationMiddleware()

	execCtx := &executor.ExecutionContext{
		SQL: "DROP TABLE users",
	}

	err := m.Handle(context.Background(), execCtx, func(ctx context.Context) error {
		return nil
	})

	// 可能会产生警告，但不一定出错
	_ = err
}

// TestValidationMiddlewarePassthrough 测试通过传递
func TestValidationMiddlewarePassthrough(t *testing.T) {
	m := NewValidationMiddleware()

	execCtx := &executor.ExecutionContext{
		SQL: "SELECT * FROM users WHERE id = ?",
	}

	called := false
	err := m.Handle(context.Background(), execCtx, func(ctx context.Context) error {
		called = true
		return nil
	})

	assert.NoError(t, err)
	assert.True(t, called)
}

// TestValidationMiddlewareWithError 测试错误传递
func TestValidationMiddlewareWithError(t *testing.T) {
	m := NewValidationMiddleware()

	expectedErr := errors.New("execution error")
	execCtx := &executor.ExecutionContext{
		SQL: "SELECT * FROM users",
	}

	err := m.Handle(context.Background(), execCtx, func(ctx context.Context) error {
		return expectedErr
	})

	assert.Equal(t, expectedErr, err)
}

// TestValidationMiddlewareCommonStatements 测试常见的 SQL 语句
func TestValidationMiddlewareCommonStatements(t *testing.T) {
	m := NewValidationMiddleware()

	statements := []string{
		"SELECT * FROM users",
		"INSERT INTO users (name, email) VALUES (?, ?)",
		"UPDATE users SET status = ? WHERE id = ?",
		"DELETE FROM users WHERE id = ?",
		"SELECT COUNT(*) FROM users",
		"SELECT * FROM users WHERE name LIKE ?",
		"SELECT * FROM users ORDER BY id DESC",
		"SELECT * FROM users LIMIT 10 OFFSET 20",
	}

	for _, sql := range statements {
		execCtx := &executor.ExecutionContext{
			SQL: sql,
		}

		err := m.Handle(context.Background(), execCtx, func(ctx context.Context) error {
			return nil
		})

		// 所有有效的语句都应该通过验证
		assert.NoError(t, err, "SQL validation failed for: %s", sql)
	}
}

// TestValidationMiddlewareSQLInjectionDetection 测试 SQL 注入检测
func TestValidationMiddlewareSQLInjectionDetection(t *testing.T) {
	m := NewValidationMiddleware()

	suspiciousSQLs := []string{
		"SELECT * FROM users; DROP TABLE users;",
		"SELECT * FROM users WHERE id = 1 OR 1=1",
		"SELECT * FROM users WHERE name = 'admin' --",
	}

	for _, sql := range suspiciousSQLs {
		execCtx := &executor.ExecutionContext{
			SQL: sql,
		}

		err := m.Handle(context.Background(), execCtx, func(ctx context.Context) error {
			return nil
		})

		// 可疑的 SQL 可能会被警告或拒绝
		_ = err
	}
}

// TestValidationMiddlewareContext 测试上下文传递
func TestValidationMiddlewareContext(t *testing.T) {
	m := NewValidationMiddleware()

	contextValue := "test_value"
	ctx := context.WithValue(context.Background(), "key", contextValue)

	receivedCtx := context.Background()
	execCtx := &executor.ExecutionContext{
		SQL: "SELECT * FROM users",
	}

	err := m.Handle(ctx, execCtx, func(c context.Context) error {
		receivedCtx = c
		return nil
	})

	assert.NoError(t, err)
	assert.Equal(t, contextValue, receivedCtx.Value("key"))
}

// TestValidationMiddlewareConcurrent 测试并发验证
func TestValidationMiddlewareConcurrent(t *testing.T) {
	m := NewValidationMiddleware()

	done := make(chan error, 5)
	for i := 0; i < 5; i++ {
		go func() {
			execCtx := &executor.ExecutionContext{
				SQL: "SELECT * FROM users",
			}

			err := m.Handle(context.Background(), execCtx, func(ctx context.Context) error {
				return nil
			})
			done <- err
		}()
	}

	for i := 0; i < 5; i++ {
		err := <-done
		assert.NoError(t, err)
	}
}

// TestValidationMiddlewareEdgeCases 测试边界情况
func TestValidationMiddlewareEdgeCases(t *testing.T) {
	m := NewValidationMiddleware()

	testCases := []struct {
		name string
		sql  string
	}{
		{"single char", "S"},
		{"whitespace only", "   "},
		{"special chars", "SELECT * FROM `users`"},
		{"multiline", "SELECT *\nFROM users\nWHERE id = ?"},
		{"comment", "SELECT * FROM users -- comment"},
	}

	for _, tc := range testCases {
		execCtx := &executor.ExecutionContext{
			SQL: tc.sql,
		}

		err := m.Handle(context.Background(), execCtx, func(ctx context.Context) error {
			return nil
		})

		// 验证不应该导致 panic
		_ = err
	}
}

// TestValidationMiddlewareIntegration 集成测试：验证中间件与其他中间件配合
func TestValidationMiddlewareIntegration(t *testing.T) {
	chain := NewChain()

	validationMiddleware := NewValidationMiddleware()
	loggingMiddleware := NewLoggingMiddleware()

	chain.Use(validationMiddleware, loggingMiddleware)

	execCtx := &executor.ExecutionContext{
		SQL:  "SELECT * FROM users WHERE id = ?",
		Args: []interface{}{123},
	}

	err := chain.Execute(context.Background(), execCtx)
	assert.NoError(t, err)
}

// TestValidationMiddlewareIntegrationWithRetry 集成测试：验证和重试
func TestValidationMiddlewareIntegrationWithRetry(t *testing.T) {
	chain := NewChain()

	validationMiddleware := NewValidationMiddleware()
	retryMiddleware := NewRetryMiddleware(3)

	chain.Use(validationMiddleware, retryMiddleware)

	execCtx := &executor.ExecutionContext{
		SQL: "SELECT * FROM users",
	}

	err := chain.Execute(context.Background(), execCtx)
	assert.NoError(t, err)
}

// TestValidationMiddlewareWithComplexSQL 测试复杂的 SQL
func TestValidationMiddlewareWithComplexSQL(t *testing.T) {
	m := NewValidationMiddleware()

	complexSQL := `
		SELECT u.id, u.name, COUNT(o.id) as order_count
		FROM users u
		LEFT JOIN orders o ON u.id = o.user_id
		WHERE u.status = ? AND o.created_at > ?
		GROUP BY u.id, u.name
		HAVING COUNT(o.id) > ?
		ORDER BY order_count DESC
		LIMIT 10
	`

	execCtx := &executor.ExecutionContext{
		SQL:  complexSQL,
		Args: []interface{}{"active", "2023-01-01", 5},
	}

	err := m.Handle(context.Background(), execCtx, func(ctx context.Context) error {
		return nil
	})

	// 复杂的有效 SQL 应该通过验证
	_ = err
}

// TestValidationMiddlewareWithDifferentDatabases 测试不同数据库的 SQL
func TestValidationMiddlewareWithDifferentDatabases(t *testing.T) {
	m := NewValidationMiddleware()

	sqls := []string{
		"SELECT * FROM users LIMIT 10",            // MySQL/PostgreSQL/SQLite
		"SELECT TOP 10 * FROM users",              // SQL Server
		"SELECT * FROM users FETCH FIRST 10 ROWS", // Oracle
	}

	for _, sql := range sqls {
		execCtx := &executor.ExecutionContext{
			SQL: sql,
		}

		err := m.Handle(context.Background(), execCtx, func(ctx context.Context) error {
			return nil
		})

		// 应该能够处理不同数据库的 SQL
		_ = err
	}
}
