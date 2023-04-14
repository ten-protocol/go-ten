package services

import (
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/enclave/components"
	"github.com/obscuronet/go-obscuro/go/enclave/core"
)

type ObscuroActor interface {
	//	SubmitTx() //todo
	ReceiveBlock(*common.BlockAndReceipts, bool) (*components.BlockIngestionType, error)
	SubmitTransaction(*common.L2Tx) error
}

type Sequencer interface {
	CreateBatch(*common.L1Block) (*core.Batch, error)
	CreateRollup() (*common.ExtRollup, error)

	ObscuroActor
}

type ObsValidator interface {
	ValidateAndStoreBatch(*core.Batch) error
	ObscuroActor
}
