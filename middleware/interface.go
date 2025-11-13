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

	"github.com/kamalyes/go-sqlbuilder/executor"
)

// Middleware 定义中间件接口
// Middleware 在执行查询前后进行处理
type Middleware interface {
	// Name 返回中间件的名称
	Name() string

	// Handle 处理请求
	// 返回下一个中间件处理的 Next 函数
	Handle(ctx context.Context, execCtx *executor.ExecutionContext, next Next) error
}

// Next 定义下一个中间件处理函数
type Next func(ctx context.Context) error

// ExecutionChain 定义中间件链接口
type ExecutionChain interface {
	// Use 添加中间件
	Use(middleware ...Middleware) ExecutionChain

	// Execute 执行中间件链
	// 依次执行所有中间件的 Handle 方法
	Execute(ctx context.Context, execCtx *executor.ExecutionContext) error

	// Remove 移除指定名称的中间件
	Remove(name string) ExecutionChain

	// Clear 清除所有中间件
	Clear() ExecutionChain

	// List 获取所有中间件名称列表
	List() []string

	// HasMiddleware 检查是否存在指定名称的中间件
	HasMiddleware(name string) bool
}

// Builder 定义中间件链构建器接口
type Builder interface {
	// WithLogging 添加日志中间件
	WithLogging() Builder

	// WithMetrics 添加指标中间件
	WithMetrics() Builder

	// WithRetry 添加重试中间件
	WithRetry(maxAttempts int) Builder

	// WithTimeout 添加超时中间件
	WithTimeout() Builder

	// WithValidation 添加验证中间件
	WithValidation() Builder

	// WithCustom 添加自定义中间件
	WithCustom(middleware Middleware) Builder

	// Build 构建中间件链
	Build() ExecutionChain
}
