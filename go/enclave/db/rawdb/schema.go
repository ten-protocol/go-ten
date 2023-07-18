package rawdb

import (
	gethcommon "github.com/ethereum/go-ethereum/common"
)

var contractReceiptPrefix = []byte("ocr") // contractReceiptPrefix + address -> tx hash

func contractReceiptKey(contractAddress gethcommon.Address) []byte {
	return append(contractReceiptPrefix, contractAddress.Bytes()...)
}

// for storing the contract creation count
func contractCreationCountKey() []byte {
	return []byte("contractCreationCountKey")
}
