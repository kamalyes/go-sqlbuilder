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
	"github.com/kamalyes/go-sqlbuilder/logger"
)

// LoggingMiddleware 日志中间件
// 记录查询的 SQL、参数、执行时间等信息
type LoggingMiddleware struct {
	name   string
	logger logger.Logger
	level  string // debug, info, warn, error
}

// NewLoggingMiddleware 创建日志中间件
func NewLoggingMiddleware() Middleware {
	return &LoggingMiddleware{
		name:   "logging",
		logger: logger.NewNoOpLogger(),
		level:  "info",
	}
}

// NewLoggingMiddlewareWithLogger 创建日志中间件，指定 logger
func NewLoggingMiddlewareWithLogger(l logger.Logger) Middleware {
	return &LoggingMiddleware{
		name:   "logging",
		logger: l,
		level:  "info",
	}
}

// NewLoggingMiddlewareWithLevel 创建日志中间件，指定日志级别
func NewLoggingMiddlewareWithLevel(level string) Middleware {
	return &LoggingMiddleware{
		name:   "logging",
		logger: logger.NewNoOpLogger(),
		level:  level,
	}
}

// Name 返回中间件名称
func (m *LoggingMiddleware) Name() string {
	return m.name
}

// Handle 处理请求
func (m *LoggingMiddleware) Handle(ctx context.Context, execCtx *executor.ExecutionContext, next Next) error {
	if m.logger == nil {
		return next(ctx)
	}

	// 记录执行前的信息
	startTime := time.Now()

	// 执行下一个中间件
	err := next(ctx)

	// 记录执行后的信息
	duration := time.Since(startTime)

	// 构建日志字段
	fields := map[string]interface{}{
		"duration_ms": duration.Milliseconds(),
		"duration":    duration.String(),
		"sql":         execCtx.SQL,
	}

	if len(execCtx.Args) > 0 {
		fields["args"] = execCtx.Args
	}

	// 根据错误情况记录日志
	if err != nil {
		fields["error"] = err.Error()
		m.logger.Errorf("SQL execution failed: %v, SQL: %s, Duration: %s",
			err, execCtx.SQL, duration.String())
	} else {
		message := fmt.Sprintf("SQL executed successfully, Duration: %s", duration.String())
		m.logWithLevel(message, fields)
	}

	return err
}

// logWithLevel 根据级别记录日志
func (m *LoggingMiddleware) logWithLevel(message string, fields map[string]interface{}) {
	if m.logger == nil {
		return
	}

	switch m.level {
	case "debug":
		m.logger.Debugf("%s: %v", message, fields)
	case "warn":
		m.logger.Warnf("%s: %v", message, fields)
	case "error":
		m.logger.Errorf("%s: %v", message, fields)
	default: // info
		m.logger.Infof("%s: %v", message, fields)
	}
}

// SetLevel 设置日志级别
func (m *LoggingMiddleware) SetLevel(level string) {
	m.level = level
}

// GetLevel 获取日志级别
func (m *LoggingMiddleware) GetLevel() string {
	return m.level
}

// SQLLoggingMiddleware SQL 专用日志中间件
// 只记录 SQL 和执行时间，不记录其他信息
type SQLLoggingMiddleware struct {
	name   string
	logger logger.Logger
}

// NewSQLLoggingMiddleware 创建 SQL 专用日志中间件
func NewSQLLoggingMiddleware() Middleware {
	return &SQLLoggingMiddleware{
		name:   "sql_logging",
		logger: logger.NewNoOpLogger(),
	}
}

// Name 返回中间件名称
func (m *SQLLoggingMiddleware) Name() string {
	return m.name
}

// Handle 处理请求
func (m *SQLLoggingMiddleware) Handle(ctx context.Context, execCtx *executor.ExecutionContext, next Next) error {
	if m.logger == nil {
		return next(ctx)
	}

	startTime := time.Now()
	err := next(ctx)
	duration := time.Since(startTime)

	if err != nil {
		m.logger.Errorf("[SQL FAILED] %s (duration: %s, error: %v)", execCtx.SQL, duration.String(), err)
	} else {
		m.logger.Infof("[SQL OK] %s (duration: %s)", execCtx.SQL, duration.String())
	}

	return err
}

// SlowQueryLoggingMiddleware 慢查询日志中间件
// 只记录执行时间超过阈值的查询
type SlowQueryLoggingMiddleware struct {
	name      string
	logger    logger.Logger
	threshold time.Duration
}

// NewSlowQueryLoggingMiddleware 创建慢查询日志中间件
// threshold: 超过此时长的查询会被记录为慢查询
func NewSlowQueryLoggingMiddleware(threshold time.Duration) Middleware {
	if threshold <= 0 {
		threshold = 100 * time.Millisecond
	}
	return &SlowQueryLoggingMiddleware{
		name:      "slow_query_logging",
		logger:    logger.NewNoOpLogger(),
		threshold: threshold,
	}
}

// Name 返回中间件名称
func (m *SlowQueryLoggingMiddleware) Name() string {
	return m.name
}

// Handle 处理请求
func (m *SlowQueryLoggingMiddleware) Handle(ctx context.Context, execCtx *executor.ExecutionContext, next Next) error {
	if m.logger == nil {
		return next(ctx)
	}

	startTime := time.Now()
	err := next(ctx)
	duration := time.Since(startTime)

	// 只记录超过阈值的查询
	if duration > m.threshold {
		fields := map[string]interface{}{
			"duration_ms":  duration.Milliseconds(),
			"sql":          execCtx.SQL,
			"threshold_ms": m.threshold.Milliseconds(),
		}

		if len(execCtx.Args) > 0 {
			fields["args"] = execCtx.Args
		}

		if err != nil {
			fields["error"] = err.Error()
			m.logger.Warnf("SLOW QUERY DETECTED (exceeded %dms): SQL: %s, Duration: %s, Error: %v",
				m.threshold.Milliseconds(), execCtx.SQL, duration.String(), err)
		} else {
			m.logger.Warnf("SLOW QUERY DETECTED (exceeded %dms): SQL: %s, Duration: %s",
				m.threshold.Milliseconds(), execCtx.SQL, duration.String())
		}
	}

	return err
}

// SetThreshold 设置慢查询阈值
func (m *SlowQueryLoggingMiddleware) SetThreshold(threshold time.Duration) {
	if threshold > 0 {
		m.threshold = threshold
	}
}

// GetThreshold 获取慢查询阈值
func (m *SlowQueryLoggingMiddleware) GetThreshold() time.Duration {
	return m.threshold
}
