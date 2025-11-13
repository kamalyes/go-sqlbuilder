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

	"github.com/kamalyes/go-sqlbuilder/executor"
	"github.com/stretchr/testify/assert"
)

// TestNewChain 测试创建新的中间件链
func TestNewChain(t *testing.T) {
	chain := NewChain()
	assert.NotNil(t, chain)
	assert.Equal(t, 0, len(chain.List()))
}

// TestChainUse 测试添加中间件
func TestChainUse(t *testing.T) {
	chain := NewChain()
	m1 := NewSimpleMiddleware("m1", func(ctx context.Context, execCtx *executor.ExecutionContext, next Next) error {
		return next(ctx)
	})
	m2 := NewSimpleMiddleware("m2", func(ctx context.Context, execCtx *executor.ExecutionContext, next Next) error {
		return next(ctx)
	})

	result := chain.Use(m1, m2)
	assert.Equal(t, chain, result) // 返回自己以支持链式调用
	assert.Equal(t, 2, len(chain.List()))
	assert.Equal(t, "m1", chain.List()[0])
	assert.Equal(t, "m2", chain.List()[1])
}

// TestChainUseWithNil 测试添加 nil 中间件
func TestChainUseWithNil(t *testing.T) {
	chain := NewChain()
	m1 := NewSimpleMiddleware("m1", func(ctx context.Context, execCtx *executor.ExecutionContext, next Next) error {
		return next(ctx)
	})

	chain.Use(m1, nil, m1)
	assert.Equal(t, 2, len(chain.List()))
}

// TestChainExecute 测试执行中间件链
func TestChainExecute(t *testing.T) {
	chain := NewChain()
	execOrder := []string{}

	m1 := NewSimpleMiddleware("m1", func(ctx context.Context, execCtx *executor.ExecutionContext, next Next) error {
		execOrder = append(execOrder, "m1-before")
		err := next(ctx)
		execOrder = append(execOrder, "m1-after")
		return err
	})

	m2 := NewSimpleMiddleware("m2", func(ctx context.Context, execCtx *executor.ExecutionContext, next Next) error {
		execOrder = append(execOrder, "m2-before")
		err := next(ctx)
		execOrder = append(execOrder, "m2-after")
		return err
	})

	chain.Use(m1, m2)

	execCtx := &executor.ExecutionContext{
		SQL: "SELECT * FROM users",
	}

	err := chain.Execute(context.Background(), execCtx)
	assert.NoError(t, err)
	assert.Equal(t, []string{"m1-before", "m2-before", "m2-after", "m1-after"}, execOrder)
}

// TestChainExecuteWithError 测试中间件链错误处理
func TestChainExecuteWithError(t *testing.T) {
	chain := NewChain()
	expectedErr := errors.New("middleware error")

	m1 := NewSimpleMiddleware("m1", func(ctx context.Context, execCtx *executor.ExecutionContext, next Next) error {
		return next(ctx)
	})

	m2 := NewSimpleMiddleware("m2", func(ctx context.Context, execCtx *executor.ExecutionContext, next Next) error {
		return expectedErr
	})

	chain.Use(m1, m2)

	execCtx := &executor.ExecutionContext{
		SQL: "SELECT * FROM users",
	}

	err := chain.Execute(context.Background(), execCtx)
	assert.Equal(t, expectedErr, err)
}

// TestChainExecuteEmpty 测试执行空链
func TestChainExecuteEmpty(t *testing.T) {
	chain := NewChain()
	execCtx := &executor.ExecutionContext{
		SQL: "SELECT * FROM users",
	}

	err := chain.Execute(context.Background(), execCtx)
	assert.NoError(t, err)
}

// TestChainRemove 测试移除中间件
func TestChainRemove(t *testing.T) {
	chain := NewChain()
	m1 := NewSimpleMiddleware("m1", func(ctx context.Context, execCtx *executor.ExecutionContext, next Next) error {
		return next(ctx)
	})
	m2 := NewSimpleMiddleware("m2", func(ctx context.Context, execCtx *executor.ExecutionContext, next Next) error {
		return next(ctx)
	})

	chain.Use(m1, m2)
	assert.Equal(t, 2, len(chain.List()))

	result := chain.Remove("m1")
	assert.Equal(t, chain, result)
	assert.Equal(t, 1, len(chain.List()))
	assert.Equal(t, "m2", chain.List()[0])
}

// TestChainRemoveNonExistent 测试移除不存在的中间件
func TestChainRemoveNonExistent(t *testing.T) {
	chain := NewChain()
	m1 := NewSimpleMiddleware("m1", func(ctx context.Context, execCtx *executor.ExecutionContext, next Next) error {
		return next(ctx)
	})

	chain.Use(m1)
	chain.Remove("nonexistent")

	assert.Equal(t, 1, len(chain.List()))
	assert.Equal(t, "m1", chain.List()[0])
}

// TestChainClear 测试清除所有中间件
func TestChainClear(t *testing.T) {
	chain := NewChain()
	m1 := NewSimpleMiddleware("m1", func(ctx context.Context, execCtx *executor.ExecutionContext, next Next) error {
		return next(ctx)
	})
	m2 := NewSimpleMiddleware("m2", func(ctx context.Context, execCtx *executor.ExecutionContext, next Next) error {
		return next(ctx)
	})

	chain.Use(m1, m2)
	assert.Equal(t, 2, len(chain.List()))

	result := chain.Clear()
	assert.Equal(t, chain, result)
	assert.Equal(t, 0, len(chain.List()))
}

// TestChainList 测试获取中间件列表
func TestChainList(t *testing.T) {
	chain := NewChain()
	m1 := NewSimpleMiddleware("m1", func(ctx context.Context, execCtx *executor.ExecutionContext, next Next) error {
		return next(ctx)
	})
	m2 := NewSimpleMiddleware("m2", func(ctx context.Context, execCtx *executor.ExecutionContext, next Next) error {
		return next(ctx)
	})
	m3 := NewSimpleMiddleware("m3", func(ctx context.Context, execCtx *executor.ExecutionContext, next Next) error {
		return next(ctx)
	})

	chain.Use(m1, m2, m3)

	list := chain.List()
	assert.Equal(t, 3, len(list))
	assert.Equal(t, "m1", list[0])
	assert.Equal(t, "m2", list[1])
	assert.Equal(t, "m3", list[2])
}

// TestChainHasMiddleware 测试检查中间件是否存在
func TestChainHasMiddleware(t *testing.T) {
	chain := NewChain()
	m1 := NewSimpleMiddleware("m1", func(ctx context.Context, execCtx *executor.ExecutionContext, next Next) error {
		return next(ctx)
	})

	assert.False(t, chain.HasMiddleware("m1"))

	chain.Use(m1)
	assert.True(t, chain.HasMiddleware("m1"))
	assert.False(t, chain.HasMiddleware("m2"))
}

// TestNewBuilder 测试创建构建器
func TestNewBuilder(t *testing.T) {
	builder := NewBuilder()
	assert.NotNil(t, builder)
}

// TestBuilderWithCustom 测试添加自定义中间件
func TestBuilderWithCustom(t *testing.T) {
	builder := NewBuilder()
	m1 := NewSimpleMiddleware("custom", func(ctx context.Context, execCtx *executor.ExecutionContext, next Next) error {
		return next(ctx)
	})

	result := builder.WithCustom(m1)
	assert.Equal(t, builder, result)

	chain := builder.Build()
	assert.NotNil(t, chain)
	assert.True(t, chain.HasMiddleware("custom"))
}

// TestBuilderBuild 测试构建链
func TestBuilderBuild(t *testing.T) {
	builder := NewBuilder()
	chain := builder.Build()
	assert.NotNil(t, chain)
	assert.Equal(t, 0, len(chain.List()))
}

// TestSimpleMiddleware 测试简单中间件
func TestSimpleMiddleware(t *testing.T) {
	called := false
	m := NewSimpleMiddleware("test", func(ctx context.Context, execCtx *executor.ExecutionContext, next Next) error {
		called = true
		return next(ctx)
	})

	assert.Equal(t, "test", m.Name())

	execCtx := &executor.ExecutionContext{}
	err := m.Handle(context.Background(), execCtx, func(ctx context.Context) error {
		return nil
	})

	assert.NoError(t, err)
	assert.True(t, called)
}

// TestSimpleMiddlewareWithNilHandle 测试 handle 为 nil 的简单中间件
func TestSimpleMiddlewareWithNilHandle(t *testing.T) {
	m := NewSimpleMiddleware("test", nil)

	execCtx := &executor.ExecutionContext{}
	err := m.Handle(context.Background(), execCtx, func(ctx context.Context) error {
		return nil
	})

	assert.NoError(t, err)
}

// TestWrapMiddleware 测试包装中间件
func TestWrapMiddleware(t *testing.T) {
	called := false
	m := WrapMiddleware("wrapped", func(ctx context.Context, execCtx *executor.ExecutionContext, next Next) error {
		called = true
		return next(ctx)
	})

	assert.Equal(t, "wrapped", m.Name())

	execCtx := &executor.ExecutionContext{}
	err := m.Handle(context.Background(), execCtx, func(ctx context.Context) error {
		return nil
	})

	assert.NoError(t, err)
	assert.True(t, called)
}

// TestPredicateMiddleware 测试条件中间件
func TestPredicateMiddleware(t *testing.T) {
	called := false

	m := NewPredicateMiddleware(
		"predicate",
		func(ctx context.Context, execCtx *executor.ExecutionContext) bool {
			return execCtx.SQL == "SELECT"
		},
		func(ctx context.Context, execCtx *executor.ExecutionContext, next Next) error {
			called = true
			return next(ctx)
		},
	)

	// 条件为真
	execCtx := &executor.ExecutionContext{SQL: "SELECT"}
	err := m.Handle(context.Background(), execCtx, func(ctx context.Context) error {
		return nil
	})

	assert.NoError(t, err)
	assert.True(t, called)

	// 条件为假
	called = false
	execCtx = &executor.ExecutionContext{SQL: "UPDATE"}
	err = m.Handle(context.Background(), execCtx, func(ctx context.Context) error {
		return nil
	})

	assert.NoError(t, err)
	assert.False(t, called)
}

// TestPredicateMiddlewareWithNilPredicate 测试 predicate 为 nil 的条件中间件
func TestPredicateMiddlewareWithNilPredicate(t *testing.T) {
	called := false

	m := NewPredicateMiddleware(
		"predicate",
		nil,
		func(ctx context.Context, execCtx *executor.ExecutionContext, next Next) error {
			called = true
			return next(ctx)
		},
	)

	execCtx := &executor.ExecutionContext{}
	err := m.Handle(context.Background(), execCtx, func(ctx context.Context) error {
		return nil
	})

	assert.NoError(t, err)
	assert.False(t, called)
}

// TestPredicateMiddlewareWithNilHandle 测试 handle 为 nil 的条件中间件
func TestPredicateMiddlewareWithNilHandle(t *testing.T) {
	m := NewPredicateMiddleware(
		"predicate",
		func(ctx context.Context, execCtx *executor.ExecutionContext) bool {
			return true
		},
		nil,
	)

	execCtx := &executor.ExecutionContext{}
	err := m.Handle(context.Background(), execCtx, func(ctx context.Context) error {
		return nil
	})

	assert.NoError(t, err)
}

// TestErrorHandlerMiddleware 测试错误处理中间件
func TestErrorHandlerMiddleware(t *testing.T) {
	recoveredErr := ""

	m := NewErrorHandlerMiddleware(
		"error_handler",
		func(ctx context.Context, execCtx *executor.ExecutionContext, err error) error {
			recoveredErr = err.Error()
			return nil // 恢复错误
		},
	)

	execCtx := &executor.ExecutionContext{}
	err := m.Handle(context.Background(), execCtx, func(ctx context.Context) error {
		return errors.New("test error")
	})

	assert.NoError(t, err)
	assert.Equal(t, "test error", recoveredErr)
}

// TestErrorHandlerMiddlewareWithNoError 测试无错误的错误处理中间件
func TestErrorHandlerMiddlewareWithNoError(t *testing.T) {
	called := false

	m := NewErrorHandlerMiddleware(
		"error_handler",
		func(ctx context.Context, execCtx *executor.ExecutionContext, err error) error {
			called = true
			return err
		},
	)

	execCtx := &executor.ExecutionContext{}
	err := m.Handle(context.Background(), execCtx, func(ctx context.Context) error {
		return nil
	})

	assert.NoError(t, err)
	assert.False(t, called)
}

// TestErrorHandlerMiddlewareWithNilHandler 测试 handler 为 nil 的错误处理中间件
func TestErrorHandlerMiddlewareWithNilHandler(t *testing.T) {
	m := NewErrorHandlerMiddleware("error_handler", nil)

	execCtx := &executor.ExecutionContext{}
	err := m.Handle(context.Background(), execCtx, func(ctx context.Context) error {
		return errors.New("test error")
	})

	assert.Error(t, err)
	assert.Equal(t, "test error", err.Error())
}

// TestChainError 测试链错误类型
func TestChainError(t *testing.T) {
	baseErr := errors.New("base error")
	chainErr := &ChainError{
		MiddlewareName: "test_middleware",
		Err:            baseErr,
	}

	assert.Equal(t, "middleware test_middleware error: base error", chainErr.Error())
	assert.Equal(t, baseErr, chainErr.Unwrap())
}

// TestConcurrentChain 测试并发安全性
func TestConcurrentChain(t *testing.T) {
	chain := NewChain()

	// 创建多个中间件
	for i := 0; i < 10; i++ {
		idx := i
		m := NewSimpleMiddleware("m"+string(rune('0'+idx)), func(ctx context.Context, execCtx *executor.ExecutionContext, next Next) error {
			return next(ctx)
		})
		chain.Use(m)
	}

	// 并发执行链
	done := make(chan error, 5)
	for i := 0; i < 5; i++ {
		go func() {
			execCtx := &executor.ExecutionContext{SQL: "SELECT"}
			done <- chain.Execute(context.Background(), execCtx)
		}()
	}

	for i := 0; i < 5; i++ {
		err := <-done
		assert.NoError(t, err)
	}

	assert.Equal(t, 10, len(chain.List()))
}

// TestChainContextPropagation 测试上下文传播
func TestChainContextPropagation(t *testing.T) {
	chain := NewChain()
	contextValue := "test_value"

	m := NewSimpleMiddleware("context_check", func(ctx context.Context, execCtx *executor.ExecutionContext, next Next) error {
		// 验证上下文值是否传递
		val := ctx.Value("key")
		assert.Equal(t, contextValue, val)
		return next(ctx)
	})

	chain.Use(m)

	ctx := context.WithValue(context.Background(), "key", contextValue)
	execCtx := &executor.ExecutionContext{SQL: "SELECT"}

	err := chain.Execute(ctx, execCtx)
	assert.NoError(t, err)
}
