/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-11 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-13 11:02:21
 * @FilePath: \go-sqlbuilder\compiler\compiler.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package compiler

import (
	"fmt"
	"strings"
	"sync"

	"github.com/kamalyes/go-sqlbuilder/constant"
	"github.com/kamalyes/go-sqlbuilder/executor"
)

// DefaultCompiler 默认SQL编译器实现
type DefaultCompiler struct {
	dialect string
	config  *CompilerConfig
	mu      sync.RWMutex
}

// NewDefaultCompiler 创建默认编译器
func NewDefaultCompiler(dialect string) SQLCompiler {
	return &DefaultCompiler{
		dialect: dialect,
		config: &CompilerConfig{
			Options: CompilerOptions{
				Dialect:                dialect,
				EnableOptimization:     true,
				EnablePlanning:         true,
				MaxQueryComplexity:     10,
				EnableParameterization: true,
				EnableQueryCache:       false,
			},
			Optimizers:         make([]Optimizer, 0),
			CustomTransformers: make(map[string]func(string) string),
		},
	}
}

// Compile 编译SQL查询
func (c *DefaultCompiler) Compile(execCtx *executor.ExecutionContext) (string, []interface{}, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if execCtx == nil {
		return "", nil, fmt.Errorf("execution context cannot be nil")
	}

	sql := execCtx.SQL
	args := execCtx.Args

	if sql == "" {
		return "", nil, fmt.Errorf("SQL cannot be empty")
	}

	// 规范化SQL
	sql = c.normalizSQL(sql)

	// 应用方言转换
	sql = c.transformForDialect(sql)

	// 应用自定义转换
	sql = c.applyCustomTransformers(sql)

	return sql, args, nil
}

// GetDialect 获取方言
func (c *DefaultCompiler) GetDialect() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.dialect
}

// SetDialect 设置方言
func (c *DefaultCompiler) SetDialect(dialect string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !isValidDialect(dialect) {
		return fmt.Errorf("invalid dialect: %s", dialect)
	}

	c.dialect = dialect
	c.config.Options.Dialect = dialect
	return nil
}

// normalizSQL 规范化SQL
func (c *DefaultCompiler) normalizSQL(sql string) string {
	// 移除多余空格
	sql = strings.TrimSpace(sql)

	// 统一换行符
	sql = strings.ReplaceAll(sql, "\r\n", "\n")
	sql = strings.ReplaceAll(sql, "\r", "\n")

	// 移除多余的空行
	lines := strings.Split(sql, "\n")
	var result []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return strings.Join(result, " ")
}

// transformForDialect 根据方言转换SQL
func (c *DefaultCompiler) transformForDialect(sql string) string {
	switch c.dialect {
	case constant.DialectMySQL:
		return c.transformMySQL(sql)
	case constant.DialectPostgres:
		return c.transformPostgreSQL(sql)
	case constant.DialectSQLite:
		return c.transformSQLite(sql)
	case constant.DialectSQLServer:
		return c.transformSQLServer(sql)
	default:
		return sql
	}
}

// transformMySQL MySQL特定转换
func (c *DefaultCompiler) transformMySQL(sql string) string {
	// MySQL特定的转换规则
	// 例如：使用 LIMIT 而不是 OFFSET
	return sql
}

// transformPostgreSQL PostgreSQL特定转换
func (c *DefaultCompiler) transformPostgreSQL(sql string) string {
	// PostgreSQL特定的转换规则
	// 例如：使用 OFFSET/LIMIT 语法
	return sql
}

// transformSQLite SQLite特定转换
func (c *DefaultCompiler) transformSQLite(sql string) string {
	// SQLite特定的转换规则
	return sql
}

// transformSQLServer SQL Server特定转换
func (c *DefaultCompiler) transformSQLServer(sql string) string {
	// SQL Server特定的转换规则
	// 例如：使用 TOP 而不是 LIMIT
	return sql
}

// transformOracle Oracle特定转换
func (c *DefaultCompiler) transformOracle(sql string) string {
	// Oracle特定的转换规则
	return sql
}

// applyCustomTransformers 应用自定义转换器
func (c *DefaultCompiler) applyCustomTransformers(sql string) string {
	for _, transformer := range c.config.CustomTransformers {
		sql = transformer(sql)
	}
	return sql
}

// AddOptimizer 添加优化器
func (c *DefaultCompiler) AddOptimizer(optimizer Optimizer) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.config.Optimizers = append(c.config.Optimizers, optimizer)
}

// AddCustomTransformer 添加自定义转换器
func (c *DefaultCompiler) AddCustomTransformer(name string, transformer func(string) string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.config.CustomTransformers[name] = transformer
}

// isValidDialect 检查是否为有效的方言
func isValidDialect(dialect string) bool {
	switch dialect {
	case constant.DialectMySQL, constant.DialectPostgres, constant.DialectSQLite, constant.DialectSQLServer:
		return true
	default:
		return false
	}
}

// DialectCompiler 方言编译器工厂
type DialectCompiler struct {
	compilers map[string]SQLCompiler
	mu        sync.RWMutex
}

// NewDialectCompiler 创建方言编译器工厂
func NewDialectCompiler() *DialectCompiler {
	dc := &DialectCompiler{
		compilers: make(map[string]SQLCompiler),
	}

	// 初始化所有方言的编译器
	dc.compilers[constant.DialectMySQL] = NewDefaultCompiler(constant.DialectMySQL)
	dc.compilers[constant.DialectPostgres] = NewDefaultCompiler(constant.DialectPostgres)
	dc.compilers[constant.DialectSQLite] = NewDefaultCompiler(constant.DialectSQLite)
	dc.compilers[constant.DialectSQLServer] = NewDefaultCompiler(constant.DialectSQLServer)

	return dc
}

// GetCompiler 获取指定方言的编译器
func (dc *DialectCompiler) GetCompiler(dialect string) (SQLCompiler, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()

	compiler, ok := dc.compilers[dialect]
	if !ok {
		return nil, fmt.Errorf("compiler not found for dialect: %s", dialect)
	}

	return compiler, nil
}

// RegisterCompiler 注册自定义编译器
func (dc *DialectCompiler) RegisterCompiler(dialect string, compiler SQLCompiler) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()

	if !isValidDialect(dialect) {
		return fmt.Errorf("invalid dialect: %s", dialect)
	}

	dc.compilers[dialect] = compiler
	return nil
}

// Compile 使用默认编译器编译
func (dc *DialectCompiler) Compile(dialect string, execCtx *executor.ExecutionContext) (string, []interface{}, error) {
	compiler, err := dc.GetCompiler(dialect)
	if err != nil {
		return "", nil, err
	}

	return compiler.Compile(execCtx)
}
