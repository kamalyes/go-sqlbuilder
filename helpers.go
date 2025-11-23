/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-23 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-23 22:48:38
 * @FilePath: \go-sqlbuilder\helpers.go
 * @Description: 仓储辅助工具 - 软删除、查询辅助等功能
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package sqlbuilder

import (
	"context"
	"gorm.io/gorm"
	"strings"
	"time"
)

// === 独立辅助函数 ===

// GetDeleted 获取已删除的记录（使用deleted_at字段）
func GetDeleted[T any](ctx context.Context, db *gorm.DB, query *Query) ([]*T, error) {
	if query == nil {
		query = NewQuery()
	}

	// 添加deleted_at不为空的条件
	query.AddFilter(&Filter{
		Field:    "deleted_at",
		Operator: OP_IS_NOT_NULL,
	})

	var results []*T
	dbQuery := db.WithContext(ctx).Model(new(T))

	// 应用查询条件（包括过滤器、分页、排序）
	dbQuery = applyQuery(dbQuery, query)

	err := dbQuery.Find(&results).Error
	return results, err
}

// GetNonDeleted 获取未删除的记录（使用deleted_at字段）
func GetNonDeleted[T any](ctx context.Context, db *gorm.DB, query *Query) ([]*T, error) {
	if query == nil {
		query = NewQuery()
	}

	// 添加deleted_at为空的条件
	query.AddFilter(&Filter{
		Field:    "deleted_at",
		Operator: OP_IS_NULL,
	})

	var results []*T
	dbQuery := db.WithContext(ctx).Model(new(T))

	// 应用查询条件（包括过滤器、分页、排序）
	dbQuery = applyQuery(dbQuery, query)

	err := dbQuery.Find(&results).Error
	return results, err
}

// RestoreDeleted 恢复单个已删除记录（将deleted_at设为NULL）
func RestoreDeleted[T any](ctx context.Context, db *gorm.DB, id interface{}) error {
	return db.WithContext(ctx).Model(new(T)).Where("id = ?", id).Update("deleted_at", nil).Error
}

// RestoreDeletedBatch 批量恢复已删除记录（将deleted_at设为NULL）
func RestoreDeletedBatch[T any](ctx context.Context, db *gorm.DB, ids []interface{}) error {
	if len(ids) == 0 {
		return nil
	}
	return db.WithContext(ctx).Model(new(T)).Where("id IN ?", ids).Update("deleted_at", nil).Error
}

// PermanentlyDelete 永久删除记录（从数据库中完全删除）
func PermanentlyDelete[T any](ctx context.Context, db *gorm.DB, id interface{}) error {
	return db.WithContext(ctx).Unscoped().Where("id = ?", id).Delete(new(T)).Error
}

// PermanentlyDeleteBatch 批量永久删除记录（从数据库中完全删除）
func PermanentlyDeleteBatch[T any](ctx context.Context, db *gorm.DB, ids []interface{}) error {
	if len(ids) == 0 {
		return nil
	}
	return db.WithContext(ctx).Unscoped().Where("id IN ?", ids).Delete(new(T)).Error
}

// applyQuery 应用查询条件的辅助函数
func applyQuery(dbQuery *gorm.DB, query *Query) *gorm.DB {
	// 应用过滤条件
	if query.FilterGroup != nil {
		conditions, args := buildGroupCondition(query.FilterGroup)
		if conditions != "" {
			dbQuery = dbQuery.Where(conditions, args...)
		}
	}
	for _, filter := range query.Filters {
		condition, arg := buildFilterCondition(filter)
		if condition != "" {
			if arg != nil {
				dbQuery = dbQuery.Where(condition, arg)
			} else {
				dbQuery = dbQuery.Where(condition)
			}
		}
	}

	// 应用排序
	for _, order := range query.Orders {
		if order.Field != "" {
			dbQuery = dbQuery.Order(order.Field + " " + order.Direction)
		}
	}

	// 应用分页
	if query.Pagination != nil && query.Pagination.Page > 0 && query.Pagination.PageSize > 0 {
		offset := (int(query.Pagination.Page) - 1) * int(query.Pagination.PageSize)
		dbQuery = dbQuery.Offset(offset).Limit(int(query.Pagination.PageSize))
	}

	// 应用限制和偏移
	if query.LimitValue != nil && *query.LimitValue > 0 {
		dbQuery = dbQuery.Limit(*query.LimitValue)
	}
	if query.OffsetValue != nil && *query.OffsetValue > 0 {
		dbQuery = dbQuery.Offset(*query.OffsetValue)
	}

	// 应用去重
	if query.Distinct {
		dbQuery = dbQuery.Distinct()
	}

	// 应用分组
	if len(query.GroupBy) > 0 {
		dbQuery = dbQuery.Group(strings.Join(query.GroupBy, ", "))
	}

	// 应用HAVING条件
	for _, havingFilter := range query.Having {
		condition, arg := buildFilterCondition(havingFilter)
		if condition != "" {
			if arg != nil {
				dbQuery = dbQuery.Having(condition, arg)
			} else {
				dbQuery = dbQuery.Having(condition)
			}
		}
	}

	return dbQuery
}

// RepositoryWithSoftDelete 带软删除功能的仓储
type RepositoryWithSoftDelete[T any] struct {
	*BaseRepository[T]
}

// NewRepositoryWithSoftDelete 创建带软删除功能的仓储
func NewRepositoryWithSoftDelete[T any](repo *BaseRepository[T]) *RepositoryWithSoftDelete[T] {
	return &RepositoryWithSoftDelete[T]{
		BaseRepository: repo,
	}
}

// === deleted_at 字段软删除 ===

// SoftDeleteWithDeletedAt 使用 deleted_at 字段软删除
func (r *RepositoryWithSoftDelete[T]) SoftDeleteWithDeletedAt(ctx context.Context, id interface{}) error {
	return r.SoftDelete(ctx, id, "deleted_at", time.Now())
}

// SoftDeleteBatchWithDeletedAt 使用 deleted_at 字段批量软删除
func (r *RepositoryWithSoftDelete[T]) SoftDeleteBatchWithDeletedAt(ctx context.Context, ids []interface{}) error {
	return r.SoftDeleteBatch(ctx, ids, "deleted_at", time.Now())
}

// SoftDeleteByFiltersWithDeletedAt 使用 deleted_at 字段按条件软删除
func (r *RepositoryWithSoftDelete[T]) SoftDeleteByFiltersWithDeletedAt(ctx context.Context, filters ...*Filter) error {
	return r.SoftDeleteByFilters(ctx, "deleted_at", time.Now(), filters...)
}

// RestoreWithDeletedAt 使用 deleted_at 字段恢复
func (r *RepositoryWithSoftDelete[T]) RestoreWithDeletedAt(ctx context.Context, id interface{}) error {
	return r.Restore(ctx, id, "deleted_at", nil)
}

// RestoreBatchWithDeletedAt 使用 deleted_at 字段批量恢复
func (r *RepositoryWithSoftDelete[T]) RestoreBatchWithDeletedAt(ctx context.Context, ids []interface{}) error {
	return r.RestoreBatch(ctx, ids, "deleted_at", nil)
}

// === is_deleted 字段软删除 ===

// SoftDeleteWithIsDeleted 使用 is_deleted 字段软删除
func (r *RepositoryWithSoftDelete[T]) SoftDeleteWithIsDeleted(ctx context.Context, id interface{}) error {
	return r.SoftDelete(ctx, id, "is_deleted", 1)
}

// SoftDeleteBatchWithIsDeleted 使用 is_deleted 字段批量软删除
func (r *RepositoryWithSoftDelete[T]) SoftDeleteBatchWithIsDeleted(ctx context.Context, ids []interface{}) error {
	return r.SoftDeleteBatch(ctx, ids, "is_deleted", 1)
}

// SoftDeleteByFiltersWithIsDeleted 使用 is_deleted 字段按条件软删除
func (r *RepositoryWithSoftDelete[T]) SoftDeleteByFiltersWithIsDeleted(ctx context.Context, filters ...*Filter) error {
	return r.SoftDeleteByFilters(ctx, "is_deleted", 1, filters...)
}

// RestoreWithIsDeleted 使用 is_deleted 字段恢复
func (r *RepositoryWithSoftDelete[T]) RestoreWithIsDeleted(ctx context.Context, id interface{}) error {
	return r.Restore(ctx, id, "is_deleted", 0)
}

// RestoreBatchWithIsDeleted 使用 is_deleted 字段批量恢复
func (r *RepositoryWithSoftDelete[T]) RestoreBatchWithIsDeleted(ctx context.Context, ids []interface{}) error {
	return r.RestoreBatch(ctx, ids, "is_deleted", 0)
}

// === 查询辅助方法 ===

// ListNotDeleted 查询未删除的记录（deleted_at 字段）
func (r *RepositoryWithSoftDelete[T]) ListNotDeleted(ctx context.Context, query *Query) ([]*T, error) {
	if query == nil {
		query = NewQuery()
	}
	query.AddFilter(NewIsNullFilter("deleted_at"))
	return r.List(ctx, query)
}

// ListNotDeletedByIsDeleted 查询未删除的记录（is_deleted 字段）
func (r *RepositoryWithSoftDelete[T]) ListNotDeletedByIsDeleted(ctx context.Context, query *Query) ([]*T, error) {
	if query == nil {
		query = NewQuery()
	}
	query.AddFilter(NewEqFilter("is_deleted", 0))
	return r.List(ctx, query)
}

// ListDeleted 查询已删除的记录（deleted_at 字段）
func (r *RepositoryWithSoftDelete[T]) ListDeleted(ctx context.Context, query *Query) ([]*T, error) {
	if query == nil {
		query = NewQuery()
	}
	query.AddFilter(NewIsNotNullFilter("deleted_at"))
	return r.List(ctx, query)
}

// ListDeletedByIsDeleted 查询已删除的记录（is_deleted 字段）
func (r *RepositoryWithSoftDelete[T]) ListDeletedByIsDeleted(ctx context.Context, query *Query) ([]*T, error) {
	if query == nil {
		query = NewQuery()
	}
	query.AddFilter(NewEqFilter("is_deleted", 1))
	return r.List(ctx, query)
}
