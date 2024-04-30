package enclavedb

import gethcommon "github.com/ethereum/go-ethereum/common"

func truncTo4(hash gethcommon.Hash) []byte {
	return truncBTo4(hash.Bytes())
}

func truncBTo4(bytes []byte) []byte {
	if len(bytes) == 0 {
		return bytes
	}
	b := bytes[0:4]
	c := make([]byte, 4)
	copy(c, b)
	return c
}
