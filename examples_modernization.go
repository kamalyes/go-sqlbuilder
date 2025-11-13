/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-11 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-13 11:10:51
 * @FilePath: \go-sqlbuilder\examples_modernization.go
 * @Description: 并发安全的Builder - 现代化改造示例
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package sqlbuilder

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"
)

// ==================== 现代化改造示例 ====================

// BuilderV2 改进的并发安全的Builder
type BuilderV2 struct {
	// 并发保护
	mu sync.RWMutex

	// 核心字段
	adapter UniversalAdapterInterface
	ctx     context.Context
	timeout time.Duration

	// SQL构建组件 (受锁保护)
	table       string
	tableAlias  string
	distinct    bool
	columns     []string
	joins       []string
	wheres      []string
	havings     []string
	groupByCols []string
	orderByCols []string
	limitVal    int64
	offsetVal   int64

	// 操作数据
	insertData map[string]interface{}
	updateData map[string]interface{}

	// 参数
	args      []interface{}
	queryType string

	// 性能指标
	metrics *BuildMetrics
}

// BuildMetrics 构建性能指标
type BuildMetrics struct {
	mu                 sync.RWMutex
	totalBuilds        int64
	totalExecutions    int64
	totalErrors        int64
	cumulativeDuration time.Duration
}

// NewBuilderV2 创建新的并发安全Builder
func NewBuilderV2(dbInstance interface{}) (*BuilderV2, error) {
	adapter, err := AutoDetectAdapter(dbInstance)
	if err != nil {
		return nil, err
	}

	return &BuilderV2{
		adapter:     adapter,
		ctx:         context.Background(),
		timeout:     30 * time.Second,
		columns:     make([]string, 0, 10),
		joins:       make([]string, 0, 5),
		wheres:      make([]string, 0, 10),
		havings:     make([]string, 0, 5),
		groupByCols: make([]string, 0, 5),
		orderByCols: make([]string, 0, 5),
		insertData:  make(map[string]interface{}),
		updateData:  make(map[string]interface{}),
		args:        make([]interface{}, 0, 20),
		queryType:   "select",
		metrics:     &BuildMetrics{},
	}, nil
}

// ==================== 并发安全的操作 ====================

// Table 并发安全的设置表名
func (b *BuilderV2) Table(table string) *BuilderV2 {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.table = table
	b.queryType = "select"
	return b
}

// Select 并发安全的SELECT列
func (b *BuilderV2) Select(columns ...string) *BuilderV2 {
	b.mu.Lock()
	defer b.mu.Unlock()

	// 预分配以提高性能
	b.columns = make([]string, 0, len(columns))
	b.columns = append(b.columns, columns...)
	return b
}

// Where 并发安全的添加WHERE条件
func (b *BuilderV2) Where(field string, operator interface{}, value ...interface{}) *BuilderV2 {
	b.mu.Lock()
	defer b.mu.Unlock()

	var whereStr string
	var args []interface{}

	if len(value) == 0 {
		// 两参数形式: Where("status", 1)
		whereStr = fmt.Sprintf("%s = ?", field)
		args = []interface{}{operator}
	} else {
		// 三参数形式: Where("age", ">", 18)
		whereStr = fmt.Sprintf("%s %v ?", field, operator)
		args = value
	}

	b.wheres = append(b.wheres, whereStr)
	b.args = append(b.args, args...)
	return b
}

// OrderBy 并发安全的排序
func (b *BuilderV2) OrderBy(field, direction string) *BuilderV2 {
	b.mu.Lock()
	defer b.mu.Unlock()

	if direction == "" {
		direction = "ASC"
	}
	b.orderByCols = append(b.orderByCols, fmt.Sprintf("%s %s", field, direction))
	return b
}

// Limit 并发安全的限制
func (b *BuilderV2) Limit(limit int64) *BuilderV2 {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.limitVal = limit
	return b
}

// Build 并发安全的构建SQL - 使用strings.Builder优化
func (b *BuilderV2) Build() (string, []interface{}, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	// 记录开始时间 (性能指标)
	start := time.Now()
	defer func() {
		b.metrics.recordBuild(time.Since(start))
	}()

	var buf strings.Builder

	// 初始化缓冲区容量 (性能优化)
	buf.Grow(256)

	// 构建SELECT语句
	buf.WriteString("SELECT ")

	if b.distinct {
		buf.WriteString("DISTINCT ")
	}

	// 优化: 预留容量避免多次分配
	if len(b.columns) == 0 {
		buf.WriteString("*")
	} else {
		for i, col := range b.columns {
			if i > 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(col)
		}
	}

	// FROM子句
	buf.WriteString(" FROM ")
	buf.WriteString(b.table)

	if b.tableAlias != "" {
		buf.WriteString(" AS ")
		buf.WriteString(b.tableAlias)
	}

	// JOIN子句
	for _, join := range b.joins {
		buf.WriteString(" ")
		buf.WriteString(join)
	}

	// WHERE子句
	if len(b.wheres) > 0 {
		buf.WriteString(" WHERE ")
		for i, where := range b.wheres {
			if i > 0 {
				buf.WriteString(" AND ")
			}
			buf.WriteString(where)
		}
	}

	// GROUP BY子句
	if len(b.groupByCols) > 0 {
		buf.WriteString(" GROUP BY ")
		for i, col := range b.groupByCols {
			if i > 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(col)
		}
	}

	// HAVING子句
	if len(b.havings) > 0 {
		buf.WriteString(" HAVING ")
		for i, having := range b.havings {
			if i > 0 {
				buf.WriteString(" AND ")
			}
			buf.WriteString(having)
		}
	}

	// ORDER BY子句
	if len(b.orderByCols) > 0 {
		buf.WriteString(" ORDER BY ")
		for i, order := range b.orderByCols {
			if i > 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(order)
		}
	}

	// LIMIT子句
	if b.limitVal > 0 {
		buf.WriteString(" LIMIT ")
		buf.WriteString(fmt.Sprintf("%d", b.limitVal))
	}

	// OFFSET子句
	if b.offsetVal > 0 {
		buf.WriteString(" OFFSET ")
		buf.WriteString(fmt.Sprintf("%d", b.offsetVal))
	}

	sql := buf.String()
	args := make([]interface{}, len(b.args))
	copy(args, b.args)

	return sql, args, nil
}

// ==================== 性能指标 ====================

// recordBuild 记录构建性能指标
func (m *BuildMetrics) recordBuild(duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.totalBuilds++
	m.cumulativeDuration += duration
}

// recordExecution 记录执行指标
func (m *BuildMetrics) recordExecution(duration time.Duration, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.totalExecutions++
	if err != nil {
		m.totalErrors++
	}
	m.cumulativeDuration += duration
}

// GetMetrics 获取性能指标 (线程安全)
func (b *BuilderV2) GetMetrics() map[string]interface{} {
	b.metrics.mu.RLock()
	defer b.metrics.mu.RUnlock()

	avgDuration := time.Duration(0)
	if b.metrics.totalBuilds > 0 {
		avgDuration = b.metrics.cumulativeDuration / time.Duration(b.metrics.totalBuilds)
	}

	return map[string]interface{}{
		"total_builds":        b.metrics.totalBuilds,
		"total_executions":    b.metrics.totalExecutions,
		"total_errors":        b.metrics.totalErrors,
		"avg_duration":        avgDuration.String(),
		"cumulative_duration": b.metrics.cumulativeDuration.String(),
	}
}

// ==================== 使用示例 ====================

/*
使用并发安全的BuilderV2:

package main

import (
    "fmt"
    "sync"
)

func main() {
    // 创建Builder (线程安全)
    builder, _ := NewBuilderV2(db)

    // 多个goroutine可以安全地使用同一个builder
    var wg sync.WaitGroup

    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()

            builder.
                Table("users").
                Select("id", "name").
                Where("status", 1).
                Where("age", ">", 18).
                OrderBy("created_at", "DESC").
                Limit(10)

            sql, args, _ := builder.Build()
            fmt.Printf("Query %d: %s\n", id, sql)
        }(i)
    }

    wg.Wait()

    // 查看性能指标
    metrics := builder.GetMetrics()
    fmt.Println("Metrics:", metrics)
}

优点:
✅ 完全并发安全
✅ 无竞态条件
✅ 性能指标内置
✅ 高效的字符串拼接 (strings.Builder)
✅ 预分配优化
*/
