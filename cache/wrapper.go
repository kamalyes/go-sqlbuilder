/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-17 16:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-17 16:00:00
 * @FilePath: \go-sqlbuilder\cache\wrapper.go
 * @Description: Repository 缓存包装器
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package cache

import (
	"context"
	"time"
	"unsafe"

	jsoniter "github.com/json-iterator/go"
)

// Cacheable 可缓存的模型接口
type Cacheable interface {
	CacheKey() string        // 缓存键
	CacheTTL() time.Duration // 缓存过期时间
}

// Wrapper 缓存包装器，用于 Repository 层的缓存操作
type Wrapper struct {
	store     Store
	keyPrefix string
}

// NewWrapper 创建缓存包装器
func NewWrapper(store Store, keyPrefix string) *Wrapper {
	if keyPrefix == "" {
		keyPrefix = "repo:"
	}
	return &Wrapper{
		store:     store,
		keyPrefix: keyPrefix,
	}
}

// buildKey 构建完整的缓存键
func (w *Wrapper) buildKey(key string) string {
	return w.keyPrefix + key
}

// Get 从缓存获取并反序列化到目标对象
func (w *Wrapper) Get(ctx context.Context, key string, target interface{}) (bool, error) {
	fullKey := w.buildKey(key)
	value, err := w.store.Get(ctx, fullKey)
	if err != nil || value == "" {
		return false, err
	}

	err = jsoniter.UnmarshalFromString(value, target)
	if err != nil {
		return false, err
	}

	return true, nil
}

// Set 序列化并设置缓存
func (w *Wrapper) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := jsoniter.MarshalToString(value)
	if err != nil {
		return err
	}

	fullKey := w.buildKey(key)
	return w.store.Set(ctx, fullKey, data, ttl)
}

// Delete 删除缓存
func (w *Wrapper) Delete(ctx context.Context, keys ...string) error {
	for _, key := range keys {
		fullKey := w.buildKey(key)
		if err := w.store.Delete(ctx, fullKey); err != nil {
			return err
		}
	}
	return nil
}

// DeleteByModels 根据可缓存模型删除缓存
func (w *Wrapper) DeleteByModels(ctx context.Context, models ...Cacheable) error {
	keys := make([]string, 0, len(models))
	for _, model := range models {
		if model != nil {
			keys = append(keys, model.CacheKey())
		}
	}
	return w.Delete(ctx, keys...)
}

// GetOrLoad 获取缓存，如果不存在则执行加载函数
func (w *Wrapper) GetOrLoad(ctx context.Context, key string, ttl time.Duration, target interface{}, loadFn func() (interface{}, error)) error {
	// 尝试从缓存读取
	hit, err := w.Get(ctx, key, target)
	if err == nil && hit {
		return nil
	}

	// 缓存未命中，执行加载函数
	data, err := loadFn()
	if err != nil {
		return err
	}

	// 异步设置缓存
	go func() {
		_ = w.Set(context.Background(), key, data, ttl)
	}()

	// 将加载的数据赋值给 target
	if data != nil {
		bytes, err := jsoniter.Marshal(data)
		if err == nil {
			_ = jsoniter.Unmarshal(bytes, target)
		}
	}

	return nil
}

// Exists 检查缓存是否存在
func (w *Wrapper) Exists(ctx context.Context, key string) (bool, error) {
	fullKey := w.buildKey(key)
	return w.store.Exists(ctx, fullKey)
}

// Clear 清除指定前缀的所有缓存
func (w *Wrapper) Clear(ctx context.Context, prefix string) error {
	fullPrefix := w.buildKey(prefix)
	return w.store.Clear(ctx, fullPrefix)
}

// GetBytes 获取原始字节数据
func (w *Wrapper) GetBytes(ctx context.Context, key string) ([]byte, error) {
	fullKey := w.buildKey(key)
	value, err := w.store.Get(ctx, fullKey)
	if err != nil {
		return nil, err
	}
	return UnsafeBytes(value), nil
}

// SetBytes 设置原始字节数据
func (w *Wrapper) SetBytes(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	fullKey := w.buildKey(key)
	return w.store.Set(ctx, fullKey, UnsafeString(value), ttl)
}

// UnsafeBytes 字符串转字节切片（零拷贝）
func UnsafeBytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// UnsafeString 字节切片转字符串（零拷贝）
func UnsafeString(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}
