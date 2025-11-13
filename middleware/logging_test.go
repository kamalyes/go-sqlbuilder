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
	"github.com/kamalyes/go-sqlbuilder/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockLogger 模拟日志记录器
type MockLogger struct {
	mock.Mock
	logs []string
}

func NewMockLogger() *MockLogger {
	return &MockLogger{
		logs: make([]string, 0),
	}
}

func (m *MockLogger) Debug(msg string, fields ...interface{}) {
	m.Called(msg, fields)
}

func (m *MockLogger) Info(msg string, fields ...interface{}) {
	m.Called(msg, fields)
}

func (m *MockLogger) Warn(msg string, fields ...interface{}) {
	m.Called(msg, fields)
}

func (m *MockLogger) Error(msg string, fields ...interface{}) {
	m.Called(msg, fields)
}

func (m *MockLogger) Fatal(msg string, fields ...interface{}) {
	m.Called(msg, fields)
}

func (m *MockLogger) Debugf(format string, args ...interface{}) {
	m.Called(format, args)
}

func (m *MockLogger) Infof(format string, args ...interface{}) {
	m.Called(format, args)
}

func (m *MockLogger) Warnf(format string, args ...interface{}) {
	m.Called(format, args)
}

func (m *MockLogger) Errorf(format string, args ...interface{}) {
	m.Called(format, args)
}

func (m *MockLogger) Fatalf(format string, args ...interface{}) {
	m.Called(format, args)
}

func (m *MockLogger) With(key string, value interface{}) logger.Logger {
	return m
}

func (m *MockLogger) WithContext(ctx interface{}) logger.Logger {
	return m
}

func (m *MockLogger) SetLevel(level string) logger.Logger {
	return m
}

func (m *MockLogger) GetLevel() string {
	return logger.LevelInfo
}

func (m *MockLogger) StartTimer() logger.Timer {
	return &simpleTimer{}
}

func (m *MockLogger) LogDuration(operation string, duration interface{}, fields ...interface{}) {
}

type simpleTimer struct{}

func (t *simpleTimer) Stop()                   {}
func (t *simpleTimer) Duration() time.Duration { return 0 }

// TestNewLoggingMiddleware 测试创建日志中间件
func TestNewLoggingMiddleware(t *testing.T) {
	m := NewLoggingMiddleware()
	assert.NotNil(t, m)
	assert.Equal(t, "logging", m.Name())
}

// TestNewLoggingMiddlewareWithLogger 测试使用自定义 logger 创建日志中间件
func TestNewLoggingMiddlewareWithLogger(t *testing.T) {
	mockLogger := NewMockLogger()
	m := NewLoggingMiddlewareWithLogger(mockLogger)

	assert.NotNil(t, m)
	assert.Equal(t, "logging", m.Name())
}

// TestNewLoggingMiddlewareWithLevel 测试指定日志级别
func TestNewLoggingMiddlewareWithLevel(t *testing.T) {
	m := NewLoggingMiddlewareWithLevel("debug")
	assert.NotNil(t, m)
	assert.Equal(t, "logging", m.Name())

	lm := m.(*LoggingMiddleware)
	assert.Equal(t, "debug", lm.GetLevel())
}

// TestLoggingMiddlewareSuccess 测试成功执行的日志记录
func TestLoggingMiddlewareSuccess(t *testing.T) {
	mockLogger := NewMockLogger()
	mockLogger.On("Infof", mock.MatchedBy(func(s string) bool {
		return true
	}), mock.Anything).Return()

	m := NewLoggingMiddlewareWithLogger(mockLogger)
	lm := m.(*LoggingMiddleware)
	lm.level = "info"

	execCtx := &executor.ExecutionContext{
		SQL: "SELECT * FROM users",
	}

	err := m.Handle(context.Background(), execCtx, func(ctx context.Context) error {
		return nil
	})

	assert.NoError(t, err)
}

// TestLoggingMiddlewareError 测试错误执行的日志记录
func TestLoggingMiddlewareError(t *testing.T) {
	mockLogger := NewMockLogger()
	mockLogger.On("Errorf", mock.MatchedBy(func(s string) bool {
		return true
	}), mock.Anything).Return()

	m := NewLoggingMiddlewareWithLogger(mockLogger)

	expectedErr := errors.New("query failed")
	execCtx := &executor.ExecutionContext{
		SQL: "SELECT * FROM users",
	}

	err := m.Handle(context.Background(), execCtx, func(ctx context.Context) error {
		return expectedErr
	})

	assert.Equal(t, expectedErr, err)
}

// TestLoggingMiddlewareWithArgs 测试记录查询参数
func TestLoggingMiddlewareWithArgs(t *testing.T) {
	mockLogger := NewMockLogger()
	mockLogger.On("Infof", mock.MatchedBy(func(s string) bool {
		return true
	}), mock.Anything).Return()

	m := NewLoggingMiddlewareWithLogger(mockLogger)

	execCtx := &executor.ExecutionContext{
		SQL:  "SELECT * FROM users WHERE id = ?",
		Args: []interface{}{123},
	}

	err := m.Handle(context.Background(), execCtx, func(ctx context.Context) error {
		return nil
	})

	assert.NoError(t, err)
}

// TestLoggingMiddlewareWithMetrics 测试记录执行指标
func TestLoggingMiddlewareWithMetrics(t *testing.T) {
	mockLogger := NewMockLogger()
	mockLogger.On("Infof", mock.MatchedBy(func(s string) bool {
		return true
	}), mock.Anything).Return()

	m := NewLoggingMiddlewareWithLogger(mockLogger)

	execCtx := &executor.ExecutionContext{
		SQL: "INSERT INTO users VALUES (?)",
	}

	err := m.Handle(context.Background(), execCtx, func(ctx context.Context) error {
		return nil
	})

	assert.NoError(t, err)
}

// TestLoggingMiddlewareSetLevel 测试设置日志级别
func TestLoggingMiddlewareSetLevel(t *testing.T) {
	m := NewLoggingMiddleware()
	lm := m.(*LoggingMiddleware)

	assert.Equal(t, "info", lm.GetLevel())

	lm.SetLevel("debug")
	assert.Equal(t, "debug", lm.GetLevel())

	lm.SetLevel("warn")
	assert.Equal(t, "warn", lm.GetLevel())

	lm.SetLevel("error")
	assert.Equal(t, "error", lm.GetLevel())
}

// TestLoggingMiddlewareWithNilLogger 测试 logger 为 nil
func TestLoggingMiddlewareWithNilLogger(t *testing.T) {
	m := NewLoggingMiddleware()
	lm := m.(*LoggingMiddleware)
	lm.logger = nil

	execCtx := &executor.ExecutionContext{
		SQL: "SELECT * FROM users",
	}

	err := m.Handle(context.Background(), execCtx, func(ctx context.Context) error {
		return nil
	})

	assert.NoError(t, err)
}

// TestNewSQLLoggingMiddleware 测试创建 SQL 日志中间件
func TestNewSQLLoggingMiddleware(t *testing.T) {
	m := NewSQLLoggingMiddleware()
	assert.NotNil(t, m)
	assert.Equal(t, "sql_logging", m.Name())
}

// TestSQLLoggingMiddlewareSuccess 测试 SQL 成功日志
func TestSQLLoggingMiddlewareSuccess(t *testing.T) {
	mockLogger := NewMockLogger()
	mockLogger.On("Infof", mock.MatchedBy(func(s string) bool {
		return true
	}), mock.Anything).Return()

	m := NewSQLLoggingMiddleware()
	sqlm := m.(*SQLLoggingMiddleware)
	sqlm.logger = mockLogger

	execCtx := &executor.ExecutionContext{
		SQL: "SELECT * FROM users",
	}

	err := m.Handle(context.Background(), execCtx, func(ctx context.Context) error {
		return nil
	})

	assert.NoError(t, err)
}

// TestSQLLoggingMiddlewareError 测试 SQL 失败日志
func TestSQLLoggingMiddlewareError(t *testing.T) {
	mockLogger := NewMockLogger()
	mockLogger.On("Errorf", mock.MatchedBy(func(s string) bool {
		return true
	}), mock.Anything).Return()

	m := NewSQLLoggingMiddleware()
	sqlm := m.(*SQLLoggingMiddleware)
	sqlm.logger = mockLogger

	expectedErr := errors.New("query failed")
	execCtx := &executor.ExecutionContext{
		SQL: "SELECT * FROM users",
	}

	err := m.Handle(context.Background(), execCtx, func(ctx context.Context) error {
		return expectedErr
	})

	assert.Equal(t, expectedErr, err)
}

// TestSQLLoggingMiddlewareWithNilLogger 测试 logger 为 nil
func TestSQLLoggingMiddlewareWithNilLogger(t *testing.T) {
	m := NewSQLLoggingMiddleware()
	sqlm := m.(*SQLLoggingMiddleware)
	sqlm.logger = nil

	execCtx := &executor.ExecutionContext{
		SQL: "SELECT * FROM users",
	}

	err := m.Handle(context.Background(), execCtx, func(ctx context.Context) error {
		return nil
	})

	assert.NoError(t, err)
}

// TestNewSlowQueryLoggingMiddleware 测试创建慢查询日志中间件
func TestNewSlowQueryLoggingMiddleware(t *testing.T) {
	m := NewSlowQueryLoggingMiddleware(100 * time.Millisecond)
	assert.NotNil(t, m)
	assert.Equal(t, "slow_query_logging", m.Name())

	sqlm := m.(*SlowQueryLoggingMiddleware)
	assert.Equal(t, 100*time.Millisecond, sqlm.GetThreshold())
}

// TestNewSlowQueryLoggingMiddlewareWithDefault 测试默认阈值
func TestNewSlowQueryLoggingMiddlewareWithDefault(t *testing.T) {
	m := NewSlowQueryLoggingMiddleware(0)
	sqlm := m.(*SlowQueryLoggingMiddleware)
	assert.Equal(t, 100*time.Millisecond, sqlm.GetThreshold())

	m = NewSlowQueryLoggingMiddleware(-1)
	sqlm = m.(*SlowQueryLoggingMiddleware)
	assert.Equal(t, 100*time.Millisecond, sqlm.GetThreshold())
}

// TestSlowQueryLoggingMiddlewareSlowQuery 测试检测慢查询
func TestSlowQueryLoggingMiddlewareSlowQuery(t *testing.T) {
	mockLogger := NewMockLogger()
	mockLogger.On("Warnf", mock.MatchedBy(func(s string) bool {
		return true
	}), mock.Anything).Return()

	m := NewSlowQueryLoggingMiddleware(50 * time.Millisecond)
	sqlm := m.(*SlowQueryLoggingMiddleware)
	sqlm.logger = mockLogger

	execCtx := &executor.ExecutionContext{
		SQL: "SELECT * FROM users WHERE id = ?",
	}

	err := m.Handle(context.Background(), execCtx, func(ctx context.Context) error {
		time.Sleep(100 * time.Millisecond)
		return nil
	})

	assert.NoError(t, err)
}

// TestSlowQueryLoggingMiddlewareFastQuery 测试快速查询不被记录
func TestSlowQueryLoggingMiddlewareFastQuery(t *testing.T) {
	mockLogger := NewMockLogger()
	m := NewSlowQueryLoggingMiddleware(100 * time.Millisecond)
	sqlm := m.(*SlowQueryLoggingMiddleware)
	sqlm.logger = mockLogger

	execCtx := &executor.ExecutionContext{
		SQL: "SELECT * FROM users WHERE id = ?",
	}

	err := m.Handle(context.Background(), execCtx, func(ctx context.Context) error {
		time.Sleep(10 * time.Millisecond)
		return nil
	})

	assert.NoError(t, err)
}

// TestSlowQueryLoggingMiddlewareWithError 测试慢查询错误
func TestSlowQueryLoggingMiddlewareWithError(t *testing.T) {
	mockLogger := NewMockLogger()
	mockLogger.On("Warnf", mock.MatchedBy(func(s string) bool {
		return true
	}), mock.Anything).Return()

	m := NewSlowQueryLoggingMiddleware(50 * time.Millisecond)
	sqlm := m.(*SlowQueryLoggingMiddleware)
	sqlm.logger = mockLogger

	expectedErr := errors.New("query failed")
	execCtx := &executor.ExecutionContext{
		SQL: "SELECT * FROM users WHERE id = ?",
	}

	err := m.Handle(context.Background(), execCtx, func(ctx context.Context) error {
		time.Sleep(100 * time.Millisecond)
		return expectedErr
	})

	assert.Equal(t, expectedErr, err)
}

// TestSlowQueryLoggingMiddlewareSetThreshold 测试设置阈值
func TestSlowQueryLoggingMiddlewareSetThreshold(t *testing.T) {
	m := NewSlowQueryLoggingMiddleware(100 * time.Millisecond)
	sqlm := m.(*SlowQueryLoggingMiddleware)

	assert.Equal(t, 100*time.Millisecond, sqlm.GetThreshold())

	sqlm.SetThreshold(200 * time.Millisecond)
	assert.Equal(t, 200*time.Millisecond, sqlm.GetThreshold())
}

// TestSlowQueryLoggingMiddlewareSetThresholdInvalid 测试设置无效阈值
func TestSlowQueryLoggingMiddlewareSetThresholdInvalid(t *testing.T) {
	m := NewSlowQueryLoggingMiddleware(100 * time.Millisecond)
	sqlm := m.(*SlowQueryLoggingMiddleware)

	original := sqlm.GetThreshold()
	sqlm.SetThreshold(0)
	assert.Equal(t, original, sqlm.GetThreshold())

	sqlm.SetThreshold(-100 * time.Millisecond)
	assert.Equal(t, original, sqlm.GetThreshold())
}

// TestSlowQueryLoggingMiddlewareWithNilLogger 测试 logger 为 nil
func TestSlowQueryLoggingMiddlewareWithNilLogger(t *testing.T) {
	m := NewSlowQueryLoggingMiddleware(50 * time.Millisecond)
	sqlm := m.(*SlowQueryLoggingMiddleware)
	sqlm.logger = nil

	execCtx := &executor.ExecutionContext{
		SQL: "SELECT * FROM users WHERE id = ?",
	}

	err := m.Handle(context.Background(), execCtx, func(ctx context.Context) error {
		time.Sleep(100 * time.Millisecond)
		return nil
	})

	assert.NoError(t, err)
}

// TestLoggingMiddlewareIntegration 集成测试：多个日志中间件
func TestLoggingMiddlewareIntegration(t *testing.T) {
	chain := NewChain()

	m1 := NewLoggingMiddleware()
	m2 := NewSQLLoggingMiddleware()
	m3 := NewSlowQueryLoggingMiddleware(50 * time.Millisecond)

	chain.Use(m1, m2, m3)

	execCtx := &executor.ExecutionContext{
		SQL:  "SELECT * FROM users WHERE id = ?",
		Args: []interface{}{123},
	}

	err := chain.Execute(context.Background(), execCtx)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(chain.List()))
}
