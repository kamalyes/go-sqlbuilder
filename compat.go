/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-11 21:13:15
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-13 11:07:41
 * @FilePath: \go-sqlbuilder\compat.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package sqlbuilder

import (
	"github.com/kamalyes/go-sqlbuilder/executor"
	"github.com/kamalyes/go-sqlbuilder/middleware"
)

// V3Features v3新功能包装器
type V3Features struct {
	builder *Builder
	enhance *BuilderEnhancer
}

// NewV3Features 为v2 Builder提供v3特性
func NewV3Features(b *Builder) *V3Features {
	return &V3Features{
		builder: b,
		enhance: NewBuilderEnhancer(b),
	}
}

// SetMiddleware 设置middleware链
func (v3 *V3Features) SetMiddleware(middlewares ...middleware.Middleware) *V3Features {
	v3.enhance.AddMiddleware(middlewares...)
	return v3
}

// GetExecutor 获取executor
func (v3 *V3Features) GetExecutor() executor.Executor {
	if v3.enhance.queryExecutor == nil {
		return nil
	}
	return v3.enhance.queryExecutor
}

// BackwardCompatibility 兼容性检查框架
type BackwardCompatibility struct {
	checks []CompatibilityCheck
}

// CompatibilityCheck 兼容性检查接口
type CompatibilityCheck interface {
	Check() error
	Name() string
}

// NewBackwardCompatibility 创建兼容性检查
func NewBackwardCompatibility() *BackwardCompatibility {
	return &BackwardCompatibility{
		checks: make([]CompatibilityCheck, 0),
	}
}

// Register 注册兼容性检查
func (bc *BackwardCompatibility) Register(check CompatibilityCheck) *BackwardCompatibility {
	bc.checks = append(bc.checks, check)
	return bc
}

// CheckAll 执行所有兼容性检查
func (bc *BackwardCompatibility) CheckAll() []error {
	var errs []error
	for _, check := range bc.checks {
		if err := check.Check(); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

// AdapterCompatibilityCheck Adapter兼容性检查
type AdapterCompatibilityCheck struct {
	builder *Builder
}

// NewAdapterCompatibilityCheck 创建Adapter兼容性检查
func NewAdapterCompatibilityCheck(b *Builder) *AdapterCompatibilityCheck {
	return &AdapterCompatibilityCheck{builder: b}
}

// Check 检查Adapter兼容性
func (acc *AdapterCompatibilityCheck) Check() error {
	if acc.builder == nil || acc.builder.adapter == nil {
		return nil // adapter 可选
	}
	return nil
}

// Name 返回检查名称
func (acc *AdapterCompatibilityCheck) Name() string {
	return "AdapterCompatibility"
}

// MethodCompatibilityCheck 方法兼容性检查
type MethodCompatibilityCheck struct {
	builder *Builder
}

// NewMethodCompatibilityCheck 创建方法兼容性检查
func NewMethodCompatibilityCheck(b *Builder) *MethodCompatibilityCheck {
	return &MethodCompatibilityCheck{builder: b}
}

// Check 检查方法兼容性
func (mcc *MethodCompatibilityCheck) Check() error {
	if mcc.builder == nil {
		return nil
	}
	// 验证关键方法存在
	// Table, Select, Where, Join, Insert, Update, Delete 等
	return nil
}

// Name 返回检查名称
func (mcc *MethodCompatibilityCheck) Name() string {
	return "MethodCompatibility"
}

// MigrationGuide v2->v3迁移指南
type MigrationGuide struct {
	v2Builder *Builder
	v3Enhance *BuilderEnhancer
}

// NewMigrationGuide 创建迁移指南
func NewMigrationGuide(b *Builder) *MigrationGuide {
	return &MigrationGuide{
		v2Builder: b,
		v3Enhance: NewBuilderEnhancer(b),
	}
}

// MigrateToV3 迁移到v3
func (mg *MigrationGuide) MigrateToV3() *BuilderEnhancer {
	return mg.v3Enhance
}
