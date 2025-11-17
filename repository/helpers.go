/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-17 15:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-17 15:30:00
 * @FilePath: \go-sqlbuilder\repository\helpers.go
 * @Description: 仓储辅助方法
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package repository

import (
	"context"
	"time"
)

// SoftDeleteHelper 软删除辅助接口
type SoftDeleteHelper[T any] interface {
	// 使用 deleted_at 字段的软删除
	SoftDeleteWithDeletedAt(ctx context.Context, id interface{}) error
	SoftDeleteBatchWithDeletedAt(ctx context.Context, ids []interface{}) error
	SoftDeleteByFiltersWithDeletedAt(ctx context.Context, filters ...*Filter) error
	RestoreWithDeletedAt(ctx context.Context, id interface{}) error
	RestoreBatchWithDeletedAt(ctx context.Context, ids []interface{}) error

	// 使用 is_deleted 字段的软删除
	SoftDeleteWithIsDeleted(ctx context.Context, id interface{}) error
	SoftDeleteBatchWithIsDeleted(ctx context.Context, ids []interface{}) error
	SoftDeleteByFiltersWithIsDeleted(ctx context.Context, filters ...*Filter) error
	RestoreWithIsDeleted(ctx context.Context, id interface{}) error
	RestoreBatchWithIsDeleted(ctx context.Context, ids []interface{}) error
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
