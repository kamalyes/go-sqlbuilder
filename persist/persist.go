/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-11 21:18:55
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-13 11:07:30
 * @FilePath: \go-sqlbuilder\persist\persist.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package persist

import (
	"context"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/kamalyes/go-sqlbuilder/meta"
	"gorm.io/gorm"
)

type Model interface {
	CacheKey() string
	CacheTTL() time.Duration
}

func Create[T any](ctx context.Context, tx DBHandler, tCreates ...*T) (err error) {
	if len(tCreates) == 0 {
		return nil
	}
	err = tx.DB().WithContext(ctx).Create(&tCreates).Error
	return
}

func Get[T Model](ctx context.Context, tx DBHandler, t *T) (tGet *T, err error) {
	err = tx.DB().WithContext(ctx).Where(t).First(&tGet).Error
	return
}

func GetByFilters[T Model](ctx context.Context, tx DBHandler, filters Filters, orders ...Order) (tGet *T, err error) {
	var opts []Option
	if orders != nil {
		opts = append(opts, WithOrders(orders...))
	}

	param := NewQueryParam(filters, nil, opts...)
	err = tx.Query(param).WithContext(ctx).First(&tGet).Error
	if err != nil {
		return
	}
	return
}

func List[T Model](ctx context.Context, tx DBHandler, filters Filters, page *meta.Paging, orders ...Order) (tList []*T, err error) {
	var opts []Option
	if orders != nil {
		opts = append(opts, WithOrders(orders...))
	}

	param := NewQueryParam(filters, page, opts...)
	err = tx.Query(param).Model(new(T)).WithContext(ctx).Find(&tList).Error
	if err != nil {
		return
	}
	if page != nil {
		param = NewQueryParam(filters, nil)
		err = tx.Query(param).Model(new(T)).WithContext(ctx).Count(&page.Total).Error
	}
	return
}

func Update[T Model](ctx context.Context, tx DBHandler, tUpdates ...*T) (err error) {
	if len(tUpdates) == 0 {
		return nil
	}

	// 使用事务确保批量更新的一致性
	err = tx.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, item := range tUpdates {
			if err := tx.Updates(item).Error; err != nil {
				return err
			}
		}
		return nil
	})
	return
}

func Save[T Model](ctx context.Context, tx DBHandler, tSaves ...*T) (err error) {
	if len(tSaves) == 0 {
		return nil
	}

	err = tx.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, item := range tSaves {
			if err := tx.Save(item).Error; err != nil {
				return err
			}
		}
		return nil
	})
	return
}

func Delete[T Model](ctx context.Context, tx DBHandler, tDeletes ...*T) (err error) {
	if len(tDeletes) == 0 {
		return nil
	}

	err = tx.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, item := range tDeletes {
			if err := tx.Delete(item).Error; err != nil {
				return err
			}
		}
		return nil
	})
	return
}

// 缓存相关操作

// CreateWithCache 创建记录并删除缓存
func CreateWithCache[T Model](ctx context.Context, tx DBHandler, cache CacheHandler, tCreates ...*T) (err error) {
	if len(tCreates) == 0 {
		return nil
	}

	// 先收集所有缓存键
	cacheKeys := make([][]byte, 0, len(tCreates))
	for _, item := range tCreates {
		cacheKeys = append(cacheKeys, UnsafeBytes((*item).CacheKey()))
	}

	err = Create(ctx, tx, tCreates...)
	if err != nil {
		return
	}

	return cache.Del(cacheKeys...)
}

// GetWithCache 获取记录，优先从缓存读取
func GetWithCache[T Model](ctx context.Context, tx DBHandler, cache CacheHandler, t *T) (tGet *T, err error) {
	cacheKey := (*t).CacheKey()

	// 尝试从缓存读取
	v, err := cache.Get(UnsafeBytes(cacheKey))
	if err == nil && v != nil {
		// 缓存命中，反序列化数据
		tGet = new(T)
		err = jsoniter.Unmarshal(v, tGet)
		if err == nil {
			return tGet, nil
		}
		// 反序列化失败，继续从数据库读取
	}

	// 缓存未命中或读取失败，从数据库读取
	tGet, err = Get(ctx, tx, t)
	if err != nil {
		return nil, err
	}

	if tGet != nil {
		// 写入缓存
		v, err = jsoniter.Marshal(tGet)
		if err != nil {
			// 记录日志但不影响主流程
			return tGet, nil
		}
		// 异步设置缓存，不阻塞主流程
		go func() {
			_ = cache.SetWithTTL(UnsafeBytes((*tGet).CacheKey()), v, (*tGet).CacheTTL())
		}()
	}

	return tGet, nil
}

// UpdateWithCache 更新记录并删除缓存
func UpdateWithCache[T Model](ctx context.Context, tx DBHandler, cache CacheHandler, tUpdates ...*T) (err error) {
	if len(tUpdates) == 0 {
		return nil
	}

	// 先收集所有缓存键
	cacheKeys := make([][]byte, 0, len(tUpdates))
	for _, item := range tUpdates {
		cacheKeys = append(cacheKeys, UnsafeBytes((*item).CacheKey()))
	}

	err = Update(ctx, tx, tUpdates...)
	if err != nil {
		return
	}

	// 批量删除缓存
	return cache.Del(cacheKeys...)
}

// SaveWithCache 保存记录并删除缓存
func SaveWithCache[T Model](ctx context.Context, tx DBHandler, cache CacheHandler, tSaves ...*T) (err error) {
	if len(tSaves) == 0 {
		return nil
	}

	// 先收集所有缓存键
	cacheKeys := make([][]byte, 0, len(tSaves))
	for _, item := range tSaves {
		cacheKeys = append(cacheKeys, UnsafeBytes((*item).CacheKey()))
	}

	err = Save(ctx, tx, tSaves...)
	if err != nil {
		return
	}

	// 批量删除缓存
	return cache.Del(cacheKeys...)
}

// DeleteWithCache 删除记录并删除缓存
func DeleteWithCache[T Model](ctx context.Context, tx DBHandler, cache CacheHandler, tDeletes ...*T) (err error) {
	if len(tDeletes) == 0 {
		return nil
	}

	// 先收集所有缓存键
	cacheKeys := make([][]byte, 0, len(tDeletes))
	for _, item := range tDeletes {
		cacheKeys = append(cacheKeys, UnsafeBytes((*item).CacheKey()))
	}

	err = Delete(ctx, tx, tDeletes...)
	if err != nil {
		return
	}

	// 批量删除缓存
	return cache.Del(cacheKeys...)
}
