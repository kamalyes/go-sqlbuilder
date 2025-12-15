/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-11 21:13:15
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-04 09:15:32
 * @FilePath: \go-sqlbuilder\repository\base.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package repository

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/kamalyes/go-logger"
	"github.com/kamalyes/go-sqlbuilder/constants"
	"github.com/kamalyes/go-sqlbuilder/db"
	"github.com/kamalyes/go-sqlbuilder/errors"
	"github.com/kamalyes/go-toolbox/pkg/errorx"
	"gorm.io/gorm"
)

// ContextFieldExtractor context字段提取器函数类型
// 用于从context中提取需要记录到日志的字段
type ContextFieldExtractor func(ctx context.Context, log logger.ILogger) logger.ILogger

// BaseRepository 基础仓储实现，包含通用的 CRUD 操作
type BaseRepository[T any] struct {
	db                    db.Handler
	table                 string
	batchSize             int            // 批处理大小
	timeout               int            // 查询超时时间（秒）
	readOnly              bool           // 只读模式
	preloads              []string       // 默认预加载关联
	defaultOrder          string         // 默认排序
	logger                logger.ILogger // 日志记录器
	primaryKeyIndexes     []int          // 主键字段索引缓存
	autoCreateTimeIndexes []int          // 创建时间字段索引缓存
	autoUpdateTimeIndexes []int          // 更新时间字段索引缓存
	modelFields           []string       // 模型字段缓存（用于自动字段选择）
	autoFields            bool           // 是否启用自动字段模式
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
func WithLogger[T any](log logger.ILogger) RepositoryOption[T] {
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
func NewBaseRepository[T any](dbHandler db.Handler, logger logger.ILogger, table string, options ...RepositoryOption[T]) *BaseRepository[T] {
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

	// 初始化字段索引缓存
	r.initFieldIndexes()

	return r
}

// initFieldIndexes 初始化字段索引缓存（主键、创建时间、更新时间）
func (r *BaseRepository[T]) initFieldIndexes() {
	var model T
	entityType := reflect.TypeOf(model)
	if entityType.Kind() == reflect.Ptr {
		entityType = entityType.Elem()
	}

	for i := 0; i < entityType.NumField(); i++ {
		field := entityType.Field(i)
		gormTag := field.Tag.Get("gorm")

		if strings.Contains(gormTag, "primaryKey") || strings.Contains(gormTag, "primary_key") {
			r.primaryKeyIndexes = append(r.primaryKeyIndexes, i)
		}
		if strings.Contains(gormTag, "autoCreateTime") {
			r.autoCreateTimeIndexes = append(r.autoCreateTimeIndexes, i)
		}
		if strings.Contains(gormTag, "autoUpdateTime") {
			r.autoUpdateTimeIndexes = append(r.autoUpdateTimeIndexes, i)
		}
	}
}

// ========== 辅助方法 ==========

// newDB 创建带上下文和表名的DB实例
func (r *BaseRepository[T]) newDB(ctx context.Context) *gorm.DB {
	return r.db.GetDB().WithContext(ctx).Table(r.table)
}

// checkReadOnly 检查是否为只读模式
func (r *BaseRepository[T]) checkReadOnly() error {
	if r.readOnly {
		return errorx.NewError(errors.ErrorCodeForbidden)
	}
	return nil
}

// checkEntity 检查实体是否为空
func (r *BaseRepository[T]) checkEntity(entity *T) error {
	if entity == nil {
		return errorx.NewError(errors.ErrorCodeInvalidInput)
	}
	return nil
}

// ApplyFilters 批量应用过滤器到查询
func ApplyFilters(db *gorm.DB, filters []*Filter) *gorm.DB {
	for _, filter := range filters {
		db = ApplyFilter(db, filter)
	}
	return db
}

// ApplyOrders 批量应用排序条件到查询
func ApplyOrders(db *gorm.DB, orders []Order) *gorm.DB {
	for _, order := range orders {
		db = db.Order(order.Field + " " + order.Direction)
	}
	return db
}

// ApplyPreloads 应用预加载关联
func ApplyPreloads(db *gorm.DB, preloads []string) *gorm.DB {
	for _, preload := range preloads {
		db = db.Preload(preload)
	}
	return db
}

// filterNilEntities 过滤nil实体
func filterNilEntities[T any](entities []*T) []*T {
	if len(entities) == 0 {
		return nil
	}
	valid := make([]*T, 0, len(entities))
	for _, entity := range entities {
		if entity != nil {
			valid = append(valid, entity)
		}
	}
	return valid
}

// copyFieldsByIndexes 按索引复制字段值
func (r *BaseRepository[T]) copyFieldsByIndexes(src, dst reflect.Value, indexes []int) {
	for _, idx := range indexes {
		srcField := src.Field(idx)
		dstField := dst.Field(idx)
		if srcField.IsValid() && dstField.CanSet() {
			dstField.Set(srcField)
		}
	}
}

// buildFiltersFromFields 从字段map构建过滤器
func (r *BaseRepository[T]) buildFiltersFromFields(entity *T, uniqueFields []string) []*Filter {
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
	return filters
}

// resetAutoUpdateTime 重置自动更新时间字段
func (r *BaseRepository[T]) resetAutoUpdateTime(entity *T) {
	if len(r.autoUpdateTimeIndexes) == 0 {
		return
	}
	entityValue := reflect.ValueOf(entity).Elem()
	entityType := entityValue.Type()
	for _, idx := range r.autoUpdateTimeIndexes {
		field := entityType.Field(idx)
		entityField := entityValue.Field(idx)
		if entityField.CanSet() {
			entityField.Set(reflect.Zero(field.Type))
		}
	}
}

// ========== CRUD 操作 ==========

// Create 创建单个记录
func (r *BaseRepository[T]) Create(ctx context.Context, entity *T) (*T, error) {
	if err := r.checkReadOnly(); err != nil {
		return nil, err
	}
	if err := r.checkEntity(entity); err != nil {
		return nil, err
	}

	if result := r.newDB(ctx).Create(entity); result.Error != nil {
		return nil, r.handleErrorWithContext(ctx, result.Error, "create")
	}
	return entity, nil
}

// CreateIfNotExists 如果不存在则创建
func (r *BaseRepository[T]) CreateIfNotExists(ctx context.Context, entity *T, uniqueFields ...string) (*T, bool, error) {
	if err := r.checkReadOnly(); err != nil {
		return nil, false, err
	}
	if entity == nil || len(uniqueFields) == 0 {
		return nil, false, errorx.NewError(errors.ErrorCodeInvalidInput)
	}

	// 构建查询条件检查是否存在
	filters := r.buildFiltersFromFields(entity, uniqueFields)

	// 检查是否存在
	if exists, err := r.Exists(ctx, filters...); err != nil {
		return nil, false, err
	} else if exists {
		existingEntity, err := r.GetByFilters(ctx, filters...)
		return existingEntity, false, err
	}

	// 不存在则创建
	createdEntity, err := r.Create(ctx, entity)
	return createdEntity, true, err
}

// CreateOrUpdate 创建或更新记录
func (r *BaseRepository[T]) CreateOrUpdate(ctx context.Context, entity *T, uniqueFields ...string) (*T, bool, error) {
	if err := r.checkReadOnly(); err != nil {
		return nil, false, err
	}

	existing, created, err := r.CreateIfNotExists(ctx, entity, uniqueFields...)
	if err != nil || created {
		return existing, created, err
	}

	// 复制关键字段（主键、创建时间）避免 GORM 误判
	if len(r.primaryKeyIndexes) > 0 || len(r.autoCreateTimeIndexes) > 0 {
		existingValue := reflect.ValueOf(existing).Elem()
		entityValue := reflect.ValueOf(entity).Elem()
		r.copyFieldsByIndexes(existingValue, entityValue, r.primaryKeyIndexes)
		r.copyFieldsByIndexes(existingValue, entityValue, r.autoCreateTimeIndexes)
	}

	updatedEntity, err := r.Update(ctx, entity)
	return updatedEntity, false, err
}

// CreateBatch 批量创建记录
func (r *BaseRepository[T]) CreateBatch(ctx context.Context, entities ...*T) error {
	if err := r.checkReadOnly(); err != nil {
		return err
	}
	if len(entities) == 0 {
		return nil
	}

	result := r.newDB(ctx).CreateInBatches(entities, r.batchSize)
	return r.handleErrorWithContext(ctx, result.Error, "create batch")
}

// BulkCreate 高性能批量创建
func (r *BaseRepository[T]) BulkCreate(ctx context.Context, entities []*T, batchSize ...int) error {
	if err := r.checkReadOnly(); err != nil {
		return err
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
		end := min(i+size, len(entities))
		batch := entities[i:end]
		if result := r.newDB(ctx).Create(&batch); result.Error != nil {
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
	query := ApplyPreloads(r.newDB(ctx), r.preloads)
	if result := query.Where("id = ?", id).First(&entity); result.Error != nil {
		return nil, r.handleErrorWithContext(ctx, result.Error, "get by id")
	}
	return &entity, nil
}

// GetWithPreloads 获取单个记录并指定预加载关联
func (r *BaseRepository[T]) GetWithPreloads(ctx context.Context, id interface{}, preloads ...string) (*T, error) {
	var entity T
	query := ApplyPreloads(r.newDB(ctx), preloads)
	if result := query.Where("id = ?", id).First(&entity); result.Error != nil {
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
	return r.GetByFilters(ctx, filter)
}

// GetByFilters 按多个过滤条件获取记录
func (r *BaseRepository[T]) GetByFilters(ctx context.Context, filters ...*Filter) (*T, error) {
	if len(filters) == 0 {
		return nil, errorx.NewError(errors.ErrorCodeInvalidInput)
	}
	var entity T
	query := ApplyFilters(r.newDB(ctx), filters)
	if result := query.First(&entity); result.Error != nil {
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
	db = ApplyFieldSelection(db, query.SelectFields, query.OmitFields, r.modelFields, r.autoFields)

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
		db = ApplyFilter(db, having)
	}

	// 应用排序
	db = ApplyOrdering(db, query.Orders, r.defaultOrder)

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
	db = ApplyFieldSelection(db, query.SelectFields, query.OmitFields, r.modelFields, r.autoFields)

	// 应用指定的预加载
	for _, preload := range preloads {
		db = db.Preload(preload)
	}

	// 应用过滤条件和其他操作
	db = r.applyFilters(db, query)
	db = ApplyOrdering(db, query.Orders, r.defaultOrder)

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
		db = ApplyFilter(db, filter)
	}

	// 计算总数
	var total int64
	countDb := db
	countDb.Model(new(T)).Count(&total)
	page.Total = total

	// 如果没有数据，直接返回空结果
	if total == 0 {
		return []*T{}, page, nil
	}

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
	if err := r.checkEntity(entity); err != nil {
		return nil, err
	}

	// 重置自动更新时间字段让 GORM 自动填充
	r.resetAutoUpdateTime(entity)

	if result := r.newDB(ctx).Save(entity); result.Error != nil {
		return nil, result.Error
	}
	return entity, nil
}

// UpdateBatch 批量更新记录
func (r *BaseRepository[T]) UpdateBatch(ctx context.Context, entities ...*T) error {
	if len(entities) == 0 {
		return nil
	}
	return r.db.GetDB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, entity := range entities {
			if entity != nil {
				if err := tx.Table(r.table).Save(entity).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

// UpdateByFilters 按过滤条件更新记录
func (r *BaseRepository[T]) UpdateByFilters(ctx context.Context, entity *T, filters ...*Filter) error {
	if err := r.checkEntity(entity); err != nil {
		return err
	}
	if len(filters) == 0 {
		return errorx.NewError(errors.ErrorCodeInvalidInput)
	}
	return ApplyFilters(r.newDB(ctx), filters).Updates(entity).Error
}

// Delete 删除单个记录
func (r *BaseRepository[T]) Delete(ctx context.Context, id interface{}) error {
	return r.newDB(ctx).Where("id = ?", id).Delete(new(T)).Error
}

// DeleteBatch 批量删除记录
func (r *BaseRepository[T]) DeleteBatch(ctx context.Context, ids ...interface{}) error {
	if len(ids) == 0 {
		return nil
	}
	return r.newDB(ctx).Where("id IN ?", ids).Delete(new(T)).Error
}

// DeleteByFilters 按过滤条件删除记录
func (r *BaseRepository[T]) DeleteByFilters(ctx context.Context, filters ...*Filter) error {
	if len(filters) == 0 {
		return errorx.NewError(errors.ErrorCodeInvalidInput)
	}
	return ApplyFilters(r.newDB(ctx), filters).Delete(new(T)).Error
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
	return count, ApplyFilters(r.newDB(ctx), filters).Count(&count).Error
}

// Exists 检查记录是否存在
func (r *BaseRepository[T]) Exists(ctx context.Context, filters ...*Filter) (bool, error) {
	count, err := r.Count(ctx, filters...)
	return count > 0, err
}

// GetAll 获取所有记录（不分页）
func (r *BaseRepository[T]) GetAll(ctx context.Context) ([]*T, error) {
	var entities []*T
	if result := r.newDB(ctx).Find(&entities); result.Error != nil {
		return nil, result.Error
	}
	return entities, nil
}

// First 获取第一条记录
func (r *BaseRepository[T]) First(ctx context.Context, filters ...*Filter) (*T, error) {
	var entity T
	if result := ApplyFilters(r.newDB(ctx), filters).First(&entity); result.Error != nil {
		return nil, result.Error
	}
	return &entity, nil
}

// Last 获取最后一条记录
func (r *BaseRepository[T]) Last(ctx context.Context, filters ...*Filter) (*T, error) {
	var entity T
	if result := ApplyFilters(r.newDB(ctx), filters).Last(&entity); result.Error != nil {
		return nil, result.Error
	}
	return &entity, nil
}

// FindOne 查找单条记录（不存在返回 nil）
func (r *BaseRepository[T]) FindOne(ctx context.Context, filters ...*Filter) (*T, error) {
	var entity T
	result := ApplyFilters(r.newDB(ctx), filters).Limit(1).Find(&entity)
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
	return r.newDB(ctx).Where("id = ?", id).Updates(fields).Error
}

// UpdateFieldsByFilters 按过滤条件更新指定字段
func (r *BaseRepository[T]) UpdateFieldsByFilters(ctx context.Context, fields map[string]interface{}, filters ...*Filter) error {
	if len(fields) == 0 {
		return nil
	}
	if len(filters) == 0 {
		return errorx.NewError(errors.ErrorCodeInvalidInput)
	}
	return ApplyFilters(r.newDB(ctx), filters).Updates(fields).Error
}

// SoftDelete 软删除（需要指定删除标记字段和值）
// field: 软删除字段名，如 "deleted_at", "is_deleted" 等
// value: 软删除标记值，如 time.Now(), 1 等
func (r *BaseRepository[T]) SoftDelete(ctx context.Context, id interface{}, field string, value interface{}) error {
	return r.newDB(ctx).Where("id = ?", id).Update(field, value).Error
}

// SoftDeleteBatch 批量软删除
// field: 软删除字段名，如 "deleted_at", "is_deleted" 等
// value: 软删除标记值，如 time.Now(), 1 等
func (r *BaseRepository[T]) SoftDeleteBatch(ctx context.Context, ids []interface{}, field string, value interface{}) error {
	if len(ids) == 0 {
		return nil
	}
	return r.newDB(ctx).Where("id IN ?", ids).Update(field, value).Error
}

// SoftDeleteByFilters 按过滤条件软删除
// field: 软删除字段名，如 "deleted_at", "is_deleted" 等
// value: 软删除标记值，如 time.Now(), 1 等
func (r *BaseRepository[T]) SoftDeleteByFilters(ctx context.Context, field string, value interface{}, filters ...*Filter) error {
	if len(filters) == 0 {
		return errorx.NewError(errors.ErrorCodeInvalidInput)
	}
	return ApplyFilters(r.newDB(ctx), filters).Update(field, value).Error
}

// Restore 恢复软删除的记录
// field: 软删除字段名，如 "deleted_at", "is_deleted" 等
// restoreValue: 恢复时的值，如 nil, 0 等
func (r *BaseRepository[T]) Restore(ctx context.Context, id interface{}, field string, restoreValue interface{}) error {
	return r.newDB(ctx).Where("id = ?", id).Update(field, restoreValue).Error
}

// RestoreBatch 批量恢复软删除的记录
// field: 软删除字段名，如 "deleted_at", "is_deleted" 等
// restoreValue: 恢复时的值，如 nil, 0 等
func (r *BaseRepository[T]) RestoreBatch(ctx context.Context, ids []interface{}, field string, restoreValue interface{}) error {
	if len(ids) == 0 {
		return nil
	}
	return r.newDB(ctx).Where("id IN ?", ids).Update(field, restoreValue).Error
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
	err := ApplyFilters(r.newDB(ctx), filters).Pluck(field, &values).Error
	return values, err
}

// Distinct 获取去重后的字段值列表
func (r *BaseRepository[T]) Distinct(ctx context.Context, field string, filters ...*Filter) ([]interface{}, error) {
	var values []interface{}
	err := ApplyFilters(r.newDB(ctx).Distinct(field), filters).Pluck(field, &values).Error
	return values, err
}

// ApplyFilter 应用单个过滤条件到 GORM 查询
func ApplyFilter(dbQuery *gorm.DB, filter *Filter) *gorm.DB {
	if filter == nil {
		return dbQuery
	}

	// 检查是否为子查询
	if subQuery, ok := filter.Value.(*SubQuery); ok {
		// 处理子查询情况
		return dbQuery.Where(fmt.Sprintf("%s %s (%s)", filter.Field, string(filter.Operator), subQuery.SQL), subQuery.Args...)
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
	validEntities := filterNilEntities(entities)
	if len(validEntities) == 0 {
		return nil
	}
	return t.db.GetDB().WithContext(ctx).Table(t.table).CreateInBatches(validEntities, len(validEntities)).Error
}

// Update 在事务中更新
func (t *transactionWrapper[T]) Update(ctx context.Context, entity *T) error {
	return t.db.GetDB().WithContext(ctx).Table(t.table).Save(entity).Error
}

// UpdateBatch 在事务中批量更新
func (t *transactionWrapper[T]) UpdateBatch(ctx context.Context, entities ...*T) error {
	validEntities := filterNilEntities(entities)
	if len(validEntities) == 0 {
		return nil
	}
	db := t.db.GetDB().WithContext(ctx)
	for _, entity := range validEntities {
		if err := db.Table(t.table).Save(entity).Error; err != nil {
			return err
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
	validEntities := filterNilEntities(entities)
	if len(validEntities) == 0 {
		return nil
	}
	db := t.db.GetDB().WithContext(ctx)
	for _, entity := range validEntities {
		if err := db.Table(t.table).Delete(entity).Error; err != nil {
			return err
		}
	}
	return nil
}

// applyFilters 应用过滤条件到查询
func (r *BaseRepository[T]) applyFilters(db *gorm.DB, query *Query) *gorm.DB {
	// 应用简单过滤条件
	for _, filter := range query.Filters {
		db = ApplyFilter(db, filter)
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

	if group.LogicOp == constants.LOGIC_OR {
		return r.applyOrFilterGroup(db, group)
	}
	return r.applyAndFilterGroup(db, group)
}

// applyOrFilterGroup 应用 OR 逻辑的过滤组
func (r *BaseRepository[T]) applyOrFilterGroup(db *gorm.DB, group *FilterGroup) *gorm.DB {
	var conditions []string
	var args []interface{}

	// 收集所有条件（过滤条件 + 子组）
	collectFilterConditions(group.Filters, &conditions, &args)
	collectSubGroupConditions(group.Groups, &conditions, &args)

	if len(conditions) == 0 {
		return db
	}

	combinedCondition := strings.Join(conditions, " OR ")
	return db.Where(combinedCondition, args...)
}

// applyAndFilterGroup 应用 AND 逻辑的过滤组
func (r *BaseRepository[T]) applyAndFilterGroup(db *gorm.DB, group *FilterGroup) *gorm.DB {
	// AND 逻辑：逐个应用过滤条件
	for _, filter := range group.Filters {
		if filter != nil {
			db = ApplyFilter(db, filter)
		}
	}

	// 递归处理子组
	for _, subGroup := range group.Groups {
		if subGroup != nil {
			db = r.applyFilterGroup(db, subGroup)
		}
	}

	return db
}

// collectFilterConditions 收集过滤条件（通用函数）
func collectFilterConditions(filters []*Filter, conditions *[]string, args *[]interface{}) {
	for _, filter := range filters {
		if filter == nil {
			continue
		}
		condition, arg := buildFilterCondition(filter)
		if condition != "" {
			*conditions = append(*conditions, condition)
			if arg != nil {
				*args = append(*args, arg)
			}
		}
	}
}

// collectSubGroupConditions 收集子组条件（通用函数）
func collectSubGroupConditions(groups []*FilterGroup, conditions *[]string, args *[]interface{}) {
	for _, subGroup := range groups {
		if subGroup == nil || subGroup.IsEmpty() {
			continue
		}
		subCondition, subArgs := buildGroupCondition(subGroup)
		if subCondition != "" {
			*conditions = append(*conditions, "("+subCondition+")")
			*args = append(*args, subArgs...)
		}
	}
}

// buildFilterCondition 构建单个过滤条件的SQL和参数
func buildFilterCondition(filter *Filter) (string, interface{}) {
	if filter == nil {
		return "", nil
	}

	// 处理特殊操作符
	switch filter.Operator {
	case constants.OP_IS_NULL, constants.OP_IS_NOT_NULL:
		return buildNullCondition(filter)
	case constants.OP_BETWEEN:
		return buildBetweenCondition(filter)
	case constants.OP_STARTS_WITH:
		return buildLikeCondition(filter, false, true) // 前缀匹配
	case constants.OP_ENDS_WITH:
		return buildLikeCondition(filter, true, false) // 后缀匹配
	case constants.OP_CONTAINS:
		return buildLikeCondition(filter, true, true) // 包含匹配
	case constants.OP_FIND_IN_SET:
		return buildFindInSetCondition(filter)
	}

	// 通用处理：使用 map 查找模板
	if template, ok := constants.OperatorTemplateMap[filter.Operator]; ok {
		return fmt.Sprintf(template, filter.Field), filter.Value
	}

	return "", nil
}

// buildNullCondition 构建 NULL 判断条件
func buildNullCondition(filter *Filter) (string, interface{}) {
	if template, ok := constants.OperatorTemplateMap[filter.Operator]; ok {
		return fmt.Sprintf(template, filter.Field), nil
	}
	return "", nil
}

// buildBetweenCondition 构建 BETWEEN 条件
func buildBetweenCondition(filter *Filter) (string, interface{}) {
	values, ok := filter.Value.([]interface{})
	if !ok || len(values) != 2 {
		return "", nil
	}
	return fmt.Sprintf(constants.SQL_BETWEEN, filter.Field), values
}

// buildLikeCondition 构建 LIKE 条件
// prefix: 是否在前面加通配符, suffix: 是否在后面加通配符
func buildLikeCondition(filter *Filter, prefix, suffix bool) (string, interface{}) {
	valueStr, ok := filter.Value.(string)
	if !ok {
		return "", nil
	}

	var pattern string
	if prefix && suffix {
		pattern = constants.SQL_WILDCARD_ANY + valueStr + constants.SQL_WILDCARD_ANY
	} else if prefix {
		pattern = constants.SQL_WILDCARD_ANY + valueStr
	} else if suffix {
		pattern = valueStr + constants.SQL_WILDCARD_ANY
	} else {
		pattern = valueStr
	}

	return fmt.Sprintf(constants.SQL_LIKE, filter.Field), pattern
}

// buildFindInSetCondition 构建 FIND_IN_SET 条件
func buildFindInSetCondition(filter *Filter) (string, interface{}) {
	template, ok := constants.OperatorTemplateMap[filter.Operator]
	if !ok {
		return "", nil
	}
	return fmt.Sprintf(template, filter.Field) + " > 0", filter.Value
}

// buildGroupCondition 递归构建过滤组的条件
func buildGroupCondition(group *FilterGroup) (string, []interface{}) {
	if group == nil || group.IsEmpty() {
		return "", nil
	}

	var conditions []string
	var args []interface{}

	// 收集过滤条件和子组条件
	collectFilterConditionsWithArgs(group.Filters, &conditions, &args)
	collectSubGroupConditions(group.Groups, &conditions, &args)

	if len(conditions) == 0 {
		return "", nil
	}

	return strings.Join(conditions, fmt.Sprintf(" %s ", group.LogicOp.String())), args
}

// collectFilterConditionsWithArgs 收集过滤条件（处理 BETWEEN 等多参数情况）
func collectFilterConditionsWithArgs(filters []*Filter, conditions *[]string, args *[]interface{}) {
	for _, filter := range filters {
		if filter == nil {
			continue
		}
		condition, arg := buildFilterCondition(filter)
		if condition != "" {
			*conditions = append(*conditions, condition)
			if arg != nil {
				// 处理BETWEEN操作的多个参数
				if values, ok := arg.([]interface{}); ok {
					*args = append(*args, values...)
				} else {
					*args = append(*args, arg)
				}
			}
		}
	}
}

// ApplyOrdering 应用排序条件
// orders: 排序条件列表
// defaultOrder: 默认排序（当orders为空时使用），可为空字符串
func ApplyOrdering(db *gorm.DB, orders []Order, defaultOrder string) *gorm.DB {
	// 应用排序条件
	for _, order := range orders {
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
	if len(orders) == 0 && defaultOrder != "" {
		db = db.Order(defaultOrder)
	}

	return db
}

// ApplyFieldSelection 应用字段选择
// selectFields: 要选择的字段列表
// omitFields: 要排除的字段列表
// modelFields: 模型的所有字段（用于自动字段模式）
// autoFields: 是否启用自动字段模式
func ApplyFieldSelection(db *gorm.DB, selectFields, omitFields, modelFields []string, autoFields bool) *gorm.DB {
	// 如果同时指定了 Select 和 Omit，优先使用 Select
	if len(selectFields) > 0 {
		db = db.Select(selectFields)
		return db
	}

	// 如果指定了 Omit
	if len(omitFields) > 0 {
		// 如果启用了自动字段模式，从模型字段中排除
		if autoFields && len(modelFields) > 0 {
			selectedFields := FilterFields(modelFields, nil, omitFields)
			if len(selectedFields) > 0 {
				db = db.Select(selectedFields)
			}
			return db
		}
		// 否则使用GORM的Omit
		db = db.Omit(omitFields...)
		return db
	}

	// 如果启用了自动字段模式且没有指定任何字段选择，使用缓存的模型字段
	if autoFields && len(modelFields) > 0 {
		db = db.Select(modelFields)
		return db
	}

	// 没有指定字段选择，返回原查询（SELECT *）
	return db
}

// ========== 并发查询支持 ==========

// ExecuteConcurrentQueries 执行并发查询任务
// 这是一个便捷方法，直接在 repository 上执行并发查询
func (r *BaseRepository[T]) ExecuteConcurrentQueries(
	ctx context.Context,
	tasks []ConcurrentQueryTask[int64],
	opts ...ConcurrentQueryOption,
) ([]ConcurrentQueryResult[int64], bool) {
	executor := NewConcurrentQueryExecutor(r.db.GetDB()).WithLogger(r.logger)
	for _, opt := range opts {
		opt(executor)
	}
	return ExecuteConcurrentQuery(executor, ctx, tasks)
}

// ConcurrentQuery 简化的并发查询接口
// 直接传入查询函数的 map，返回结果的 map
func (r *BaseRepository[T]) ConcurrentQuery(
	ctx context.Context,
	queries map[string]func(ctx context.Context) (int64, error),
	opts ...ConcurrentQueryOption,
) (map[string]int64, bool) {
	executor := NewConcurrentQueryExecutor(r.db.GetDB()).WithLogger(r.logger)
	for _, opt := range opts {
		opt(executor)
	}
	return ConcurrentSimpleQuery(executor, ctx, queries)
}
