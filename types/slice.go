/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-11 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-08 15:00:00
 * @FilePath: \go-sqlbuilder\types\slice.go
 * @Description: Slice - 泛型切片类型，支持数据库 JSON 序列化
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package types

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/kamalyes/go-toolbox/pkg/serializer"
)

// Slice 泛型切片类型
type Slice[T any] []T

func (s *Slice[T]) Scan(value interface{}) error {
	if value == nil {
		*s = []T{}
		return nil
	}
	b, err := toBytes(value)
	if err != nil {
		return fmt.Errorf("failed to unmarshal Slice: %w", err)
	}
	if len(b) == 0 {
		*s = []T{}
		return nil
	}
	b = normalizeSliceJSON(b)
	if err := serializer.JSONUnmarshal(b, s); err != nil {
		// protojson 反序列化失败时回退到标准 json
		// 兼容旧格式数据（如 Timestamp 对象格式、枚举数字值等 generated proto json tag 可解析的格式）
		return json.Unmarshal(b, s)
	}
	return nil
}

func normalizeSliceJSON(data []byte) []byte {
	data = bytes.TrimSpace(data)
	if len(data) == 0 || data[0] == '[' || bytes.EqualFold(data, []byte("null")) {
		return data
	}
	wrapped := make([]byte, 0, len(data)+2)
	wrapped = append(wrapped, '[')
	wrapped = append(wrapped, data...)
	wrapped = append(wrapped, ']')
	return wrapped
}

func (s Slice[T]) Value() (driver.Value, error) {
	if s == nil {
		return nil, nil
	}
	// 返回 string 而不是 []byte
	//
	// driver.Value 里的 []byte 会被驱动当成二进制：pgx 在简单协议下
	// （go-rpc-gateway 对 postgres 固定开 PreferSimpleProtocol）把它编码成
	// bytea 的十六进制字面量 '\x5b7b...'，写进 json/jsonb 列直接报
	// `invalid input syntax for type json`。扩展协议下靠参数类型推断侥幸能过，
	// 于是这个 bug 只在「关掉预编译」之后才出现，看起来像是别的改动引发的。
	//
	// string 在两种协议下都按文本送，json 列能正确解析。
	//
	// [EN] A []byte in a driver.Value is treated as binary: under the simple protocol, which
	// go-rpc-gateway forces for postgres, pgx encodes it as a bytea hex literal and a json or
	// jsonb column rejects it outright. The extended protocol gets away with it through
	// parameter type inference, so the bug only appears once prepared statements are turned off
	// and looks like it came from some other change. A string is sent as text under both.
	data, err := serializer.JSONMarshal(s)
	if err != nil {
		return nil, err
	}
	return string(data), nil
}

// Len 获取长度
func (s Slice[T]) Len() int {
	return len(s)
}

// Append 追加元素
func (s *Slice[T]) Append(items ...T) Slice[T] {
	*s = append(*s, items...)
	return *s
}

// Filter 过滤元素
func (s Slice[T]) Filter(fn func(T) bool) Slice[T] {
	result := make(Slice[T], 0)
	for _, v := range s {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}

// Map 映射转换
func Map[T any, R any](s Slice[T], fn func(T) R) Slice[R] {
	result := make(Slice[R], len(s))
	for i, v := range s {
		result[i] = fn(v)
	}
	return result
}

// Clone 克隆切片
func (s Slice[T]) Clone() Slice[T] {
	clone := make(Slice[T], len(s))
	copy(clone, s)
	return clone
}
