/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-23 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-23 22:50:00
 * @FilePath: \go-sqlbuilder\mapopt.go
 * @Description: Map类型扩展 - MapAny、MapString、StringSlice的数据库序列化
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package sqlbuilder

import (
	"database/sql/driver"
	"encoding/json"
)

type MapAny map[string]any

func (m *MapAny) Scan(value interface{}) error {
	bytesValue, _ := value.([]byte)
	return json.Unmarshal(bytesValue, m)
}

func (m MapAny) Value() (driver.Value, error) {
	return json.Marshal(m)
}

type MapString map[string]string

func (m *MapString) Scan(value interface{}) error {
	bytesValue, _ := value.([]byte)
	return json.Unmarshal(bytesValue, m)
}

func (m MapString) Value() (driver.Value, error) {
	return json.Marshal(m)
}

type StringSlice []string

func (s *StringSlice) Scan(value interface{}) error {
	bytesValue, _ := value.([]byte)
	return json.Unmarshal(bytesValue, s)
}

func (s StringSlice) Value() (driver.Value, error) {
	value, err := json.Marshal(s)
	return value, err
}
