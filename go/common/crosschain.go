package common

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type CrossChainRootHashes [][]byte

type ExtCrossChainBundle struct {
	LastBatchHash        gethcommon.Hash
	Signature            []byte
	L1BlockHash          gethcommon.Hash      // The block hash that's expected to be canonical on signature submission
	L1BlockNum           *big.Int             // The number of the block that has the block hash. This is used to verify the block hash.
	CrossChainRootHashes CrossChainRootHashes // The CrossChainRoots of the batches that are being submitted
}

func (hashes CrossChainRootHashes) ToHexString() string {
	hexStrings := make([]string, 0, len(hashes))

	for _, hash := range hashes {
		hexStrings = append(hexStrings, gethcommon.BytesToHash(hash).Hex())
	}

	return strings.Join(hexStrings, ",")
}

func (bundle ExtCrossChainBundle) HashPacked() common.Hash {
	uint256type, _ := abi.NewType("uint256", "", nil)
	bytes32type, _ := abi.NewType("bytes32", "", nil)
	bytesТype, _ := abi.NewType("bytes[]", "", nil)

	args := abi.Arguments{
		{
			Type: bytes32type,
		},
		{
			Type: bytes32type,
		},
		{
			Type: uint256type,
		},
		{
			Type: bytesТype,
		},
	}

	bytes, err := args.Pack(bundle.LastBatchHash, bundle.L1BlockHash, bundle.L1BlockNum, bundle.CrossChainRootHashes)
	if err != nil {
		panic(err)
	}

	hash := crypto.Keccak256Hash(bytes)
	return hash
}
