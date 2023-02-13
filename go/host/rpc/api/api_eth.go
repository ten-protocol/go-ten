package api

import (
	"context"

	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/host"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"
)

// EthereumAPI implements a subset of the Ethereum JSON RPC operations.
// It reuses the RPC API implemented by the Sequencer and allows to add additional endpoints.
type EthereumAPI struct {
	SequencerEthAPI
}

func NewValidatorEthAPI(host host.Host, logger gethlog.Logger) *EthereumAPI {
	return &EthereumAPI{
		SequencerEthAPI{
			host:   host,
			logger: logger,
		},
	}
}

// SendRawTransaction sends the encrypted transaction.
func (api *EthereumAPI) SendRawTransaction(_ context.Context, encryptedParams common.EncryptedParamsSendRawTx) (string, error) {
	encryptedResponse, err := api.host.SubmitAndBroadcastTx(encryptedParams)
	if err != nil {
		return "", err
	}
	return gethcommon.Bytes2Hex(encryptedResponse), nil
}
