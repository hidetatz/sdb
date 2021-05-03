package ssdb

import (
	"crypto/sha256"
	"fmt"

	"github.com/dty1er/sdb/lru"
)

// BufferPool manages pages, indices, and files on the disk.
type BufferPool struct {
	// lru cache element type is *PageDescriptor. BufferPool never manages Page directly.
	frames *lru.Cache
	// TODO: keep b-tree index here
}

func NewBufferPool(entryCount int) (*BufferPool, error) {
	frames := lru.New(lru.WithCap(entryCount))
	return &BufferPool{
		frames: frames,
	}, nil
}
