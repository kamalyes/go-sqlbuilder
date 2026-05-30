/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-11 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-08 15:15:19
 * @FilePath: \go-sqlbuilder\types\map_any.go
 * @Description: MapAny - 任意类型 Map，支持数据库 JSON 序列化
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package types

import (
	"database/sql/driver"
	"fmt"

	"github.com/kamalyes/go-toolbox/pkg/serializer"
)

// MapAny 任意类型的 Map，支持数据库 JSON 序列化
type MapAny map[string]any

func (m *MapAny) Scan(value interface{}) error {
	if value == nil {
		*m = make(MapAny)
		return nil
	}
	b, err := toBytes(value)
	if err != nil {
		return fmt.Errorf("failed to unmarshal MapAny: %w", err)
	}
	if len(b) == 0 {
		*m = make(MapAny)
		return nil
	}
	return serializer.JSONUnmarshal(b, m)
}

func (m MapAny) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil
	}
	return serializer.JSONMarshal(m)
}

// Get 获取指定 key 的值，支持默认值
func (m MapAny) Get(key string, defaultValue ...any) any {
	if val, ok := m[key]; ok {
		return val
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return nil
}

// GetString 获取字符串值
func (m MapAny) GetString(key string, defaultValue ...string) string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return ""
}

// GetInt 获取整数值
func (m MapAny) GetInt(key string, defaultValue ...int) int {
	if val, ok := m[key]; ok {
		switch v := val.(type) {
		case int:
			return v
		case int32:
			return int(v)
		case int64:
			return int(v)
		case float64:
			return int(v)
		}
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return 0
}

// GetInt32 获取 int32 值
func (m MapAny) GetInt32(key string, defaultValue ...int32) int32 {
	if val, ok := m[key]; ok {
		switch v := val.(type) {
		case int32:
			return v
		case int:
			return int32(v)
		case int64:
			return int32(v)
		case float64:
			return int32(v)
		}
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return 0
}

// GetBool 获取布尔值
func (m MapAny) GetBool(key string, defaultValue ...bool) bool {
	if val, ok := m[key]; ok {
		if b, ok := val.(bool); ok {
			return b
		}
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return false
}

// GetMap 获取 MapAny 值
func (m MapAny) GetMap(key string, defaultValue ...MapAny) MapAny {
	if val, ok := m[key]; ok {
		switch v := val.(type) {
		case MapAny:
			return v
		case map[string]interface{}:
			return MapAny(v)
		}
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return MapAny{}
}

// Set 设置值
func (m MapAny) Set(key string, value any) MapAny {
	m[key] = value
	return m
}

// Has 检查 key 是否存在
func (m MapAny) Has(key string) bool {
	_, ok := m[key]
	return ok
}

// Delete 删除指定 key
func (m MapAny) Delete(key string) MapAny {
	delete(m, key)
	return m
}

// Keys 获取所有 key
func (m MapAny) Keys() []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// Values 获取所有 value
func (m MapAny) Values() []any {
	values := make([]any, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}

// Merge 合并另一个 Map
func (m MapAny) Merge(other MapAny) MapAny {
	for k, v := range other {
		m[k] = v
	}
	return m
}

// Clone 克隆一个新的 Map
func (m MapAny) Clone() MapAny {
	clone := make(MapAny, len(m))
	for k, v := range m {
		clone[k] = v
	}
	return clone
}
