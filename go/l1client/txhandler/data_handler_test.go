package txhandler

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
)

func TestSerialization(t *testing.T) {
	// TODO would be nice to have a rollup with randomized data instead of mostly nil data
	rol := nodecommon.Rollup{
		Header: &nodecommon.Header{
			ParentHash:  obscurocommon.GenesisBlock.Hash(),
			Agg:         common.Address{},
			Nonce:       2,
			L1Proof:     obscurocommon.L1RootHash{},
			State:       nodecommon.StateRoot{},
			Height:      0,
			Withdrawals: nil,
		},
		Transactions: nil,
	}

	serializedRollup := EncodeToString(nodecommon.EncodeRollup(&rol))
	deserializedRollup := DecodeFromString(serializedRollup)
	newRollup, err := nodecommon.DecodeRollup(deserializedRollup)
	if err != nil {
		t.Error(err)
	}

	if rol.Hash() != newRollup.Hash() {
		t.Errorf("unexpected hashes when converting")
	}
}

func TestCompression(t *testing.T) {
	rol := nodecommon.Rollup{
		Header: &nodecommon.Header{
			ParentHash:  obscurocommon.GenesisBlock.Hash(),
			Agg:         common.Address{},
			Nonce:       2,
			L1Proof:     obscurocommon.L1RootHash{},
			State:       nodecommon.StateRoot{},
			Height:      0,
			Withdrawals: nil,
		},
		Transactions: nil,
	}

	compressedRollup := Compress(nodecommon.EncodeRollup(&rol))
	serializedRollup := EncodeToString(compressedRollup)
	deserializedRollup := DecodeFromString(serializedRollup)
	decompressedRollup := Decompress(deserializedRollup)
	newRollup, err := nodecommon.DecodeRollup(decompressedRollup)
	if err != nil {
		t.Error(err)
	}

	if rol.Hash() != newRollup.Hash() {
		t.Errorf("unexpected hashes when converting")
	}
}
