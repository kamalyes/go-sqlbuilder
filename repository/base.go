/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-11 21:13:15
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-13 11:07:41
 * @FilePath: \go-sqlbuilder\persist\query_param.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package repository

import (
	"context"
	"fmt"

	"github.com/kamalyes/go-sqlbuilder/constant"
	"github.com/kamalyes/go-sqlbuilder/db"
	"github.com/kamalyes/go-sqlbuilder/errors"
	"github.com/kamalyes/go-sqlbuilder/meta"
	"gorm.io/gorm"
)

// BaseRepository 基础仓储实现，包含通用的 CRUD 操作
type BaseRepository[T any] struct {
	db    db.Handler
	table string
}

// NewBaseRepository 创建基础仓储
func NewBaseRepository[T any](dbHandler db.Handler, table string) *BaseRepository[T] {
	return &BaseRepository[T]{
		db:    dbHandler,
		table: table,
	}
}

// Create 创建单个记录
func (r *BaseRepository[T]) Create(ctx context.Context, entity *T) (*T, error) {
	if entity == nil {
		return nil, errors.NewError(errors.ErrorCodeInvalidInput, errors.MsgEntityCannotBeNil)
	}

	result := r.db.DB().WithContext(ctx).Table(r.table).Create(entity)
	if result.Error != nil {
		return nil, result.Error
	}

	return entity, nil
}

// CreateBatch 批量创建记录
func (r *BaseRepository[T]) CreateBatch(ctx context.Context, entities ...*T) error {
	if len(entities) == 0 {
		return nil
	}

	result := r.db.DB().WithContext(ctx).Table(r.table).CreateInBatches(entities, 100)
	return result.Error
}

// Get 获取单个记录
func (r *BaseRepository[T]) Get(ctx context.Context, id interface{}) (*T, error) {
	var entity T
	result := r.db.DB().WithContext(ctx).Table(r.table).Where("id = ?", id).First(&entity)
	if result.Error != nil {
		return nil, result.Error
	}

	return &entity, nil
}

// GetByFilter 按单个过滤条件获取记录
func (r *BaseRepository[T]) GetByFilter(ctx context.Context, filter *Filter) (*T, error) {
	if filter == nil {
		return nil, errors.NewError(errors.ErrorCodeInvalidInput, errors.MsgFilterCannotBeNil)
	}

	var entity T
	query := r.db.DB().WithContext(ctx).Table(r.table)
	query = applyFilter(query, filter)

	result := query.First(&entity)
	if result.Error != nil {
		return nil, result.Error
	}

	return &entity, nil
}

// GetByFilters 按多个过滤条件获取记录
func (r *BaseRepository[T]) GetByFilters(ctx context.Context, filters ...*Filter) (*T, error) {
	if len(filters) == 0 {
		return nil, errors.NewError(errors.ErrorCodeInvalidInput, errors.MsgAtLeastOneFilterRequired)
	}

	var entity T
	query := r.db.DB().WithContext(ctx).Table(r.table)
	for _, filter := range filters {
		query = applyFilter(query, filter)
	}

	result := query.First(&entity)
	if result.Error != nil {
		return nil, result.Error
	}

	return &entity, nil
}

// List 列表查询
func (r *BaseRepository[T]) List(ctx context.Context, query *Query) ([]*T, error) {
	if query == nil {
		query = NewQuery()
	}

	var entities []*T
	db := r.db.DB().WithContext(ctx).Table(r.table)

	// 应用过滤条件
	for _, filter := range query.Filters {
		db = applyFilter(db, filter)
	}

	// 应用排序
	for _, order := range query.Orders {
		db = db.Order(order.Field + " " + order.Direction)
	}

	result := db.Find(&entities)
	if result.Error != nil {
		return nil, result.Error
	}

	return entities, nil
}

// ListWithPagination 分页列表查询
func (r *BaseRepository[T]) ListWithPagination(ctx context.Context, query *Query, page *meta.Paging) ([]*T, *meta.Paging, error) {
	if query == nil {
		query = NewQuery()
	}

	if page == nil {
		page = &meta.Paging{Page: 1, PageSize: 10}
	}

	var entities []*T
	db := r.db.DB().WithContext(ctx).Table(r.table)

	// 应用过滤条件
	for _, filter := range query.Filters {
		db = applyFilter(db, filter)
	}

	// 计算总数
	var total int64
	countDb := db
	countDb.Model(new(T)).Count(&total)
	page.Total = total

	// 应用排序
	for _, order := range query.Orders {
		db = db.Order(order.Field + " " + order.Direction)
	}

	// 应用分页
	offset := (int(page.Page) - 1) * int(page.PageSize)
	result := db.Offset(offset).Limit(int(page.PageSize)).Find(&entities)
	if result.Error != nil {
		return nil, nil, result.Error
	}

	return entities, page, nil
}

// Update 更新单个记录
func (r *BaseRepository[T]) Update(ctx context.Context, entity *T) (*T, error) {
	if entity == nil {
		return nil, errors.NewError(errors.ErrorCodeInvalidInput, errors.MsgEntityCannotBeNil)
	}

	result := r.db.DB().WithContext(ctx).Table(r.table).Save(entity)
	if result.Error != nil {
		return nil, result.Error
	}

	return entity, nil
}

// UpdateBatch 批量更新记录
func (r *BaseRepository[T]) UpdateBatch(ctx context.Context, entities ...*T) error {
	if len(entities) == 0 {
		return nil
	}

	// 使用事务确保批量更新的一致性
	return r.db.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, entity := range entities {
			if err := tx.Table(r.table).Save(entity).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// UpdateByFilters 按过滤条件更新记录
func (r *BaseRepository[T]) UpdateByFilters(ctx context.Context, entity *T, filters ...*Filter) error {
	if entity == nil {
		return errors.NewError(errors.ErrorCodeInvalidInput, errors.MsgEntityCannotBeNil)
	}

	if len(filters) == 0 {
		return errors.NewError(errors.ErrorCodeInvalidInput, errors.MsgAtLeastOneFilterRequired)
	}

	db := r.db.DB().WithContext(ctx).Table(r.table)
	for _, filter := range filters {
		db = applyFilter(db, filter)
	}

	result := db.Updates(entity)
	return result.Error
}

// Delete 删除单个记录
func (r *BaseRepository[T]) Delete(ctx context.Context, id interface{}) error {
	result := r.db.DB().WithContext(ctx).Table(r.table).Where("id = ?", id).Delete(new(T))
	return result.Error
}

// DeleteBatch 批量删除记录
func (r *BaseRepository[T]) DeleteBatch(ctx context.Context, ids ...interface{}) error {
	if len(ids) == 0 {
		return nil
	}

	result := r.db.DB().WithContext(ctx).Table(r.table).Where("id IN ?", ids).Delete(new(T))
	return result.Error
}

// DeleteByFilters 按过滤条件删除记录
func (r *BaseRepository[T]) DeleteByFilters(ctx context.Context, filters ...*Filter) error {
	if len(filters) == 0 {
		return errors.NewError(errors.ErrorCodeInvalidInput, errors.MsgAtLeastOneFilterRequired)
	}

	db := r.db.DB().WithContext(ctx).Table(r.table)
	for _, filter := range filters {
		db = applyFilter(db, filter)
	}

	result := db.Delete(new(T))
	return result.Error
}

// Transaction 事务支持
func (r *BaseRepository[T]) Transaction(ctx context.Context, fn func(tx Transaction) error) error {
	return r.db.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txWrapper := &transactionWrapper{db: db.NewGormHandler(tx)}
		return fn(txWrapper)
	})
}

// Count 计数
func (r *BaseRepository[T]) Count(ctx context.Context, filters ...*Filter) (int64, error) {
	var count int64
	db := r.db.DB().WithContext(ctx).Table(r.table)

	for _, filter := range filters {
		db = applyFilter(db, filter)
	}

	result := db.Count(&count)
	return count, result.Error
}

// Exists 检查记录是否存在
func (r *BaseRepository[T]) Exists(ctx context.Context, filters ...*Filter) (bool, error) {
	count, err := r.Count(ctx, filters...)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// applyFilter 应用单个过滤条件到 GORM 查询
func applyFilter(dbQuery *gorm.DB, filter *Filter) *gorm.DB {
	if filter == nil {
		return dbQuery
	}

	switch filter.Operator {
	case "=":
		return dbQuery.Where(fmt.Sprintf("%s = ?", filter.Field), filter.Value)
	case ">":
		return dbQuery.Where(fmt.Sprintf("%s > ?", filter.Field), filter.Value)
	case "<":
		return dbQuery.Where(fmt.Sprintf("%s < ?", filter.Field), filter.Value)
	case ">=":
		return dbQuery.Where(fmt.Sprintf("%s >= ?", filter.Field), filter.Value)
	case "<=":
		return dbQuery.Where(fmt.Sprintf("%s <= ?", filter.Field), filter.Value)
	case "!=":
		return dbQuery.Where(fmt.Sprintf("%s != ?", filter.Field), filter.Value)
	case string(constant.OP_IN):
		return dbQuery.Where(fmt.Sprintf("%s IN ?", filter.Field), filter.Value)
	case string(constant.OP_LIKE):
		return dbQuery.Where(fmt.Sprintf("%s LIKE ?", filter.Field), filter.Value)
	case string(constant.OP_BETWEEN):
		if values, ok := filter.Value.([]interface{}); ok && len(values) == 2 {
			return dbQuery.Where(fmt.Sprintf("%s BETWEEN ? AND ?", filter.Field), values[0], values[1])
		}
	}

	return dbQuery
}

// transactionWrapper 事务包装器
type transactionWrapper struct {
	db db.Handler
}

// Create 在事务中创建
func (t *transactionWrapper) Create(ctx context.Context, entity interface{}) error {
	return t.db.DB().WithContext(ctx).Create(entity).Error
}

// CreateBatch 在事务中批量创建
func (t *transactionWrapper) CreateBatch(ctx context.Context, entities ...interface{}) error {
	return t.db.DB().WithContext(ctx).CreateInBatches(entities, 100).Error
}

// Update 在事务中更新
func (t *transactionWrapper) Update(ctx context.Context, entity interface{}) error {
	return t.db.DB().WithContext(ctx).Save(entity).Error
}

// UpdateBatch 在事务中批量更新
func (t *transactionWrapper) UpdateBatch(ctx context.Context, entities ...interface{}) error {
	for _, entity := range entities {
		if err := t.db.DB().WithContext(ctx).Save(entity).Error; err != nil {
			return err
		}
	}
	return nil
}

// Delete 在事务中删除
func (t *transactionWrapper) Delete(ctx context.Context, entity interface{}) error {
	return t.db.DB().WithContext(ctx).Delete(entity).Error
}

// DeleteBatch 在事务中批量删除
func (t *transactionWrapper) DeleteBatch(ctx context.Context, entities ...interface{}) error {
	for _, entity := range entities {
		if err := t.db.DB().WithContext(ctx).Delete(entity).Error; err != nil {
			return err
		}
	}
	return nil
}
