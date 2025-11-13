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

// TestNewRetryMiddleware 测试创建重试中间件
func TestNewRetryMiddleware(t *testing.T) {
	m := NewRetryMiddleware(3)
	assert.NotNil(t, m)
	assert.Equal(t, "retry", m.Name())
}

// TestRetryMiddlewareSuccess 测试一次成功的重试
func TestRetryMiddlewareSuccess(t *testing.T) {
	m := NewRetryMiddleware(3)

	callCount := 0
	execCtx := &executor.ExecutionContext{
		SQL: "SELECT * FROM users",
	}

	err := m.Handle(context.Background(), execCtx, func(ctx context.Context) error {
		callCount++
		return nil
	})

	assert.NoError(t, err)
	assert.Equal(t, 1, callCount) // 应该只调用一次
}

// TestRetryMiddlewareRetryOnce 测试重试一次
func TestRetryMiddlewareRetryOnce(t *testing.T) {
	m := NewRetryMiddleware(3)

	callCount := 0
	execCtx := &executor.ExecutionContext{
		SQL: "SELECT * FROM users",
	}

	err := m.Handle(context.Background(), execCtx, func(ctx context.Context) error {
		callCount++
		if callCount == 1 {
			return errors.New("first attempt failed")
		}
		return nil
	})

	assert.NoError(t, err)
	assert.Equal(t, 2, callCount) // 调用两次，第二次成功
}

// TestRetryMiddlewareExhausted 测试重试次数耗尽
func TestRetryMiddlewareExhausted(t *testing.T) {
	m := NewRetryMiddleware(2)

	callCount := 0
	expectedErr := errors.New("persistent error")
	execCtx := &executor.ExecutionContext{
		SQL: "SELECT * FROM users",
	}

	err := m.Handle(context.Background(), execCtx, func(ctx context.Context) error {
		callCount++
		return expectedErr
	})

	assert.Error(t, err)
	assert.Equal(t, 2, callCount) // 调用两次，都失败
	assert.Equal(t, expectedErr, err)
}

// TestRetryMiddlewareMultipleRetries 测试多次重试
func TestRetryMiddlewareMultipleRetries(t *testing.T) {
	m := NewRetryMiddleware(5)

	callCount := 0
	execCtx := &executor.ExecutionContext{
		SQL: "SELECT * FROM users",
	}

	err := m.Handle(context.Background(), execCtx, func(ctx context.Context) error {
		callCount++
		if callCount < 4 {
			return errors.New("not yet")
		}
		return nil
	})

	assert.NoError(t, err)
	assert.Equal(t, 4, callCount)
}

// TestRetryMiddlewareZeroAttempts 测试零次重试
func TestRetryMiddlewareZeroAttempts(t *testing.T) {
	m := NewRetryMiddleware(0)

	callCount := 0
	execCtx := &executor.ExecutionContext{
		SQL: "SELECT * FROM users",
	}

	err := m.Handle(context.Background(), execCtx, func(ctx context.Context) error {
		callCount++
		return errors.New("error")
	})

	assert.Error(t, err)
	// 当 maxAttempts <= 0 时，实现将其设置为默认值 3
	// 所以会尝试 3 次
	assert.Equal(t, 3, callCount)
}

// TestRetryMiddlewareNegativeAttempts 测试负数次重试
func TestRetryMiddlewareNegativeAttempts(t *testing.T) {
	m := NewRetryMiddleware(-1)

	callCount := 0
	execCtx := &executor.ExecutionContext{
		SQL: "SELECT * FROM users",
	}

	err := m.Handle(context.Background(), execCtx, func(ctx context.Context) error {
		callCount++
		return errors.New("error")
	})

	assert.Error(t, err)
	// 当 maxAttempts < 0 时，实现将其设置为默认值 3
	// 所以会尝试 3 次
	assert.Equal(t, 3, callCount)
}

// TestRetryMiddlewareContextCancellation 测试上下文取消
func TestRetryMiddlewareContextCancellation(t *testing.T) {
	m := NewRetryMiddleware(5)

	callCount := 0
	execCtx := &executor.ExecutionContext{
		SQL: "SELECT * FROM users",
	}

	// 创建可取消的上下文
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		// 在第二次调用后取消上下文
		select {
		case <-time.After(50 * time.Millisecond):
			cancel()
		}
	}()

	err := m.Handle(ctx, execCtx, func(ctx context.Context) error {
		callCount++
		if callCount <= 3 {
			return errors.New("error")
		}
		return nil
	})

	// 由于上下文被取消，应该在有限的调用后停止
	assert.Error(t, err)
	assert.True(t, callCount <= 5)
}

// TestRetryMiddlewareWithDifferentErrors 测试不同的错误类型
func TestRetryMiddlewareWithDifferentErrors(t *testing.T) {
	m := NewRetryMiddleware(3)

	callCount := 0
	execCtx := &executor.ExecutionContext{
		SQL: "SELECT * FROM users",
	}

	err := m.Handle(context.Background(), execCtx, func(ctx context.Context) error {
		callCount++
		if callCount == 1 {
			return errors.New("connection error")
		}
		if callCount == 2 {
			return errors.New("timeout")
		}
		return nil
	})

	assert.NoError(t, err)
	assert.Equal(t, 3, callCount)
}

// TestRetryMiddlewarePreservesError 测试保留最后的错误
func TestRetryMiddlewarePreservesError(t *testing.T) {
	m := NewRetryMiddleware(2)

	lastErr := errors.New("final error")
	execCtx := &executor.ExecutionContext{
		SQL: "SELECT * FROM users",
	}

	err := m.Handle(context.Background(), execCtx, func(ctx context.Context) error {
		return lastErr
	})

	assert.Error(t, err)
	assert.Equal(t, lastErr, err)
}

// TestNewRetryPolicy 测试创建重试策略
func TestNewRetryPolicy(t *testing.T) {
	backoff := NewDefaultBackoff()
	assert.NotNil(t, backoff)
}

// TestDefaultBackoffNextBackoff 测试默认退避策略
func TestDefaultBackoffNextBackoff(t *testing.T) {
	backoff := NewDefaultBackoff()
	assert.NotNil(t, backoff)

	// 测试指数退避
	d1 := backoff.NextBackoff(1)
	d2 := backoff.NextBackoff(2)
	assert.True(t, d2 > d1) // 第二次延迟应该更长
}

// TestLinearBackoffNextBackoff 测试线性退避策略
func TestLinearBackoffNextBackoff(t *testing.T) {
	backoff := NewLinearBackoff(50 * time.Millisecond)
	assert.NotNil(t, backoff)

	d1 := backoff.NextBackoff(1)
	d2 := backoff.NextBackoff(2)
	assert.Equal(t, 50*time.Millisecond, d1)
	assert.Equal(t, 100*time.Millisecond, d2)
}

// TestNewRetryMiddlewareWithConfig 测试使用配置创建重试中间件
func TestNewRetryMiddlewareWithConfig(t *testing.T) {
	backoff := NewLinearBackoff(50 * time.Millisecond)
	retryable := &DefaultRetryableChecker{}
	m := NewRetryMiddlewareWithConfig(3, backoff, retryable)

	assert.NotNil(t, m)
	assert.Equal(t, "retry", m.Name())
}

// TestRetryMiddlewareIntegration 集成测试：重试中间件与其他中间件配合
func TestRetryMiddlewareIntegration(t *testing.T) {
	chain := NewChain()

	retryMiddleware := NewRetryMiddleware(2)
	loggingMiddleware := NewLoggingMiddleware()

	chain.Use(retryMiddleware, loggingMiddleware)

	execCtx := &executor.ExecutionContext{
		SQL: "SELECT * FROM users",
	}

	err := chain.Execute(context.Background(), execCtx)
	assert.NoError(t, err)
}

// TestRetryMiddlewareWithMetrics 集成测试：重试和指标中间件
func TestRetryMiddlewareWithMetrics(t *testing.T) {
	chain := NewChain()

	retryMiddleware := NewRetryMiddleware(3)
	metricsMiddleware := NewMetricsMiddleware()

	chain.Use(retryMiddleware, metricsMiddleware)

	execCtx := &executor.ExecutionContext{
		SQL: "SELECT * FROM users",
	}

	err := chain.Execute(context.Background(), execCtx)
	assert.NoError(t, err)
}

// TestRetryMiddlewareMaxAttempts 测试最大尝试次数的边界情况
func TestRetryMiddlewareMaxAttempts(t *testing.T) {
	tests := []struct {
		name        string
		maxAttempts int
		callCount   int
		wantErr     bool
	}{
		{
			name:        "single attempt",
			maxAttempts: 1,
			callCount:   0,
			wantErr:     false,
		},
		{
			name:        "three attempts",
			maxAttempts: 3,
			callCount:   0,
			wantErr:     false,
		},
		{
			name:        "large number of attempts",
			maxAttempts: 100,
			callCount:   0,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewRetryMiddleware(tt.maxAttempts)
			execCtx := &executor.ExecutionContext{
				SQL: "SELECT * FROM users",
			}

			err := m.Handle(context.Background(), execCtx, func(ctx context.Context) error {
				return nil
			})

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestRetryMiddlewareWithExecutionContext 测试执行上下文在重试中保持
func TestRetryMiddlewareWithExecutionContext(t *testing.T) {
	m := NewRetryMiddleware(3)

	execCtx := &executor.ExecutionContext{
		SQL:  "SELECT * FROM users WHERE id = ?",
		Args: []interface{}{123},
	}

	callCount := 0
	err := m.Handle(context.Background(), execCtx, func(ctx context.Context) error {
		callCount++
		if callCount == 1 {
			return errors.New("retry")
		}
		return nil
	})

	assert.NoError(t, err)
	assert.Equal(t, "SELECT * FROM users WHERE id = ?", execCtx.SQL)
	assert.Equal(t, []interface{}{123}, execCtx.Args)
}
