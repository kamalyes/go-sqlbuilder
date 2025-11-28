/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-11 21:13:15
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-28 13:15:15
 * @FilePath: \go-sqlbuilder\repository\base.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package repository

import (
	"context"
	"fmt"
	"github.com/kamalyes/go-logger"
	gologger "github.com/kamalyes/go-logger"
	"github.com/kamalyes/go-sqlbuilder/constants"
	"github.com/kamalyes/go-sqlbuilder/db"
	"github.com/kamalyes/go-sqlbuilder/errors"
	"github.com/kamalyes/go-toolbox/pkg/errorx"
	"gorm.io/gorm"
	"reflect"
	"strings"
	"time"
)

// ContextFieldExtractor context字段提取器函数类型
// 用于从context中提取需要记录到日志的字段
type ContextFieldExtractor func(ctx context.Context, log logger.ILogger) logger.ILogger

// BaseRepository 基础仓储实现，包含通用的 CRUD 操作
type BaseRepository[T any] struct {
	db               db.Handler
	table            string
	batchSize        int                   // 批处理大小
	timeout          int                   // 查询超时时间（秒）
	readOnly         bool                  // 只读模式
	preloads         []string              // 默认预加载关联
	defaultOrder     string                // 默认排序
	logger           gologger.ILogger      // 日志记录器
	contextExtractor ContextFieldExtractor // context字段提取器
	modelFields      []string              // 模型字段缓存（用于自动字段选择）
	autoFields       bool                  // 是否启用自动字段模式
}

// 编译时检查 - 确保 BaseRepository 实现了 Repository 接口
var _ Repository[any] = (*BaseRepository[any])(nil)

// RepositoryOption 仓储配置选项
type RepositoryOption[T any] func(*BaseRepository[T])

// WithBatchSize 设置批处理大小
func WithBatchSize[T any](size int) RepositoryOption[T] {
	return func(r *BaseRepository[T]) {
		if size > 0 {
			r.batchSize = size
		}
	}
}

// WithTimeout 设置查询超时时间
func WithTimeout[T any](seconds int) RepositoryOption[T] {
	return func(r *BaseRepository[T]) {
		if seconds > 0 {
			r.timeout = seconds
		}
	}
}

// WithReadOnly 设置为只读模式
func WithReadOnly[T any]() RepositoryOption[T] {
	return func(r *BaseRepository[T]) {
		r.readOnly = true
	}
}

// WithDefaultPreloads 设置默认预加载关联
func WithDefaultPreloads[T any](preloads ...string) RepositoryOption[T] {
	return func(r *BaseRepository[T]) {
		r.preloads = preloads
	}
}

// WithDefaultOrder 设置默认排序
func WithDefaultOrder[T any](order string) RepositoryOption[T] {
	return func(r *BaseRepository[T]) {
		r.defaultOrder = order
	}
}

// WithLogger 设置日志记录器
func WithLogger[T any](log gologger.ILogger) RepositoryOption[T] {
	return func(r *BaseRepository[T]) {
		r.logger = log
	}
}

// WithAutoFields 启用自动字段模式（根据model自动提取字段）
func WithAutoFields[T any]() RepositoryOption[T] {
	return func(r *BaseRepository[T]) {
		r.autoFields = true
		// 初始化时提取并缓存字段
		var model T
		r.modelFields = GetStructFields(model)
	}
}

// NewBaseRepository 创建基础仓储
func NewBaseRepository[T any](dbHandler db.Handler, logger gologger.ILogger, table string, options ...RepositoryOption[T]) *BaseRepository[T] {
	r := &BaseRepository[T]{
		db:         dbHandler,
		table:      table,
		batchSize:  constants.DefaultBatchSize,
		timeout:    constants.DefaultQueryTimeout,
		logger:     logger,
		autoFields: false,
	}

	// 应用配置选项
	for _, option := range options {
		option(r)
	}

	return r
}

// Create 创建单个记录
func (r *BaseRepository[T]) Create(ctx context.Context, entity *T) (*T, error) {
	if r.readOnly {
		return nil, errorx.NewError(errors.ErrorCodeForbidden)
	}

	if entity == nil {
		return nil, errorx.NewError(errors.ErrorCodeInvalidInput)
	}

	result := r.db.GetDB().WithContext(ctx).Table(r.table).Create(entity)

	if result.Error != nil {
		return nil, r.handleErrorWithContext(ctx, result.Error, "create")
	}

	return entity, nil
}

// CreateIfNotExists 如果不存在则创建
func (r *BaseRepository[T]) CreateIfNotExists(ctx context.Context, entity *T, uniqueFields ...string) (*T, bool, error) {
	if r.readOnly {
		return nil, false, errorx.NewError(errors.ErrorCodeForbidden)
	}

	if entity == nil || len(uniqueFields) == 0 {
		return nil, false, errorx.NewError(errors.ErrorCodeInvalidInput)
	}

	// 构建查询条件检查是否存在
	filters := make([]*Filter, 0, len(uniqueFields))
	entityValue := reflect.ValueOf(entity).Elem()
	entityType := entityValue.Type()

	for _, field := range uniqueFields {
		for i := 0; i < entityType.NumField(); i++ {
			structField := entityType.Field(i)
			if structField.Name == field || structField.Tag.Get("json") == field {
				fieldValue := entityValue.Field(i).Interface()
				filters = append(filters, NewEqFilter(field, fieldValue))
				break
			}
		}
	}

	// 检查是否存在
	exists, err := r.Exists(ctx, filters...)
	if err != nil {
		return nil, false, err
	}

	if exists {
		// 返回现有记录
		existingEntity, err := r.GetByFilters(ctx, filters...)
		return existingEntity, false, err
	}

	// 不存在则创建
	createdEntity, err := r.Create(ctx, entity)
	return createdEntity, true, err
}

// CreateOrUpdate 创建或更新记录
func (r *BaseRepository[T]) CreateOrUpdate(ctx context.Context, entity *T, uniqueFields ...string) (*T, bool, error) {
	if r.readOnly {
		return nil, false, errorx.NewError(errors.ErrorCodeForbidden)
	}

	existing, created, err := r.CreateIfNotExists(ctx, entity, uniqueFields...)
	if err != nil || created {
		return existing, created, err
	}

	// 如果存在则更新
	updatedEntity, err := r.Update(ctx, entity)
	return updatedEntity, false, err
}

// CreateBatch 批量创建记录
func (r *BaseRepository[T]) CreateBatch(ctx context.Context, entities ...*T) error {
	if r.readOnly {
		return errorx.NewError(errors.ErrorCodeForbidden)
	}

	if len(entities) == 0 {
		return nil
	}

	result := r.db.GetDB().WithContext(ctx).Table(r.table).CreateInBatches(entities, r.batchSize)
	return r.handleErrorWithContext(ctx, result.Error, "create batch")
}

// BulkCreate 高性能批量创建
func (r *BaseRepository[T]) BulkCreate(ctx context.Context, entities []*T, batchSize ...int) error {
	if r.readOnly {
		return errorx.NewError(errors.ErrorCodeForbidden)
	}

	if len(entities) == 0 {
		return nil
	}

	size := r.batchSize
	if len(batchSize) > 0 && batchSize[0] > 0 {
		size = batchSize[0]
	}

	// 分批处理
	for i := 0; i < len(entities); i += size {
		end := i + size
		if end > len(entities) {
			end = len(entities)
		}

		batch := entities[i:end]
		result := r.db.GetDB().WithContext(ctx).Table(r.table).Create(&batch)
		if result.Error != nil {
			return r.handleErrorWithContext(ctx, result.Error, fmt.Sprintf("bulk create batch %d-%d", i, end-1))
		}
	}

	return nil
}

// handleErrorWithContext 带上下文的错误处理
func (r *BaseRepository[T]) handleErrorWithContext(ctx context.Context, err error, operation string) error {
	if err == nil {
		return nil
	}

	// 使用 logger.WithContext 自动提取上下文字段
	contextLogger := r.logger.WithContext(ctx)

	// 添加操作相关字段
	fields := map[string]interface{}{
		"table":      r.table,
		"operation":  operation,
		"error_type": fmt.Sprintf("%T", err),
		"timestamp":  time.Now().Format(time.RFC3339),
	}

	// 区分错误级别：record not found 是INFO级别，其他是ERROR级别
	if err == gorm.ErrRecordNotFound {
		contextLogger.WithFields(fields).Infof("Record not found: %v", err)
	} else {
		contextLogger.WithFields(fields).Errorf("Repository operation failed: %v", err)
	}

	return err
}

// Get 获取单个记录
func (r *BaseRepository[T]) Get(ctx context.Context, id interface{}) (*T, error) {
	var entity T

	query := r.db.GetDB().WithContext(ctx).Table(r.table)

	// 应用默认预加载
	for _, preload := range r.preloads {
		query = query.Preload(preload)
	}

	result := query.Where("id = ?", id).First(&entity)

	if result.Error != nil {
		return nil, r.handleErrorWithContext(ctx, result.Error, "get by id")
	}

	return &entity, nil
}

// GetWithPreloads 获取单个记录并指定预加载关联
func (r *BaseRepository[T]) GetWithPreloads(ctx context.Context, id interface{}, preloads ...string) (*T, error) {
	var entity T

	query := r.db.GetDB().WithContext(ctx).Table(r.table)

	// 应用指定的预加载
	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	result := query.Where("id = ?", id).First(&entity)
	if result.Error != nil {
		return nil, r.handleErrorWithContext(ctx, result.Error, "get with preloads")
	}

	return &entity, nil
}

// GetByFields 根据多个字段获取记录
func (r *BaseRepository[T]) GetByFields(ctx context.Context, fields map[string]interface{}) (*T, error) {
	if len(fields) == 0 {
		return nil, errorx.NewError(errors.ErrorCodeInvalidInput)
	}

	filters := make([]*Filter, 0, len(fields))
	for field, value := range fields {
		filters = append(filters, NewEqFilter(field, value))
	}

	return r.GetByFilters(ctx, filters...)
}

// GetByFilter 按单个过滤条件获取记录
func (r *BaseRepository[T]) GetByFilter(ctx context.Context, filter *Filter) (*T, error) {
	if filter == nil {
		return nil, errorx.NewError(errors.ErrorCodeInvalidInput)
	}

	var entity T
	query := r.db.GetDB().WithContext(ctx).Table(r.table)
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
		return nil, errorx.NewError(errors.ErrorCodeInvalidInput)
	}

	var entity T
	query := r.db.GetDB().WithContext(ctx).Table(r.table)
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

	db := r.db.GetDB().WithContext(ctx).Table(r.table)

	// 应用字段选择
	db = r.applyFieldSelection(db, query)

	// 应用默认预加载
	for _, preload := range r.preloads {
		db = db.Preload(preload)
	}

	// 应用去重
	if query.Distinct {
		db = db.Distinct()
	}

	// 应用过滤条件
	db = r.applyFilters(db, query)

	// 应用分组
	for _, groupBy := range query.GroupBy {
		db = db.Group(groupBy)
	}

	// 应用HAVING条件
	for _, having := range query.Having {
		db = applyFilter(db, having)
	}

	// 应用排序
	db = r.applyOrdering(db, query)

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
		return nil, r.handleErrorWithContext(ctx, result.Error, "list")
	}

	return entities, nil
}

// ListWithPreloads 列表查询并指定预加载关联
func (r *BaseRepository[T]) ListWithPreloads(ctx context.Context, query *Query, preloads ...string) ([]*T, error) {
	if query == nil {
		query = NewQuery()
	}

	var entities []*T

	db := r.db.GetDB().WithContext(ctx).Table(r.table)

	// 应用字段选择
	db = r.applyFieldSelection(db, query)

	// 应用指定的预加载
	for _, preload := range preloads {
		db = db.Preload(preload)
	}

	// 应用过滤条件和其他操作
	db = r.applyFilters(db, query)
	db = r.applyOrdering(db, query)

	if query.LimitValue != nil {
		db = db.Limit(*query.LimitValue)
	}
	if query.OffsetValue != nil {
		db = db.Offset(*query.OffsetValue)
	}

	result := db.Find(&entities)
	if result.Error != nil {
		return nil, r.handleErrorWithContext(ctx, result.Error, "list with preloads")
	}

	return entities, nil
}

// ListWithPagination 分页列表查询
func (r *BaseRepository[T]) ListWithPagination(ctx context.Context, query *Query, page *Pagination) ([]*T, *Pagination, error) {
	if query == nil {
		query = NewQuery()
	}

	if page == nil {
		page = &Pagination{
			Page:     constants.DefaultPage,
			PageSize: constants.DefaultPageSize,
		}
	}

	// 参数校验和安全限制
	if page.Page <= 0 {
		page.Page = constants.DefaultPage
	}
	if page.PageSize <= 0 {
		page.PageSize = constants.DefaultPageSize
	}
	if page.PageSize < constants.MinPageSize {
		page.PageSize = constants.MinPageSize
	}
	if page.PageSize > constants.MaxPageSize {
		page.PageSize = constants.MaxPageSize
	}

	var entities []*T
	db := r.db.GetDB().WithContext(ctx).Table(r.table)

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
		return nil, errorx.NewError(errors.ErrorCodeInvalidInput)
	}

	result := r.db.GetDB().WithContext(ctx).Table(r.table).Save(entity)

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
	return r.db.GetDB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
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
		return errorx.NewError(errors.ErrorCodeInvalidInput)
	}

	if len(filters) == 0 {
		return errorx.NewError(errors.ErrorCodeInvalidInput)
	}

	db := r.db.GetDB().WithContext(ctx).Table(r.table)
	for _, filter := range filters {
		db = applyFilter(db, filter)
	}

	result := db.Updates(entity)
	return result.Error
}

// Delete 删除单个记录
func (r *BaseRepository[T]) Delete(ctx context.Context, id interface{}) error {
	result := r.db.GetDB().WithContext(ctx).Table(r.table).Where("id = ?", id).Delete(new(T))
	return result.Error
}

// DeleteBatch 批量删除记录
func (r *BaseRepository[T]) DeleteBatch(ctx context.Context, ids ...interface{}) error {
	if len(ids) == 0 {
		return nil
	}

	result := r.db.GetDB().WithContext(ctx).Table(r.table).Where("id IN ?", ids).Delete(new(T))
	return result.Error
}

// DeleteByFilters 按过滤条件删除记录
func (r *BaseRepository[T]) DeleteByFilters(ctx context.Context, filters ...*Filter) error {
	if len(filters) == 0 {
		return errorx.NewError(errors.ErrorCodeInvalidInput)
	}

	db := r.db.GetDB().WithContext(ctx).Table(r.table)
	for _, filter := range filters {
		db = applyFilter(db, filter)
	}

	result := db.Delete(new(T))
	return result.Error
}

// Transaction 事务支持
func (r *BaseRepository[T]) Transaction(ctx context.Context, fn func(tx Transaction[T]) error) error {
	return r.db.GetDB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txHandler, err := db.NewGormHandler(tx)
		if err != nil {
			return err
		}
		txWrapper := &transactionWrapper[T]{db: txHandler, table: r.table}
		return fn(txWrapper)
	})
}

// TransactionWithRawDB 通用事务支持（可以操作任意模型）
func (r *BaseRepository[T]) TransactionWithRawDB(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return r.db.GetDB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(tx.WithContext(ctx))
	})
}

// Count 计数
func (r *BaseRepository[T]) Count(ctx context.Context, filters ...*Filter) (int64, error) {
	var count int64
	db := r.db.GetDB().WithContext(ctx).Table(r.table)

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
	result := r.db.GetDB().WithContext(ctx).Table(r.table).Find(&entities)
	if result.Error != nil {
		return nil, result.Error
	}
	return entities, nil
}

// First 获取第一条记录
func (r *BaseRepository[T]) First(ctx context.Context, filters ...*Filter) (*T, error) {
	var entity T
	db := r.db.GetDB().WithContext(ctx).Table(r.table)

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
	db := r.db.GetDB().WithContext(ctx).Table(r.table)

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
	db := r.db.GetDB().WithContext(ctx).Table(r.table)

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

	result := r.db.GetDB().WithContext(ctx).Table(r.table).Where("id = ?", id).Updates(fields)
	return result.Error
}

// UpdateFieldsByFilters 按过滤条件更新指定字段
func (r *BaseRepository[T]) UpdateFieldsByFilters(ctx context.Context, fields map[string]interface{}, filters ...*Filter) error {
	if len(fields) == 0 {
		return nil
	}

	if len(filters) == 0 {
		return errorx.NewError(errors.ErrorCodeInvalidInput)
	}

	db := r.db.GetDB().WithContext(ctx).Table(r.table)
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
	result := r.db.GetDB().WithContext(ctx).Table(r.table).Where("id = ?", id).Update(field, value)
	return result.Error
}

// SoftDeleteBatch 批量软删除
// field: 软删除字段名，如 "deleted_at", "is_deleted" 等
// value: 软删除标记值，如 time.Now(), 1 等
func (r *BaseRepository[T]) SoftDeleteBatch(ctx context.Context, ids []interface{}, field string, value interface{}) error {
	if len(ids) == 0 {
		return nil
	}

	result := r.db.GetDB().WithContext(ctx).Table(r.table).Where("id IN ?", ids).Update(field, value)
	return result.Error
}

// SoftDeleteByFilters 按过滤条件软删除
// field: 软删除字段名，如 "deleted_at", "is_deleted" 等
// value: 软删除标记值，如 time.Now(), 1 等
func (r *BaseRepository[T]) SoftDeleteByFilters(ctx context.Context, field string, value interface{}, filters ...*Filter) error {
	if len(filters) == 0 {
		return errorx.NewError(errors.ErrorCodeInvalidInput)
	}

	db := r.db.GetDB().WithContext(ctx).Table(r.table)
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
	result := r.db.GetDB().WithContext(ctx).Table(r.table).Where("id = ?", id).Update(field, restoreValue)
	return result.Error
}

// RestoreBatch 批量恢复软删除的记录
// field: 软删除字段名，如 "deleted_at", "is_deleted" 等
// restoreValue: 恢复时的值，如 nil, 0 等
func (r *BaseRepository[T]) RestoreBatch(ctx context.Context, ids []interface{}, field string, restoreValue interface{}) error {
	if len(ids) == 0 {
		return nil
	}

	result := r.db.GetDB().WithContext(ctx).Table(r.table).Where("id IN ?", ids).Update(field, restoreValue)
	return result.Error
}

// CountByField 按字段计数（GROUP BY）
func (r *BaseRepository[T]) CountByField(ctx context.Context, field string) (map[interface{}]int64, error) {
	rows, err := r.db.GetDB().WithContext(ctx).Table(r.table).
		Select(field + ", COUNT(*) as count").
		Group(field).Rows()
	if err != nil {
		return nil, r.handleErrorWithContext(ctx, err, "count by field")
	}
	defer rows.Close()

	countMap := make(map[interface{}]int64)
	for rows.Next() {
		var fieldValue interface{}
		var count int64
		if err := rows.Scan(&fieldValue, &count); err != nil {
			return nil, r.handleErrorWithContext(ctx, err, "scan count result")
		}
		countMap[fieldValue] = count
	}

	if err := rows.Err(); err != nil {
		return nil, r.handleErrorWithContext(ctx, err, "iterate count results")
	}

	return countMap, nil
}

// Pluck 提取单个字段的值列表
func (r *BaseRepository[T]) Pluck(ctx context.Context, field string, filters ...*Filter) ([]interface{}, error) {
	var values []interface{}
	db := r.db.GetDB().WithContext(ctx).Table(r.table)

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
	db := r.db.GetDB().WithContext(ctx).Table(r.table).Distinct(field)

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
	case constants.OP_EQ, constants.OP_NEQ, constants.OP_GT, constants.OP_GTE,
		constants.OP_LT, constants.OP_LTE, constants.OP_IN, constants.OP_NOT_IN,
		constants.OP_LIKE, constants.OP_NOT_LIKE:
		return dbQuery.Where(fmt.Sprintf("%s %s ?", filter.Field, string(filter.Operator)), filter.Value)
	case constants.OP_BETWEEN:
		values, ok := filter.Value.([]interface{})
		if ok && len(values) == 2 {
			return dbQuery.Where(fmt.Sprintf("%s BETWEEN ? AND ?", filter.Field), values[0], values[1])
		}
	case constants.OP_IS_NULL:
		return dbQuery.Where(fmt.Sprintf("%s IS NULL", filter.Field))
	case constants.OP_IS_NOT_NULL:
		return dbQuery.Where(fmt.Sprintf("%s IS NOT NULL", filter.Field))
	case constants.OP_FIND_IN_SET:
		// 修复参数顺序：FIND_IN_SET(value, field_list)
		return dbQuery.Where("FIND_IN_SET(?, ?)", filter.Value, filter.Field)
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

// EnableAutoFields 启用自动字段模式
func (r *BaseRepository[T]) EnableAutoFields() {
	if !r.autoFields {
		r.autoFields = true
		// 提取并缓存字段（如果尚未缓存）
		if len(r.modelFields) == 0 {
			var model T
			r.modelFields = GetStructFields(model)
		}
	}
}

// DisableAutoFields 禁用自动字段模式
func (r *BaseRepository[T]) DisableAutoFields() {
	r.autoFields = false
}

// IsAutoFieldsEnabled 检查是否启用自动字段模式
func (r *BaseRepository[T]) IsAutoFieldsEnabled() bool {
	return r.autoFields
}

// GetModelFields 获取缓存的模型字段
func (r *BaseRepository[T]) GetModelFields() []string {
	if len(r.modelFields) == 0 {
		var model T
		r.modelFields = GetStructFields(model)
	}
	return r.modelFields
}

// transactionWrapper 事务包装器
type transactionWrapper[T any] struct {
	db    db.Handler
	table string
}

// Create 在事务中创建
func (t *transactionWrapper[T]) Create(ctx context.Context, entity *T) error {
	return t.db.GetDB().WithContext(ctx).Table(t.table).Create(entity).Error
}

// CreateBatch 在事务中批量创建
func (t *transactionWrapper[T]) CreateBatch(ctx context.Context, entities ...*T) error {
	if len(entities) == 0 {
		return nil
	}

	// 过滤掉nil实体
	var validEntities []*T
	for _, entity := range entities {
		if entity != nil {
			validEntities = append(validEntities, entity)
		}
	}

	if len(validEntities) == 0 {
		return nil
	}

	// 使用GORM的CreateInBatches进行真正的批量插入
	return t.db.GetDB().WithContext(ctx).Table(t.table).CreateInBatches(validEntities, len(validEntities)).Error
}

// Update 在事务中更新
func (t *transactionWrapper[T]) Update(ctx context.Context, entity *T) error {
	return t.db.GetDB().WithContext(ctx).Table(t.table).Save(entity).Error
}

// UpdateBatch 在事务中批量更新
func (t *transactionWrapper[T]) UpdateBatch(ctx context.Context, entities ...*T) error {
	if len(entities) == 0 {
		return nil
	}

	// 在事务中执行批量更新
	db := t.db.GetDB().WithContext(ctx)
	for _, entity := range entities {
		if entity != nil {
			if err := db.Table(t.table).Save(entity).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

// Delete 在事务中删除
func (t *transactionWrapper[T]) Delete(ctx context.Context, entity *T) error {
	return t.db.GetDB().WithContext(ctx).Table(t.table).Delete(entity).Error
}

// DeleteBatch 在事务中批量删除
func (t *transactionWrapper[T]) DeleteBatch(ctx context.Context, entities ...*T) error {
	if len(entities) == 0 {
		return nil
	}

	// 使用GORM的Delete方法进行批量删除
	// GORM会自动根据主键构建WHERE条件
	db := t.db.GetDB().WithContext(ctx)
	for _, entity := range entities {
		if entity != nil {
			if err := db.Table(t.table).Delete(entity).Error; err != nil {
				return err
			}
		}
	}
	return nil
} // applyFilters 应用过滤条件到查询
func (r *BaseRepository[T]) applyFilters(db *gorm.DB, query *Query) *gorm.DB {
	// 应用简单过滤条件
	for _, filter := range query.Filters {
		db = applyFilter(db, filter)
	}

	// 应用复合过滤条件组
	if query.FilterGroup != nil {
		db = r.applyFilterGroup(db, query.FilterGroup)
	}

	return db
}

// applyFilterGroup 应用过滤条件组
func (r *BaseRepository[T]) applyFilterGroup(db *gorm.DB, group *FilterGroup) *gorm.DB {
	if group == nil || group.IsEmpty() {
		return db
	}

	isOr := group.LogicOp == constants.LOGIC_OR

	// 构建条件表达式
	if isOr {
		// OR 逻辑：使用Where("condition1 OR condition2 OR ...", args...)
		var conditions []string
		var args []interface{}

		// 处理过滤条件
		for _, filter := range group.Filters {
			if filter != nil {
				condition, arg := buildFilterCondition(filter)
				if condition != "" {
					conditions = append(conditions, condition)
					if arg != nil {
						args = append(args, arg)
					}
				}
			}
		}

		// 处理子组
		for _, subGroup := range group.Groups {
			if subGroup != nil && !subGroup.IsEmpty() {
				// 递归构建子组条件
				subConditions, subArgs := buildGroupCondition(subGroup)
				if subConditions != "" {
					conditions = append(conditions, "("+subConditions+")")
					args = append(args, subArgs...)
				}
			}
		}

		if len(conditions) > 0 {
			combinedCondition := strings.Join(conditions, " OR ")
			db = db.Where(combinedCondition, args...)
		}
	} else {
		// AND 逻辑：逐个应用条件
		for _, filter := range group.Filters {
			if filter != nil {
				db = applyFilter(db, filter)
			}
		}

		// 递归处理子组
		for _, subGroup := range group.Groups {
			if subGroup != nil {
				db = r.applyFilterGroup(db, subGroup)
			}
		}
	}

	return db
}

// buildFilterCondition 构建单个过滤条件的SQL和参数
func buildFilterCondition(filter *Filter) (string, interface{}) {
	if filter == nil {
		return "", nil
	}

	// 处理特殊情况
	switch filter.Operator {
	case constants.OP_IS_NULL, constants.OP_IS_NOT_NULL:
		if template, ok := constants.OperatorTemplateMap[filter.Operator]; ok {
			return fmt.Sprintf(template, filter.Field), nil
		}
	case constants.OP_BETWEEN:
		if values, ok := filter.Value.([]interface{}); ok && len(values) == 2 {
			return fmt.Sprintf(constants.SQL_BETWEEN, filter.Field), values
		}
		return "", nil
	case constants.OP_STARTS_WITH:
		// STARTS_WITH 转换为 LIKE 'value%'
		if valueStr, ok := filter.Value.(string); ok {
			return fmt.Sprintf(constants.SQL_LIKE, filter.Field), valueStr + constants.SQL_WILDCARD_ANY
		}
		return "", nil
	case constants.OP_ENDS_WITH:
		// ENDS_WITH 转换为 LIKE '%value'
		if valueStr, ok := filter.Value.(string); ok {
			return fmt.Sprintf(constants.SQL_LIKE, filter.Field), constants.SQL_WILDCARD_ANY + valueStr
		}
		return "", nil
	case constants.OP_CONTAINS:
		// CONTAINS 转换为 LIKE '%value%'
		if valueStr, ok := filter.Value.(string); ok {
			return fmt.Sprintf(constants.SQL_LIKE, filter.Field), constants.SQL_WILDCARD_ANY + valueStr + constants.SQL_WILDCARD_ANY
		}
		return "", nil
	case constants.OP_FIND_IN_SET:
		// FIND_IN_SET 需要特殊处理参数位置
		if template, ok := constants.OperatorTemplateMap[filter.Operator]; ok {
			return fmt.Sprintf(template, filter.Field) + " > 0", filter.Value
		}
		return "", nil
	}

	// 通用处理：使用 map 查找模板
	if template, ok := constants.OperatorTemplateMap[filter.Operator]; ok {
		return fmt.Sprintf(template, filter.Field), filter.Value
	}

	return "", nil
}

// buildGroupCondition 递归构建过滤组的条件
func buildGroupCondition(group *FilterGroup) (string, []interface{}) {
	if group == nil || group.IsEmpty() {
		return "", nil
	}

	var conditions []string
	var args []interface{}

	// 处理过滤条件
	for _, filter := range group.Filters {
		if filter != nil {
			condition, arg := buildFilterCondition(filter)
			if condition != "" {
				conditions = append(conditions, condition)
				if arg != nil {
					// 处理BETWEEN操作的多个参数
					if values, ok := arg.([]interface{}); ok {
						args = append(args, values...)
					} else {
						args = append(args, arg)
					}
				}
			}
		}
	}

	// 递归处理子组
	for _, subGroup := range group.Groups {
		if subGroup != nil && !subGroup.IsEmpty() {
			subCondition, subArgs := buildGroupCondition(subGroup)
			if subCondition != "" {
				conditions = append(conditions, "("+subCondition+")")
				args = append(args, subArgs...)
			}
		}
	}

	if len(conditions) == 0 {
		return "", nil
	}

	// 根据逻辑操作符连接条件
	separator := " AND "
	if group.LogicOp == constants.LOGIC_OR {
		separator = " OR "
	}

	return strings.Join(conditions, separator), args
}

// applyOrdering 应用排序条件
func (r *BaseRepository[T]) applyOrdering(db *gorm.DB, query *Query) *gorm.DB {
	// 应用查询中的排序条件
	for _, order := range query.Orders {
		if order.Field != "" {
			orderClause := func() string {
				if order.Direction != "" {
					return order.Field + " " + order.Direction
				}
				return order.Field
			}()
			db = db.Order(orderClause)
		}
	}

	// 如果没有排序条件且有默认排序，应用默认排序
	if len(query.Orders) == 0 && r.defaultOrder != "" {
		db = db.Order(r.defaultOrder)
	}

	return db
}

// applyFieldSelection 应用字段选择（Select/Omit）
func (r *BaseRepository[T]) applyFieldSelection(db *gorm.DB, query *Query) *gorm.DB {
	// 如果同时指定了 Select 和 Omit，优先使用 Select
	if len(query.SelectFields) > 0 {
		db = db.Select(query.SelectFields)
		return db
	}

	// 如果指定了 Omit
	if len(query.OmitFields) > 0 {
		// 如果启用了自动字段模式，从模型字段中排除
		if r.autoFields && len(r.modelFields) > 0 {
			selectedFields := FilterFields(r.modelFields, nil, query.OmitFields)
			if len(selectedFields) > 0 {
				db = db.Select(selectedFields)
			}
			return db
		}
		// 否则使用GORM的Omit
		db = db.Omit(query.OmitFields...)
		return db
	}

	// 如果启用了自动字段模式且没有指定任何字段选择，使用缓存的模型字段
	if r.autoFields && len(r.modelFields) > 0 {
		db = db.Select(r.modelFields)
		return db
	}

	// 没有指定字段选择，返回原查询（SELECT *）
	return db
}
