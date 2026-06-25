/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-11 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-08 15:05:00
 * @FilePath: \go-sqlbuilder\types\json.go
 * @Description: JSON - 泛型 JSON 类型，支持任意可序列化类型的数据库存储
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package types

import (
	"database/sql/driver"
	"fmt"
	"reflect"

	"github.com/kamalyes/go-toolbox/pkg/serializer"
)

// JSON 泛型 JSON 类型，支持任意可序列化类型
type JSON[T any] struct {
	Data T
}

// Len 获取数据长度（仅对切片类型有效）
func (j JSON[T]) Len() int {
	rv := reflect.ValueOf(j.Data)
	if rv.Kind() == reflect.Slice {
		return rv.Len()
	}
	return 0
}

// IsEmpty 判断是否为空（切片长度为0或数据为零值）
func (j JSON[T]) IsEmpty() bool {
	rv := reflect.ValueOf(j.Data)
	if rv.Kind() == reflect.Slice {
		return rv.Len() == 0
	}
	return rv.IsZero()
}

// Append 追加元素（仅对切片类型有效）
func (j *JSON[T]) Append(items ...any) {
	rv := reflect.ValueOf(&j.Data)
	if rv.Elem().Kind() != reflect.Slice {
		return
	}
	for _, item := range items {
		rv.Elem().Set(reflect.Append(rv.Elem(), reflect.ValueOf(item)))
	}
}

// Clone 克隆数据
func (j JSON[T]) Clone() JSON[T] {
	rv := reflect.ValueOf(j.Data)
	clone := reflect.New(rv.Type()).Elem()
	clone.Set(rv)
	return JSON[T]{Data: clone.Interface().(T)}
}

func (j *JSON[T]) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	b, err := toBytes(value)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}
	if len(b) == 0 {
		return nil
	}
	return serializer.JSONUnmarshal(b, &j.Data)
}

func (j JSON[T]) Value() (driver.Value, error) {
	return serializer.JSONMarshal(j.Data)
}

// Get 获取数据
func (j *JSON[T]) Get() T {
	return j.Data
}

// Set 设置数据
func (j *JSON[T]) Set(data T) {
	j.Data = data
}
