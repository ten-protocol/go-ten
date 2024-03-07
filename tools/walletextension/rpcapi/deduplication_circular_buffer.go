package rpcapi

import "github.com/ethereum/go-ethereum/common"

// LogKey uniquely represents a log (consists of BlockHash, TxHash, and Index)
type LogKey struct {
	BlockHash common.Hash // Not necessary, but can be helpful in edge case of block reorg.
	TxHash    common.Hash
	Index     uint
}

// CircularBuffer is a data structure that uses a single, fixed-size buffer as if it was connected end-to-end.
type CircularBuffer struct {
	data []LogKey
	size int
	end  int
}

// NewCircularBuffer initializes a new CircularBuffer of the given size.
func NewCircularBuffer(size int) *CircularBuffer {
	return &CircularBuffer{
		data: make([]LogKey, size),
		size: size,
		end:  0,
	}
}

// Push adds a new LogKey to the end of the buffer. If the buffer is full,
// it overwrites the oldest data with the new LogKey.
func (cb *CircularBuffer) Push(key LogKey) {
	cb.data[cb.end] = key
	cb.end = (cb.end + 1) % cb.size
}

// Contains checks if the given LogKey exists in the buffer
func (cb *CircularBuffer) Contains(key LogKey) bool {
	for _, item := range cb.data {
		if item == key {
			return true
		}
	}
	return false
}
