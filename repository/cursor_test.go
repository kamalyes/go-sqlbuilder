/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-12-11 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-16 11:15:27
 * @FilePath: \go-sqlbuilder\repository\cursor_test.go
 * @Description: 游标工具（用于复合排序分页）
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package repository

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimeEntityCursor(t *testing.T) {
	// 测试创建游标
	t.Run("NewTimeEntityCursor", func(t *testing.T) {
		now := time.Now()
		cursor := NewTimeEntityCursor(&now, "123")
		assert.NotNil(t, cursor)
		assert.Equal(t, now.Format(TimeFormat), cursor.Timestamp)
		assert.Equal(t, "123", cursor.EntityID)
	})

	// 测试游标字符串表示
	t.Run("String", func(t *testing.T) {
		now := time.Now()
		cursor := NewTimeEntityCursor(&now, "123")
		assert.Equal(t, cursor.String(), cursor.Timestamp+CursorSeparator+cursor.EntityID)
	})

	// 测试游标有效性
	t.Run("IsValid", func(t *testing.T) {
		cursor := NewTimeEntityCursor(nil, "")
		assert.Error(t, cursor.IsValid())

		cursor = NewTimeEntityCursor(timePtr(time.Now()), "123")
		assert.NoError(t, cursor.IsValid())
	})

	// 测试游标克隆
	t.Run("Clone", func(t *testing.T) {
		cursor := NewTimeEntityCursor(timePtr(time.Now()), "123")
		clone := cursor.Clone()
		assert.Equal(t, cursor, clone)
		assert.NotSame(t, cursor, clone) // 确保是不同的实例
	})

	// 测试游标比较
	t.Run("Compare", func(t *testing.T) {
		cursor1 := NewTimeEntityCursor(timePtr(time.Now().Add(-time.Minute)), "123")
		cursor2 := NewTimeEntityCursor(timePtr(time.Now()), "123")
		result, err := cursor1.Compare(cursor2)
		assert.NoError(t, err)
		assert.Equal(t, -1, result)

		cursor3 := NewTimeEntityCursor(timePtr(time.Now()), "124")
		result, err = cursor2.Compare(cursor3)
		assert.NoError(t, err)
		assert.Equal(t, -1, result)

		cursor4 := NewTimeEntityCursor(timePtr(time.Now()), "123")
		result, err = cursor2.Compare(cursor4)
		assert.NoError(t, err)
		assert.Equal(t, 0, result)
	})

	// 测试游标序列化
	t.Run("ToJSON", func(t *testing.T) {
		cursor := NewTimeEntityCursor(timePtr(time.Now()), "123")
		jsonStr, err := cursor.ToJSON()
		assert.NoError(t, err)
		assert.Contains(t, jsonStr, cursor.EntityID)
		assert.Contains(t, jsonStr, cursor.Timestamp)
	})

	// 测试游标反序列化
	t.Run("ParseTimeEntityCursor", func(t *testing.T) {
		cursor := NewTimeEntityCursor(timePtr(time.Now()), "123")
		cursorStr := cursor.String()
		parsedCursor := ParseTimeEntityCursor(cursorStr)
		assert.Equal(t, cursor, parsedCursor)
	})

	// 测试Base64编码
	t.Run("ToBase64", func(t *testing.T) {
		cursor := NewTimeEntityCursor(timePtr(time.Now()), "123")
		base64Str := cursor.ToBase64()
		assert.NotEmpty(t, base64Str)
		decodedCursor := ParseTimeEntityCursor(base64Str)
		assert.Equal(t, cursor, decodedCursor)
	})

	// 测试构建器模式
	t.Run("Builder", func(t *testing.T) {
		builder := NewCursorBuilder().
			WithTime(timePtr(time.Now())).
			WithEntityID("123")
		cursor, err := builder.Build()
		assert.NoError(t, err)
		assert.Equal(t, "123", cursor.EntityID)
	})

	// 测试游标反序列化JSON
	t.Run("ParseTimeEntityCursorFromJSON", func(t *testing.T) {
		cursor := NewTimeEntityCursor(timePtr(time.Now()), "123")
		jsonStr, _ := cursor.ToJSON()
		parsedCursor, err := ParseTimeEntityCursorFromJSON(jsonStr)
		assert.NoError(t, err)
		assert.Equal(t, cursor, parsedCursor)
	})
}

func timePtr(t time.Time) *time.Time {
	return &t
}
