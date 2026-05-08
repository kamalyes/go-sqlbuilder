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
	"encoding/json"
	"fmt"
)

// JSON 泛型 JSON 类型，支持任意可序列化类型
type JSON[T any] struct {
	Data T
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
	return json.Unmarshal(b, &j.Data)
}

func (j JSON[T]) Value() (driver.Value, error) {
	return json.Marshal(j.Data)
}

// Get 获取数据
func (j *JSON[T]) Get() T {
	return j.Data
}

// Set 设置数据
func (j *JSON[T]) Set(data T) {
	j.Data = data
}
