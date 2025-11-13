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
	"time"

	"github.com/kamalyes/go-sqlbuilder/executor"
	"github.com/stretchr/testify/assert"
)

// TestNewTimeoutMiddleware 测试创建超时中间件
func TestNewTimeoutMiddleware(t *testing.T) {
	m := NewTimeoutMiddleware()
	assert.NotNil(t, m)
	assert.Equal(t, "timeout", m.Name())
}

// TestNewTimeoutMiddlewareWithDuration 测试指定超时时间
func TestNewTimeoutMiddlewareWithDuration(t *testing.T) {
	m := NewTimeoutMiddlewareWithDuration(500 * time.Millisecond)
	assert.NotNil(t, m)
	assert.Equal(t, "timeout", m.Name())
}

// TestTimeoutMiddlewareSuccess 测试成功执行
func TestTimeoutMiddlewareSuccess(t *testing.T) {
	m := NewTimeoutMiddlewareWithDuration(1000 * time.Millisecond)

	execCtx := &executor.ExecutionContext{
		SQL: "SELECT * FROM users",
	}

	err := m.Handle(context.Background(), execCtx, func(ctx context.Context) error {
		time.Sleep(100 * time.Millisecond)
		return nil
	})

	assert.NoError(t, err)
}

// TestTimeoutMiddlewareTimeout 测试超时
func TestTimeoutMiddlewareTimeout(t *testing.T) {
	m := NewTimeoutMiddlewareWithDuration(100 * time.Millisecond)

	execCtx := &executor.ExecutionContext{
		SQL: "SELECT * FROM users",
	}

	err := m.Handle(context.Background(), execCtx, func(ctx context.Context) error {
		time.Sleep(500 * time.Millisecond)
		return nil
	})

	assert.Error(t, err)
	// 超时会返回一个包装后的错误，包含"timeout"关键字
	assert.Contains(t, err.Error(), "timeout")
}

// TestTimeoutMiddlewareEdgeCase 测试边界情况（刚好在超时时）
func TestTimeoutMiddlewareEdgeCase(t *testing.T) {
	m := NewTimeoutMiddlewareWithDuration(50 * time.Millisecond)

	execCtx := &executor.ExecutionContext{
		SQL: "SELECT * FROM users",
	}

	err := m.Handle(context.Background(), execCtx, func(ctx context.Context) error {
		time.Sleep(40 * time.Millisecond)
		return nil
	})

	assert.NoError(t, err)
}

// TestTimeoutMiddlewareWithError 测试超时前发生错误
func TestTimeoutMiddlewareWithError(t *testing.T) {
	m := NewTimeoutMiddlewareWithDuration(1000 * time.Millisecond)

	expectedErr := errors.New("query error")
	execCtx := &executor.ExecutionContext{
		SQL: "SELECT * FROM users",
	}

	err := m.Handle(context.Background(), execCtx, func(ctx context.Context) error {
		return expectedErr
	})

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

// TestTimeoutMiddlewareContextCancellation 测试上下文取消
func TestTimeoutMiddlewareContextCancellation(t *testing.T) {
	m := NewTimeoutMiddlewareWithDuration(1000 * time.Millisecond)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // 立即取消

	execCtx := &executor.ExecutionContext{
		SQL: "SELECT * FROM users",
	}

	err := m.Handle(ctx, execCtx, func(ctx context.Context) error {
		time.Sleep(100 * time.Millisecond)
		return nil
	})

	assert.Error(t, err)
}

// TestTimeoutMiddlewareZeroDuration 测试零超时时间
func TestTimeoutMiddlewareZeroDuration(t *testing.T) {
	m := NewTimeoutMiddlewareWithDuration(0)
	tm := m.(*TimeoutMiddleware)

	// 应该使用默认超时时间
	assert.True(t, tm.GetTimeout() > 0)
}

// TestTimeoutMiddlewareNegativeDuration 测试负数超时时间
func TestTimeoutMiddlewareNegativeDuration(t *testing.T) {
	m := NewTimeoutMiddlewareWithDuration(-100 * time.Millisecond)
	tm := m.(*TimeoutMiddleware)

	// 应该使用默认超时时间
	assert.True(t, tm.GetTimeout() > 0)
}

// TestTimeoutMiddlewareSetTimeout 测试设置超时时间
func TestTimeoutMiddlewareSetTimeout(t *testing.T) {
	m := NewTimeoutMiddleware()
	tm := m.(*TimeoutMiddleware)

	original := tm.GetTimeout()
	tm.SetTimeout(500 * time.Millisecond)
	assert.Equal(t, 500*time.Millisecond, tm.GetTimeout())
	assert.NotEqual(t, original, tm.GetTimeout())
}

// TestTimeoutMiddlewareGetTimeout 测试获取超时时间
func TestTimeoutMiddlewareGetTimeout(t *testing.T) {
	m := NewTimeoutMiddlewareWithDuration(300 * time.Millisecond)
	tm := m.(*TimeoutMiddleware)

	assert.Equal(t, 300*time.Millisecond, tm.GetTimeout())
}

// TestTimeoutMiddlewareConcurrent 测试并发执行
func TestTimeoutMiddlewareConcurrent(t *testing.T) {
	m := NewTimeoutMiddlewareWithDuration(500 * time.Millisecond)

	execCtx := &executor.ExecutionContext{
		SQL: "SELECT * FROM users",
	}

	done := make(chan error, 5)
	for i := 0; i < 5; i++ {
		go func() {
			err := m.Handle(context.Background(), execCtx, func(ctx context.Context) error {
				time.Sleep(100 * time.Millisecond)
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

// TestTimeoutMiddlewareWithExistingDeadline 测试已有截止时间的上下文
func TestTimeoutMiddlewareWithExistingDeadline(t *testing.T) {
	m := NewTimeoutMiddlewareWithDuration(1000 * time.Millisecond)

	// 创建已有截止时间的上下文（更短的）
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	execCtx := &executor.ExecutionContext{
		SQL: "SELECT * FROM users",
	}

	err := m.Handle(ctx, execCtx, func(ctx context.Context) error {
		time.Sleep(500 * time.Millisecond)
		return nil
	})

	assert.Error(t, err)
	// 超时会返回一个包装后的错误或上下文取消错误
	assert.NotNil(t, err)
}

// TestTimeoutMiddlewareIntegration 集成测试：超时中间件与其他中间件配合
func TestTimeoutMiddlewareIntegration(t *testing.T) {
	chain := NewChain()

	timeoutMiddleware := NewTimeoutMiddlewareWithDuration(500 * time.Millisecond)
	retryMiddleware := NewRetryMiddleware(2)

	chain.Use(timeoutMiddleware, retryMiddleware)

	execCtx := &executor.ExecutionContext{
		SQL: "SELECT * FROM users",
	}

	err := chain.Execute(context.Background(), execCtx)
	assert.NoError(t, err)
}

// TestTimeoutMiddlewareIntegrationWithLogging 集成测试：超时、重试和日志中间件
func TestTimeoutMiddlewareIntegrationWithLogging(t *testing.T) {
	chain := NewChain()

	timeoutMiddleware := NewTimeoutMiddlewareWithDuration(500 * time.Millisecond)
	retryMiddleware := NewRetryMiddleware(2)
	loggingMiddleware := NewLoggingMiddleware()

	chain.Use(timeoutMiddleware, retryMiddleware, loggingMiddleware)

	execCtx := &executor.ExecutionContext{
		SQL: "SELECT * FROM users WHERE id = ?",
	}

	err := chain.Execute(context.Background(), execCtx)
	assert.NoError(t, err)
}

// TestMultipleTimeoutMiddleware 测试多个超时中间件
func TestMultipleTimeoutMiddleware(t *testing.T) {
	m1 := NewTimeoutMiddlewareWithDuration(200 * time.Millisecond)

	execCtx := &executor.ExecutionContext{
		SQL: "SELECT * FROM users",
	}

	// 第一个中间件更严格
	err := m1.Handle(context.Background(), execCtx, func(ctx context.Context) error {
		time.Sleep(100 * time.Millisecond)
		return nil
	})

	assert.NoError(t, err)
}

// TestTimeoutMiddlewareVeryShortTimeout 测试很短的超时
func TestTimeoutMiddlewareVeryShortTimeout(t *testing.T) {
	m := NewTimeoutMiddlewareWithDuration(1 * time.Millisecond)

	execCtx := &executor.ExecutionContext{
		SQL: "SELECT * FROM users",
	}

	err := m.Handle(context.Background(), execCtx, func(ctx context.Context) error {
		time.Sleep(100 * time.Millisecond)
		return nil
	})

	assert.Error(t, err)
	// 超时会返回一个包装后的错误，包含"timeout"关键字
	assert.Contains(t, err.Error(), "timeout")
}

// TestTimeoutMiddlewareVeryLongTimeout 测试很长的超时
func TestTimeoutMiddlewareVeryLongTimeout(t *testing.T) {
	m := NewTimeoutMiddlewareWithDuration(10 * time.Second)

	execCtx := &executor.ExecutionContext{
		SQL: "SELECT * FROM users",
	}

	err := m.Handle(context.Background(), execCtx, func(ctx context.Context) error {
		time.Sleep(100 * time.Millisecond)
		return nil
	})

	assert.NoError(t, err)
}
