/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-11 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-13 11:02:21
 * @FilePath: \go-sqlbuilder\cache\handler.go
 * @Description: 缓存处理器接口 - 从 go-data-repository 迁移
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package cache

import "time"

// Handler 缓存处理器接口
type Handler interface {
	Get(key []byte) ([]byte, error)
	Set(key, value []byte, ttl time.Duration) error
	SetWithTTL(key, value []byte, ttl time.Duration) error
	Del(keys ...[]byte) error
	Exists(key []byte) (bool, error)
	Close() error
}

// NoCacheHandler 无缓存实现 - 所有操作都直接返回
type NoCacheHandler struct{}

func (n *NoCacheHandler) Get(key []byte) ([]byte, error) {
	return nil, nil
}

func (n *NoCacheHandler) Set(key, value []byte, ttl time.Duration) error {
	return nil
}

func (n *NoCacheHandler) SetWithTTL(key, value []byte, ttl time.Duration) error {
	return nil
}

func (n *NoCacheHandler) Del(keys ...[]byte) error {
	return nil
}

func (n *NoCacheHandler) Exists(key []byte) (bool, error) {
	return false, nil
}

func (n *NoCacheHandler) Close() error {
	return nil
}
