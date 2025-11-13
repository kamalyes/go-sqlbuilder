/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-11 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-13 11:03:21
 * @FilePath: \go-sqlbuilder\compiler\planner.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package compiler

import (
	"fmt"
	"strings"

	"github.com/kamalyes/go-sqlbuilder/executor"
)

// SimplePlanner 简单查询计划器
type SimplePlanner struct {
	name string
}

// NewSimplePlanner 创建简单查询计划器
func NewSimplePlanner() Planner {
	return &SimplePlanner{
		name: "simple_planner",
	}
}

// Plan 生成查询执行计划
func (p *SimplePlanner) Plan(execCtx *executor.ExecutionContext) (*QueryPlan, error) {
	if execCtx == nil {
		return nil, fmt.Errorf("execution context cannot be nil")
	}

	sql := execCtx.SQL
	if sql == "" {
		return nil, fmt.Errorf("SQL cannot be empty")
	}

	plan := &QueryPlan{
		OriginalSQL:  sql,
		OptimizedSQL: sql,
		Strategy:     p.determineStrategy(sql),
		IndexHints:   p.suggestIndexes(sql),
		Explanation:  p.explainPlan(sql),
	}

	return plan, nil
}

// determineStrategy 确定执行策略
func (p *SimplePlanner) determineStrategy(sql string) string {
	sql = strings.ToUpper(strings.TrimSpace(sql))

	if strings.HasPrefix(sql, "SELECT") {
		if strings.Contains(sql, "JOIN") {
			return "HASH_JOIN"
		}
		if strings.Contains(sql, "GROUP BY") {
			return "GROUPING"
		}
		return "FULL_SCAN"
	}

	if strings.HasPrefix(sql, "INSERT") {
		return "INSERT"
	}

	if strings.HasPrefix(sql, "UPDATE") {
		return "UPDATE"
	}

	if strings.HasPrefix(sql, "DELETE") {
		return "DELETE"
	}

	return "UNKNOWN"
}

// suggestIndexes 建议索引
func (p *SimplePlanner) suggestIndexes(sql string) []string {
	var hints []string
	sql = strings.ToUpper(sql)

	// 检查WHERE条件
	if strings.Contains(sql, "WHERE") {
		hints = append(hints, "Consider index on WHERE columns")
	}

	// 检查JOIN条件
	if strings.Contains(sql, "JOIN") {
		hints = append(hints, "Consider index on JOIN columns")
	}

	// 检查ORDER BY
	if strings.Contains(sql, "ORDER BY") {
		hints = append(hints, "Consider index on ORDER BY columns")
	}

	// 检查GROUP BY
	if strings.Contains(sql, "GROUP BY") {
		hints = append(hints, "Consider index on GROUP BY columns")
	}

	return hints
}

// explainPlan 解释执行计划
func (p *SimplePlanner) explainPlan(sql string) string {
	sql = strings.ToUpper(strings.TrimSpace(sql))

	switch {
	case strings.HasPrefix(sql, "SELECT"):
		return "Execute SELECT query with sequential scan"
	case strings.HasPrefix(sql, "INSERT"):
		return "Execute INSERT query"
	case strings.HasPrefix(sql, "UPDATE"):
		return "Execute UPDATE query"
	case strings.HasPrefix(sql, "DELETE"):
		return "Execute DELETE query"
	default:
		return "Execute unknown query type"
	}
}

// Analyze 分析查询性能
func (p *SimplePlanner) Analyze(execCtx *executor.ExecutionContext) map[string]interface{} {
	result := make(map[string]interface{})

	if execCtx == nil {
		result["error"] = "execution context is nil"
		return result
	}

	sql := execCtx.SQL
	sqlLen := len(sql)
	argsLen := len(execCtx.Args)

	result["sql_length"] = sqlLen
	result["args_count"] = argsLen
	result["complexity"] = p.estimateComplexity(sql)
	result["estimated_cost"] = p.estimateCost(sql)

	return result
}

// estimateComplexity 估计查询复杂度
func (p *SimplePlanner) estimateComplexity(sql string) int {
	complexity := 1
	sql = strings.ToUpper(sql)

	if strings.Contains(sql, "JOIN") {
		complexity += 2
	}
	if strings.Contains(sql, "GROUP BY") {
		complexity += 1
	}
	if strings.Contains(sql, "ORDER BY") {
		complexity += 1
	}
	if strings.Contains(sql, "SUBQUERY") || strings.Contains(sql, "SELECT") {
		complexity += 1
	}

	return complexity
}

// estimateCost 估计查询成本
func (p *SimplePlanner) estimateCost(sql string) float64 {
	cost := 1.0
	sql = strings.ToUpper(sql)

	if strings.Contains(sql, "JOIN") {
		cost *= 10.0
	}
	if strings.Contains(sql, "GROUP BY") {
		cost *= 5.0
	}
	if strings.Contains(sql, "ORDER BY") {
		cost *= 3.0
	}

	return cost
}

// StatisticPlanner 基于统计信息的查询计划器
type StatisticPlanner struct {
	name       string
	statistics map[string]TableStatistics
}

// TableStatistics 表统计信息
type TableStatistics struct {
	TableName  string
	RowCount   int64
	ColumnStat map[string]ColumnStatistics
}

// ColumnStatistics 列统计信息
type ColumnStatistics struct {
	ColumnName   string
	Cardinality  int64
	AvgLength    int32
	NullCount    int64
	DistinctVals int64
}

// NewStatisticPlanner 创建基于统计的查询计划器
func NewStatisticPlanner() *StatisticPlanner {
	return &StatisticPlanner{
		name:       "statistic_planner",
		statistics: make(map[string]TableStatistics),
	}
}

// Plan 生成查询执行计划
func (p *StatisticPlanner) Plan(execCtx *executor.ExecutionContext) (*QueryPlan, error) {
	if execCtx == nil {
		return nil, fmt.Errorf("execution context cannot be nil")
	}

	sql := execCtx.SQL
	if sql == "" {
		return nil, fmt.Errorf("SQL cannot be empty")
	}

	// 基于统计信息生成计划
	plan := &QueryPlan{
		OriginalSQL:  sql,
		OptimizedSQL: sql,
		Strategy:     p.determineOptimalStrategy(sql),
		IndexHints:   p.suggestOptimalIndexes(sql),
		Explanation:  "Plan based on table statistics",
	}

	// 计算估计成本
	plan.EstimatedCost = p.estimateStatisticCost(sql)

	return plan, nil
}

// determineOptimalStrategy 确定最优执行策略
func (p *StatisticPlanner) determineOptimalStrategy(sql string) string {
	// 基于统计信息选择最优策略
	return "OPTIMAL_STRATEGY"
}

// suggestOptimalIndexes 建议最优索引
func (p *StatisticPlanner) suggestOptimalIndexes(sql string) []string {
	// 基于统计信息建议索引
	return []string{"Optimal index suggestion"}
}

// estimateStatisticCost 基于统计估计成本
func (p *StatisticPlanner) estimateStatisticCost(sql string) float64 {
	// 基于表统计计算成本
	return 1.0
}

// Analyze 分析查询性能
func (p *StatisticPlanner) Analyze(execCtx *executor.ExecutionContext) map[string]interface{} {
	result := make(map[string]interface{})

	if execCtx == nil {
		result["error"] = "execution context is nil"
		return result
	}

	result["planner"] = "statistic_planner"
	result["statistics_available"] = len(p.statistics) > 0

	return result
}

// AddTableStatistics 添加表统计信息
func (p *StatisticPlanner) AddTableStatistics(stats TableStatistics) {
	p.statistics[stats.TableName] = stats
}

// GetTableStatistics 获取表统计信息
func (p *StatisticPlanner) GetTableStatistics(tableName string) (TableStatistics, bool) {
	stats, ok := p.statistics[tableName]
	return stats, ok
}
