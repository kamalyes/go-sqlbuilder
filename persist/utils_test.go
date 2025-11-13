/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-11 21:13:15
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-13 11:07:41
 * @FilePath: \go-sqlbuilder\persist\query_param.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package persist

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnsafeBytes(t *testing.T) {
	tests := []struct {
		input    string
		expected []byte
	}{
		{"hello", []byte{104, 101, 108, 108, 111}}, // ASCII for "hello"
		{"", []byte{}}, // Empty string
		{"Go语言", []byte{0x47, 0x6f, 0xe8, 0xaf, 0xad, 0xe8, 0xa8, 0x80}}, // UTF-8 for "Go语言"
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result := UnsafeBytes(test.input)
			assert.Equal(t, []byte(test.input), result)
			assert.Equal(t, test.expected, result, "For input '%s', expected '%v' but got '%v'", test.input, test.expected, result)
		})
	}
}
