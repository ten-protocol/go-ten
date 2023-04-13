package actors

import (
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/enclave/components"
	"github.com/obscuronet/go-obscuro/go/enclave/core"
)

type ObscuroActor interface {
	//	SubmitTx() //todo
	ReceiveBlock(*common.BlockAndReceipts, bool) (*components.BlockIngestionType, error)
}

type Sequencer interface {
	CreateBatch() (*core.Batch, error)
	CreateRollup() (*common.ExtRollup, error)
	IsReady() bool

	ObscuroActor
}

type ObsValidator interface {
	ValidateAndStoreBatch(*core.Batch) error
	ReceiveBlock(*common.BlockAndReceipts, bool) (*components.BlockIngestionType, error)
	GetLatestHead() (*core.Batch, error)

	ObscuroActor
}
