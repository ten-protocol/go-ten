package hostdb

import gethcommon "github.com/ethereum/go-ethereum/common"

const truncHash = 16

func truncTo16(hash gethcommon.Hash) []byte {
	return truncLastTo16(hash.Bytes())
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

func truncLastTo16(bytes []byte) []byte {
	if len(bytes) == 0 {
		return bytes
	}
	start := len(bytes) - truncHash
	if start < 0 {
		start = 0
	}
	b := bytes[start:]
	c := make([]byte, truncHash)
	copy(c, b)
	return c
}
