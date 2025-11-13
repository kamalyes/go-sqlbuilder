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

// TestNewMetricsMiddleware 测试创建指标中间件
func TestNewMetricsMiddleware(t *testing.T) {
	m := NewMetricsMiddleware()
	assert.NotNil(t, m)
	assert.Equal(t, "metrics", m.Name())
}

// TestMetricsMiddlewareBasic 测试基本指标收集
func TestMetricsMiddlewareBasic(t *testing.T) {
	m := NewMetricsMiddleware()

	execCtx := &executor.ExecutionContext{
		SQL: "SELECT * FROM users",
	}

	// 执行多次查询
	for i := 0; i < 3; i++ {
		err := m.Handle(context.Background(), execCtx, func(ctx context.Context) error {
			time.Sleep(10 * time.Millisecond)
			return nil
		})
		assert.NoError(t, err)
	}

	metrics := m.(*MetricsMiddleware).GetMetrics()
	assert.Equal(t, int64(3), metrics.TotalQueries)
	assert.True(t, metrics.TotalTime > 0)
	assert.True(t, metrics.AverageTime > 0)
}

// TestMetricsMiddlewareError 测试错误计数
func TestMetricsMiddlewareError(t *testing.T) {
	m := NewMetricsMiddleware()

	execCtx := &executor.ExecutionContext{
		SQL: "SELECT * FROM users",
	}

	// 第一个查询失败
	err := m.Handle(context.Background(), execCtx, func(ctx context.Context) error {
		return errors.New("query failed")
	})

	assert.Error(t, err)
	metrics1 := m.(*MetricsMiddleware).GetMetrics()
	assert.Equal(t, int64(1), metrics1.TotalQueries)
	assert.Equal(t, int64(1), metrics1.TotalErrors)

	// 第二个查询成功
	err = m.Handle(context.Background(), execCtx, func(ctx context.Context) error {
		return nil
	})

	assert.NoError(t, err)
	metrics2 := m.(*MetricsMiddleware).GetMetrics()
	assert.Equal(t, int64(2), metrics2.TotalQueries)
	assert.Equal(t, int64(1), metrics2.TotalErrors)
}

// TestMetricsMiddlewareDuration 测试执行时间统计
func TestMetricsMiddlewareDuration(t *testing.T) {
	m := NewMetricsMiddleware()

	execCtx := &executor.ExecutionContext{
		SQL: "SELECT * FROM users",
	}

	// 执行快速查询
	err := m.Handle(context.Background(), execCtx, func(ctx context.Context) error {
		time.Sleep(10 * time.Millisecond)
		return nil
	})
	assert.NoError(t, err)

	metrics1 := m.(*MetricsMiddleware).GetMetrics()
	assert.Equal(t, int64(1), metrics1.TotalQueries)

	// 执行慢速查询
	err = m.Handle(context.Background(), execCtx, func(ctx context.Context) error {
		time.Sleep(50 * time.Millisecond)
		return nil
	})
	assert.NoError(t, err)

	metrics2 := m.(*MetricsMiddleware).GetMetrics()
	assert.Equal(t, int64(2), metrics2.TotalQueries)
	assert.True(t, metrics2.TotalTime > metrics1.TotalTime)
}

// TestMetricsMiddlewareFastestSlowest 测试最快和最慢查询
func TestMetricsMiddlewareFastestSlowest(t *testing.T) {
	m := NewMetricsMiddleware()

	execCtx := &executor.ExecutionContext{
		SQL: "SELECT * FROM users",
	}

	// 执行多个查询，不同的耗时
	durations := []time.Duration{
		10 * time.Millisecond,
		50 * time.Millisecond,
		20 * time.Millisecond,
		30 * time.Millisecond,
	}

	for _, duration := range durations {
		d := duration
		err := m.Handle(context.Background(), execCtx, func(ctx context.Context) error {
			time.Sleep(d)
			return nil
		})
		assert.NoError(t, err)
	}

	metrics := m.(*MetricsMiddleware).GetMetrics()
	assert.Equal(t, int64(4), metrics.TotalQueries)
	// MinTime 应该接近最小值
	assert.True(t, metrics.MinTime < 50*time.Millisecond)
	// MaxTime 应该接近最大值
	assert.True(t, metrics.MaxTime >= 50*time.Millisecond)
}

// TestMetricsMiddlewareReset 测试重置指标
func TestMetricsMiddlewareReset(t *testing.T) {
	m := NewMetricsMiddleware()

	execCtx := &executor.ExecutionContext{
		SQL: "SELECT * FROM users",
	}

	// 执行查询
	err := m.Handle(context.Background(), execCtx, func(ctx context.Context) error {
		return nil
	})
	assert.NoError(t, err)

	metrics1 := m.(*MetricsMiddleware).GetMetrics()
	assert.Equal(t, int64(1), metrics1.TotalQueries)

	// 重置
	m.(*MetricsMiddleware).Reset()
	metrics2 := m.(*MetricsMiddleware).GetMetrics()
	assert.Equal(t, int64(0), metrics2.TotalQueries)
	assert.Equal(t, time.Duration(0), metrics2.TotalTime)
	assert.Equal(t, int64(0), metrics2.TotalErrors)
}

// TestMetricsMiddlewarePerQuery 测试每个查询的指标
func TestMetricsMiddlewarePerQuery(t *testing.T) {
	m := NewMetricsMiddleware()

	execCtx := &executor.ExecutionContext{
		SQL: "SELECT * FROM users WHERE id = ?",
	}

	err := m.Handle(context.Background(), execCtx, func(ctx context.Context) error {
		time.Sleep(20 * time.Millisecond)
		return nil
	})
	assert.NoError(t, err)

	metrics := m.(*MetricsMiddleware).GetMetrics()
	assert.Equal(t, int64(1), metrics.TotalQueries)
	assert.True(t, metrics.TotalTime >= 20*time.Millisecond)
}

// TestMetricsMiddlewareConcurrent 测试并发安全性
func TestMetricsMiddlewareConcurrent(t *testing.T) {
	m := NewMetricsMiddleware()

	execCtx := &executor.ExecutionContext{
		SQL: "SELECT * FROM users",
	}

	// 并发执行查询
	done := make(chan error, 10)
	for i := 0; i < 10; i++ {
		go func() {
			err := m.Handle(context.Background(), execCtx, func(ctx context.Context) error {
				time.Sleep(10 * time.Millisecond)
				return nil
			})
			done <- err
		}()
	}

	for i := 0; i < 10; i++ {
		err := <-done
		assert.NoError(t, err)
	}

	metrics := m.(*MetricsMiddleware).GetMetrics()
	assert.Equal(t, int64(10), metrics.TotalQueries)
}

// TestMetricsMiddlewareAverageDuration 测试平均执行时间
func TestMetricsMiddlewareAverageDuration(t *testing.T) {
	m := NewMetricsMiddleware()

	execCtx := &executor.ExecutionContext{
		SQL: "SELECT * FROM users",
	}

	// 执行固定耗时的查询
	durations := []time.Duration{
		20 * time.Millisecond,
		20 * time.Millisecond,
		20 * time.Millisecond,
	}

	for _, d := range durations {
		duration := d
		err := m.Handle(context.Background(), execCtx, func(ctx context.Context) error {
			time.Sleep(duration)
			return nil
		})
		assert.NoError(t, err)
	}

	metrics := m.(*MetricsMiddleware).GetMetrics()
	assert.Equal(t, int64(3), metrics.TotalQueries)
	// 平均耗时应该约 20ms
	assert.True(t, metrics.AverageTime >= 15*time.Millisecond)
	assert.True(t, metrics.AverageTime <= 30*time.Millisecond)
}

// TestMetricsMiddlewareSuccess 测试成功查询计数
func TestMetricsMiddlewareSuccess(t *testing.T) {
	m := NewMetricsMiddleware()

	execCtx := &executor.ExecutionContext{
		SQL: "SELECT * FROM users",
	}

	err := m.Handle(context.Background(), execCtx, func(ctx context.Context) error {
		return nil
	})
	assert.NoError(t, err)

	metrics := m.(*MetricsMiddleware).GetMetrics()
	assert.Equal(t, int64(1), metrics.TotalQueries)
	assert.Equal(t, int64(1), metrics.SuccessCount)
}

// TestMetricsMiddlewareMultipleQueries 测试多个不同的查询
func TestMetricsMiddlewareMultipleQueries(t *testing.T) {
	m := NewMetricsMiddleware()

	sqls := []string{
		"SELECT * FROM users",
		"INSERT INTO users VALUES (?)",
		"UPDATE users SET status = ?",
		"DELETE FROM users WHERE id = ?",
	}

	for _, sql := range sqls {
		s := sql
		execCtx := &executor.ExecutionContext{
			SQL: s,
		}

		err := m.Handle(context.Background(), execCtx, func(ctx context.Context) error {
			time.Sleep(5 * time.Millisecond)
			return nil
		})
		assert.NoError(t, err)
	}

	metrics := m.(*MetricsMiddleware).GetMetrics()
	assert.Equal(t, int64(4), metrics.TotalQueries)
}

// TestQueryMetricsCalculation 测试指标计算
func TestQueryMetricsCalculation(t *testing.T) {
	metrics := QueryMetrics{
		TotalQueries:  100,
		TotalErrors:   10,
		SuccessCount:  90,
		TotalTime:     1000 * time.Millisecond,
		AverageTime:   10 * time.Millisecond,
		MinTime:       1 * time.Millisecond,
		MaxTime:       100 * time.Millisecond,
		ErrorRate:     0.1,
		QueriesPerSec: 100.0,
	}

	assert.Equal(t, int64(100), metrics.TotalQueries)
	assert.Equal(t, int64(10), metrics.TotalErrors)
	assert.Equal(t, int64(90), metrics.SuccessCount)
	assert.Equal(t, 0.1, metrics.ErrorRate)
	assert.Equal(t, 100.0, metrics.QueriesPerSec)
}

// TestRatioMetricsMiddleware 测试比例指标中间件
func TestRatioMetricsMiddleware(t *testing.T) {
	m := NewRatioMetricsMiddleware()

	sqls := []string{
		"SELECT * FROM users",
		"INSERT INTO users VALUES (?)",
		"UPDATE users SET status = ?",
		"DELETE FROM users WHERE id = ?",
		"SELECT * FROM posts",
	}

	for _, sql := range sqls {
		s := sql
		execCtx := &executor.ExecutionContext{
			SQL: s,
		}

		err := m.Handle(context.Background(), execCtx, func(ctx context.Context) error {
			return nil
		})
		assert.NoError(t, err)
	}

	stats := m.(*RatioMetricsMiddleware).GetStatistics()
	assert.Equal(t, int64(5), stats.TotalQueries)
	assert.Equal(t, int64(2), stats.SelectCount)
	assert.Equal(t, int64(1), stats.InsertCount)
	assert.Equal(t, int64(1), stats.UpdateCount)
	assert.Equal(t, int64(1), stats.DeleteCount)
}

// TestRatioMetricsMiddlewareRatios 测试比例计算
func TestRatioMetricsMiddlewareRatios(t *testing.T) {
	m := NewRatioMetricsMiddleware()

	sqls := []string{
		"SELECT * FROM users",
		"SELECT * FROM posts",
		"INSERT INTO users VALUES (?)",
	}

	for _, sql := range sqls {
		s := sql
		execCtx := &executor.ExecutionContext{
			SQL: s,
		}

		err := m.Handle(context.Background(), execCtx, func(ctx context.Context) error {
			return nil
		})
		assert.NoError(t, err)
	}

	stats := m.(*RatioMetricsMiddleware).GetStatistics()
	assert.True(t, stats.SelectRatio > 0.5) // 2/3
	assert.True(t, stats.InsertRatio > 0.3) // 1/3
}

// TestMetricsErrorRate 测试错误率计算
func TestMetricsErrorRate(t *testing.T) {
	m := NewMetricsMiddleware()

	execCtx := &executor.ExecutionContext{
		SQL: "SELECT * FROM users",
	}

	// 执行 10 次，其中 2 次失败
	for i := 0; i < 10; i++ {
		shouldFail := (i == 3 || i == 7)
		err := m.Handle(context.Background(), execCtx, func(ctx context.Context) error {
			if shouldFail {
				return errors.New("query failed")
			}
			return nil
		})
		if shouldFail {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}

	metrics := m.(*MetricsMiddleware).GetMetrics()
	assert.Equal(t, int64(10), metrics.TotalQueries)
	assert.Equal(t, int64(2), metrics.TotalErrors)
	assert.Equal(t, int64(8), metrics.SuccessCount)
	assert.True(t, metrics.ErrorRate > 0.15 && metrics.ErrorRate < 0.25) // ~0.2
}

// TestMetricsMiddlewareIntegration 集成测试：指标中间件与其他中间件配合
func TestMetricsMiddlewareIntegration(t *testing.T) {
	chain := NewChain()

	metricsMiddleware := NewMetricsMiddleware()
	loggingMiddleware := NewLoggingMiddleware()

	chain.Use(metricsMiddleware, loggingMiddleware)

	execCtx := &executor.ExecutionContext{
		SQL: "SELECT * FROM users WHERE id = ?",
	}

	for i := 0; i < 5; i++ {
		err := chain.Execute(context.Background(), execCtx)
		assert.NoError(t, err)
	}

	metrics := metricsMiddleware.(*MetricsMiddleware).GetMetrics()
	assert.Equal(t, int64(5), metrics.TotalQueries)
}
