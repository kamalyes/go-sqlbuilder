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
	"sync"
	"sync/atomic"
	"time"

	"github.com/kamalyes/go-sqlbuilder/executor"
)

// MetricsMiddleware 指标中间件
// 收集查询的统计信息，如执行次数、总耗时、错误次数等
type MetricsMiddleware struct {
	name         string
	totalQueries int64
	totalErrors  int64
	totalTime    int64 // 纳秒
	minTime      int64 // 纳秒
	maxTime      int64 // 纳秒
	mu           sync.RWMutex
	initialized  bool
}

// NewMetricsMiddleware 创建指标中间件
func NewMetricsMiddleware() Middleware {
	return &MetricsMiddleware{
		name:        "metrics",
		minTime:     int64(^uint64(0) >> 1), // max int64
		maxTime:     0,
		initialized: true,
	}
}

// Name 返回中间件名称
func (m *MetricsMiddleware) Name() string {
	return m.name
}

// Handle 处理请求
func (m *MetricsMiddleware) Handle(ctx context.Context, execCtx *executor.ExecutionContext, next Next) error {
	startTime := time.Now()

	// 执行下一个中间件
	err := next(ctx)

	// 更新指标
	duration := time.Since(startTime)
	durationNano := duration.Nanoseconds()

	atomic.AddInt64(&m.totalQueries, 1)
	atomic.AddInt64(&m.totalTime, durationNano)

	if err != nil {
		atomic.AddInt64(&m.totalErrors, 1)
	}

	// 更新最小/最大时间
	m.updateMinMaxTime(durationNano)

	return err
}

// updateMinMaxTime 更新最小和最大执行时间
func (m *MetricsMiddleware) updateMinMaxTime(duration int64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if duration < m.minTime {
		m.minTime = duration
	}
	if duration > m.maxTime {
		m.maxTime = duration
	}
}

// GetMetrics 获取统计指标
func (m *MetricsMiddleware) GetMetrics() QueryMetrics {
	m.mu.RLock()
	defer m.mu.RUnlock()

	total := atomic.LoadInt64(&m.totalQueries)
	totalTime := atomic.LoadInt64(&m.totalTime)
	totalErrors := atomic.LoadInt64(&m.totalErrors)

	avgTime := int64(0)
	if total > 0 {
		avgTime = totalTime / total
	}

	return QueryMetrics{
		TotalQueries:  total,
		TotalErrors:   totalErrors,
		SuccessCount:  total - totalErrors,
		TotalTime:     time.Duration(totalTime),
		AverageTime:   time.Duration(avgTime),
		MinTime:       time.Duration(m.minTime),
		MaxTime:       time.Duration(m.maxTime),
		ErrorRate:     float64(totalErrors) / float64(total),
		QueriesPerSec: calculateQPS(total, totalTime),
	}
}

// QueryMetrics 查询指标
type QueryMetrics struct {
	TotalQueries  int64
	TotalErrors   int64
	SuccessCount  int64
	TotalTime     time.Duration
	AverageTime   time.Duration
	MinTime       time.Duration
	MaxTime       time.Duration
	ErrorRate     float64
	QueriesPerSec float64
}

// Reset 重置所有指标
func (m *MetricsMiddleware) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	atomic.StoreInt64(&m.totalQueries, 0)
	atomic.StoreInt64(&m.totalErrors, 0)
	atomic.StoreInt64(&m.totalTime, 0)
	m.minTime = int64(^uint64(0) >> 1)
	m.maxTime = 0
}

// calculateQPS 计算每秒查询数
func calculateQPS(totalQueries int64, totalTimeNano int64) float64 {
	if totalTimeNano <= 0 {
		return 0
	}
	seconds := float64(totalTimeNano) / 1e9
	return float64(totalQueries) / seconds
}

// RatioMetricsMiddleware 比例指标中间件
// 统计不同查询类型的比例
type RatioMetricsMiddleware struct {
	name         string
	selectCount  int64
	insertCount  int64
	updateCount  int64
	deleteCount  int64
	otherCount   int64
	totalQueries int64
	mu           sync.RWMutex
}

// NewRatioMetricsMiddleware 创建比例指标中间件
func NewRatioMetricsMiddleware() Middleware {
	return &RatioMetricsMiddleware{
		name: "ratio_metrics",
	}
}

// Name 返回中间件名称
func (m *RatioMetricsMiddleware) Name() string {
	return m.name
}

// Handle 处理请求
func (m *RatioMetricsMiddleware) Handle(ctx context.Context, execCtx *executor.ExecutionContext, next Next) error {
	// 识别查询类型
	queryType := m.identifyQueryType(execCtx.SQL)

	err := next(ctx)

	// 更新计数器
	atomic.AddInt64(&m.totalQueries, 1)
	switch queryType {
	case "SELECT":
		atomic.AddInt64(&m.selectCount, 1)
	case "INSERT":
		atomic.AddInt64(&m.insertCount, 1)
	case "UPDATE":
		atomic.AddInt64(&m.updateCount, 1)
	case "DELETE":
		atomic.AddInt64(&m.deleteCount, 1)
	default:
		atomic.AddInt64(&m.otherCount, 1)
	}

	return err
}

// identifyQueryType 识别查询类型
func (m *RatioMetricsMiddleware) identifyQueryType(sql string) string {
	if len(sql) == 0 {
		return "UNKNOWN"
	}

	// 查找第一个字母
	for i, c := range sql {
		if (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') {
			// 读取接下来的 6 个字符
			end := i + 6
			if end > len(sql) {
				end = len(sql)
			}
			keyword := sql[i:end]

			switch {
			case len(keyword) >= 6 && keyword[:6] == "SELECT":
				return "SELECT"
			case len(keyword) >= 6 && keyword[:6] == "INSERT":
				return "INSERT"
			case len(keyword) >= 6 && keyword[:6] == "UPDATE":
				return "UPDATE"
			case len(keyword) >= 6 && keyword[:6] == "DELETE":
				return "DELETE"
			}
			break
		}
	}

	return "OTHER"
}

// GetStatistics 获取统计信息
func (m *RatioMetricsMiddleware) GetStatistics() QueryStatistics {
	m.mu.RLock()
	defer m.mu.RUnlock()

	total := atomic.LoadInt64(&m.totalQueries)
	selectCount := atomic.LoadInt64(&m.selectCount)
	insertCount := atomic.LoadInt64(&m.insertCount)
	updateCount := atomic.LoadInt64(&m.updateCount)
	deleteCount := atomic.LoadInt64(&m.deleteCount)
	otherCount := atomic.LoadInt64(&m.otherCount)

	var selectRatio, insertRatio, updateRatio, deleteRatio, otherRatio float64
	if total > 0 {
		selectRatio = float64(selectCount) / float64(total)
		insertRatio = float64(insertCount) / float64(total)
		updateRatio = float64(updateCount) / float64(total)
		deleteRatio = float64(deleteCount) / float64(total)
		otherRatio = float64(otherCount) / float64(total)
	}

	return QueryStatistics{
		TotalQueries: total,
		SelectCount:  selectCount,
		InsertCount:  insertCount,
		UpdateCount:  updateCount,
		DeleteCount:  deleteCount,
		OtherCount:   otherCount,
		SelectRatio:  selectRatio,
		InsertRatio:  insertRatio,
		UpdateRatio:  updateRatio,
		DeleteRatio:  deleteRatio,
		OtherRatio:   otherRatio,
	}
}

// QueryStatistics 查询统计信息
type QueryStatistics struct {
	TotalQueries int64
	SelectCount  int64
	InsertCount  int64
	UpdateCount  int64
	DeleteCount  int64
	OtherCount   int64
	SelectRatio  float64
	InsertRatio  float64
	UpdateRatio  float64
	DeleteRatio  float64
	OtherRatio   float64
}

// Reset 重置所有统计
func (m *RatioMetricsMiddleware) Reset() {
	atomic.StoreInt64(&m.selectCount, 0)
	atomic.StoreInt64(&m.insertCount, 0)
	atomic.StoreInt64(&m.updateCount, 0)
	atomic.StoreInt64(&m.deleteCount, 0)
	atomic.StoreInt64(&m.otherCount, 0)
	atomic.StoreInt64(&m.totalQueries, 0)
}
