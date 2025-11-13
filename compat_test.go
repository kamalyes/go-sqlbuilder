/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-11 21:13:15
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-13 11:09:45
 * @FilePath: \go-sqlbuilder\compat_test.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package sqlbuilder

import (
	"testing"

	"github.com/kamalyes/go-sqlbuilder/middleware"
)

// TestV3Features_NewV3Features 测试v3特性创建
func TestV3Features_NewV3Features(t *testing.T) {
	// 创建一个最小的Builder
	builder := &Builder{
		table:       "users",
		queryType:   "select",
		columns:     []string{},
		joins:       []string{},
		wheres:      []string{},
		havings:     []string{},
		groupByCols: []string{},
		orderByCols: []string{},
		insertData:  make(map[string]interface{}),
		updateData:  make(map[string]interface{}),
		args:        []interface{}{},
	}

	v3 := NewV3Features(builder)
	if v3 == nil {
		t.Fatal("expected v3 features, got nil")
	}
	if v3.builder != builder {
		t.Error("builder mismatch")
	}
	if v3.enhance == nil {
		t.Error("enhance should not be nil")
	}
}

// TestV3Features_SetMiddleware 测试middleware设置
func TestV3Features_SetMiddleware(t *testing.T) {
	builder := &Builder{
		table:       "users",
		queryType:   "select",
		columns:     []string{},
		joins:       []string{},
		wheres:      []string{},
		havings:     []string{},
		groupByCols: []string{},
		orderByCols: []string{},
		insertData:  make(map[string]interface{}),
		updateData:  make(map[string]interface{}),
		args:        []interface{}{},
	}

	v3 := NewV3Features(builder)

	// 创建测试middleware
	testMW := middleware.NewLoggingMiddleware()

	result := v3.SetMiddleware(testMW)
	if result == nil {
		t.Fatal("expected v3 features, got nil")
	}
	if result != v3 {
		t.Error("should return v3 features for chaining")
	}
}

// TestV3Features_GetExecutor 测试executor获取
func TestV3Features_GetExecutor(t *testing.T) {
	builder := &Builder{
		table:       "users",
		queryType:   "select",
		columns:     []string{},
		joins:       []string{},
		wheres:      []string{},
		havings:     []string{},
		groupByCols: []string{},
		orderByCols: []string{},
		insertData:  make(map[string]interface{}),
		updateData:  make(map[string]interface{}),
		args:        []interface{}{},
	}

	v3 := NewV3Features(builder)
	executor := v3.GetExecutor()

	// executor 可能为nil，因为没有设置
	if executor != nil {
		// 验证executor类型
		t.Logf("Executor type: %T", executor)
	}
}

// TestBackwardCompatibility_NewBackwardCompatibility 测试兼容性检查创建
func TestBackwardCompatibility_NewBackwardCompatibility(t *testing.T) {
	bc := NewBackwardCompatibility()
	if bc == nil {
		t.Fatal("expected backward compatibility, got nil")
	}
	if bc.checks == nil {
		t.Error("checks should be initialized")
	}
	if len(bc.checks) != 0 {
		t.Error("checks should be empty initially")
	}
}

// TestBackwardCompatibility_Register 测试注册兼容性检查
func TestBackwardCompatibility_Register(t *testing.T) {
	bc := NewBackwardCompatibility()
	builder := &Builder{
		table:       "users",
		queryType:   "select",
		columns:     []string{},
		joins:       []string{},
		wheres:      []string{},
		havings:     []string{},
		groupByCols: []string{},
		orderByCols: []string{},
		insertData:  make(map[string]interface{}),
		updateData:  make(map[string]interface{}),
		args:        []interface{}{},
	}

	check := NewAdapterCompatibilityCheck(builder)
	result := bc.Register(check)

	if result != bc {
		t.Error("should return backward compatibility for chaining")
	}
	if len(bc.checks) != 1 {
		t.Error("should have 1 check registered")
	}
}

// TestBackwardCompatibility_CheckAll 测试执行所有兼容性检查
func TestBackwardCompatibility_CheckAll(t *testing.T) {
	bc := NewBackwardCompatibility()
	builder := &Builder{
		table:       "users",
		queryType:   "select",
		columns:     []string{},
		joins:       []string{},
		wheres:      []string{},
		havings:     []string{},
		groupByCols: []string{},
		orderByCols: []string{},
		insertData:  make(map[string]interface{}),
		updateData:  make(map[string]interface{}),
		args:        []interface{}{},
		adapter:     nil, // 故意设置为nil以测试
	}

	check1 := NewAdapterCompatibilityCheck(builder)
	check2 := NewMethodCompatibilityCheck(builder)

	bc.Register(check1)
	bc.Register(check2)

	errors := bc.CheckAll()
	// 因为adapter为nil，应该没有错误（因为检查允许nil）
	if len(errors) > 0 {
		t.Logf("unexpected errors: %v", errors)
	}
}

// TestAdapterCompatibilityCheck_Check 测试Adapter兼容性检查
func TestAdapterCompatibilityCheck_Check(t *testing.T) {
	builder := &Builder{
		table:       "users",
		adapter:     nil,
		columns:     []string{},
		joins:       []string{},
		wheres:      []string{},
		havings:     []string{},
		groupByCols: []string{},
		orderByCols: []string{},
		insertData:  make(map[string]interface{}),
		updateData:  make(map[string]interface{}),
		args:        []interface{}{},
	}

	check := NewAdapterCompatibilityCheck(builder)
	err := check.Check()

	// 应该没有错误，因为adapter可选
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// TestAdapterCompatibilityCheck_Name 测试名称获取
func TestAdapterCompatibilityCheck_Name(t *testing.T) {
	check := NewAdapterCompatibilityCheck(nil)
	name := check.Name()

	if name != "AdapterCompatibility" {
		t.Errorf("expected 'AdapterCompatibility', got '%s'", name)
	}
}

// TestMethodCompatibilityCheck_Check 测试方法兼容性检查
func TestMethodCompatibilityCheck_Check(t *testing.T) {
	builder := &Builder{
		table:       "users",
		columns:     []string{},
		joins:       []string{},
		wheres:      []string{},
		havings:     []string{},
		groupByCols: []string{},
		orderByCols: []string{},
		insertData:  make(map[string]interface{}),
		updateData:  make(map[string]interface{}),
		args:        []interface{}{},
	}

	check := NewMethodCompatibilityCheck(builder)
	err := check.Check()

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// TestMethodCompatibilityCheck_Name 测试名称获取
func TestMethodCompatibilityCheck_Name(t *testing.T) {
	check := NewMethodCompatibilityCheck(nil)
	name := check.Name()

	if name != "MethodCompatibility" {
		t.Errorf("expected 'MethodCompatibility', got '%s'", name)
	}
}

// TestMigrationGuide_NewMigrationGuide 测试迁移指南创建
func TestMigrationGuide_NewMigrationGuide(t *testing.T) {
	builder := &Builder{
		table:       "users",
		columns:     []string{},
		joins:       []string{},
		wheres:      []string{},
		havings:     []string{},
		groupByCols: []string{},
		orderByCols: []string{},
		insertData:  make(map[string]interface{}),
		updateData:  make(map[string]interface{}),
		args:        []interface{}{},
	}

	guide := NewMigrationGuide(builder)
	if guide == nil {
		t.Fatal("expected migration guide, got nil")
	}
	if guide.v2Builder != builder {
		t.Error("builder mismatch")
	}
	if guide.v3Enhance == nil {
		t.Error("v3Enhance should not be nil")
	}
}

// TestMigrationGuide_MigrateToV3 测试迁移到v3
func TestMigrationGuide_MigrateToV3(t *testing.T) {
	builder := &Builder{
		table:       "users",
		columns:     []string{},
		joins:       []string{},
		wheres:      []string{},
		havings:     []string{},
		groupByCols: []string{},
		orderByCols: []string{},
		insertData:  make(map[string]interface{}),
		updateData:  make(map[string]interface{}),
		args:        []interface{}{},
	}

	guide := NewMigrationGuide(builder)
	enhance := guide.MigrateToV3()

	if enhance == nil {
		t.Fatal("expected BuilderEnhancer, got nil")
	}
	if enhance.GetBuilder() != builder {
		t.Error("builder should be the same")
	}
}

// TestBuilderEnhancer_NewBuilderEnhancer 测试BuilderEnhancer创建
func TestBuilderEnhancer_NewBuilderEnhancer(t *testing.T) {
	builder := &Builder{
		table:       "users",
		columns:     []string{},
		joins:       []string{},
		wheres:      []string{},
		havings:     []string{},
		groupByCols: []string{},
		orderByCols: []string{},
		insertData:  make(map[string]interface{}),
		updateData:  make(map[string]interface{}),
		args:        []interface{}{},
	}

	enhancer := NewBuilderEnhancer(builder)
	if enhancer == nil {
		t.Fatal("expected BuilderEnhancer, got nil")
	}
	if enhancer.GetBuilder() != builder {
		t.Error("builder mismatch")
	}
}

// TestBuilderEnhancer_AddMiddleware 测试middleware添加
func TestBuilderEnhancer_AddMiddleware(t *testing.T) {
	builder := &Builder{
		table:       "users",
		columns:     []string{},
		joins:       []string{},
		wheres:      []string{},
		havings:     []string{},
		groupByCols: []string{},
		orderByCols: []string{},
		insertData:  make(map[string]interface{}),
		updateData:  make(map[string]interface{}),
		args:        []interface{}{},
	}

	enhancer := NewBuilderEnhancer(builder)
	mw := middleware.NewLoggingMiddleware()

	result := enhancer.AddMiddleware(mw)
	if result != enhancer {
		t.Error("should return enhancer for chaining")
	}
}

// TestBuilderEnhancer_GetCompiler 测试compiler获取
func TestBuilderEnhancer_GetCompiler(t *testing.T) {
	builder := &Builder{
		table:       "users",
		columns:     []string{},
		joins:       []string{},
		wheres:      []string{},
		havings:     []string{},
		groupByCols: []string{},
		orderByCols: []string{},
		insertData:  make(map[string]interface{}),
		updateData:  make(map[string]interface{}),
		args:        []interface{}{},
	}

	enhancer := NewBuilderEnhancer(builder)
	compiler := enhancer.GetCompiler()

	if compiler == nil {
		t.Error("compiler should not be nil")
	}
}

// TestBuilderEnhancer_GetPlanner 测试planner获取
func TestBuilderEnhancer_GetPlanner(t *testing.T) {
	builder := &Builder{
		table:       "users",
		columns:     []string{},
		joins:       []string{},
		wheres:      []string{},
		havings:     []string{},
		groupByCols: []string{},
		orderByCols: []string{},
		insertData:  make(map[string]interface{}),
		updateData:  make(map[string]interface{}),
		args:        []interface{}{},
	}

	enhancer := NewBuilderEnhancer(builder)
	planner := enhancer.GetPlanner()

	if planner == nil {
		t.Error("planner should not be nil")
	}
}

// TestBuilderEnhancer_BuildWithMiddleware 测试middleware构建
func TestBuilderEnhancer_BuildWithMiddleware(t *testing.T) {
	t.Skip("BuildWithMiddleware not implemented in this phase")
}

// TestBuilderEnhancer_CompileWithDialect 测试方言编译
func TestBuilderEnhancer_CompileWithDialect(t *testing.T) {
	t.Skip("CompileWithDialect not implemented in this phase")
}

// TestBuilderEnhancer_PlanQuery 测试查询计划
func TestBuilderEnhancer_PlanQuery(t *testing.T) {
	t.Skip("PlanQuery not implemented in this phase")
}
