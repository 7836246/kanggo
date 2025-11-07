package kanggo

import (
	"unsafe"
)

// b2s 字节切片转字符串（零拷贝）
// 警告：使用时要确保字节切片不会被修改
func b2s(b []byte) string {
	if len(b) == 0 {
		return ""
	}
	return *(*string)(unsafe.Pointer(&b))
}

// s2b 字符串转字节切片（零拷贝）
// 警告：返回的切片是只读的，不可修改
func s2b(s string) []byte {
	if len(s) == 0 {
		return nil
	}
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}

// BytesToString 安全的字节切片转字符串（零拷贝）
// 注意：只在确保字节切片不会被修改时使用
func BytesToString(b []byte) string {
	return b2s(b)
}

// StringToBytes 安全的字符串转字节切片（零拷贝）
// 注意：返回的切片是只读的，不可修改
func StringToBytes(s string) []byte {
	return s2b(s)
}

// UnsafeString 不安全的字节切片转字符串（零拷贝）
// 用于内部优化，外部使用请谨慎
func UnsafeString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// UnsafeBytes 不安全的字符串转字节切片（零拷贝）
// 用于内部优化，外部使用请谨慎
func UnsafeBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}
