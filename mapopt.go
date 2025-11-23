/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-23 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-23 22:50:00
 * @FilePath: \go-sqlbuilder\mapopt.go
 * @Description: Map类型扩展 - MapAny、MapString、StringSlice的数据库序列化和泛型操作
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package sqlbuilder

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// MapAny 任意类型的 Map，支持数据库 JSON 序列化
type MapAny map[string]any

func (m *MapAny) Scan(value interface{}) error {
	if value == nil {
		*m = make(MapAny)
		return nil
	}
	bytesValue, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal MapAny value: %v", value)
	}
	return json.Unmarshal(bytesValue, m)
}

func (m MapAny) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil
	}
	return json.Marshal(m)
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

// MapString 字符串类型的 Map，支持数据库 JSON 序列化
type MapString map[string]string

func (m *MapString) Scan(value interface{}) error {
	if value == nil {
		*m = make(MapString)
		return nil
	}
	bytesValue, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal MapString value: %v", value)
	}
	return json.Unmarshal(bytesValue, m)
}

func (m MapString) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil
	}
	return json.Marshal(m)
}

// Get 获取指定 key 的值，支持默认值
func (m MapString) Get(key string, defaultValue ...string) string {
	if val, ok := m[key]; ok {
		return val
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return ""
}

// Set 设置值
func (m MapString) Set(key, value string) MapString {
	m[key] = value
	return m
}

// Has 检查 key 是否存在
func (m MapString) Has(key string) bool {
	_, ok := m[key]
	return ok
}

// Delete 删除指定 key
func (m MapString) Delete(key string) MapString {
	delete(m, key)
	return m
}

// Keys 获取所有 key
func (m MapString) Keys() []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// Values 获取所有 value
func (m MapString) Values() []string {
	values := make([]string, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}

// Merge 合并另一个 Map
func (m MapString) Merge(other MapString) MapString {
	for k, v := range other {
		m[k] = v
	}
	return m
}

// Clone 克隆一个新的 Map
func (m MapString) Clone() MapString {
	clone := make(MapString, len(m))
	for k, v := range m {
		clone[k] = v
	}
	return clone
}

// ToMapAny 转换为 MapAny
func (m MapString) ToMapAny() MapAny {
	result := make(MapAny, len(m))
	for k, v := range m {
		result[k] = v
	}
	return result
}

// StringSlice 字符串切片，支持数据库 JSON 序列化
type StringSlice []string

func (s *StringSlice) Scan(value interface{}) error {
	if value == nil {
		*s = []string{}
		return nil
	}
	bytesValue, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal StringSlice value: %v", value)
	}
	return json.Unmarshal(bytesValue, s)
}

func (s StringSlice) Value() (driver.Value, error) {
	if s == nil {
		return nil, nil
	}
	return json.Marshal(s)
}

// Contains 检查是否包含指定元素
func (s StringSlice) Contains(item string) bool {
	for _, v := range s {
		if v == item {
			return true
		}
	}
	return false
}

// IndexOf 查找元素索引，未找到返回 -1
func (s StringSlice) IndexOf(item string) int {
	for i, v := range s {
		if v == item {
			return i
		}
	}
	return -1
}

// Append 追加元素
func (s *StringSlice) Append(items ...string) StringSlice {
	*s = append(*s, items...)
	return *s
}

// Remove 移除指定元素（第一个匹配的）
func (s *StringSlice) Remove(item string) StringSlice {
	for i, v := range *s {
		if v == item {
			*s = append((*s)[:i], (*s)[i+1:]...)
			break
		}
	}
	return *s
}

// RemoveAt 移除指定索引的元素
func (s *StringSlice) RemoveAt(index int) StringSlice {
	if index >= 0 && index < len(*s) {
		*s = append((*s)[:index], (*s)[index+1:]...)
	}
	return *s
}

// Filter 过滤元素
func (s StringSlice) Filter(fn func(string) bool) StringSlice {
	result := make(StringSlice, 0)
	for _, v := range s {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}

// Map 映射转换
func (s StringSlice) Map(fn func(string) string) StringSlice {
	result := make(StringSlice, len(s))
	for i, v := range s {
		result[i] = fn(v)
	}
	return result
}

// Unique 去重
func (s StringSlice) Unique() StringSlice {
	seen := make(map[string]bool)
	result := make(StringSlice, 0)
	for _, v := range s {
		if !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}
	return result
}

// Join 连接为字符串
func (s StringSlice) Join(sep string) string {
	if len(s) == 0 {
		return ""
	}
	result := s[0]
	for i := 1; i < len(s); i++ {
		result += sep + s[i]
	}
	return result
}

// Clone 克隆一个新的切片
func (s StringSlice) Clone() StringSlice {
	clone := make(StringSlice, len(s))
	copy(clone, s)
	return clone
}

// 泛型 JSON 类型，支持任意可序列化类型
type JSONType[T any] struct {
	Data T
}

func (j *JSONType[T]) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytesValue, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONType value: %v", value)
	}
	return json.Unmarshal(bytesValue, &j.Data)
}

func (j JSONType[T]) Value() (driver.Value, error) {
	return json.Marshal(j.Data)
}

// Get 获取数据
func (j *JSONType[T]) Get() T {
	return j.Data
}

// Set 设置数据
func (j *JSONType[T]) Set(data T) {
	j.Data = data
}

// 泛型切片类型
type Slice[T any] []T

func (s *Slice[T]) Scan(value interface{}) error {
	if value == nil {
		*s = []T{}
		return nil
	}
	bytesValue, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal Slice value: %v", value)
	}
	return json.Unmarshal(bytesValue, s)
}

func (s Slice[T]) Value() (driver.Value, error) {
	if s == nil {
		return nil, nil
	}
	return json.Marshal(s)
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
