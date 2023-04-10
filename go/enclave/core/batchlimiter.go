package core

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

type BatchSizeLimiter struct {
	Size            uint64
	ContractAddress common.Address
	Topic           common.Hash
}

var ErrInsufficientSpace = errors.New("insufficient space in BatchSizeLimiter")

func NewBatchSizeLimiter(size uint64, contractAddress common.Address, topic common.Hash) *BatchSizeLimiter {
	return &BatchSizeLimiter{
		Size:            size,
		ContractAddress: contractAddress,
		Topic:           topic,
	}
}

func (l *BatchSizeLimiter) AcceptTransaction(tx *types.Transaction) error {
	if l == nil {
		return nil
	}

	rlpSize, err := getRlpSize(tx)
	if err != nil {
		return err
	}

	if uint64(rlpSize) > l.Size {
		return ErrInsufficientSpace
	}

	l.Size -= uint64(rlpSize)
	return nil
}

func (l *BatchSizeLimiter) ProcessReceipt(receipt *types.Receipt) error {
	if l == nil {
		return nil
	}

	for _, log := range receipt.Logs {
		if log.Address == l.ContractAddress {
			for _, t := range log.Topics {
				if t == l.Topic {
					rlpSize, err := getRlpSize(log)
					if err != nil {
						return err
					}

					if uint64(rlpSize) > l.Size {
						return ErrInsufficientSpace
					}

					l.Size -= uint64(rlpSize)
				}
			}
		}
	}

	return nil
}

func getRlpSize(val interface{}) (int, error) {
	enc, err := rlp.EncodeToBytes(val)
	if err != nil {
		return 0, err
	}

	return len(enc), nil
}
