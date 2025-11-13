/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-11 21:13:15
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-13 11:09:38
 * @FilePath: \go-sqlbuilder\builder_enhancer.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package sqlbuilder

import (
	"github.com/kamalyes/go-sqlbuilder/compiler"
	"github.com/kamalyes/go-sqlbuilder/executor"
	"github.com/kamalyes/go-sqlbuilder/middleware"
)

// BuilderEnhancer Builder增强器 - 为现有Builder提供新功能映射
type BuilderEnhancer struct {
	builder *Builder

	// Phase 4: Middleware
	chain middleware.ExecutionChain

	// Phase 5: Compiler
	sqlCompiler compiler.SQLCompiler

	// Phase 3: Executor
	queryExecutor executor.Executor

	// Phase 5: Planner
	planner compiler.Planner
}

// NewBuilderEnhancer 为Builder增加新功能
func NewBuilderEnhancer(b *Builder) *BuilderEnhancer {
	if b == nil {
		return nil
	}

	dialect := "mysql" // 默认方言
	if b.adapter != nil {
		dialect = b.adapter.GetDialect()
	}

	return &BuilderEnhancer{
		builder:     b,
		chain:       middleware.NewChain(),
		sqlCompiler: compiler.NewDefaultCompiler(dialect),
		planner:     compiler.NewSimplePlanner(),
	}
}

// GetMiddlewareChain 获取middleware链
func (be *BuilderEnhancer) GetMiddlewareChain() middleware.ExecutionChain {
	return be.chain
}

// AddMiddleware 添加中间件
func (be *BuilderEnhancer) AddMiddleware(middlewares ...middleware.Middleware) *BuilderEnhancer {
	be.chain.Use(middlewares...)
	return be
}

// GetCompiler 获取SQL编译器
func (be *BuilderEnhancer) GetCompiler() compiler.SQLCompiler {
	return be.sqlCompiler
}

// GetExecutor 获取查询执行器
func (be *BuilderEnhancer) GetExecutor() executor.Executor {
	return be.queryExecutor
}

// GetPlanner 获取查询计划器
func (be *BuilderEnhancer) GetPlanner() compiler.Planner {
	return be.planner
}

// GetBuilder 获取原始Builder
func (be *BuilderEnhancer) GetBuilder() *Builder {
	return be.builder
}
