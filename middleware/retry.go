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
	"time"

	"github.com/kamalyes/go-sqlbuilder/executor"
)

// RetryMiddleware 重试中间件
// 当查询失败时自动重试
type RetryMiddleware struct {
	name        string
	maxAttempts int
	backoff     BackoffStrategy
	retryable   RetryableChecker
}

// BackoffStrategy 重试等待策略
type BackoffStrategy interface {
	// NextBackoff 获取下一次重试的等待时间
	NextBackoff(attempt int) time.Duration
}

// RetryableChecker 检查错误是否可重试
type RetryableChecker interface {
	// IsRetryable 检查给定的错误是否应该重试
	IsRetryable(err error) bool
}

// DefaultBackoff 默认重试等待策略（指数退避）
type DefaultBackoff struct {
	initialDelay time.Duration
	maxDelay     time.Duration
}

// NewDefaultBackoff 创建默认重试等待策略
func NewDefaultBackoff() BackoffStrategy {
	return &DefaultBackoff{
		initialDelay: 100 * time.Millisecond,
		maxDelay:     10 * time.Second,
	}
}

// NextBackoff 获取下一次重试的等待时间
func (b *DefaultBackoff) NextBackoff(attempt int) time.Duration {
	if attempt <= 0 {
		return 0
	}

	// 指数退避：初始延迟 * 2^(attempt-1)
	delay := b.initialDelay
	for i := 1; i < attempt; i++ {
		delay *= 2
	}

	if delay > b.maxDelay {
		delay = b.maxDelay
	}

	return delay
}

// LinearBackoff 线性重试等待策略
type LinearBackoff struct {
	interval time.Duration
}

// NewLinearBackoff 创建线性重试等待策略
func NewLinearBackoff(interval time.Duration) BackoffStrategy {
	return &LinearBackoff{
		interval: interval,
	}
}

// NextBackoff 获取下一次重试的等待时间
func (b *LinearBackoff) NextBackoff(attempt int) time.Duration {
	return time.Duration(attempt) * b.interval
}

// DefaultRetryableChecker 默认的可重试错误检查器
type DefaultRetryableChecker struct{}

// IsRetryable 检查错误是否可重试
func (c *DefaultRetryableChecker) IsRetryable(err error) bool {
	if err == nil {
		return false
	}

	// 可以添加更多的可重试错误类型
	// 例如：网络错误、超时、死锁等

	return true
}

// NewRetryMiddleware 创建重试中间件
func NewRetryMiddleware(maxAttempts int) Middleware {
	if maxAttempts <= 0 {
		maxAttempts = 3
	}

	return &RetryMiddleware{
		name:        "retry",
		maxAttempts: maxAttempts,
		backoff:     NewDefaultBackoff(),
		retryable:   &DefaultRetryableChecker{},
	}
}

// NewRetryMiddlewareWithConfig 创建重试中间件，指定配置
func NewRetryMiddlewareWithConfig(maxAttempts int, backoff BackoffStrategy, retryable RetryableChecker) Middleware {
	if maxAttempts <= 0 {
		maxAttempts = 3
	}

	if backoff == nil {
		backoff = NewDefaultBackoff()
	}

	if retryable == nil {
		retryable = &DefaultRetryableChecker{}
	}

	return &RetryMiddleware{
		name:        "retry",
		maxAttempts: maxAttempts,
		backoff:     backoff,
		retryable:   retryable,
	}
}

// Name 返回中间件名称
func (m *RetryMiddleware) Name() string {
	return m.name
}

// Handle 处理请求
func (m *RetryMiddleware) Handle(ctx context.Context, execCtx *executor.ExecutionContext, next Next) error {
	var lastErr error

	for attempt := 0; attempt < m.maxAttempts; attempt++ {
		err := next(ctx)

		if err == nil {
			// 成功执行
			return nil
		}

		lastErr = err

		// 检查错误是否可重试
		if !m.retryable.IsRetryable(err) {
			return err
		}

		// 如果是最后一次尝试，直接返回错误
		if attempt == m.maxAttempts-1 {
			return err
		}

		// 等待后重试
		backoffDuration := m.backoff.NextBackoff(attempt + 1)
		select {
		case <-time.After(backoffDuration):
			// 继续重试
		case <-ctx.Done():
			// 上下文已取消
			return ctx.Err()
		}
	}

	return lastErr
}

// CircuitBreakerRetryMiddleware 断路器重试中间件
// 当错误率过高时，停止重试
type CircuitBreakerRetryMiddleware struct {
	name             string
	maxAttempts      int
	backoff          BackoffStrategy
	retryable        RetryableChecker
	failureThreshold int64
	successThreshold int64
	timeout          time.Duration
	state            CircuitState
	failureCount     int64
	successCount     int64
	lastFailureTime  time.Time
	lastSuccessTime  time.Time
}

// CircuitState 断路器状态
type CircuitState string

const (
	// CircuitStateClosed 闭合状态（正常）
	CircuitStateClosed CircuitState = "closed"
	// CircuitStateOpen 开启状态（断路）
	CircuitStateOpen CircuitState = "open"
	// CircuitStateHalfOpen 半开状态（试探）
	CircuitStateHalfOpen CircuitState = "half_open"
)

// NewCircuitBreakerRetryMiddleware 创建断路器重试中间件
func NewCircuitBreakerRetryMiddleware(maxAttempts int) Middleware {
	if maxAttempts <= 0 {
		maxAttempts = 3
	}

	return &CircuitBreakerRetryMiddleware{
		name:             "circuit_breaker_retry",
		maxAttempts:      maxAttempts,
		backoff:          NewDefaultBackoff(),
		retryable:        &DefaultRetryableChecker{},
		failureThreshold: 5,
		successThreshold: 2,
		timeout:          30 * time.Second,
		state:            CircuitStateClosed,
		failureCount:     0,
		successCount:     0,
	}
}

// Name 返回中间件名称
func (m *CircuitBreakerRetryMiddleware) Name() string {
	return m.name
}

// Handle 处理请求
func (m *CircuitBreakerRetryMiddleware) Handle(ctx context.Context, execCtx *executor.ExecutionContext, next Next) error {
	// 检查断路器状态
	if m.state == CircuitStateOpen {
		// 检查是否应该进入半开状态
		if time.Since(m.lastFailureTime) > m.timeout {
			m.state = CircuitStateHalfOpen
			m.failureCount = 0
			m.successCount = 0
		} else {
			// 断路器开启，直接返回错误
			return errors.New("circuit breaker is open")
		}
	}

	var lastErr error

	for attempt := 0; attempt < m.maxAttempts; attempt++ {
		err := next(ctx)

		if err == nil {
			// 成功执行
			m.recordSuccess()
			return nil
		}

		lastErr = err
		m.recordFailure()

		// 检查错误是否可重试
		if !m.retryable.IsRetryable(err) {
			return err
		}

		// 如果是最后一次尝试，直接返回错误
		if attempt == m.maxAttempts-1 {
			return err
		}

		// 等待后重试
		backoffDuration := m.backoff.NextBackoff(attempt + 1)
		select {
		case <-time.After(backoffDuration):
			// 继续重试
		case <-ctx.Done():
			// 上下文已取消
			return ctx.Err()
		}
	}

	return lastErr
}

// recordSuccess 记录成功
func (m *CircuitBreakerRetryMiddleware) recordSuccess() {
	m.successCount++
	m.failureCount = 0
	m.lastSuccessTime = time.Now()

	// 如果在半开状态下连续成功，则关闭断路器
	if m.state == CircuitStateHalfOpen && m.successCount >= m.successThreshold {
		m.state = CircuitStateClosed
		m.failureCount = 0
		m.successCount = 0
	}
}

// recordFailure 记录失败
func (m *CircuitBreakerRetryMiddleware) recordFailure() {
	m.failureCount++
	m.successCount = 0
	m.lastFailureTime = time.Now()

	// 如果失败次数超过阈值，则打开断路器
	if m.failureCount >= m.failureThreshold {
		m.state = CircuitStateOpen
		m.failureCount = 0
		m.successCount = 0
	}
}

// GetState 获取断路器状态
func (m *CircuitBreakerRetryMiddleware) GetState() CircuitState {
	return m.state
}

// Reset 重置断路器
func (m *CircuitBreakerRetryMiddleware) Reset() {
	m.state = CircuitStateClosed
	m.failureCount = 0
	m.successCount = 0
}
