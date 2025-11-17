/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-17 15:40:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-17 16:00:00
 * @FilePath: \go-sqlbuilder\repository\cache.go
 * @Description: Repository 缓存支持
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package repository

import (
	"context"
	"time"

	"github.com/kamalyes/go-sqlbuilder/cache"
)

// RepositoryWithCache 带缓存的仓储
type RepositoryWithCache[T cache.Cacheable] struct {
	*BaseRepository[T]
	cache *cache.Wrapper
}

// NewRepositoryWithCache 创建带缓存的仓储
func NewRepositoryWithCache[T cache.Cacheable](repo *BaseRepository[T], cacheWrapper *cache.Wrapper) *RepositoryWithCache[T] {
	return &RepositoryWithCache[T]{
		BaseRepository: repo,
		cache:          cacheWrapper,
	}
}

// GetWithCache 获取记录，优先从缓存读取
func (r *RepositoryWithCache[T]) GetWithCache(ctx context.Context, id interface{}) (*T, error) {
	// 先尝试从数据库获取以构建缓存键
	entity, err := r.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, nil
	}

	cacheKey := (*entity).CacheKey()

	// 尝试从缓存读取
	cached := new(T)
	hit, err := r.cache.Get(ctx, cacheKey, cached)
	if err == nil && hit {
		return cached, nil
	}

	// 异步设置缓存
	go func() {
		_ = r.cache.Set(context.Background(), (*entity).CacheKey(), entity, (*entity).CacheTTL())
	}()

	return entity, nil
}

// CreateWithCache 创建记录并删除相关缓存
func (r *RepositoryWithCache[T]) CreateWithCache(ctx context.Context, entity *T) (*T, error) {
	result, err := r.Create(ctx, entity)
	if err != nil {
		return nil, err
	}

	// 删除缓存
	if result != nil {
		_ = r.cache.Delete(ctx, (*result).CacheKey())
	}

	return result, nil
}

// UpdateWithCache 更新记录并删除缓存
func (r *RepositoryWithCache[T]) UpdateWithCache(ctx context.Context, entity *T) (*T, error) {
	result, err := r.Update(ctx, entity)
	if err != nil {
		return nil, err
	}

	// 删除缓存
	if result != nil {
		_ = r.cache.Delete(ctx, (*result).CacheKey())
	}

	return result, nil
}

// DeleteWithCache 删除记录并删除缓存
func (r *RepositoryWithCache[T]) DeleteWithCache(ctx context.Context, id interface{}) error {
	// 先获取实体以得到缓存键
	entity, err := r.Get(ctx, id)
	if err != nil {
		return err
	}

	err = r.Delete(ctx, id)
	if err != nil {
		return err
	}

	// 删除缓存
	if entity != nil {
		_ = r.cache.Delete(ctx, (*entity).CacheKey())
	}

	return nil
}

// ListWithCache 查询列表，支持缓存
func (r *RepositoryWithCache[T]) ListWithCache(ctx context.Context, cacheKey string, ttl time.Duration, query *Query) ([]*T, error) {
	// 尝试从缓存读取
	var cached []*T
	hit, err := r.cache.Get(ctx, cacheKey, &cached)
	if err == nil && hit {
		return cached, nil
	}

	// 缓存未命中，从数据库查询
	entities, err := r.List(ctx, query)
	if err != nil {
		return nil, err
	}

	// 异步设置缓存
	go func() {
		_ = r.cache.Set(context.Background(), cacheKey, entities, ttl)
	}()

	return entities, nil
}

// InvalidateCache 使缓存失效
func (r *RepositoryWithCache[T]) InvalidateCache(ctx context.Context, keys ...string) error {
	return r.cache.Delete(ctx, keys...)
}

// InvalidateCacheByEntities 根据实体使缓存失效
func (r *RepositoryWithCache[T]) InvalidateCacheByEntities(ctx context.Context, entities ...*T) error {
	models := make([]cache.Cacheable, 0, len(entities))
	for _, entity := range entities {
		if entity != nil {
			models = append(models, *entity)
		}
	}
	return r.cache.DeleteByModels(ctx, models...)
}

// GetOrLoad 获取缓存，如果不存在则执行加载函数
func (r *RepositoryWithCache[T]) GetOrLoad(ctx context.Context, cacheKey string, ttl time.Duration, loadFn func() (*T, error)) (*T, error) {
	var result T
	err := r.cache.GetOrLoad(ctx, cacheKey, ttl, &result, func() (interface{}, error) {
		return loadFn()
	})
	if err != nil {
		return nil, err
	}
	return &result, nil
}
