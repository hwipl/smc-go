package util

import (
	"bytes"
	"sync"
)

// Buffer is a bytes.Buffer protected by a mutex
type Buffer struct {
	lock   sync.Mutex
	buffer bytes.Buffer
}

// Write writes p to the buffer
func (b *Buffer) Write(p []byte) (n int, err error) {
	b.lock.Lock()
	defer b.lock.Unlock()
	return b.buffer.Write(p)
}

// CopyBuffer copies the underlying bytes.Buffer and returns it
func (b *Buffer) CopyBuffer() *bytes.Buffer {
	b.lock.Lock()
	defer b.lock.Unlock()
	oldBuf := b.buffer.Bytes()
	newBuf := make([]byte, len(oldBuf))
	copy(newBuf, oldBuf)
	return bytes.NewBuffer(newBuf)
}

// Reset removes everything from the underlying bytes.Buffer
func (b *Buffer) Reset() {
	b.lock.Lock()
	defer b.lock.Unlock()
	b.buffer = bytes.Buffer{}
}
