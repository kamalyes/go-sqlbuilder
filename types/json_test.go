/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-11 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-08 15:05:00
 * @FilePath: \go-sqlbuilder\types\json_test.go
 * @Description: JSON 测试
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package types

import (
	"database/sql/driver"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSON(t *testing.T) {
	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	// Test Scan and Value
	j := &JSON[Person]{}
	jsonData := []byte(`{"name":"Alice","age":30}`)
	err := j.Scan(jsonData)
	assert.NoError(t, err)
	assert.Equal(t, "Alice", j.Data.Name)
	assert.Equal(t, 30, j.Data.Age)

	// Test Get
	person := j.Get()
	assert.Equal(t, "Alice", person.Name)

	// Test Set
	j.Set(Person{Name: "Bob", Age: 25})
	assert.Equal(t, "Bob", j.Data.Name)

	// Test Value
	val, err := j.Value()
	assert.NoError(t, err)
	assert.NotNil(t, val)

	// Test Scan with nil
	j2 := &JSON[Person]{}
	err = j2.Scan(nil)
	assert.NoError(t, err)

	// Test Scan with invalid type
	j3 := &JSON[Person]{}
	err = j3.Scan(123)
	assert.Error(t, err)

	// Test with driver.Valuer interface
	var _ driver.Valuer = (*JSON[Person])(nil)
}
