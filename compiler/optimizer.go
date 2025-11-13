/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-11 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-13 11:03:14
 * @FilePath: \go-sqlbuilder\compiler\optimizer.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package compiler

import (
	"strings"

	"github.com/kamalyes/go-sqlbuilder/executor"
)

// SelectOptimizer SELECT语句优化器
type SelectOptimizer struct {
	name string
}

// NewSelectOptimizer 创建SELECT优化器
func NewSelectOptimizer() Optimizer {
	return &SelectOptimizer{
		name: "select_optimizer",
	}
}

// Optimize 优化SELECT语句
func (o *SelectOptimizer) Optimize(execCtx *executor.ExecutionContext) (*executor.ExecutionContext, error) {
	if execCtx == nil {
		return execCtx, nil
	}

	// 检查是否是SELECT语句
	sql := strings.ToUpper(strings.TrimSpace(execCtx.SQL))
	if !strings.HasPrefix(sql, "SELECT") {
		return execCtx, nil
	}

	// 优化逻辑
	// 1. 消除不必要的DISTINCT
	// 2. 列出所需的列而不是 *
	// 3. 推送谓词下降

	return execCtx, nil
}

// GetName 获取优化器名称
func (o *SelectOptimizer) GetName() string {
	return o.name
}

// JoinOptimizer JOIN优化器
type JoinOptimizer struct {
	name string
}

// NewJoinOptimizer 创建JOIN优化器
func NewJoinOptimizer() Optimizer {
	return &JoinOptimizer{
		name: "join_optimizer",
	}
}

// Optimize 优化JOIN操作
func (o *JoinOptimizer) Optimize(execCtx *executor.ExecutionContext) (*executor.ExecutionContext, error) {
	if execCtx == nil {
		return execCtx, nil
	}

	// 检查是否包含JOIN
	sql := strings.ToUpper(execCtx.SQL)
	if !strings.Contains(sql, "JOIN") {
		return execCtx, nil
	}

	// JOIN优化逻辑
	// 1. 调整JOIN顺序
	// 2. 选择最优JOIN算法
	// 3. 推送谓词到JOIN之前

	return execCtx, nil
}

// GetName 获取优化器名称
func (o *JoinOptimizer) GetName() string {
	return o.name
}

// WhereOptimizer WHERE子句优化器
type WhereOptimizer struct {
	name string
}

// NewWhereOptimizer 创建WHERE优化器
func NewWhereOptimizer() Optimizer {
	return &WhereOptimizer{
		name: "where_optimizer",
	}
}

// Optimize 优化WHERE子句
func (o *WhereOptimizer) Optimize(execCtx *executor.ExecutionContext) (*executor.ExecutionContext, error) {
	if execCtx == nil {
		return execCtx, nil
	}

	// 检查是否包含WHERE
	sql := strings.ToUpper(execCtx.SQL)
	if !strings.Contains(sql, "WHERE") {
		return execCtx, nil
	}

	// WHERE优化逻辑
	// 1. 消除冗余条件
	// 2. 简化布尔表达式
	// 3. 优化索引使用

	return execCtx, nil
}

// GetName 获取优化器名称
func (o *WhereOptimizer) GetName() string {
	return o.name
}

// GroupByOptimizer GROUP BY优化器
type GroupByOptimizer struct {
	name string
}

// NewGroupByOptimizer 创建GROUP BY优化器
func NewGroupByOptimizer() Optimizer {
	return &GroupByOptimizer{
		name: "groupby_optimizer",
	}
}

// Optimize 优化GROUP BY操作
func (o *GroupByOptimizer) Optimize(execCtx *executor.ExecutionContext) (*executor.ExecutionContext, error) {
	if execCtx == nil {
		return execCtx, nil
	}

	// 检查是否包含GROUP BY
	sql := strings.ToUpper(execCtx.SQL)
	if !strings.Contains(sql, "GROUP BY") {
		return execCtx, nil
	}

	// GROUP BY优化逻辑
	// 1. 调整分组列顺序
	// 2. 使用索引优化分组
	// 3. 推送聚合函数

	return execCtx, nil
}

// GetName 获取优化器名称
func (o *GroupByOptimizer) GetName() string {
	return o.name
}

// OrderByOptimizer ORDER BY优化器
type OrderByOptimizer struct {
	name string
}

// NewOrderByOptimizer 创建ORDER BY优化器
func NewOrderByOptimizer() Optimizer {
	return &OrderByOptimizer{
		name: "orderby_optimizer",
	}
}

// Optimize 优化ORDER BY操作
func (o *OrderByOptimizer) Optimize(execCtx *executor.ExecutionContext) (*executor.ExecutionContext, error) {
	if execCtx == nil {
		return execCtx, nil
	}

	// 检查是否包含ORDER BY
	sql := strings.ToUpper(execCtx.SQL)
	if !strings.Contains(sql, "ORDER BY") {
		return execCtx, nil
	}

	// ORDER BY优化逻辑
	// 1. 使用索引排序
	// 2. 避免额外排序
	// 3. 推送排序下降

	return execCtx, nil
}

// GetName 获取优化器名称
func (o *OrderByOptimizer) GetName() string {
	return o.name
}

// CompositeOptimizer 组合优化器
type CompositeOptimizer struct {
	name       string
	optimizers []Optimizer
}

// NewCompositeOptimizer 创建组合优化器
func NewCompositeOptimizer(optimizers ...Optimizer) Optimizer {
	return &CompositeOptimizer{
		name:       "composite_optimizer",
		optimizers: optimizers,
	}
}

// Optimize 应用所有优化器
func (o *CompositeOptimizer) Optimize(execCtx *executor.ExecutionContext) (*executor.ExecutionContext, error) {
	var err error
	for _, optimizer := range o.optimizers {
		execCtx, err = optimizer.Optimize(execCtx)
		if err != nil {
			return execCtx, err
		}
	}
	return execCtx, nil
}

// GetName 获取优化器名称
func (o *CompositeOptimizer) GetName() string {
	return o.name
}

// AddOptimizer 添加优化器
func (o *CompositeOptimizer) AddOptimizer(optimizer Optimizer) {
	o.optimizers = append(o.optimizers, optimizer)
}
