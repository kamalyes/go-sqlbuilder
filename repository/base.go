/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-11 21:13:15
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-20 13:33:13
 * @FilePath: \go-sqlbuilder\repository\base.go
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

	// 应用 Limit
	if query.LimitValue != nil {
		db = db.Limit(*query.LimitValue)
	}

	// 应用 Offset
	if query.OffsetValue != nil {
		db = db.Offset(*query.OffsetValue)
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

// Find 通用查询方法，兼容旧的API调用方式
func (r *BaseRepository[T]) Find(ctx context.Context, options *FindOptions) ([]*T, error) {
	if options == nil {
		return r.GetAll(ctx)
	}

	query := NewQuery()

	// 转换条件
	for _, condition := range options.Conditions {
		filter := &Filter{
			Field:    condition.Field,
			Operator: condition.Op,
			Value:    condition.Value,
		}
		query.AddFilter(filter)
	}

	// 转换排序
	for _, order := range options.Orders {
		query.AddOrder(order.Field, order.Direction)
	}

	// 应用限制
	if options.Limit > 0 {
		query.Limit(options.Limit)
	}

	// 应用偏移量
	if options.Offset > 0 {
		query.Offset(options.Offset)
	}

	return r.List(ctx, query)
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

// GetAll 获取所有记录（不分页）
func (r *BaseRepository[T]) GetAll(ctx context.Context) ([]*T, error) {
	var entities []*T
	result := r.db.DB().WithContext(ctx).Table(r.table).Find(&entities)
	if result.Error != nil {
		return nil, result.Error
	}
	return entities, nil
}

// First 获取第一条记录
func (r *BaseRepository[T]) First(ctx context.Context, filters ...*Filter) (*T, error) {
	var entity T
	db := r.db.DB().WithContext(ctx).Table(r.table)

	for _, filter := range filters {
		db = applyFilter(db, filter)
	}

	result := db.First(&entity)
	if result.Error != nil {
		return nil, result.Error
	}

	return &entity, nil
}

// Last 获取最后一条记录
func (r *BaseRepository[T]) Last(ctx context.Context, filters ...*Filter) (*T, error) {
	var entity T
	db := r.db.DB().WithContext(ctx).Table(r.table)

	for _, filter := range filters {
		db = applyFilter(db, filter)
	}

	result := db.Last(&entity)
	if result.Error != nil {
		return nil, result.Error
	}

	return &entity, nil
}

// FindOne 查找单条记录（不存在返回 nil）
func (r *BaseRepository[T]) FindOne(ctx context.Context, filters ...*Filter) (*T, error) {
	var entity T
	db := r.db.DB().WithContext(ctx).Table(r.table)

	for _, filter := range filters {
		db = applyFilter(db, filter)
	}

	result := db.Limit(1).Find(&entity)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	return &entity, nil
}

// UpdateFields 更新指定字段
func (r *BaseRepository[T]) UpdateFields(ctx context.Context, id interface{}, fields map[string]interface{}) error {
	if len(fields) == 0 {
		return nil
	}

	result := r.db.DB().WithContext(ctx).Table(r.table).Where("id = ?", id).Updates(fields)
	return result.Error
}

// UpdateFieldsByFilters 按过滤条件更新指定字段
func (r *BaseRepository[T]) UpdateFieldsByFilters(ctx context.Context, fields map[string]interface{}, filters ...*Filter) error {
	if len(fields) == 0 {
		return nil
	}

	if len(filters) == 0 {
		return errors.NewError(errors.ErrorCodeInvalidInput, errors.MsgAtLeastOneFilterRequired)
	}

	db := r.db.DB().WithContext(ctx).Table(r.table)
	for _, filter := range filters {
		db = applyFilter(db, filter)
	}

	result := db.Updates(fields)
	return result.Error
}

// SoftDelete 软删除（需要指定删除标记字段和值）
// field: 软删除字段名，如 "deleted_at", "is_deleted" 等
// value: 软删除标记值，如 time.Now(), 1 等
func (r *BaseRepository[T]) SoftDelete(ctx context.Context, id interface{}, field string, value interface{}) error {
	result := r.db.DB().WithContext(ctx).Table(r.table).Where("id = ?", id).Update(field, value)
	return result.Error
}

// SoftDeleteBatch 批量软删除
// field: 软删除字段名，如 "deleted_at", "is_deleted" 等
// value: 软删除标记值，如 time.Now(), 1 等
func (r *BaseRepository[T]) SoftDeleteBatch(ctx context.Context, ids []interface{}, field string, value interface{}) error {
	if len(ids) == 0 {
		return nil
	}

	result := r.db.DB().WithContext(ctx).Table(r.table).Where("id IN ?", ids).Update(field, value)
	return result.Error
}

// SoftDeleteByFilters 按过滤条件软删除
// field: 软删除字段名，如 "deleted_at", "is_deleted" 等
// value: 软删除标记值，如 time.Now(), 1 等
func (r *BaseRepository[T]) SoftDeleteByFilters(ctx context.Context, field string, value interface{}, filters ...*Filter) error {
	if len(filters) == 0 {
		return errors.NewError(errors.ErrorCodeInvalidInput, errors.MsgAtLeastOneFilterRequired)
	}

	db := r.db.DB().WithContext(ctx).Table(r.table)
	for _, filter := range filters {
		db = applyFilter(db, filter)
	}

	result := db.Update(field, value)
	return result.Error
}

// Restore 恢复软删除的记录
// field: 软删除字段名，如 "deleted_at", "is_deleted" 等
// restoreValue: 恢复时的值，如 nil, 0 等
func (r *BaseRepository[T]) Restore(ctx context.Context, id interface{}, field string, restoreValue interface{}) error {
	result := r.db.DB().WithContext(ctx).Table(r.table).Where("id = ?", id).Update(field, restoreValue)
	return result.Error
}

// RestoreBatch 批量恢复软删除的记录
// field: 软删除字段名，如 "deleted_at", "is_deleted" 等
// restoreValue: 恢复时的值，如 nil, 0 等
func (r *BaseRepository[T]) RestoreBatch(ctx context.Context, ids []interface{}, field string, restoreValue interface{}) error {
	if len(ids) == 0 {
		return nil
	}

	result := r.db.DB().WithContext(ctx).Table(r.table).Where("id IN ?", ids).Update(field, restoreValue)
	return result.Error
}

// CountByField 按字段计数（GROUP BY）
func (r *BaseRepository[T]) CountByField(ctx context.Context, field string) (map[interface{}]int64, error) {
	type Result struct {
		Field interface{}
		Count int64
	}

	var results []Result
	query := fmt.Sprintf("%s as field, COUNT(*) as count", field)
	db := r.db.DB().WithContext(ctx).Table(r.table).Select(query).Group(field)

	if err := db.Scan(&results).Error; err != nil {
		return nil, err
	}

	countMap := make(map[interface{}]int64)
	for _, result := range results {
		countMap[result.Field] = result.Count
	}

	return countMap, nil
}

// Pluck 提取单个字段的值列表
func (r *BaseRepository[T]) Pluck(ctx context.Context, field string, filters ...*Filter) ([]interface{}, error) {
	var values []interface{}
	db := r.db.DB().WithContext(ctx).Table(r.table)

	for _, filter := range filters {
		db = applyFilter(db, filter)
	}

	if err := db.Pluck(field, &values).Error; err != nil {
		return nil, err
	}

	return values, nil
}

// Distinct 获取去重后的字段值列表
func (r *BaseRepository[T]) Distinct(ctx context.Context, field string, filters ...*Filter) ([]interface{}, error) {
	var values []interface{}
	db := r.db.DB().WithContext(ctx).Table(r.table).Distinct(field)

	for _, filter := range filters {
		db = applyFilter(db, filter)
	}

	if err := db.Pluck(field, &values).Error; err != nil {
		return nil, err
	}

	return values, nil
}

// applyFilter 应用单个过滤条件到 GORM 查询
func applyFilter(dbQuery *gorm.DB, filter *Filter) *gorm.DB {
	if filter == nil {
		return dbQuery
	}

	switch filter.Operator {
	case constant.OP_EQ:
		return dbQuery.Where(fmt.Sprintf("%s = ?", filter.Field), filter.Value)
	case constant.OP_GT:
		return dbQuery.Where(fmt.Sprintf("%s > ?", filter.Field), filter.Value)
	case constant.OP_LT:
		return dbQuery.Where(fmt.Sprintf("%s < ?", filter.Field), filter.Value)
	case constant.OP_GTE:
		return dbQuery.Where(fmt.Sprintf("%s >= ?", filter.Field), filter.Value)
	case constant.OP_LTE:
		return dbQuery.Where(fmt.Sprintf("%s <= ?", filter.Field), filter.Value)
	case constant.OP_NEQ:
		return dbQuery.Where(fmt.Sprintf("%s != ?", filter.Field), filter.Value)
	case constant.OP_IN:
		return dbQuery.Where(fmt.Sprintf("%s IN ?", filter.Field), filter.Value)
	case constant.OP_NOT_IN:
		return dbQuery.Where(fmt.Sprintf("%s NOT IN ?", filter.Field), filter.Value)
	case constant.OP_LIKE:
		return dbQuery.Where(fmt.Sprintf("%s LIKE ?", filter.Field), filter.Value)
	case constant.OP_BETWEEN:
		if values, ok := filter.Value.([]interface{}); ok && len(values) == 2 {
			return dbQuery.Where(fmt.Sprintf("%s BETWEEN ? AND ?", filter.Field), values[0], values[1])
		}
	case constant.OP_IS_NULL:
		return dbQuery.Where(fmt.Sprintf("%s IS NULL", filter.Field))
	case constant.OP_IS_NOT_NULL:
		return dbQuery.Where(fmt.Sprintf("%s IS NOT NULL", filter.Field))
	}

	return dbQuery
}

// DBHandler 获取数据库处理器
func (r *BaseRepository[T]) DBHandler() db.Handler {
	return r.db
}

// Table 获取表名
func (r *BaseRepository[T]) Table() string {
	return r.table
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
