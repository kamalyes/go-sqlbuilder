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
	"unsafe"
)

// UnsafeBytes 将字符串转换为字节切片
func UnsafeBytes(s string) []byte {
	if s == "" {
		return []byte{} // 返回空切片而不是 nil
	}
	return *(*[]byte)(unsafe.Pointer(&s))
}
