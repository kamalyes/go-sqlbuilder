/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-11 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-08 15:00:00
 * @FilePath: \go-sqlbuilder\types\map_any_test.go
 * @Description: MapAny 测试
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package types

import (
	"database/sql/driver"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapAny_ScanAndValue(t *testing.T) {
	m := MapAny{}

	// Test Scan
	jsonData := []byte(`{"name":"test","age":25,"active":true}`)
	err := m.Scan(jsonData)
	assert.NoError(t, err)
	assert.Equal(t, "test", m["name"])
	assert.Equal(t, float64(25), m["age"]) // JSON 数字默认为 float64
	assert.Equal(t, true, m["active"])

	// Test Scan with nil
	m2 := MapAny{}
	err = m2.Scan(nil)
	assert.NoError(t, err)
	assert.NotNil(t, m2)
	assert.Equal(t, 0, len(m2))

	// Test Scan with invalid type
	m3 := MapAny{}
	err = m3.Scan(123)
	assert.Error(t, err)

	// Test Value
	val, err := m.Value()
	assert.NoError(t, err)
	assert.NotNil(t, val)

	// Test Value with nil
	var m4 MapAny
	val, err = m4.Value()
	assert.NoError(t, err)
	assert.Nil(t, val)
}

func TestMapAny_Operations(t *testing.T) {
	m := MapAny{
		"name":   "Alice",
		"age":    30,
		"active": true,
	}

	// Test Get
	assert.Equal(t, "Alice", m.Get("name"))
	assert.Equal(t, 30, m.Get("age"))
	assert.Nil(t, m.Get("nonexist"))
	assert.Equal(t, "default", m.Get("nonexist", "default"))

	// Test GetString
	assert.Equal(t, "Alice", m.GetString("name"))
	assert.Equal(t, "", m.GetString("age")) // 类型不匹配
	assert.Equal(t, "default", m.GetString("nonexist", "default"))

	// Test GetInt
	assert.Equal(t, 30, m.GetInt("age"))
	assert.Equal(t, 0, m.GetInt("name")) // 类型不匹配
	assert.Equal(t, 99, m.GetInt("nonexist", 99))

	// Test GetInt with different number types
	m["int64val"] = int64(100)
	m["float64val"] = float64(200)
	assert.Equal(t, 100, m.GetInt("int64val"))
	assert.Equal(t, 200, m.GetInt("float64val"))

	// Test GetBool
	assert.Equal(t, true, m.GetBool("active"))
	assert.Equal(t, false, m.GetBool("name")) // 类型不匹配
	assert.Equal(t, true, m.GetBool("nonexist", true))

	// Test Set
	m.Set("city", "Beijing")
	assert.Equal(t, "Beijing", m["city"])

	// Test Has
	assert.True(t, m.Has("name"))
	assert.False(t, m.Has("nonexist"))

	// Test Delete
	m.Delete("city")
	assert.False(t, m.Has("city"))

	// Test Keys
	keys := m.Keys()
	assert.Contains(t, keys, "name")
	assert.Contains(t, keys, "age")
	assert.Contains(t, keys, "active")

	// Test Values
	values := m.Values()
	assert.Contains(t, values, "Alice")
	assert.Contains(t, values, 30)
	assert.Contains(t, values, true)

	// Test Merge
	other := MapAny{"email": "alice@test.com", "age": 31}
	m.Merge(other)
	assert.Equal(t, "alice@test.com", m["email"])
	assert.Equal(t, 31, m["age"]) // 覆盖

	// Test Clone
	clone := m.Clone()
	clone["name"] = "Bob"
	assert.Equal(t, "Alice", m["name"]) // 原数据不变
	assert.Equal(t, "Bob", clone["name"])
}

func TestMapAny_ValuerInterface(t *testing.T) {
	var _ driver.Valuer = MapAny(nil)
}
