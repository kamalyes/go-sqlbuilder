/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-11 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-13 11:03:06
 * @FilePath: \go-sqlbuilder\compiler\interface.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package compiler

import (
	"github.com/kamalyes/go-sqlbuilder/executor"
)

// SQLCompiler SQL编译器接口
// 将ExecutionContext编译成特定方言的SQL语句
type SQLCompiler interface {
	// Compile 编译执行上下文为SQL语句
	Compile(execCtx *executor.ExecutionContext) (string, []interface{}, error)

	// GetDialect 获取当前方言
	GetDialect() string

	// SetDialect 设置方言
	SetDialect(dialect string) error
}

// Optimizer SQL优化器接口
// 优化SQL查询性能
type Optimizer interface {
	// Optimize 优化SQL执行上下文
	Optimize(execCtx *executor.ExecutionContext) (*executor.ExecutionContext, error)

	// GetName 获取优化器名称
	GetName() string
}

// QueryPlan 查询计划
type QueryPlan struct {
	// 原始SQL
	OriginalSQL string

	// 优化后的SQL
	OptimizedSQL string

	// 执行策略
	Strategy string

	// 估计成本
	EstimatedCost float64

	// 索引建议
	IndexHints []string

	// 优化说明
	Explanation string
}

// Planner 查询计划器接口
type Planner interface {
	// Plan 生成查询执行计划
	Plan(execCtx *executor.ExecutionContext) (*QueryPlan, error)

	// Analyze 分析查询性能
	Analyze(execCtx *executor.ExecutionContext) map[string]interface{}
}

// CompilerOptions 编译器选项
type CompilerOptions struct {
	// 方言
	Dialect string

	// 是否启用优化
	EnableOptimization bool

	// 是否启用查询计划
	EnablePlanning bool

	// 最大查询复杂度
	MaxQueryComplexity int

	// 是否启用参数化查询
	EnableParameterization bool

	// 是否启用查询缓存
	EnableQueryCache bool
}

// CompilerConfig 编译器配置
type CompilerConfig struct {
	// 选项
	Options CompilerOptions

	// 优化器列表
	Optimizers []Optimizer

	// 查询计划器
	Planner Planner

	// 自定义SQL转换函数
	CustomTransformers map[string]func(string) string
}
