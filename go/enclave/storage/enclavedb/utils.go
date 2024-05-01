package enclavedb

import (
	"strings"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

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

func repeat(token string, sep string, count int) string {
	elems := make([]string, count)
	for i := 0; i < count; i++ {
		elems[i] = token
	}
	return strings.Join(elems, sep)
}
