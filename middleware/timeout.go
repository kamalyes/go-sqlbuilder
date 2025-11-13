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
	"fmt"
	"time"

	"github.com/kamalyes/go-sqlbuilder/executor"
)

// TimeoutMiddleware 超时中间件
// 限制查询执行时间，超时则取消查询
type TimeoutMiddleware struct {
	name    string
	timeout time.Duration
}

// NewTimeoutMiddleware 创建超时中间件
// 使用默认超时时间（30 秒）
func NewTimeoutMiddleware() Middleware {
	return &TimeoutMiddleware{
		name:    "timeout",
		timeout: 30 * time.Second,
	}
}

// NewTimeoutMiddlewareWithDuration 创建超时中间件，指定超时时间
func NewTimeoutMiddlewareWithDuration(timeout time.Duration) Middleware {
	if timeout <= 0 {
		timeout = 30 * time.Second
	}

	return &TimeoutMiddleware{
		name:    "timeout",
		timeout: timeout,
	}
}

// Name 返回中间件名称
func (m *TimeoutMiddleware) Name() string {
	return m.name
}

// Handle 处理请求
func (m *TimeoutMiddleware) Handle(ctx context.Context, execCtx *executor.ExecutionContext, next Next) error {
	// 创建超时上下文
	timeoutCtx, cancel := context.WithTimeout(ctx, m.timeout)
	defer cancel()

	// 创建错误通道
	errChan := make(chan error, 1)

	// 在 goroutine 中执行下一个中间件
	go func() {
		errChan <- next(timeoutCtx)
	}()

	// 等待执行完成或超时
	select {
	case err := <-errChan:
		return err
	case <-timeoutCtx.Done():
		if timeoutCtx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("query execution timeout: %s", m.timeout.String())
		}
		return timeoutCtx.Err()
	}
}

// SetTimeout 设置超时时间
func (m *TimeoutMiddleware) SetTimeout(timeout time.Duration) {
	if timeout > 0 {
		m.timeout = timeout
	}
}

// GetTimeout 获取超时时间
func (m *TimeoutMiddleware) GetTimeout() time.Duration {
	return m.timeout
}

// SetDuration 设置超时时间（别名）
func (m *TimeoutMiddleware) SetDuration(duration time.Duration) {
	if duration > 0 {
		m.timeout = duration
	}
}

// GetDuration 获取超时时间（别名）
func (m *TimeoutMiddleware) GetDuration() time.Duration {
	return m.timeout
}

// AdaptiveTimeoutMiddleware 自适应超时中间件
// 根据历史查询耗时自动调整超时时间
type AdaptiveTimeoutMiddleware struct {
	name            string
	baseTimeout     time.Duration
	percentile      float64 // 99th percentile
	maxTimeout      time.Duration
	minTimeout      time.Duration
	durations       []time.Duration
	totalExecutions int64
	timeoutExceeds  int64
}

// NewAdaptiveTimeoutMiddleware 创建自适应超时中间件
func NewAdaptiveTimeoutMiddleware(baseTimeout time.Duration) Middleware {
	if baseTimeout <= 0 {
		baseTimeout = 30 * time.Second
	}

	return &AdaptiveTimeoutMiddleware{
		name:            "adaptive_timeout",
		baseTimeout:     baseTimeout,
		percentile:      0.99, // 99th percentile
		maxTimeout:      5 * time.Minute,
		minTimeout:      100 * time.Millisecond,
		durations:       make([]time.Duration, 0),
		totalExecutions: 0,
	}
}

// Name 返回中间件名称
func (m *AdaptiveTimeoutMiddleware) Name() string {
	return m.name
}

// Handle 处理请求
func (m *AdaptiveTimeoutMiddleware) Handle(ctx context.Context, execCtx *executor.ExecutionContext, next Next) error {
	// 获取当前的有效超时时间
	timeout := m.getEffectiveTimeout()

	// 创建超时上下文
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	startTime := time.Now()

	// 创建错误通道
	errChan := make(chan error, 1)

	// 在 goroutine 中执行下一个中间件
	go func() {
		errChan <- next(timeoutCtx)
	}()

	// 等待执行完成或超时
	var err error
	select {
	case err = <-errChan:
		// 正常完成
	case <-timeoutCtx.Done():
		m.recordTimeoutExceed()
		if timeoutCtx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("query execution timeout: %s", timeout.String())
		}
		return timeoutCtx.Err()
	}

	// 记录执行时间
	duration := time.Since(startTime)
	m.recordDuration(duration)

	return err
}

// recordDuration 记录查询执行时间
func (m *AdaptiveTimeoutMiddleware) recordDuration(duration time.Duration) {
	m.durations = append(m.durations, duration)
	m.totalExecutions++

	// 只保留最后 1000 次执行的时间
	if len(m.durations) > 1000 {
		m.durations = m.durations[1:]
	}
}

// recordTimeoutExceed 记录超时
func (m *AdaptiveTimeoutMiddleware) recordTimeoutExceed() {
	m.timeoutExceeds++
}

// getEffectiveTimeout 获取有效的超时时间
func (m *AdaptiveTimeoutMiddleware) getEffectiveTimeout() time.Duration {
	if len(m.durations) < 10 {
		// 数据不足，使用基础超时
		return m.baseTimeout
	}

	// 计算百分位数
	percentileDuration := m.calculatePercentile(m.percentile)

	// 使用百分位数 + 20% 的缓冲
	timeout := time.Duration(float64(percentileDuration) * 1.2)

	// 限制在最大和最小值之间
	if timeout > m.maxTimeout {
		timeout = m.maxTimeout
	}
	if timeout < m.minTimeout {
		timeout = m.minTimeout
	}

	return timeout
}

// calculatePercentile 计算百分位数
func (m *AdaptiveTimeoutMiddleware) calculatePercentile(percentile float64) time.Duration {
	if len(m.durations) == 0 {
		return m.baseTimeout
	}

	// 简单的百分位数计算
	// 更复杂的实现可以使用排序和插值
	index := int(float64(len(m.durations)-1) * percentile)
	if index >= len(m.durations) {
		index = len(m.durations) - 1
	}

	return m.durations[index]
}

// GetStatistics 获取统计信息
func (m *AdaptiveTimeoutMiddleware) GetStatistics() TimeoutStatistics {
	maxDuration := time.Duration(0)
	minDuration := time.Duration(^uint64(0) >> 1)
	totalDuration := time.Duration(0)

	for _, d := range m.durations {
		totalDuration += d
		if d > maxDuration {
			maxDuration = d
		}
		if d < minDuration {
			minDuration = d
		}
	}

	avgDuration := time.Duration(0)
	if len(m.durations) > 0 {
		avgDuration = totalDuration / time.Duration(len(m.durations))
	}

	return TimeoutStatistics{
		TotalExecutions:   m.totalExecutions,
		TimeoutExceeds:    m.timeoutExceeds,
		AverageDuration:   avgDuration,
		MaxDuration:       maxDuration,
		MinDuration:       minDuration,
		CurrentTimeout:    m.getEffectiveTimeout(),
		TimeoutExceedRate: float64(m.timeoutExceeds) / float64(m.totalExecutions),
	}
}

// TimeoutStatistics 超时统计信息
type TimeoutStatistics struct {
	TotalExecutions   int64
	TimeoutExceeds    int64
	AverageDuration   time.Duration
	MaxDuration       time.Duration
	MinDuration       time.Duration
	CurrentTimeout    time.Duration
	TimeoutExceedRate float64
}

// Reset 重置统计信息
func (m *AdaptiveTimeoutMiddleware) Reset() {
	m.durations = make([]time.Duration, 0)
	m.totalExecutions = 0
	m.timeoutExceeds = 0
}

// SetPercentile 设置百分位数
func (m *AdaptiveTimeoutMiddleware) SetPercentile(percentile float64) {
	if percentile > 0 && percentile < 1 {
		m.percentile = percentile
	}
}

// SetMaxTimeout 设置最大超时时间
func (m *AdaptiveTimeoutMiddleware) SetMaxTimeout(maxTimeout time.Duration) {
	if maxTimeout > 0 {
		m.maxTimeout = maxTimeout
	}
}

// SetMinTimeout 设置最小超时时间
func (m *AdaptiveTimeoutMiddleware) SetMinTimeout(minTimeout time.Duration) {
	if minTimeout > 0 {
		m.minTimeout = minTimeout
	}
}
