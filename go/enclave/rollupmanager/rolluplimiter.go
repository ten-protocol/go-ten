package rollupmanager

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/rlp"
)

var (
	ErrFailedToEncode   = errors.New("failed to encode data")
	ErrSizeExceedsLimit = errors.New("data size exceeds remaining limit")
)

type RollupLimiter uint64

func (rl *RollupLimiter) Consume(encodable interface{}) error {
	encodedData, err := rlp.EncodeToBytes(encodable)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrFailedToEncode, err)
	}

	encodedSize := uint64(len(encodedData))
	if encodedSize > uint64(*rl) {
		return fmt.Errorf("%w: data size %d, remaining limit %d", ErrSizeExceedsLimit, encodedSize, *rl)
	}

	*rl -= RollupLimiter(encodedSize)
	return nil
}
