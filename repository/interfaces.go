/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-10 01:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-11 09:20:54
 * @FilePath: \go-sqlbuilder\repository\interfaces.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package repository

import (
	"context"
	"github.com/kamalyes/go-sqlbuilder/cache"
	"github.com/kamalyes/go-sqlbuilder/constant"
	"github.com/kamalyes/go-sqlbuilder/core"
	"github.com/kamalyes/go-sqlbuilder/meta"
	"time"
)

// Filter 过滤条件（复用 core.Filter）
type Filter = core.Filter

// Order 排序条件
type Order struct {
	Field     string
	Direction string // "ASC" or "DESC"
}

// Query 查询条件
type Query struct {
	Filters     []*Filter
	Orders      []Order
	Pagination  *meta.Paging
	LimitValue  *int
	OffsetValue *int
}

// NewEqFilter 创建等于过滤条件
func NewEqFilter(field string, value interface{}) *Filter {
	return &Filter{Field: field, Operator: constant.OP_EQ, Value: value}
}

// NewGtFilter 创建大于过滤条件
func NewGtFilter(field string, value interface{}) *Filter {
	return &Filter{Field: field, Operator: constant.OP_GT, Value: value}
}

// NewLtFilter 创建小于过滤条件
func NewLtFilter(field string, value interface{}) *Filter {
	return &Filter{Field: field, Operator: constant.OP_LT, Value: value}
}

// NewGteFilter 创建大于等于过滤条件
func NewGteFilter(field string, value interface{}) *Filter {
	return &Filter{Field: field, Operator: constant.OP_GTE, Value: value}
}

// NewLteFilter 创建小于等于过滤条件
func NewLteFilter(field string, value interface{}) *Filter {
	return &Filter{Field: field, Operator: constant.OP_LTE, Value: value}
}

// NewInFilter 创建 IN 过滤条件
func NewInFilter(field string, values ...interface{}) *Filter {
	return &Filter{Field: field, Operator: constant.OP_IN, Value: values}
}

// NewLikeFilter 创建 LIKE 过滤条件
func NewLikeFilter(field string, value string) *Filter {
	return &Filter{Field: field, Operator: constant.OP_LIKE, Value: "%" + value + "%"}
}

// NewNeqFilter 创建不等于过滤条件
func NewNeqFilter(field string, value interface{}) *Filter {
	return &Filter{Field: field, Operator: constant.OP_NEQ, Value: value}
}

// NewBetweenFilter 创建 BETWEEN 过滤条件
func NewBetweenFilter(field string, min, max interface{}) *Filter {
	return &Filter{Field: field, Operator: constant.OP_BETWEEN, Value: []interface{}{min, max}}
}

// NewNotInFilter 创建 NOT IN 过滤条件
func NewNotInFilter(field string, values ...interface{}) *Filter {
	return &Filter{Field: field, Operator: constant.OP_NOT_IN, Value: values}
}

// NewIsNullFilter 创建 IS NULL 过滤条件
func NewIsNullFilter(field string) *Filter {
	return &Filter{Field: field, Operator: constant.OP_IS_NULL, Value: nil}
}

// NewIsNotNullFilter 创建 IS NOT NULL 过滤条件
func NewIsNotNullFilter(field string) *Filter {
	return &Filter{Field: field, Operator: constant.OP_IS_NOT_NULL, Value: nil}
}

// NewStartsWithFilter 创建前缀匹配过滤条件
func NewStartsWithFilter(field string, value string) *Filter {
	return &Filter{Field: field, Operator: constant.OP_LIKE, Value: value + "%"}
}

// NewEndsWithFilter 创建后缀匹配过滤条件
func NewEndsWithFilter(field string, value string) *Filter {
	return &Filter{Field: field, Operator: constant.OP_LIKE, Value: "%" + value}
}

// NewContainsFilter 创建包含匹配过滤条件（与 NewLikeFilter 相同）
func NewContainsFilter(field string, value string) *Filter {
	return NewLikeFilter(field, value)
}

// NewQuery 创建查询条件
func NewQuery() *Query {
	return &Query{
		Filters: make([]*Filter, 0),
		Orders:  make([]Order, 0),
	}
}

// FindOptions 兼容旧API的查询选项结构
type FindOptions struct {
	Conditions []Condition
	Orders     []OrderBy
	Limit      int
	Offset     int
}

// Condition 查询条件结构
type Condition struct {
	Field string
	Op    constant.Operator
	Value interface{}
}

// OrderBy 排序条件结构
type OrderBy struct {
	Field     string
	Direction string
}

// AddFilter 添加过滤条件
func (q *Query) AddFilter(filter *Filter) *Query {
	if filter != nil {
		q.Filters = append(q.Filters, filter)
	}
	return q
}

// AddFilters 批量添加过滤条件
func (q *Query) AddFilters(filters ...*Filter) *Query {
	for _, f := range filters {
		q.AddFilter(f)
	}
	return q
}

// AddOrder 添加排序条件
func (q *Query) AddOrder(field, direction string) *Query {
	q.Orders = append(q.Orders, Order{Field: field, Direction: direction})
	return q
}

// WithPaging 设置分页条件
func (q *Query) WithPaging(page, pageSize int) *Query {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	q.Pagination = &meta.Paging{
		Page:     int32(page),
		PageSize: int32(pageSize),
	}
	return q
}

// Limit 设置查询限制数量
func (q *Query) Limit(limit int) *Query {
	q.LimitValue = &limit
	return q
}

// Offset 设置查询偏移量
func (q *Query) Offset(offset int) *Query {
	q.OffsetValue = &offset
	return q
}

// Transaction 事务接口
type Transaction interface {
	// 执行事务内的操作
	Create(ctx context.Context, entity interface{}) error
	CreateBatch(ctx context.Context, entities ...interface{}) error
	Update(ctx context.Context, entity interface{}) error
	UpdateBatch(ctx context.Context, entities ...interface{}) error
	Delete(ctx context.Context, entity interface{}) error
	DeleteBatch(ctx context.Context, entities ...interface{}) error
}

// CacheConfig 缓存配置
type CacheConfig struct {
	Enabled bool
	TTL     time.Duration
}

// Repository 通用仓储接口
type Repository[T any] interface {
	// 创建
	Create(ctx context.Context, entity *T) (*T, error)
	CreateBatch(ctx context.Context, entities ...*T) error

	// 读取
	Get(ctx context.Context, id interface{}) (*T, error)
	GetByFilter(ctx context.Context, filter *Filter) (*T, error)
	GetByFilters(ctx context.Context, filters ...*Filter) (*T, error)
	GetAll(ctx context.Context) ([]*T, error)
	First(ctx context.Context, filters ...*Filter) (*T, error)
	Last(ctx context.Context, filters ...*Filter) (*T, error)
	FindOne(ctx context.Context, filters ...*Filter) (*T, error)
	List(ctx context.Context, query *Query) ([]*T, error)
	ListWithPagination(ctx context.Context, query *Query, page *meta.Paging) ([]*T, *meta.Paging, error)

	// 更新
	Update(ctx context.Context, entity *T) (*T, error)
	UpdateBatch(ctx context.Context, entities ...*T) error
	UpdateByFilters(ctx context.Context, entity *T, filters ...*Filter) error
	UpdateFields(ctx context.Context, id interface{}, fields map[string]interface{}) error
	UpdateFieldsByFilters(ctx context.Context, fields map[string]interface{}, filters ...*Filter) error

	// 删除
	Delete(ctx context.Context, id interface{}) error
	DeleteBatch(ctx context.Context, ids ...interface{}) error
	DeleteByFilters(ctx context.Context, filters ...*Filter) error
	SoftDelete(ctx context.Context, id interface{}, field string, value interface{}) error
	SoftDeleteBatch(ctx context.Context, ids []interface{}, field string, value interface{}) error
	SoftDeleteByFilters(ctx context.Context, field string, value interface{}, filters ...*Filter) error
	Restore(ctx context.Context, id interface{}, field string, restoreValue interface{}) error
	RestoreBatch(ctx context.Context, ids []interface{}, field string, restoreValue interface{}) error

	// 事务
	Transaction(ctx context.Context, fn func(tx Transaction) error) error

	// 工具方法
	Count(ctx context.Context, filters ...*Filter) (int64, error)
	Exists(ctx context.Context, filters ...*Filter) (bool, error)
	CountByField(ctx context.Context, field string) (map[interface{}]int64, error)
	Pluck(ctx context.Context, field string, filters ...*Filter) ([]interface{}, error)
	Distinct(ctx context.Context, field string, filters ...*Filter) ([]interface{}, error)
}

// CacheManager 缓存管理接口
type CacheManager interface {
	// 获取
	Get(ctx context.Context, key string) (interface{}, error)

	// 设置
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error

	// 删除
	Delete(ctx context.Context, keys ...string) error

	// 删除匹配的键
	DeletePattern(ctx context.Context, pattern string) error

	// 检查存在
	Exists(ctx context.Context, key string) (bool, error)

	// 获取 TTL
	TTL(ctx context.Context, key string) (time.Duration, error)

	// 批量操作
	BatchGet(ctx context.Context, keys ...string) (map[string]interface{}, error)
	BatchSet(ctx context.Context, items map[string]interface{}, ttl time.Duration) error
	BatchDelete(ctx context.Context, keys ...string) error
}

// CachedRepository 带缓存的仓储接口
type CachedRepository[T cache.Cacheable] interface {
	Repository[T]

	// 带缓存的操作
	GetWithCache(ctx context.Context, id interface{}) (*T, error)
	CreateWithCache(ctx context.Context, entity *T) (*T, error)
	UpdateWithCache(ctx context.Context, entity *T) (*T, error)
	DeleteWithCache(ctx context.Context, id interface{}) error
	ListWithCache(ctx context.Context, cacheKey string, ttl time.Duration, query *Query) ([]*T, error)
	GetOrLoad(ctx context.Context, cacheKey string, ttl time.Duration, loadFn func() (*T, error)) (*T, error)

	// 缓存管理
	InvalidateCache(ctx context.Context, keys ...string) error
	InvalidateCacheByEntities(ctx context.Context, entities ...*T) error
}

// RepositoryFactory 仓储工厂接口
type RepositoryFactory interface {
	// 缓存管理
	CacheManager() CacheManager

	// 关闭工厂
	Close() error
}
