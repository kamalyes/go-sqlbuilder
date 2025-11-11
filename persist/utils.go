// utils.go
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
