package kanggo

import (
	"sync"
)

// ByteBuffer 字节缓冲区
type ByteBuffer struct {
	B []byte
}

// Reset 重置缓冲区
func (b *ByteBuffer) Reset() {
	b.B = b.B[:0]
}

// Write 写入数据
func (b *ByteBuffer) Write(p []byte) (int, error) {
	b.B = append(b.B, p...)
	return len(p), nil
}

// WriteString 写入字符串
func (b *ByteBuffer) WriteString(s string) (int, error) {
	b.B = append(b.B, s...)
	return len(s), nil
}

// WriteByte 写入单个字节
func (b *ByteBuffer) WriteByte(c byte) error {
	b.B = append(b.B, c)
	return nil
}

// String 转换为字符串（零拷贝）
func (b *ByteBuffer) String() string {
	return b2s(b.B)
}

// Bytes 返回字节切片
func (b *ByteBuffer) Bytes() []byte {
	return b.B
}

// Len 返回缓冲区长度
func (b *ByteBuffer) Len() int {
	return len(b.B)
}

// Cap 返回缓冲区容量
func (b *ByteBuffer) Cap() int {
	return cap(b.B)
}

// ByteBufferPool 字节缓冲池
type ByteBufferPool struct {
	small  sync.Pool // < 1KB
	medium sync.Pool // 1KB - 8KB
	large  sync.Pool // > 8KB
}

// DefaultByteBufferPool 默认缓冲池
var DefaultByteBufferPool = &ByteBufferPool{
	small: sync.Pool{
		New: func() interface{} {
			return &ByteBuffer{
				B: make([]byte, 0, 512), // 512 字节
			}
		},
	},
	medium: sync.Pool{
		New: func() interface{} {
			return &ByteBuffer{
				B: make([]byte, 0, 4096), // 4KB
			}
		},
	},
	large: sync.Pool{
		New: func() interface{} {
			return &ByteBuffer{
				B: make([]byte, 0, 16384), // 16KB
			}
		},
	},
}

// Get 获取缓冲区
func (p *ByteBufferPool) Get(size int) *ByteBuffer {
	var pool *sync.Pool

	switch {
	case size <= 1024:
		pool = &p.small
	case size <= 8192:
		pool = &p.medium
	default:
		pool = &p.large
	}

	buf := pool.Get().(*ByteBuffer)
	buf.Reset()
	return buf
}

// Put 归还缓冲区
func (p *ByteBufferPool) Put(buf *ByteBuffer) {
	if buf == nil {
		return
	}

	size := cap(buf.B)
	var pool *sync.Pool

	switch {
	case size <= 1024:
		pool = &p.small
	case size <= 8192:
		pool = &p.medium
	default:
		pool = &p.large
	}

	// 如果缓冲区太大，不归还到池中
	if size > 65536 { // 64KB
		return
	}

	pool.Put(buf)
}

// AcquireByteBuffer 从默认池获取缓冲区
func AcquireByteBuffer(size int) *ByteBuffer {
	return DefaultByteBufferPool.Get(size)
}

// ReleaseByteBuffer 归还缓冲区到默认池
func ReleaseByteBuffer(buf *ByteBuffer) {
	DefaultByteBufferPool.Put(buf)
}
