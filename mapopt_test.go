/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-23 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-23 23:30:00
 * @FilePath: \go-sqlbuilder\mapopt_test.go
 * @Description: Map类型扩展测试
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package sqlbuilder

import (
	"database/sql/driver"
	"github.com/stretchr/testify/assert"
	"testing"
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

func TestMapString_ScanAndValue(t *testing.T) {
	m := MapString{}

	// Test Scan
	jsonData := []byte(`{"key1":"value1","key2":"value2"}`)
	err := m.Scan(jsonData)
	assert.NoError(t, err)
	assert.Equal(t, "value1", m["key1"])
	assert.Equal(t, "value2", m["key2"])

	// Test Scan with nil
	m2 := MapString{}
	err = m2.Scan(nil)
	assert.NoError(t, err)
	assert.NotNil(t, m2)

	// Test Scan with invalid type
	m3 := MapString{}
	err = m3.Scan(123)
	assert.Error(t, err)

	// Test Value
	val, err := m.Value()
	assert.NoError(t, err)
	assert.NotNil(t, val)

	// Test Value with nil
	var m4 MapString
	val, err = m4.Value()
	assert.NoError(t, err)
	assert.Nil(t, val)
}

func TestMapString_Operations(t *testing.T) {
	m := MapString{
		"name":  "Alice",
		"email": "alice@test.com",
	}

	// Test Get
	assert.Equal(t, "Alice", m.Get("name"))
	assert.Equal(t, "", m.Get("nonexist"))
	assert.Equal(t, "default", m.Get("nonexist", "default"))

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
	assert.Contains(t, keys, "email")

	// Test Values
	values := m.Values()
	assert.Contains(t, values, "Alice")
	assert.Contains(t, values, "alice@test.com")

	// Test Merge
	other := MapString{"phone": "123456", "email": "new@test.com"}
	m.Merge(other)
	assert.Equal(t, "123456", m["phone"])
	assert.Equal(t, "new@test.com", m["email"]) // 覆盖

	// Test Clone
	clone := m.Clone()
	clone["name"] = "Bob"
	assert.Equal(t, "Alice", m["name"]) // 原数据不变
	assert.Equal(t, "Bob", clone["name"])

	// Test ToMapAny
	mapAny := m.ToMapAny()
	assert.Equal(t, "Alice", mapAny["name"])
	assert.IsType(t, MapAny{}, mapAny)
}

func TestStringSlice_ScanAndValue(t *testing.T) {
	s := StringSlice{}

	// Test Scan
	jsonData := []byte(`["apple","banana","cherry"]`)
	err := s.Scan(jsonData)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(s))
	assert.Equal(t, "apple", s[0])

	// Test Scan with nil
	s2 := StringSlice{}
	err = s2.Scan(nil)
	assert.NoError(t, err)
	assert.NotNil(t, s2)
	assert.Equal(t, 0, len(s2))

	// Test Scan with invalid type
	s3 := StringSlice{}
	err = s3.Scan(123)
	assert.Error(t, err)

	// Test Value
	val, err := s.Value()
	assert.NoError(t, err)
	assert.NotNil(t, val)

	// Test Value with nil
	var s4 StringSlice
	val, err = s4.Value()
	assert.NoError(t, err)
	assert.Nil(t, val)
}

func TestStringSlice_Operations(t *testing.T) {
	s := StringSlice{"apple", "banana", "cherry"}

	// Test Contains
	assert.True(t, s.Contains("apple"))
	assert.False(t, s.Contains("orange"))

	// Test IndexOf
	assert.Equal(t, 0, s.IndexOf("apple"))
	assert.Equal(t, 1, s.IndexOf("banana"))
	assert.Equal(t, -1, s.IndexOf("orange"))

	// Test Append
	s.Append("orange", "grape")
	assert.Equal(t, 5, len(s))
	assert.Equal(t, "orange", s[3])

	// Test Remove
	s.Remove("banana")
	assert.Equal(t, 4, len(s))
	assert.False(t, s.Contains("banana"))

	// Test RemoveAt
	s.RemoveAt(0)
	assert.Equal(t, 3, len(s))
	assert.Equal(t, "cherry", s[0])

	// Test RemoveAt with invalid index
	s.RemoveAt(100)
	assert.Equal(t, 3, len(s)) // 不变

	// Test Filter
	s = StringSlice{"apple", "apricot", "banana", "avocado"}
	filtered := s.Filter(func(item string) bool {
		return item[0] == 'a'
	})
	assert.Equal(t, 3, len(filtered))
	assert.Contains(t, filtered, "apple")

	// Test Map
	mapped := s.Map(func(item string) string {
		return item + "!"
	})
	assert.Equal(t, "apple!", mapped[0])
	assert.Equal(t, "banana!", mapped[2])

	// Test Unique
	s = StringSlice{"apple", "banana", "apple", "cherry", "banana"}
	unique := s.Unique()
	assert.Equal(t, 3, len(unique))
	assert.Contains(t, unique, "apple")
	assert.Contains(t, unique, "banana")
	assert.Contains(t, unique, "cherry")

	// Test Join
	s = StringSlice{"apple", "banana", "cherry"}
	joined := s.Join(", ")
	assert.Equal(t, "apple, banana, cherry", joined)

	// Test Join empty
	s = StringSlice{}
	joined = s.Join(", ")
	assert.Equal(t, "", joined)

	// Test Clone
	s = StringSlice{"apple", "banana"}
	clone := s.Clone()
	clone[0] = "orange"
	assert.Equal(t, "apple", s[0]) // 原数据不变
	assert.Equal(t, "orange", clone[0])
}

func TestJSONType(t *testing.T) {
	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	// Test Scan and Value
	j := &JSONType[Person]{}
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
	j2 := &JSONType[Person]{}
	err = j2.Scan(nil)
	assert.NoError(t, err)

	// Test Scan with invalid type
	j3 := &JSONType[Person]{}
	err = j3.Scan(123)
	assert.Error(t, err)

	// Test with driver.Valuer interface
	var _ driver.Valuer = (*JSONType[Person])(nil)
}

func TestSlice(t *testing.T) {
	// Test Scan and Value
	s := &Slice[int]{}
	jsonData := []byte(`[1,2,3,4,5]`)
	err := s.Scan(jsonData)
	assert.NoError(t, err)
	assert.Equal(t, 5, len(*s))
	assert.Equal(t, 1, (*s)[0])

	// Test Len
	assert.Equal(t, 5, s.Len())

	// Test Append
	s.Append(6, 7)
	assert.Equal(t, 7, s.Len())

	// Test Filter
	filtered := s.Filter(func(n int) bool {
		return n%2 == 0
	})
	assert.Equal(t, 3, len(filtered)) // 2, 4, 6

	// Test Clone
	clone := s.Clone()
	clone[0] = 99
	assert.Equal(t, 1, (*s)[0]) // 原数据不变
	assert.Equal(t, 99, clone[0])

	// Test Value
	val, err := s.Value()
	assert.NoError(t, err)
	assert.NotNil(t, val)

	// Test Scan with nil
	s2 := &Slice[int]{}
	err = s2.Scan(nil)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(*s2))

	// Test Scan with invalid type
	s3 := &Slice[int]{}
	err = s3.Scan(123)
	assert.Error(t, err)

	// Test Value with nil
	var s4 Slice[int]
	val, err = s4.Value()
	assert.NoError(t, err)
	assert.Nil(t, val)
}

func TestSlice_Map(t *testing.T) {
	s := Slice[int]{1, 2, 3, 4, 5}

	// Test Map to string
	mapped := Map(s, func(n int) string {
		return string(rune('A' + n - 1))
	})
	assert.Equal(t, 5, len(mapped))
	assert.Equal(t, "A", mapped[0])
	assert.Equal(t, "E", mapped[4])

	// Test Map to double
	doubled := Map(s, func(n int) int {
		return n * 2
	})
	assert.Equal(t, 2, doubled[0])
	assert.Equal(t, 10, doubled[4])
}

func TestSlice_ComplexTypes(t *testing.T) {
	type User struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	s := &Slice[User]{}
	jsonData := []byte(`[{"id":1,"name":"Alice"},{"id":2,"name":"Bob"}]`)
	err := s.Scan(jsonData)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(*s))
	assert.Equal(t, "Alice", (*s)[0].Name)
	assert.Equal(t, "Bob", (*s)[1].Name)

	// Test Filter
	filtered := s.Filter(func(u User) bool {
		return u.ID == 1
	})
	assert.Equal(t, 1, len(filtered))
	assert.Equal(t, "Alice", filtered[0].Name)

	// Test Map
	names := Map(*s, func(u User) string {
		return u.Name
	})
	assert.Equal(t, 2, len(names))
	assert.Equal(t, "Alice", names[0])
	assert.Equal(t, "Bob", names[1])
}
