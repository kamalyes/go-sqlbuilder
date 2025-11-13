/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-13
 * @Description: 统一导出层 - 简化版本，确保可编译
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package unified

import (
	"github.com/kamalyes/go-sqlbuilder/db"
	"github.com/kamalyes/go-sqlbuilder/errors"
	"github.com/kamalyes/go-sqlbuilder/logger"
	"github.com/kamalyes/go-sqlbuilder/meta"
	"github.com/kamalyes/go-sqlbuilder/repository"
)

// ==================== Repository 相关导出 ====================

// Repository - 通用仓储接口
type Repository[T any] = repository.Repository[T]

// BaseRepository - 基础仓储实现
type BaseRepository[T any] = repository.BaseRepository[T]

// Filter - 查询过滤条件
type Filter = repository.Filter

// Query - 查询构建器
type Query = repository.Query

// Order - 排序条件
type Order = repository.Order

// Transaction - 事务接口
type Transaction = repository.Transaction

// ==================== DB 相关导出 ====================

// Handler - 数据库处理器接口
type Handler = db.Handler

// ==================== Log 相关导出 ====================

// Logger - 日志接口
type Logger = logger.Logger

// ==================== Error 相关导出 ====================

// ErrorCode - 错误码
type ErrorCode = errors.ErrorCode

// BusinessError - 业务错误
type BusinessError = errors.AppError

// ==================== Meta 相关导出 ====================

// Paging - 分页信息
type Paging = meta.Paging

// ==================== Repository 工厂函数 ====================

// NewBaseRepository - 创建基础仓储
func NewBaseRepository[T any](handler Handler, tableName string) Repository[T] {
	return repository.NewBaseRepository[T](handler, tableName)
}

// ==================== Filter 工厂函数 ====================

// Eq - 等于过滤
func Eq(field string, value interface{}) *Filter {
	return repository.NewEqFilter(field, value)
}

// Gt - 大于过滤
func Gt(field string, value interface{}) *Filter {
	return repository.NewGtFilter(field, value)
}

// Lt - 小于过滤
func Lt(field string, value interface{}) *Filter {
	return repository.NewLtFilter(field, value)
}

// Gte - 大于等于过滤
func Gte(field string, value interface{}) *Filter {
	return repository.NewGteFilter(field, value)
}

// Lte - 小于等于过滤
func Lte(field string, value interface{}) *Filter {
	return repository.NewLteFilter(field, value)
}

// Neq - 不等于过滤
func Neq(field string, value interface{}) *Filter {
	return repository.NewNeqFilter(field, value)
}

// In - IN 过滤
func In(field string, values ...interface{}) *Filter {
	return repository.NewInFilter(field, values...)
}

// Like - LIKE 过滤
func Like(field string, value string) *Filter {
	return repository.NewLikeFilter(field, value)
}

// Between - BETWEEN 过滤
func Between(field string, min, max interface{}) *Filter {
	return repository.NewBetweenFilter(field, min, max)
}

// ==================== Query 工厂函数 ====================

// NewQuery - 创建新的查询条件
func NewQuery() *repository.Query {
	return repository.NewQuery()
}

// ==================== DB 工厂函数 ====================

// NewGormHandler - 创建 GORM 处理器
// 需要导入 gorm，从 gorm.DB 创建处理器
func NewGormHandler(gormDB interface{}) Handler {
	// 需要调用方导入 gorm 并传入 *gorm.DB
	// 这里我们无法处理 interface{} 到 *gorm.DB 的转换
	// 用户应该直接使用: unified.NewGormHandler(db.NewGormHandler(gormDB).DB())
	// 或者在应用层处理类型
	if h, ok := gormDB.(Handler); ok {
		return h
	}
	panic("invalid gormDB type, expected Handler or *gorm.DB")
}

// ==================== Error 工厂函数 ====================

// NewNotFound - 创建 NotFound 错误
func NewNotFound(resource string) *BusinessError {
	return errors.NewNotFound(resource)
}

// NewDatabaseError - 创建数据库错误
func NewDatabaseError(operation string, err error) *BusinessError {
	return errors.NewDatabaseError(operation, err)
}

// NewAccessDenied - 创建访问拒绝错误
func NewAccessDenied(reason string) *BusinessError {
	return errors.NewAccessDenied(reason)
}

// NewInvalidInput - 创建无效输入错误
func NewInvalidInput(field string, reason string) *BusinessError {
	return errors.NewInvalidInput(field, reason)
}

// NewConflict - 创建冲突错误
func NewConflict(resource string) *BusinessError {
	return errors.NewConflict(resource)
}

// NewInternal - 创建内部错误
func NewInternal(message string) *BusinessError {
	return errors.NewInternal(message)
}

// ==================== Error 检查函数 ====================

// IsNotFound - 检查是否是 NotFound 错误
func IsNotFound(err error) bool {
	return errors.IsNotFound(err)
}

// IsAccessDenied - 检查是否是访问拒绝错误
func IsAccessDenied(err error) bool {
	return errors.IsAccessDenied(err)
}

// IsConflict - 检查是否是冲突错误
func IsConflict(err error) bool {
	return errors.IsConflict(err)
}

// IsAlreadyExists - 检查是否是已存在错误
func IsAlreadyExists(err error) bool {
	return errors.IsAlreadyExists(err)
}

// ==================== Paging 工厂函数 ====================

// NewPaging - 创建分页信息
func NewPaging(page, pageSize int) *Paging {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 1000 {
		pageSize = 1000
	}
	return &Paging{
		Page:     int32(page),
		PageSize: int32(pageSize),
	}
}
