package crypto

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/big"

	"github.com/ten-protocol/go-ten/go/common"

	gethlog "github.com/ethereum/go-ethereum/log"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// EvmEntropyService - generates the entropy that is injected into the EVM - unique for each transaction
type EvmEntropyService struct {
	sharedSecretService *SharedSecretService
	logger              gethlog.Logger
}

func NewEvmEntropyService(sc *SharedSecretService, logger gethlog.Logger) *EvmEntropyService {
	return &EvmEntropyService{sharedSecretService: sc, logger: logger}
}

// BatchEntropy - calculates entropy per batch
// In Ten, we use a root entropy per batch, which is then used to calculate randomness exposed to individual transactions
// The RootBatchEntropy is calculated based on the shared secret, the batch height and the timestamp
// This ensures that sibling batches will naturally use the same root entropy so that transactions will have the same results
func (ees *EvmEntropyService) BatchEntropy(batch *common.BatchHeader) gethcommon.Hash {
	if !ees.sharedSecretService.IsInitialised() {
		ees.logger.Crit("shared secret service is not initialised")
	}
	extra := batch.Number.Bytes()
	extra = append(extra, big.NewInt(int64(batch.Time)).Bytes()...)
	return gethcommon.BytesToHash(ees.sharedSecretService.ExtendEntropy(extra))
}

// TxEntropy - calculates the randomness exposed to individual transactions
// In TEN, each tx has its own independent randomness,  because otherwise a malicious transaction from the same batch
// could reveal information.
func (ees *EvmEntropyService) TxEntropy(rootBatchEntropy []byte, tCount int) gethcommon.Hash {
	return crypto.Keccak256Hash(rootBatchEntropy, intToBytes(tCount))
}

func intToBytes(val int) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, int64(val))
	if err != nil {
		panic(fmt.Sprintf("Could not convert int to bytes. Cause: %s", err))
	}
	return buf.Bytes()
}
