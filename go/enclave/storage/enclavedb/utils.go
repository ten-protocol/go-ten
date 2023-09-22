package enclavedb

import gethcommon "github.com/ethereum/go-ethereum/common"

const truncHash = 16

func truncTo16(hash gethcommon.Hash) []byte {
	return truncBTo16(hash.Bytes())
}

func truncBTo16(bytes []byte) []byte {
	if len(bytes) == 0 {
		return bytes
	}
	b := bytes[0:truncHash]
	c := make([]byte, truncHash)
	copy(c, b)
	return c
}
