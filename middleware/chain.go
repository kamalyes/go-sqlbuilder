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
	"sync"

	"github.com/kamalyes/go-sqlbuilder/executor"
)

// chain 实现 ExecutionChain 接口
type chain struct {
	middlewares []Middleware
	mu          sync.RWMutex
}

// NewChain 创建新的中间件链
func NewChain() ExecutionChain {
	return &chain{
		middlewares: make([]Middleware, 0),
	}
}

// Use 添加中间件
func (c *chain) Use(middlewares ...Middleware) ExecutionChain {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, m := range middlewares {
		if m != nil {
			c.middlewares = append(c.middlewares, m)
		}
	}
	return c
}

// Execute 执行中间件链
func (c *chain) Execute(ctx context.Context, execCtx *executor.ExecutionContext) error {
	c.mu.RLock()
	middlewares := make([]Middleware, len(c.middlewares))
	copy(middlewares, c.middlewares)
	c.mu.RUnlock()

	if len(middlewares) == 0 {
		return nil
	}

	var index int = 0

	var handle func(context.Context) error
	handle = func(ctx context.Context) error {
		if index >= len(middlewares) {
			return nil
		}

		mid := middlewares[index]
		index++

		return mid.Handle(ctx, execCtx, handle)
	}

	return handle(ctx)
}

// Remove 移除指定名称的中间件
func (c *chain) Remove(name string) ExecutionChain {
	c.mu.Lock()
	defer c.mu.Unlock()

	var newMiddlewares []Middleware
	for _, m := range c.middlewares {
		if m.Name() != name {
			newMiddlewares = append(newMiddlewares, m)
		}
	}
	c.middlewares = newMiddlewares

	return c
}

// Clear 清除所有中间件
func (c *chain) Clear() ExecutionChain {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.middlewares = make([]Middleware, 0)
	return c
}

// List 获取所有中间件名称列表
func (c *chain) List() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	names := make([]string, len(c.middlewares))
	for i, m := range c.middlewares {
		names[i] = m.Name()
	}
	return names
}

// HasMiddleware 检查是否存在指定名称的中间件
func (c *chain) HasMiddleware(name string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	for _, m := range c.middlewares {
		if m.Name() == name {
			return true
		}
	}
	return false
}

// builder 实现 Builder 接口
type builder struct {
	chain ExecutionChain
}

// NewBuilder 创建新的中间件链构建器
func NewBuilder() Builder {
	return &builder{
		chain: NewChain(),
	}
}

// WithLogging 添加日志中间件
func (b *builder) WithLogging() Builder {
	loggingMiddleware := NewLoggingMiddleware()
	b.chain.Use(loggingMiddleware)
	return b
}

// WithMetrics 添加指标中间件
func (b *builder) WithMetrics() Builder {
	metricsMiddleware := NewMetricsMiddleware()
	b.chain.Use(metricsMiddleware)
	return b
}

// WithRetry 添加重试中间件
func (b *builder) WithRetry(maxAttempts int) Builder {
	if maxAttempts <= 0 {
		maxAttempts = 3
	}
	retryMiddleware := NewRetryMiddleware(maxAttempts)
	b.chain.Use(retryMiddleware)
	return b
}

// WithTimeout 添加超时中间件
func (b *builder) WithTimeout() Builder {
	timeoutMiddleware := NewTimeoutMiddleware()
	b.chain.Use(timeoutMiddleware)
	return b
}

// WithValidation 添加验证中间件
func (b *builder) WithValidation() Builder {
	validationMiddleware := NewValidationMiddleware()
	b.chain.Use(validationMiddleware)
	return b
}

// WithCustom 添加自定义中间件
func (b *builder) WithCustom(middleware Middleware) Builder {
	if middleware != nil {
		b.chain.Use(middleware)
	}
	return b
}

// Build 构建中间件链
func (b *builder) Build() ExecutionChain {
	if b.chain == nil {
		return NewChain()
	}
	return b.chain
}

// SimpleMiddleware 定义简单的中间件实现
type SimpleMiddleware struct {
	name   string
	handle func(ctx context.Context, execCtx *executor.ExecutionContext, next Next) error
}

// NewSimpleMiddleware 创建简单中间件
func NewSimpleMiddleware(name string, handle func(ctx context.Context, execCtx *executor.ExecutionContext, next Next) error) Middleware {
	return &SimpleMiddleware{
		name:   name,
		handle: handle,
	}
}

// Name 返回中间件名称
func (m *SimpleMiddleware) Name() string {
	return m.name
}

// Handle 处理请求
func (m *SimpleMiddleware) Handle(ctx context.Context, execCtx *executor.ExecutionContext, next Next) error {
	if m.handle == nil {
		return next(ctx)
	}
	return m.handle(ctx, execCtx, next)
}

// WrapMiddleware 将简单的处理函数包装为中间件
func WrapMiddleware(name string, handle func(ctx context.Context, execCtx *executor.ExecutionContext, next Next) error) Middleware {
	return NewSimpleMiddleware(name, handle)
}

// PredicateMiddleware 条件中间件
type PredicateMiddleware struct {
	name      string
	predicate func(ctx context.Context, execCtx *executor.ExecutionContext) bool
	handle    func(ctx context.Context, execCtx *executor.ExecutionContext, next Next) error
}

// NewPredicateMiddleware 创建条件中间件
// 只有当 predicate 返回 true 时才执行 handle
func NewPredicateMiddleware(
	name string,
	predicate func(ctx context.Context, execCtx *executor.ExecutionContext) bool,
	handle func(ctx context.Context, execCtx *executor.ExecutionContext, next Next) error,
) Middleware {
	return &PredicateMiddleware{
		name:      name,
		predicate: predicate,
		handle:    handle,
	}
}

// Name 返回中间件名称
func (m *PredicateMiddleware) Name() string {
	return m.name
}

// Handle 处理请求
func (m *PredicateMiddleware) Handle(ctx context.Context, execCtx *executor.ExecutionContext, next Next) error {
	if m.predicate == nil || !m.predicate(ctx, execCtx) {
		return next(ctx)
	}

	if m.handle == nil {
		return next(ctx)
	}

	return m.handle(ctx, execCtx, next)
}

// ErrorHandlerMiddleware 错误处理中间件
type ErrorHandlerMiddleware struct {
	name    string
	handler func(ctx context.Context, execCtx *executor.ExecutionContext, err error) error
}

// NewErrorHandlerMiddleware 创建错误处理中间件
func NewErrorHandlerMiddleware(
	name string,
	handler func(ctx context.Context, execCtx *executor.ExecutionContext, err error) error,
) Middleware {
	return &ErrorHandlerMiddleware{
		name:    name,
		handler: handler,
	}
}

// Name 返回中间件名称
func (m *ErrorHandlerMiddleware) Name() string {
	return m.name
}

// Handle 处理请求
func (m *ErrorHandlerMiddleware) Handle(ctx context.Context, execCtx *executor.ExecutionContext, next Next) error {
	err := next(ctx)
	if err != nil && m.handler != nil {
		return m.handler(ctx, execCtx, err)
	}
	return err
}

// ChainError 中间件链执行错误
type ChainError struct {
	MiddlewareName string
	Err            error
}

// Error 实现 error 接口
func (e *ChainError) Error() string {
	return fmt.Sprintf("middleware %s error: %v", e.MiddlewareName, e.Err)
}

// Unwrap 返回底层错误
func (e *ChainError) Unwrap() error {
	return e.Err
}
