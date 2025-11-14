/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-14 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-14 13:30:11
 * @FilePath: \engine-im-service\go-sqlbuilder\easy_crud.go
 * @Description: 简单易用的CRUD操作 - 基于EnhancedBuilder的简化版本
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package sqlbuilder

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/kamalyes/go-sqlbuilder/constant"
)

// ==================== 简单易用的方法 ====================

// Add 添加数据 - 一行代码搞定
func (eb *EnhancedBuilder) Add(data map[string]interface{}) error {
	if len(data) == 0 {
		return errors.New("数据不能为空")
	}

	// 使用hook自动处理时间戳
	eb.setupDefaultHooks()

	ctx := context.Background()
	options := &CreateOptions{}
	_, err := eb.SmartCreate(ctx, data, options)
	if err != nil {
		return errors.New("添加数据失败")
	}

	return nil
}

// Get 获取单条数据 - 按ID查询
func (eb *EnhancedBuilder) Get(id interface{}) (map[string]interface{}, error) {
	if id == nil {
		return nil, errors.New("ID不能为空")
	}

	ctx := context.Background()
	var result map[string]interface{}
	err := eb.WithContext(ctx).Table(eb.table).Where("id", "=", id).First(&result)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("数据不存在")
		}
		return nil, errors.New("查询数据失败")
	}

	return result, nil
}

// List 获取数据列表 - 支持分页
func (eb *EnhancedBuilder) List(page, pageSize int64) ([]map[string]interface{}, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}

	ctx := context.Background()
	offset := (page - 1) * pageSize
	var results []map[string]interface{}

	query := eb.WithContext(ctx).Table(eb.table).OrderByDesc("id")

	// 如果启用软删除，过滤已删除的数据
	if eb.softDeleteEnabled {
		query = query.WhereNull(constant.FieldDeletedAt)
	}

	err := query.Limit(pageSize).Offset(offset).Get(&results)
	if err != nil {
		return nil, errors.New("查询数据失败")
	}

	return results, nil
}

// Update 更新数据 - 按ID更新
func (eb *EnhancedBuilder) Update(id interface{}, data map[string]interface{}) error {
	if id == nil {
		return errors.New("ID不能为空")
	}
	if len(data) == 0 {
		return errors.New("更新数据不能为空")
	}

	// 使用hook自动处理时间戳
	eb.setupDefaultHooks()

	ctx := context.Background()
	options := &UpdateOptions{}
	_, err := eb.SmartUpdate(ctx, id, data, options)
	if err != nil {
		return errors.New("更新数据失败")
	}

	return nil
}

// Delete 删除数据 - 按ID删除（软删除）
func (eb *EnhancedBuilder) Delete(id interface{}) error {
	if id == nil {
		return errors.New("ID不能为空")
	}

	// 使用hook自动处理时间戳
	eb.setupDefaultHooks()

	ctx := context.Background()
	options := &DeleteOptions{Force: false} // 强制软删除
	_, err := eb.SmartDelete(ctx, id, options)
	if err != nil {
		return errors.New("删除数据失败")
	}

	return nil
}

// Count 统计数据数量
func (eb *EnhancedBuilder) Count() (int64, error) {
	ctx := context.Background()
	query := eb.WithContext(ctx).Table(eb.table)

	// 如果启用软删除，只统计未删除的数据
	if eb.softDeleteEnabled {
		query = query.WhereNull(constant.FieldDeletedAt)
	}

	count, err := query.Count()
	if err != nil {
		return 0, errors.New("统计数据失败")
	}

	return count, nil
}

// Search 简单搜索 - 模糊查询
func (eb *EnhancedBuilder) Search(field, keyword string, page, pageSize int64) ([]map[string]interface{}, error) {
	if field == "" || keyword == "" {
		return eb.List(page, pageSize)
	}

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}

	ctx := context.Background()
	offset := (page - 1) * pageSize
	var results []map[string]interface{}

	query := eb.WithContext(ctx).Table(eb.table).
		WhereLike(field, "%"+keyword+"%").
		OrderByDesc("id")

	// 如果启用软删除，过滤已删除的数据
	if eb.softDeleteEnabled {
		query = query.WhereNull(constant.FieldDeletedAt)
	}

	err := query.Limit(pageSize).Offset(offset).Get(&results)
	if err != nil {
		return nil, errors.New("搜索数据失败")
	}

	return results, nil
}

// Save 保存数据（智能判断是新增还是更新）
func (eb *EnhancedBuilder) Save(data map[string]interface{}) error {
	if id, exists := data["id"]; exists && id != nil {
		// 有ID就更新
		updateData := make(map[string]interface{})
		for k, v := range data {
			if k != "id" { // 排除ID字段
				updateData[k] = v
			}
		}
		return eb.Update(id, updateData)
	} else {
		// 没ID就新增
		return eb.Add(data)
	}
}

// ==================== 别名方法 ====================

// Create 创建数据（Add的别名）
func (eb *EnhancedBuilder) Create(data map[string]interface{}) error {
	return eb.Add(data)
}

// Read 读取数据（Get的别名）
func (eb *EnhancedBuilder) Read(id interface{}) (map[string]interface{}, error) {
	return eb.Get(id)
}

// Find 查找数据（List的别名）
func (eb *EnhancedBuilder) Find(page, pageSize int64) ([]map[string]interface{}, error) {
	return eb.List(page, pageSize)
}

// ==================== 辅助方法 ====================

// setupDefaultHooks 设置默认的时间戳hooks
func (eb *EnhancedBuilder) setupDefaultHooks() {
	// 只设置一次
	if len(eb.hooks) > 0 {
		return
	}

	// 创建前自动添加时间戳
	eb.AddHook(constant.HookEventBeforeCreate, func(ctx context.Context, data interface{}) error {
		if dataMap, ok := data.(map[string]interface{}); ok {
			now := time.Now()
			dataMap[constant.FieldCreatedAt] = now
			dataMap[constant.FieldUpdatedAt] = now
		}
		return nil
	})

	// 更新前自动添加更新时间
	eb.AddHook(constant.HookEventBeforeUpdate, func(ctx context.Context, data interface{}) error {
		if dataMap, ok := data.(map[string]interface{}); ok {
			dataMap[constant.FieldUpdatedAt] = time.Now()
		}
		return nil
	})

	// 删除前自动添加删除时间
	eb.AddHook(constant.HookEventBeforeDelete, func(ctx context.Context, data interface{}) error {
		// 删除操作在SmartDelete中已经处理了软删除逻辑
		return nil
	})
}

// ==================== 快速使用函数 ====================

// NewSimple 创建一个简单易用的CRUD操作器
func NewSimple(db interface{}, tableName string) (*EnhancedBuilder, error) {
	builder, err := NewEnhanced(db)
	if err != nil {
		return nil, errors.New("数据库连接失败")
	}

	// 设置表名
	builder.Table(tableName)

	// 启用默认功能
	builder.EnableSoftDelete(true)

	// 设置默认审计字段
	builder.AddAuditFields(
		constant.FieldCreatedAt,
		constant.FieldUpdatedAt,
		constant.FieldDeletedAt,
	)

	return builder, nil
}
