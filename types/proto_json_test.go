/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-05-09 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-09 00:00:00
 * @FilePath: \go-sqlbuilder\types\proto_json_test.go
 * @Description: ProtoJSON 单元测试
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */
package types

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestProtoJSON_ProtoMessage(t *testing.T) {
	t.Run("scan and value roundtrip", func(t *testing.T) {
		original := wrapperspb.String("hello")
		pj := ProtoJSON[*wrapperspb.StringValue]{Data: original}

		val, err := pj.Value()
		require.NoError(t, err)
		assert.Equal(t, `"hello"`, val)

		var restored ProtoJSON[*wrapperspb.StringValue]
		err = restored.Scan(val)
		require.NoError(t, err)
		assert.Equal(t, "hello", restored.Data.GetValue())
	})

	t.Run("scan nil", func(t *testing.T) {
		var pj ProtoJSON[*wrapperspb.StringValue]
		err := pj.Scan(nil)
		require.NoError(t, err)
	})

	t.Run("scan empty bytes", func(t *testing.T) {
		var pj ProtoJSON[*wrapperspb.StringValue]
		err := pj.Scan([]byte{})
		require.NoError(t, err)
	})

	t.Run("scan null string", func(t *testing.T) {
		var pj ProtoJSON[*wrapperspb.StringValue]
		err := pj.Scan("null")
		require.NoError(t, err)
	})

	t.Run("scan string type", func(t *testing.T) {
		jsonStr := `"world"`
		var pj ProtoJSON[*wrapperspb.StringValue]
		err := pj.Scan(jsonStr)
		require.NoError(t, err)
		assert.Equal(t, "world", pj.Data.GetValue())
	})

	t.Run("value nil data", func(t *testing.T) {
		var pj ProtoJSON[*wrapperspb.StringValue]
		val, err := pj.Value()
		require.NoError(t, err)
		assert.Nil(t, val)
	})

	t.Run("value empty proto", func(t *testing.T) {
		pj := ProtoJSON[*wrapperspb.StringValue]{Data: &wrapperspb.StringValue{Value: ""}}
		val, err := pj.Value()
		require.NoError(t, err)
		assert.Nil(t, val)
	})

	t.Run("unsupported scan type", func(t *testing.T) {
		var pj ProtoJSON[*wrapperspb.StringValue]
		err := pj.Scan(123)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unsupported scan type")
	})
}

func TestProtoJSON_GetSet(t *testing.T) {
	t.Run("get and set", func(t *testing.T) {
		pj := ProtoJSON[*wrapperspb.StringValue]{Data: wrapperspb.String("original")}
		assert.Equal(t, "original", pj.Get().GetValue())

		pj.Set(wrapperspb.String("updated"))
		assert.Equal(t, "updated", pj.Get().GetValue())
	})
}

type TestPayload struct {
	Name   *wrapperspb.StringValue `json:"name"`
	Age    *wrapperspb.Int32Value  `json:"age"`
	Active *wrapperspb.BoolValue   `json:"active"`
}

func TestProtoJSON_Struct(t *testing.T) {
	t.Run("scan and value roundtrip", func(t *testing.T) {
		original := TestPayload{
			Name:   wrapperspb.String("test"),
			Age:    wrapperspb.Int32(25),
			Active: wrapperspb.Bool(true),
		}
		pj := ProtoJSON[TestPayload]{Data: original}

		val, err := pj.Value()
		require.NoError(t, err)
		assert.NotEmpty(t, val)

		var restored ProtoJSON[TestPayload]
		err = restored.Scan(val)
		require.NoError(t, err)
		assert.Equal(t, "test", restored.Data.Name.GetValue())
		assert.Equal(t, int32(25), restored.Data.Age.GetValue())
		assert.True(t, restored.Data.Active.GetValue())
	})

	t.Run("scan nil", func(t *testing.T) {
		var pj ProtoJSON[TestPayload]
		err := pj.Scan(nil)
		require.NoError(t, err)
	})

	t.Run("scan empty bytes", func(t *testing.T) {
		var pj ProtoJSON[TestPayload]
		err := pj.Scan([]byte{})
		require.NoError(t, err)
	})

	t.Run("scan null string", func(t *testing.T) {
		var pj ProtoJSON[TestPayload]
		err := pj.Scan("null")
		require.NoError(t, err)
	})

	t.Run("scan partial json", func(t *testing.T) {
		jsonStr := `{"name":"partial"}`
		var pj ProtoJSON[TestPayload]
		err := pj.Scan(jsonStr)
		require.NoError(t, err)
		assert.Equal(t, "partial", pj.Data.Name.GetValue())
	})

	t.Run("value empty struct", func(t *testing.T) {
		pj := ProtoJSON[TestPayload]{Data: TestPayload{}}
		val, err := pj.Value()
		require.NoError(t, err)
		assert.NotNil(t, val)
	})

	t.Run("verify json structure", func(t *testing.T) {
		original := TestPayload{
			Name: wrapperspb.String("test"),
			Age:  wrapperspb.Int32(100),
		}
		pj := ProtoJSON[TestPayload]{Data: original}

		val, err := pj.Value()
		require.NoError(t, err)

		var raw map[string]json.RawMessage
		err = json.Unmarshal([]byte(val.(string)), &raw)
		require.NoError(t, err)
		assert.Contains(t, raw, "name")
		assert.Contains(t, raw, "age")
	})
}

func TestProtoJSON_Int32Value(t *testing.T) {
	t.Run("int32 value roundtrip", func(t *testing.T) {
		pj := ProtoJSON[*wrapperspb.Int32Value]{Data: wrapperspb.Int32(42)}
		val, err := pj.Value()
		require.NoError(t, err)
		assert.Equal(t, `42`, val)

		var restored ProtoJSON[*wrapperspb.Int32Value]
		err = restored.Scan(val)
		require.NoError(t, err)
		assert.Equal(t, int32(42), restored.Data.GetValue())
	})
}

func TestProtoJSON_BoolValue(t *testing.T) {
	t.Run("bool value roundtrip", func(t *testing.T) {
		pj := ProtoJSON[*wrapperspb.BoolValue]{Data: wrapperspb.Bool(true)}
		val, err := pj.Value()
		require.NoError(t, err)
		assert.Equal(t, `true`, val)

		var restored ProtoJSON[*wrapperspb.BoolValue]
		err = restored.Scan(val)
		require.NoError(t, err)
		assert.True(t, restored.Data.GetValue())
	})
}

func TestProtoJSONMap_ScanValue(t *testing.T) {
	t.Run("scan and value roundtrip", func(t *testing.T) {
		jsonStr := `{
			"key1": "hello",
			"key2": "world"
		}`

		pjm := NewProtoJSONMap[*wrapperspb.StringValue]()
		err := pjm.Scan(jsonStr)
		require.NoError(t, err)
		assert.Len(t, pjm.Fields, 2)
		assert.Equal(t, "hello", pjm.Get("key1").GetValue())
		assert.Equal(t, "world", pjm.Get("key2").GetValue())

		val, err := pjm.Value()
		require.NoError(t, err)

		var restored map[string]string
		err = json.Unmarshal(val.([]byte), &restored)
		require.NoError(t, err)
		assert.Len(t, restored, 2)
		assert.Equal(t, `"hello"`, restored["key1"])
		assert.Equal(t, `"world"`, restored["key2"])
	})

	t.Run("scan nil", func(t *testing.T) {
		pjm := NewProtoJSONMap[*wrapperspb.StringValue]()
		err := pjm.Scan(nil)
		require.NoError(t, err)
	})

	t.Run("scan empty bytes", func(t *testing.T) {
		pjm := NewProtoJSONMap[*wrapperspb.StringValue]()
		err := pjm.Scan([]byte{})
		require.NoError(t, err)
	})

	t.Run("scan null string", func(t *testing.T) {
		pjm := NewProtoJSONMap[*wrapperspb.StringValue]()
		err := pjm.Scan("null")
		require.NoError(t, err)
	})

	t.Run("value empty map", func(t *testing.T) {
		pjm := NewProtoJSONMap[*wrapperspb.StringValue]()
		val, err := pjm.Value()
		require.NoError(t, err)
		assert.Nil(t, val)
	})

	t.Run("get set", func(t *testing.T) {
		pjm := NewProtoJSONMap[*wrapperspb.StringValue]()
		pjm.Set("key1", wrapperspb.String("test"))
		assert.Equal(t, "test", pjm.Get("key1").GetValue())
		assert.Nil(t, pjm.Get("nonexistent"))
	})

	t.Run("keys", func(t *testing.T) {
		pjm := NewProtoJSONMap[*wrapperspb.StringValue]()
		pjm.Set("a", wrapperspb.String("1"))
		pjm.Set("b", wrapperspb.String("2"))
		keys := pjm.Keys()
		assert.Len(t, keys, 2)
	})

	t.Run("unsupported scan type", func(t *testing.T) {
		pjm := NewProtoJSONMap[*wrapperspb.StringValue]()
		err := pjm.Scan(123)
		assert.Error(t, err)
	})
}
