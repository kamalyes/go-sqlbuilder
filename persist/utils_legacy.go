//go:build !go1.17
// +build !go1.17

/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-11 21:28:49
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-11 21:29:36
 * @FilePath: \go-sqlbuilder\persist\utils_legacy.go
 * @Description:
 *
 * UnsafeBytes converts a string to a byte slice without allocating new memory.
 * This function is only for Go versions below 1.17, as it uses unsafe operations
 * that may lead to undefined behavior if the string is modified after conversion.
 * Use with caution.
 *
 * Example:
 *   str := "hello"
 *   bytes := UnsafeBytes(str)
 *   fmt.Println(bytes) // Output: [104 101 108 108 111]
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package persist

import (
	"reflect"
	"unsafe"
)

// UnsafeBytes 将字符串转换为字节切片且不分配新内存（Go 1.16及以下）
func UnsafeBytes(s string) []byte {
	stringHeader := (*reflect.StringHeader)(unsafe.Pointer(&s))
	return *(*[]byte)(unsafe.Pointer(&struct {
		Data uintptr
		Len  int
	}{stringHeader.Data, len(s)}))
}
