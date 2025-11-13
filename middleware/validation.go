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
	"regexp"
	"strings"

	"github.com/kamalyes/go-sqlbuilder/executor"
)

// ValidationMiddleware 验证中间件
// 在执行查询前进行验证
type ValidationMiddleware struct {
	name       string
	validators []Validator
}

// Validator 验证器接口
type Validator interface {
	// Validate 验证执行上下文
	Validate(ctx context.Context, execCtx *executor.ExecutionContext) error
}

// NewValidationMiddleware 创建验证中间件
func NewValidationMiddleware() Middleware {
	return &ValidationMiddleware{
		name:       "validation",
		validators: make([]Validator, 0),
	}
}

// Name 返回中间件名称
func (m *ValidationMiddleware) Name() string {
	return m.name
}

// AddValidator 添加验证器
func (m *ValidationMiddleware) AddValidator(validator Validator) {
	if validator != nil {
		m.validators = append(m.validators, validator)
	}
}

// Handle 处理请求
func (m *ValidationMiddleware) Handle(ctx context.Context, execCtx *executor.ExecutionContext, next Next) error {
	// 执行所有验证器
	for _, validator := range m.validators {
		if err := validator.Validate(ctx, execCtx); err != nil {
			return fmt.Errorf("validation failed: %w", err)
		}
	}

	return next(ctx)
}

// SQLValidator SQL 验证器
type SQLValidator struct {
	name string
}

// NewSQLValidator 创建 SQL 验证器
func NewSQLValidator() Validator {
	return &SQLValidator{
		name: "sql_validator",
	}
}

// Validate 验证 SQL
func (v *SQLValidator) Validate(ctx context.Context, execCtx *executor.ExecutionContext) error {
	// 检查 SQL 是否为空
	if strings.TrimSpace(execCtx.SQL) == "" {
		return fmt.Errorf("SQL cannot be empty")
	}

	// 检查是否包含危险的 SQL 关键字
	upperSQL := strings.ToUpper(execCtx.SQL)
	dangerousKeywords := []string{"DROP", "TRUNCATE", "DELETE FROM"}

	for _, keyword := range dangerousKeywords {
		if strings.Contains(upperSQL, keyword) {
			return fmt.Errorf("dangerous SQL keyword detected: %s", keyword)
		}
	}

	return nil
}

// ArgsValidator 参数验证器
type ArgsValidator struct {
	name    string
	maxArgs int
}

// NewArgsValidator 创建参数验证器
func NewArgsValidator() Validator {
	return &ArgsValidator{
		name:    "args_validator",
		maxArgs: 1000,
	}
}

// Validate 验证参数
func (v *ArgsValidator) Validate(ctx context.Context, execCtx *executor.ExecutionContext) error {
	// 检查参数数量
	if len(execCtx.Args) > v.maxArgs {
		return fmt.Errorf("too many arguments: %d > %d", len(execCtx.Args), v.maxArgs)
	}

	// 检查参数是否包含 nil
	for i, arg := range execCtx.Args {
		if arg == nil {
			// 可以根据需要修改此行为
			// return fmt.Errorf("argument %d is nil", i)
		} else {
			// 检查参数类型
			switch v := arg.(type) {
			case string:
				if len(v) > 1000000 { // 1MB
					return fmt.Errorf("argument %d is too large: %d bytes", i, len(v))
				}
			}
		}
	}

	return nil
}

// SetMaxArgs 设置最大参数数
func (v *ArgsValidator) SetMaxArgs(maxArgs int) {
	if maxArgs > 0 {
		v.maxArgs = maxArgs
	}
}

// PatternValidator 模式验证器
// 检查 SQL 是否匹配预定义的模式
type PatternValidator struct {
	name     string
	patterns []*regexp.Regexp
	allow    bool // true: 只允许匹配的模式，false: 拒绝匹配的模式
}

// NewPatternValidator 创建模式验证器
// allow: true 表示只允许匹配的模式，false 表示拒绝匹配的模式
func NewPatternValidator(allow bool) *PatternValidator {
	return &PatternValidator{
		name:     "pattern_validator",
		patterns: make([]*regexp.Regexp, 0),
		allow:    allow,
	}
}

// AddPattern 添加模式
func (v *PatternValidator) AddPattern(pattern string) error {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("invalid pattern: %w", err)
	}
	v.patterns = append(v.patterns, re)
	return nil
}

// Validate 验证模式
func (v *PatternValidator) Validate(ctx context.Context, execCtx *executor.ExecutionContext) error {
	for _, pattern := range v.patterns {
		if pattern.MatchString(execCtx.SQL) {
			if v.allow {
				return nil // 允许该模式
			}
			return fmt.Errorf("SQL matches forbidden pattern: %s", pattern.String())
		}
	}

	if v.allow {
		return fmt.Errorf("SQL does not match any allowed pattern")
	}

	return nil
}

// LengthValidator SQL 长度验证器
type LengthValidator struct {
	name      string
	maxLength int
}

// NewLengthValidator 创建 SQL 长度验证器
func NewLengthValidator(maxLength int) Validator {
	if maxLength <= 0 {
		maxLength = 1000000 // 1MB
	}

	return &LengthValidator{
		name:      "length_validator",
		maxLength: maxLength,
	}
}

// Validate 验证 SQL 长度
func (v *LengthValidator) Validate(ctx context.Context, execCtx *executor.ExecutionContext) error {
	if len(execCtx.SQL) > v.maxLength {
		return fmt.Errorf("SQL is too long: %d > %d bytes", len(execCtx.SQL), v.maxLength)
	}

	return nil
}

// CompositeValidator 复合验证器
// 支持多个验证器的组合，所有验证器都必须通过
type CompositeValidator struct {
	name       string
	validators []Validator
}

// NewCompositeValidator 创建复合验证器
func NewCompositeValidator(validators ...Validator) Validator {
	return &CompositeValidator{
		name:       "composite_validator",
		validators: validators,
	}
}

// AddValidator 添加验证器
func (v *CompositeValidator) AddValidator(validator Validator) {
	if validator != nil {
		v.validators = append(v.validators, validator)
	}
}

// Validate 执行所有验证
func (v *CompositeValidator) Validate(ctx context.Context, execCtx *executor.ExecutionContext) error {
	for _, validator := range v.validators {
		if err := validator.Validate(ctx, execCtx); err != nil {
			return err
		}
	}

	return nil
}

// CustomValidator 自定义验证器
type CustomValidator struct {
	name     string
	validate func(ctx context.Context, execCtx *executor.ExecutionContext) error
}

// NewCustomValidator 创建自定义验证器
func NewCustomValidator(name string, validate func(ctx context.Context, execCtx *executor.ExecutionContext) error) Validator {
	return &CustomValidator{
		name:     name,
		validate: validate,
	}
}

// Validate 执行自定义验证
func (v *CustomValidator) Validate(ctx context.Context, execCtx *executor.ExecutionContext) error {
	if v.validate == nil {
		return nil
	}

	return v.validate(ctx, execCtx)
}

// SQLInjectionValidator SQL 注入防护验证器
type SQLInjectionValidator struct {
	name               string
	suspiciousPatterns []*regexp.Regexp
}

// NewSQLInjectionValidator 创建 SQL 注入防护验证器
func NewSQLInjectionValidator() Validator {
	patterns := []string{
		`(?i)(union|select|insert|update|delete|drop|create|alter|exec|execute)\s+(.*)(union|select|insert|update|delete|drop|create|alter|exec|execute)`,
		`(?i)(;|--|#)\s*`,
		`(?i)(\*|<|>|!|=)\s*`,
	}

	suspiciousPatterns := make([]*regexp.Regexp, 0)
	for _, pattern := range patterns {
		if re, err := regexp.Compile(pattern); err == nil {
			suspiciousPatterns = append(suspiciousPatterns, re)
		}
	}

	return &SQLInjectionValidator{
		name:               "sql_injection_validator",
		suspiciousPatterns: suspiciousPatterns,
	}
}

// Validate 验证 SQL 注入
func (v *SQLInjectionValidator) Validate(ctx context.Context, execCtx *executor.ExecutionContext) error {
	for _, pattern := range v.suspiciousPatterns {
		if pattern.MatchString(execCtx.SQL) {
			return fmt.Errorf("potential SQL injection detected")
		}
	}

	return nil
}
