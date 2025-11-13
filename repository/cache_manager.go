/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-10 01:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-13 07:45:23
 * @FilePath: \go-sqlbuilder\repository\cache_manager.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package repository

import (
	"context"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/kamalyes/go-sqlbuilder/cache"
	"github.com/kamalyes/go-sqlbuilder/errors"
)

// CacheManagerImpl 缓存管理器实现
type CacheManagerImpl struct {
	handler cache.Handler
}

// NewCacheManager 创建缓存管理器
func NewCacheManager(handler cache.Handler) CacheManager {
	if handler == nil {
		handler = &cache.NoCacheHandler{}
	}
	return &CacheManagerImpl{handler: handler}
}

// Get 获取缓存值
func (m *CacheManagerImpl) Get(ctx context.Context, key string) (interface{}, error) {
	data, err := m.handler.Get(stringToBytes(key))
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}
	var result interface{}
	err = jsoniter.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Set 设置缓存值
func (m *CacheManagerImpl) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	if value == nil {
		return errors.NewError(errors.ErrorCodeInvalidInput, errors.MsgCacheValueCannotBeNil)
	}

	// 序列化为 JSON
	data, err := jsoniter.Marshal(value)
	if err != nil {
		return errors.NewErrorf(errors.ErrorCodeCacheError, errors.MsgFailedToMarshalCache+": %v", err)
	}

	return m.handler.SetWithTTL(stringToBytes(key), data, ttl)
}

// Delete 删除单个或多个缓存键
func (m *CacheManagerImpl) Delete(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}

	byteKeys := make([][]byte, len(keys))
	for i, key := range keys {
		byteKeys[i] = stringToBytes(key)
	}

	return m.handler.Del(byteKeys...)
}

// DeletePattern 删除匹配模式的所有键（仅适用于支持的后端）
func (m *CacheManagerImpl) DeletePattern(ctx context.Context, pattern string) error {
	return errors.NewError(errors.ErrorCodeUnsupported, errors.MsgDeletePatternNotImplemented)
}

// Exists 检查缓存键是否存在
func (m *CacheManagerImpl) Exists(ctx context.Context, key string) (bool, error) {
	return m.handler.Exists(stringToBytes(key))
}

// TTL 获取键的剩余 TTL
func (m *CacheManagerImpl) TTL(ctx context.Context, key string) (time.Duration, error) {
	return 0, errors.NewError(errors.ErrorCodeUnsupported, errors.MsgTTLQueryNotImplemented)
}

// BatchGet 批量获取缓存值
func (m *CacheManagerImpl) BatchGet(ctx context.Context, keys ...string) (map[string]interface{}, error) {
	results := make(map[string]interface{})

	for _, key := range keys {
		value, err := m.Get(ctx, key)
		if err == nil && value != nil {
			results[key] = value
		}
	}

	return results, nil
}

// BatchSet 批量设置缓存值
func (m *CacheManagerImpl) BatchSet(ctx context.Context, items map[string]interface{}, ttl time.Duration) error {
	var errs []error

	for key, value := range items {
		v := value
		err := m.Set(ctx, key, &v, ttl)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return errors.NewErrorf(errors.ErrorCodeCacheError, errors.MsgBatchSetEncounteredErrors+": %v", errs)
	}

	return nil
}

// BatchDelete 批量删除缓存键
func (m *CacheManagerImpl) BatchDelete(ctx context.Context, keys ...string) error {
	return m.Delete(ctx, keys...)
}

// stringToBytes 将字符串转换为字节数组
func stringToBytes(s string) []byte {
	return []byte(s)
}
